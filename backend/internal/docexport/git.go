package docexport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"

	"github.com/WiseLabz/wiselabz/internal/httpx"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// GitOptions configures the Git remote target. It mirrors
// config.DocExportGitSettings; scheme validation happens there.
type GitOptions struct {
	CommitMode          string
	AuthorFromUser      bool
	MaxRevisionsPerRun  int
	Remote              string
	Branch              string
	Path                string // subdirectory of the repo the docs are written to
	AuthorName          string
	AuthorEmail         string
	Token               string // HTTPS token, sent as x-access-token basic auth; never logged
	SSHKeyPath          string
	SSHKnownHosts       string
	InsecureSkipHostKey bool
}

// gitTarget holds the resolved Git settings for an Exporter.
type gitTarget struct {
	commitMode         string
	authorFromUser     bool
	maxRevisionsPerRun int
	remote             string
	branch             plumbing.ReferenceName
	path               string // slash-separated, relative to the repo root
	author             object.Signature
	auth               transport.AuthMethod
	insecure           bool
	beforePush         func() // test hook, nil in production
}

var installHTTPSOnce sync.Once

// installHTTPS routes go-git's https transport through the shared hardened
// client (TLS 1.2+, no redirects). No client timeout: a clone or push of a
// large repo is bounded by the job's ctx instead.
func installHTTPS() {
	installHTTPSOnce.Do(func() {
		client.InstallProtocol("https", githttp.NewClient(httpx.NewClient(httpx.Options{Timeout: -1})))
	})
}

// ConfigureGit switches the Exporter to Git mode. It fails closed: an SSH
// remote without a readable key, or without known_hosts when host key
// checking isn't explicitly disabled, is an error.
func (e *Exporter) ConfigureGit(o GitOptions) error {
	if o.Remote == "" {
		return errors.New("git remote must not be empty")
	}
	if o.Branch == "" {
		o.Branch = "main"
	}
	if o.Path == "" {
		o.Path = "docs"
	}
	p := path.Clean(filepath.ToSlash(o.Path))
	if path.IsAbs(p) || p == ".." || strings.HasPrefix(p, "../") {
		return fmt.Errorf("git path %q must be relative to the repository root", o.Path)
	}
	if o.AuthorName == "" {
		o.AuthorName = "WiseLabz"
	}
	if o.AuthorEmail == "" {
		o.AuthorEmail = "wiselabz@localhost"
	}
	if o.CommitMode == "" {
		o.CommitMode = "snapshot"
	}
	if o.CommitMode != "snapshot" && o.CommitMode != "per_revision" {
		return errors.New("git commit mode must be snapshot or per_revision")
	}
	if o.MaxRevisionsPerRun == 0 {
		o.MaxRevisionsPerRun = 500
	}
	if o.MaxRevisionsPerRun < 0 {
		return errors.New("git max revisions per run must be positive")
	}

	auth, err := gitAuth(o)
	if err != nil {
		return err
	}
	e.git = &gitTarget{
		commitMode: o.CommitMode, authorFromUser: o.AuthorFromUser,
		maxRevisionsPerRun: o.MaxRevisionsPerRun,
		remote:             o.Remote,
		branch:             plumbing.NewBranchReferenceName(o.Branch),
		path:               p,
		author:             object.Signature{Name: o.AuthorName, Email: o.AuthorEmail},
		auth:               auth,
		insecure:           o.InsecureSkipHostKey,
	}
	return nil
}

