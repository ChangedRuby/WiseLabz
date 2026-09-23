package docexport_test

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"github.com/go-git/go-git/v5/plumbing/transport/server"

	"github.com/WiseLabz/wiselabz/internal/docexport"
	"github.com/WiseLabz/wiselabz/internal/notifications"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Serve file:// remotes in-process so the tests need neither network nor a
// git binary.
func init() { client.InstallProtocol("file", server.DefaultServer) }

type notifyCall struct{ eventType, severity, title string }

type fakeNotifier struct {
	mu    sync.Mutex
	calls []notifyCall
}

func (f *fakeNotifier) NotifySystemEvent(_ context.Context, eventType, severity, title, _ string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, notifyCall{eventType, severity, title})
}

func (f *fakeNotifier) snapshot() []notifyCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]notifyCall(nil), f.calls...)
}

type gitFixture struct {
	t        *testing.T
	ctx      context.Context
	store    *store.Store
	exporter *docexport.Exporter
	notifier *fakeNotifier
	bare     string
	remote   string
	workdir  string
	logger   *slog.Logger
}

func newGitFixture(t *testing.T) *gitFixture {
	t.Helper()
	bare := t.TempDir()
	if _, err := git.PlainInit(bare, true); err != nil {
		t.Fatalf("init bare remote: %v", err)
	}
	f := &gitFixture{
		t:        t,
		ctx:      context.Background(),
		store:    newTestStore(t),
		notifier: &fakeNotifier{},
		bare:     bare,
		remote:   "file://" + bare,
		workdir:  filepath.Join(t.TempDir(), "clone"),
		logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	f.exporter = docexport.NewExporter(f.store)
	if err := f.exporter.ConfigureGit(docexport.GitOptions{
		Remote: f.remote, Branch: "main", Path: "docs",
		AuthorName: "Doc Bot", AuthorEmail: "bot@example.com",
	}); err != nil {
		t.Fatalf("ConfigureGit: %v", err)
	}
	f.exporter.SetNotifier(f.notifier)
	return f
}

func (f *gitFixture) run() { docexport.RunExportOnce(f.ctx, f.exporter, f.workdir, f.logger) }

func (f *gitFixture) createDoc(id, title, content string) {
	f.t.Helper()
	if err := f.store.CreateDoc(f.ctx, &store.DocRecord{ID: id, Title: title, Content: content}); err != nil {
		f.t.Fatalf("create doc: %v", err)
	}
}

// head returns the remote's main branch head commit.
func (f *gitFixture) head() *object.Commit {
	f.t.Helper()
	repo, err := git.PlainOpen(f.bare)
	if err != nil {
		f.t.Fatalf("open bare: %v", err)
	}
	ref, err := repo.Reference(plumbing.NewBranchReferenceName("main"), true)
	if err != nil {
		f.t.Fatalf("remote main: %v", err)
	}
	c, err := repo.CommitObject(ref.Hash())
	if err != nil {
		f.t.Fatalf("head commit: %v", err)
	}
	return c
}

func (f *gitFixture) commitCount() int {
	f.t.Helper()
	n := 0
	iter := object.NewCommitPreorderIter(f.head(), nil, nil)
	if err := iter.ForEach(func(*object.Commit) error { n++; return nil }); err != nil {
		f.t.Fatal(err)
	}
	return n
}

// files returns path → content for every file on the remote's main branch.
func (f *gitFixture) files() map[string]string {
	f.t.Helper()
	tree, err := f.head().Tree()
	if err != nil {
		f.t.Fatal(err)
	}
	out := map[string]string{}
	if err := tree.Files().ForEach(func(file *object.File) error {
		c, err := file.Contents()
		out[file.Name] = c
		return err
	}); err != nil {
		f.t.Fatal(err)
	}
	return out
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// pushExternal commits file=content to the remote from an unrelated clone,
// as a human or another tool would.
func (f *gitFixture) pushExternal(file, content string) {
	f.t.Helper()
	dir := f.t.TempDir()
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{URL: f.remote, ReferenceName: plumbing.NewBranchReferenceName("main")})
	if err != nil {
		f.t.Fatalf("external clone: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, file), []byte(content), 0o644); err != nil {
		f.t.Fatal(err)
	}
	wt, _ := repo.Worktree()
	if _, err := wt.Add(file); err != nil {
		f.t.Fatal(err)
	}
	sig := &object.Signature{Name: "Human", Email: "human@example.com", When: time.Now()}
	if _, err := wt.Commit("external edit", &git.CommitOptions{Author: sig}); err != nil {
		f.t.Fatal(err)
	}
	if err := repo.Push(&git.PushOptions{}); err != nil {
		f.t.Fatalf("external push: %v", err)
	}
}

func TestGitExportLifecycle(t *testing.T) {
	f := newGitFixture(t)
	f.createDoc(id1, "Runbook", "# Runbook v1")
	f.createDoc(id2, "Lab Topology", "# Topology")

	// First run against an empty remote creates the branch.
	f.run()
	if n := f.commitCount(); n != 1 {
		t.Fatalf("after first run: %d commits, want 1", n)
	}
	head := f.head()
	if !strings.HasPrefix(head.Message, "docs: export 2 docs (+2 ~0 -0)\n") {
		t.Errorf("first commit message = %q", head.Message)
	}
	if !strings.Contains(head.Message, "- docs/runbook-0000000a.md") {
		t.Errorf("commit body missing file list: %q", head.Message)
	}
	if head.Author.Name != "Doc Bot" || head.Author.Email != "bot@example.com" {
		t.Errorf("author = %s <%s>", head.Author.Name, head.Author.Email)
	}
	want := map[string]string{
		"docs/runbook-0000000a.md":      "# Runbook v1",
		"docs/lab-topology-0000000b.md": "# Topology",
	}
	if got := f.files(); len(got) != 2 || got["docs/runbook-0000000a.md"] != want["docs/runbook-0000000a.md"] || got["docs/lab-topology-0000000b.md"] != want["docs/lab-topology-0000000b.md"] {
		t.Fatalf("remote files = %v, want %v", got, want)
	}

	// No doc changes: no new commit.
	f.run()
	if n := f.commitCount(); n != 1 {
		t.Fatalf("after unchanged run: %d commits, want 1", n)
	}

	// An edit makes exactly one commit.
	if err := f.store.UpdateDoc(f.ctx, id1, "# Runbook v2", nil); err != nil {
		t.Fatal(err)
	}
	f.run()
	if n := f.commitCount(); n != 2 {
		t.Fatalf("after edit: %d commits, want 2", n)
	}
	if msg := f.head().Message; !strings.HasPrefix(msg, "docs: export 2 docs (+0 ~1 -0)\n") {
		t.Errorf("edit commit message = %q", msg)
	}
	if got := f.files()["docs/runbook-0000000a.md"]; got != "# Runbook v2" {
		t.Errorf("edited content = %q", got)
	}

	// A human commits a README (inside the export path) and a root file.
	f.pushExternal("README.md", "hello")
	f.pushExternal("docs/README.md", "about these docs")

	// A deleted doc removes its file; the external commits survive.
	if err := f.store.DeleteDoc(f.ctx, id2); err != nil {
		t.Fatal(err)
	}
	f.run()
	if n := f.commitCount(); n != 5 {
		t.Fatalf("after delete: %d commits, want 5", n)
	}
	if msg := f.head().Message; !strings.HasPrefix(msg, "docs: export 1 docs (+0 ~0 -1)\n") {
		t.Errorf("delete commit message = %q", msg)
	}
	got := f.files()
	wantKeys := []string{"README.md", "docs/README.md", "docs/runbook-0000000a.md"}
	if strings.Join(keys(got), ",") != strings.Join(wantKeys, ",") {
		t.Fatalf("remote files = %v, want %v", keys(got), wantKeys)
	}

	if calls := f.notifier.snapshot(); len(calls) != 0 {
		t.Errorf("unexpected notifications: %+v", calls)
	}
}

func TestGitExportPushRejectionNotifiesOnTransitions(t *testing.T) {
	f := newGitFixture(t)
	f.createDoc(id1, "Runbook", "v1")
	f.run() // seed the remote
	if n := f.commitCount(); n != 1 {
		t.Fatalf("seed: %d commits, want 1", n)
	}

	// Make the remote move between our fetch and our push, twice in a row,
	// so both pushes are rejected as non-fast-forward.
	external := 0
	docexport.SetBeforePushForTest(f.exporter, func() {
		external++
		f.pushExternal("external.txt", strings.Repeat("x", external))
	})
	for i, content := range []string{"v2", "v3"} {
		if err := f.store.UpdateDoc(f.ctx, id1, content, nil); err != nil {
			t.Fatal(err)
		}
		f.run()
		if msg := f.head().Message; msg != "external edit" {
			t.Fatalf("run %d: remote head = %q, want the external commit (push must not be forced)", i, msg)
		}
	}
	calls := f.notifier.snapshot()
	if len(calls) != 1 || calls[0] != (notifyCall{"system.job_failed", "warning", "Doc export failing"}) {
		t.Fatalf("after two failures: notifications = %+v, want exactly one warning", calls)
	}

	// Next run succeeds on top of the external commits and notifies recovery once.
	docexport.SetBeforePushForTest(f.exporter, nil)
	f.run()
	f.run()
	calls = f.notifier.snapshot()
	if len(calls) != 2 || calls[1] != (notifyCall{"system.job_failed", "info", "Doc export recovered"}) {
		t.Fatalf("after recovery: notifications = %+v, want warning then one info", calls)
	}
	got := f.files()
	if got["docs/runbook-0000000a.md"] != "v3" || got["external.txt"] != "xx" {
		t.Fatalf("remote files after recovery = %v", got)
	}
}

func TestGitExportRefusesForeignDirectory(t *testing.T) {
	t.Run("not a repository", func(t *testing.T) {
		f := newGitFixture(t)
		if err := os.MkdirAll(f.workdir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(f.workdir, "important.txt"), []byte("mine"), 0o644); err != nil {
			t.Fatal(err)
		}
		f.run()
		if calls := f.notifier.snapshot(); len(calls) != 1 || calls[0].severity != "warning" {
			t.Fatalf("notifications = %+v, want one warning", calls)
		}
		if data, err := os.ReadFile(filepath.Join(f.workdir, "important.txt")); err != nil || string(data) != "mine" {
			t.Fatalf("operator file touched: %q, %v", data, err)
		}
	})
	t.Run("clone of another remote", func(t *testing.T) {
		f := newGitFixture(t)
		other := t.TempDir()
		if _, err := git.PlainInit(other, true); err != nil {
			t.Fatal(err)
		}
		repo, err := git.PlainInit(f.workdir, false)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := repo.CreateRemote(&gitconfig.RemoteConfig{Name: "origin", URLs: []string{"file://" + other}}); err != nil {
			t.Fatal(err)
		}
		f.createDoc(id1, "Runbook", "v1")
		f.run()
		if calls := f.notifier.snapshot(); len(calls) != 1 || calls[0].severity != "warning" {
			t.Fatalf("notifications = %+v, want one warning", calls)
		}
		if _, err := os.Stat(filepath.Join(f.workdir, "docs")); !os.IsNotExist(err) {
			t.Fatalf("export ran into a foreign clone: %v", err)
		}
	})
}

// The exporter talks to the real dispatcher through Notifier, and both
// packages must agree on the event name.
var _ docexport.Notifier = (*notifications.Dispatcher)(nil)

func TestEventJobFailedMatchesDispatcher(t *testing.T) {
	if docexport.EventJobFailed != notifications.EventSystemJobFailed {
		t.Fatalf("docexport.EventJobFailed = %q, notifications.EventSystemJobFailed = %q", docexport.EventJobFailed, notifications.EventSystemJobFailed)
	}
}
