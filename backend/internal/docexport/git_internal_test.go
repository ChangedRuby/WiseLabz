package docexport

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"

	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

// SetBeforePushForTest installs a hook that runs right before the push, so
// the external tests can make the remote move and the push get rejected.
func SetBeforePushForTest(e *Exporter, f func()) { e.git.beforePush = f }

func writeTestKey(t *testing.T) string {
	t.Helper()
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	block, err := ssh.MarshalPrivateKey(priv, "")
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "id_ed25519")
	if err := os.WriteFile(p, pem.EncodeToMemory(block), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestGitAuthHTTPSToken(t *testing.T) {
	auth, err := gitAuth(GitOptions{Remote: "https://git.example.com/o/r.git", Token: "s3cr3t"})
	if err != nil {
		t.Fatal(err)
	}
	basic, ok := auth.(*githttp.BasicAuth)
	if !ok || basic.Username != "x-access-token" || basic.Password != "s3cr3t" {
		t.Fatalf("auth = %#v, want BasicAuth x-access-token", auth)
	}
	if strings.Contains(basic.String(), "s3cr3t") {
		t.Error("BasicAuth.String() leaks the token")
	}
}

func TestGitAuthHTTPSNoToken(t *testing.T) {
	auth, err := gitAuth(GitOptions{Remote: "https://git.example.com/o/r.git"})
	if err != nil || auth != nil {
		t.Fatalf("auth, err = %v, %v; want nil, nil", auth, err)
	}
}

func TestGitAuthSSH(t *testing.T) {
	key := writeTestKey(t)
	knownHosts := filepath.Join(t.TempDir(), "known_hosts")
	if err := os.WriteFile(knownHosts, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("known_hosts", func(t *testing.T) {
		auth, err := gitAuth(GitOptions{Remote: "ssh://deploy@git.example.com/o/r.git", SSHKeyPath: key, SSHKnownHosts: knownHosts})
		if err != nil {
			t.Fatal(err)
		}
		keys, ok := auth.(*gitssh.PublicKeys)
		if !ok || keys.User != "deploy" || keys.HostKeyCallback == nil {
			t.Fatalf("auth = %#v, want PublicKeys for deploy with a host key callback", auth)
		}
	})
	t.Run("scp-style defaults user to git", func(t *testing.T) {
		auth, err := gitAuth(GitOptions{Remote: "git@git.example.com:o/r.git", SSHKeyPath: key, SSHKnownHosts: knownHosts})
		if err != nil {
			t.Fatal(err)
		}
		if keys := auth.(*gitssh.PublicKeys); keys.User != "git" {
			t.Fatalf("user = %q, want git", keys.User)
		}
	})
	t.Run("fails closed without known_hosts", func(t *testing.T) {
		if _, err := gitAuth(GitOptions{Remote: "ssh://git@git.example.com/o/r.git", SSHKeyPath: key}); err == nil {
			t.Fatal("expected error without known_hosts")
		}
	})
	t.Run("fails closed on unreadable known_hosts", func(t *testing.T) {
		if _, err := gitAuth(GitOptions{Remote: "ssh://git@git.example.com/o/r.git", SSHKeyPath: key, SSHKnownHosts: "/nonexistent/known_hosts"}); err == nil {
			t.Fatal("expected error for missing known_hosts file")
		}
	})
	t.Run("insecure opt-in", func(t *testing.T) {
		auth, err := gitAuth(GitOptions{Remote: "ssh://git@git.example.com/o/r.git", SSHKeyPath: key, InsecureSkipHostKey: true})
		if err != nil {
			t.Fatal(err)
		}
		if auth.(*gitssh.PublicKeys).HostKeyCallback == nil {
			t.Fatal("insecure mode should still install a (permissive) callback")
		}
	})
	t.Run("missing key", func(t *testing.T) {
		if _, err := gitAuth(GitOptions{Remote: "ssh://git@git.example.com/o/r.git", SSHKnownHosts: knownHosts}); err == nil {
			t.Fatal("expected error without key")
		}
	})
}

func TestCommitMessage(t *testing.T) {
	got := commitMessage(3, commitResult{added: []string{"docs/a-0000000a.md"}, removed: []string{"docs/b-0000000b.md"}})
	want := "docs: export 3 docs (+1 ~0 -1)\n\nAdded:\n- docs/a-0000000a.md\n\nRemoved:\n- docs/b-0000000b.md\n"
	if got != want {
		t.Fatalf("commitMessage =\n%q\nwant\n%q", got, want)
	}
}