// gitAuth picks the auth method for the remote's protocol: basic auth with
// the token for https, a key file plus host key verification for ssh, and
// none for anything else (file:// in tests).
func gitAuth(o GitOptions) (transport.AuthMethod, error) {
	ep, err := transport.NewEndpoint(o.Remote)
	if err != nil {
		return nil, fmt.Errorf("parse git remote: %w", err)
	}
	switch ep.Protocol {
	case "https":
		installHTTPS()
		if o.Token == "" {
			return nil, nil
		}
		return &githttp.BasicAuth{Username: "x-access-token", Password: o.Token}, nil
	case "ssh":
		if o.SSHKeyPath == "" {
			return nil, errors.New("ssh remote requires an ssh key path")
		}
		user := ep.User
		if user == "" {
			user = "git"
		}
		keys, err := gitssh.NewPublicKeysFromFile(user, o.SSHKeyPath, "")
		if err != nil {
			return nil, fmt.Errorf("load ssh key: %w", err)
		}
		switch {
		case o.InsecureSkipHostKey:
			keys.HostKeyCallback = ssh.InsecureIgnoreHostKey() //nolint:gosec // explicit operator opt-in, warned on every run
		case o.SSHKnownHosts == "":
			return nil, errors.New("ssh remote requires a known_hosts file (or insecure_skip_host_key)")
		default:
			cb, err := gitssh.NewKnownHostsCallback(o.SSHKnownHosts)
			if err != nil {
				return nil, fmt.Errorf("load ssh known_hosts: %w", err)
			}
			keys.HostKeyCallback = cb
		}
		return keys, nil
	default:
		return nil, nil
	}
}

// runGit is one Git-mode export: sync the clone with the remote branch,
// regenerate the docs into it, and commit + push whatever changed.
func (e *Exporter) runGit(ctx context.Context, dir string, logger *slog.Logger) error {
	g := e.git
	if g.insecure {
		logger.Warn("doc export: SSH host key verification is disabled (doc_export.git.insecure_skip_host_key)")
	}

	repo, err := g.open(dir)
	if err != nil {
		return err
	}
	remoteHead, err := g.sync(ctx, repo)
	if err != nil {
		return err
	}

	var commit commitResult
	var count int
	if g.commitMode == "per_revision" {
		commit, count, err = e.runPerRevision(ctx, dir, repo)
	} else {
		commit, count, err = e.runSnapshot(ctx, dir, repo)
	}
	if err != nil {
		return err
	}

	local, err := repo.Reference(g.branch, true)
	if err != nil {
		// Nothing was ever committed (empty remote, no docs): nothing to push.
		logger.Info("doc export: completed, nothing to commit", "count", count)
		return nil
	}
	if local.Hash() == remoteHead {
		logger.Info("doc export: completed, no changes", "count", count)
		return nil
	}

	if g.beforePush != nil {
		g.beforePush()
	}
	err = repo.PushContext(ctx, &git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []gitconfig.RefSpec{gitconfig.RefSpec(g.branch + ":" + g.branch)},
		Auth:       g.auth,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("push to origin %s rejected or failed (never forced; retrying next run): %w", g.branch.Short(), err)
	}
	logger.Info("doc export: pushed", "count", count, "branch", g.branch.Short(),
		"commit", commit.stats(), "head", local.Hash().String())
	return nil
}

func (e *Exporter) runSnapshot(ctx context.Context, dir string, repo *git.Repository) (commitResult, int, error) {
	res, err := e.ExportAll(ctx, filepath.Join(dir, filepath.FromSlash(e.git.path)))
	if err != nil {
		return commitResult{}, 0, err
	}
	commit, err := e.git.commit(repo, res.Count)
	return commit, res.Count, err
}

type exportCursor struct {
	Rev  int    `json:"rev"`
	File string `json:"file"`
}

type exportState struct {
	Version int                     `json:"version"`
	Docs    map[string]exportCursor `json:"docs"`
}

func writeExportState(file string, state exportState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(file, append(data, '\n'), 0o644)
}

