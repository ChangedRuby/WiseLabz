package notifications

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/crypto"
)

func TestSendWebhook_BlocksLoopbackAndLinkLocal(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { hits.Add(1) }))
	defer srv.Close()

	for _, target := range []string{srv.URL, "http://169.254.169.254/latest/meta-data", "http://0.0.0.0:9/"} {
		err := sendWebhook(context.Background(), target, "", map[string]string{"a": "b"})
		if err == nil || !strings.Contains(err.Error(), "blocked address") {
			t.Errorf("target %s: expected blocked-address error, got %v", target, err)
		}
	}
	if hits.Load() != 0 {
		t.Errorf("loopback server received %d requests, want 0", hits.Load())
	}
}

func TestSendWebhook_ErrorOmitsURL(t *testing.T) {
	err := sendWebhook(context.Background(), "http://127.0.0.1:1/hook?token=supersecret", "", "x")
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "supersecret") {
		t.Errorf("error leaks URL token: %v", err)
	}
}

func TestSendWebhook_DoesNotFollowRedirects(t *testing.T) {
	connector.AllowLoopbackForTest(t)
	var targetHits atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { targetHits.Add(1) }))
	defer target.Close()
	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer redirector.Close()

	err := sendWebhook(context.Background(), redirector.URL, "", "x")
	if err == nil || !strings.Contains(err.Error(), "status 302") {
		t.Errorf("expected status 302 failure, got %v", err)
	}
	if targetHits.Load() != 0 {
		t.Errorf("redirect target was contacted %d times", targetHits.Load())
	}
}

func TestSendWebhook_BoundsResponseRead(t *testing.T) {
	connector.AllowLoopbackForTest(t)
	const chunks = 512
	var written atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		chunk := make([]byte, 1<<20)
		for i := 0; i < chunks; i++ {
			if _, err := w.Write(chunk); err != nil {
				return
			}
			written.Add(1)
		}
	}))
	defer srv.Close()

	if err := sendWebhook(context.Background(), srv.URL, "", "x"); err != nil {
		t.Fatalf("send: %v", err)
	}
	srv.CloseClientConnections()
	if got := written.Load(); got >= chunks {
		t.Errorf("server wrote the full %d MiB body; client should stop reading early", got)
	}
}

func TestSendWebhook_SignatureVerifies(t *testing.T) {
	connector.AllowLoopbackForTest(t)
	const secret = "shared-secret"
	var gotSig, gotTS string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig, gotTS = r.Header.Get(signatureHeader), r.Header.Get(timestampHeader)
		gotBody, _ = io.ReadAll(r.Body)
	}))
	defer srv.Close()

	if err := sendWebhook(context.Background(), srv.URL, secret, map[string]string{"title": "t"}); err != nil {
		t.Fatalf("send: %v", err)
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(gotTS + "."))
	mac.Write(gotBody)
	if want := "sha256=" + hex.EncodeToString(mac.Sum(nil)); gotSig != want {
		t.Errorf("signature = %q, want %q", gotSig, want)
	}
}

func TestSendWebhook_UnsignedWithoutSecret(t *testing.T) {
	connector.AllowLoopbackForTest(t)
	var sig, ts string
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		sig, ts = r.Header.Get(signatureHeader), r.Header.Get(timestampHeader)
	}))
	defer srv.Close()
	if err := sendWebhook(context.Background(), srv.URL, "", "x"); err != nil {
		t.Fatal(err)
	}
	if sig != "" || ts != "" {
		t.Errorf("expected no signature headers, got %q %q", sig, ts)
	}
}

func TestDispatcher_SignsWithEncryptedChannelSecret(t *testing.T) {
	s := newTestStore(t)
	var gotSig, gotTS string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		gotSig, gotTS = r.Header.Get(signatureHeader), r.Header.Get(timestampHeader)
		gotBody, _ = io.ReadAll(r.Body)
	}))
	defer srv.Close()

	rawKey := []byte(strings.Repeat("k", 32))
	enc, err := crypto.Encrypt("known-secret", rawKey)
	if err != nil {
		t.Fatal(err)
	}
	cfgJSON := `{"channels":[{"type":"webhook","enabled":true,"config":{"url":"` + srv.URL + `","secretEncrypted":"` + enc + `"}}]}`
	if _, err := s.DB().ExecContext(context.Background(),
		`INSERT INTO notification_config (id, config_json) VALUES (1, ?) ON CONFLICT(id) DO UPDATE SET config_json = excluded.config_json`, cfgJSON); err != nil {
		t.Fatal(err)
	}

	d := NewDispatcher(s, nil)
	d.SetEncryptionKey(base64.StdEncoding.EncodeToString(rawKey))
	d.NotifyAlert("alert-1", "user-1", "alert.created", "Title", "Message")

	if gotSig == "" {
		t.Fatal("expected signed delivery")
	}
	if want := signWebhook("known-secret", gotTS, gotBody); gotSig != want {
		t.Errorf("signature = %q, want %q", gotSig, want)
	}
}