func (e *Exporter) runPerRevision(ctx context.Context, dir string, repo *git.Repository) (commitResult, int, error) {
	g := e.git
	exportDir := filepath.Join(dir, filepath.FromSlash(g.path))
	stateFile := filepath.Join(exportDir, ".wiselabz-export.json")
	docs, err := fetchAllDocs(ctx, e.store)
	if err != nil {
		return commitResult{}, 0, fmt.Errorf("fetch docs: %w", err)
	}
	current := make(map[string]store.DocRecord, len(docs))
	for _, d := range docs {
		current[d.ID] = d
	}
	data, err := os.ReadFile(stateFile)
	if errors.Is(err, fs.ErrNotExist) {
		res, err := e.ExportAll(ctx, exportDir)
		if err != nil {
			return commitResult{}, 0, err
		}
		state := exportState{Version: 1, Docs: make(map[string]exportCursor, len(docs))}
		for _, d := range docs {
			state.Docs[d.ID] = exportCursor{Rev: d.CurrentVersion, File: fileName(d)}
		}
		if err := writeExportState(stateFile, state); err != nil {
			return commitResult{}, 0, fmt.Errorf("write export state: %w", err)
		}
		commit, err := g.commit(repo, res.Count)
		return commit, res.Count, err
	}
	if err != nil {
		return commitResult{}, 0, fmt.Errorf("read export state: %w", err)
	}
	var state exportState
	if err := json.Unmarshal(data, &state); err != nil {
		return commitResult{}, 0, fmt.Errorf("parse export state: %w", err)
	}
	if state.Version != 1 || state.Docs == nil {
		return commitResult{}, 0, errors.New("invalid export state: version 1 with docs is required")
	}
	after := make(map[string]int, len(state.Docs))
	for id, cursor := range state.Docs {
		after[id] = cursor.Rev
	}
	versions, err := e.store.ListDocVersionsAfter(ctx, after, g.maxRevisionsPerRun)
	if err != nil {
		return commitResult{}, 0, err
	}
	users := make(map[string]*store.User)
	var total commitResult
	for _, v := range versions {
		d := current[v.DocID]
		name := fileName(d)
		if old := state.Docs[v.DocID].File; old != "" && old != name {
			if !IsGeneratedName(old) || filepath.Base(old) != old {
				return total, len(docs), fmt.Errorf("invalid export state filename %q", old)
			}
			if err := os.Remove(filepath.Join(exportDir, old)); err != nil && !errors.Is(err, fs.ErrNotExist) {
				return total, len(docs), err
			}
		}
		if err := os.WriteFile(filepath.Join(exportDir, name), []byte(v.Content), 0o644); err != nil {
			return total, len(docs), fmt.Errorf("write revision: %w", err)
		}
		state.Docs[v.DocID] = exportCursor{Rev: v.Rev, File: name}
		if err := writeExportState(stateFile, state); err != nil {
			return total, len(docs), err
		}
		when, err := time.Parse(time.RFC3339, v.CreatedAt)
		if err != nil {
			return total, len(docs), fmt.Errorf("parse revision time: %w", err)
		}
		author := g.author
		author.When = when
		if g.authorFromUser && v.Author != "" {
			user, found := users[v.Author]
			if !found {
				user, err = e.store.GetUserByID(ctx, v.Author)
				if err != nil && !errors.Is(err, store.ErrNotFound) {
					return total, len(docs), err
				}
				users[v.Author] = user
			}
			if user != nil && !user.Disabled {
				author.Name = user.DisplayName
				if author.Name == "" {
					author.Name = user.Username
				}
				author.Email = user.Email
				if author.Email == "" {
					author.Email = user.Username + "@users.noreply.wiselabz"
				}
			}
		}
		committer := g.author
		committer.When = time.Now()
		message := fmt.Sprintf("docs(%s): rev %d (%s)", slugify(d.Title), v.Rev, v.Trigger)
		part, err := g.commitWith(repo, 0, message, author, committer)
		if err != nil {
			return total, len(docs), err
		}
		total.added = append(total.added, part.added...)
		total.modified = append(total.modified, part.modified...)
		total.removed = append(total.removed, part.removed...)
	}
	if len(versions) == g.maxRevisionsPerRun {
		return total, len(docs), nil
	}
	res, err := e.ExportAll(ctx, exportDir)
	if err != nil {
		return total, len(docs), err
	}
	for id := range state.Docs {
		if _, ok := current[id]; !ok {
			delete(state.Docs, id)
		}
	}
	for _, d := range docs {
		state.Docs[d.ID] = exportCursor{Rev: d.CurrentVersion, File: fileName(d)}
	}
	if err := writeExportState(stateFile, state); err != nil {
		return total, len(docs), err
	}
	part, err := g.commit(repo, res.Count)
	if err != nil {
		return total, len(docs), err
	}
	total.added = append(total.added, part.added...)
	total.modified = append(total.modified, part.modified...)
	total.removed = append(total.removed, part.removed...)
	return total, len(docs), nil
}

// open returns the persistent clone in dir. An empty (or missing) dir is
// initialised with origin pointing at the remote; a non-empty dir must
// already be a repository whose origin is exactly the configured remote.
func (g *gitTarget) open(dir string) (*git.Repository, error) {
	entries, err := os.ReadDir(dir)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("read export directory: %w", err)
	}
	if len(entries) == 0 {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create export directory: %w", err)
		}
		repo, err := git.PlainInitWithOptions(dir, &git.PlainInitOptions{InitOptions: git.InitOptions{DefaultBranch: g.branch}})
		if err != nil {
			return nil, fmt.Errorf("init clone: %w", err)
		}
		if _, err := repo.CreateRemote(&gitconfig.RemoteConfig{Name: "origin", URLs: []string{g.remote}}); err != nil {
			return nil, fmt.Errorf("add origin: %w", err)
		}
		return repo, nil
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, fmt.Errorf("export directory %s is not empty and not a git repository; refusing to use it: %w", dir, err)
	}
	origin, err := repo.Remote("origin")
	if err != nil {
		return nil, fmt.Errorf("export directory %s has no origin remote; refusing to use it: %w", dir, err)
	}
	if urls := origin.Config().URLs; len(urls) == 0 || urls[0] != g.remote {
		return nil, fmt.Errorf("export directory %s is a clone of a different remote; refusing to use it", dir)
	}
	return repo, nil
}

// sync fetches the branch and hard-resets the clone to origin/<branch>, so
// commits pushed by anyone else are kept and never overwritten. It returns
// the remote branch head, or the zero hash when the remote (or the branch)
// doesn't exist yet and the first push will create it.
func (g *gitTarget) sync(ctx context.Context, repo *git.Repository) (plumbing.Hash, error) {
	if err := repo.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, g.branch)); err != nil {
		return plumbing.ZeroHash, fmt.Errorf("checkout %s: %w", g.branch.Short(), err)
	}

	// Drop any local commit a previous rejected push left behind: the export
	// is regenerated from the store every run, so local history is never
	// the source of truth, and advertising commits the remote doesn't have
	// only confuses fetch negotiation.
	tracking := plumbing.NewRemoteReferenceName("origin", g.branch.Short())
	if ref, err := repo.Reference(tracking, true); err == nil {
		err = repo.Storer.SetReference(plumbing.NewHashReference(g.branch, ref.Hash()))
		if err != nil {
			return plumbing.ZeroHash, fmt.Errorf("move %s: %w", g.branch.Short(), err)
		}
	} else if err := repo.Storer.RemoveReference(g.branch); err != nil {
		return plumbing.ZeroHash, fmt.Errorf("drop unpushed %s: %w", g.branch.Short(), err)
	}

	err := repo.FetchContext(ctx, &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []gitconfig.RefSpec{gitconfig.RefSpec("+" + g.branch + ":" + tracking)},
		Auth:       g.auth,
	})
	switch {
	case err == nil, errors.Is(err, git.NoErrAlreadyUpToDate):
	case errors.Is(err, transport.ErrEmptyRemoteRepository), errors.Is(err, git.NoMatchingRefSpecError{}):
		return plumbing.ZeroHash, nil
	default:
		return plumbing.ZeroHash, fmt.Errorf("fetch origin %s: %w", g.branch.Short(), err)
	}

	ref, err := repo.Reference(tracking, true)
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("resolve %s: %w", tracking, err)
	}
	if err := repo.Storer.SetReference(plumbing.NewHashReference(g.branch, ref.Hash())); err != nil {
		return plumbing.ZeroHash, fmt.Errorf("move %s: %w", g.branch.Short(), err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		return plumbing.ZeroHash, err
	}
	if err := wt.Reset(&git.ResetOptions{Commit: ref.Hash(), Mode: git.HardReset}); err != nil {
		return plumbing.ZeroHash, fmt.Errorf("reset to %s: %w", tracking, err)
	}
	if err := wt.Clean(&git.CleanOptions{Dir: true}); err != nil {
		return plumbing.ZeroHash, fmt.Errorf("clean worktree: %w", err)
	}
	return ref.Hash(), nil
}

// commitResult lists the paths a commit touched, by kind.
type commitResult struct {
	added, modified, removed []string
}

func (c commitResult) stats() string {
	return fmt.Sprintf("+%d ~%d -%d", len(c.added), len(c.modified), len(c.removed))
}

func (c commitResult) empty() bool {
	return len(c.added)+len(c.modified)+len(c.removed) == 0
}

// commit stages every change under the export path (adds, edits and
// removals) and commits it. A clean tree makes no commit.
func (g *gitTarget) commit(repo *git.Repository, count int) (commitResult, error) {
	sig := g.author
	sig.When = time.Now()
	return g.commitWith(repo, count, "", sig, sig)
}

func (g *gitTarget) commitWith(repo *git.Repository, count int, message string, author, committer object.Signature) (commitResult, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return commitResult{}, err
	}
	status, err := wt.Status()
	if err != nil {
		return commitResult{}, fmt.Errorf("worktree status: %w", err)
	}

	var res commitResult
	for p, st := range status {
		if !g.underPath(p) || (st.Worktree == git.Unmodified && st.Staging == git.Unmodified) {
			continue
		}
		switch {
		case st.Worktree == git.Deleted:
			if _, err := wt.Remove(p); err != nil {
				return commitResult{}, fmt.Errorf("stage removal of %s: %w", p, err)
			}
			res.removed = append(res.removed, p)
		case st.Staging == git.Untracked || st.Worktree == git.Untracked || st.Staging == git.Added:
			if _, err := wt.Add(p); err != nil {
				return commitResult{}, fmt.Errorf("stage %s: %w", p, err)
			}
			res.added = append(res.added, p)
		default:
			if _, err := wt.Add(p); err != nil {
				return commitResult{}, fmt.Errorf("stage %s: %w", p, err)
			}
			res.modified = append(res.modified, p)
		}
	}
	if res.empty() {
		return res, nil
	}
	sort.Strings(res.added)
	sort.Strings(res.modified)
	sort.Strings(res.removed)

	if message == "" {
		message = commitMessage(count, res)
	}
	if _, err := wt.Commit(message, &git.CommitOptions{Author: &author, Committer: &committer}); err != nil {
		return commitResult{}, fmt.Errorf("commit: %w", err)
	}
	return res, nil
}

func (g *gitTarget) underPath(p string) bool {
	return g.path == "." || p == g.path || strings.HasPrefix(p, g.path+"/")
}

// commitMessage is "docs: export N docs (+a ~m -r)" with the file lists in
// the body.
func commitMessage(count int, c commitResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "docs: export %d docs (%s)\n", count, c.stats())
	for _, sec := range []struct {
		title string
		files []string
	}{{"Added", c.added}, {"Modified", c.modified}, {"Removed", c.removed}} {
		if len(sec.files) == 0 {
			continue
		}
		fmt.Fprintf(&b, "\n%s:\n", sec.title)
		for _, f := range sec.files {
			fmt.Fprintf(&b, "- %s\n", f)
		}
	}
	return b.String()
}
