// Package connectortest checks connector behavior at common failure boundaries.
package connectortest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// Factory constructs a connector targeting the supplied test server URL.
type Factory func(url string) (connector.Connector, map[string]any, error)

// Run checks both Validate and Fetch against rejected credentials, a deadline,
// and malformed upstream data. A Fetch may return an error or a degraded
// snapshot, as both are supported by the connector contract.
func Run(t *testing.T, factory Factory, opaqueBody bool) {
	t.Helper()
	connector.AllowLoopbackForTest(t)
	var malformed atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if malformed.Load() {
			_, _ = w.Write([]byte(`{malformed`))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid credentials"}`))
	}))
	defer server.Close()
	makeConnector := func() (connector.Connector, map[string]any) {
		t.Helper()
		c, cfg, err := factory(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		return c, cfg
	}
	t.Run("bad credentials", func(t *testing.T) {
		c, cfg := makeConnector()
		var auth *connector.AuthError
		if err := c.Validate(context.Background(), cfg); !errors.As(err, &auth) {
			t.Errorf("Validate error = %v, want AuthError", err)
		}
		assertFailedFetch(context.Background(), t, c, cfg, "auth error")
	})
	t.Run("timeout", func(t *testing.T) {
		c, cfg := makeConnector()
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
		defer cancel()
		var timeout *connector.TimeoutError
		if err := c.Validate(ctx, cfg); !errors.As(err, &timeout) {
			t.Errorf("Validate error = %v, want TimeoutError", err)
		}
		assertFailedFetch(ctx, t, c, cfg, "timeout")
	})
	malformed.Store(true)
	t.Run("malformed body", func(t *testing.T) {
		c, cfg := makeConnector()
		// Some connectors validate reachability only; others parse the body.
		// Both must remain safe on malformed upstream data.
		_ = c.Validate(context.Background(), cfg)
		if opaqueBody {
			// Custom HTTP intentionally preserves arbitrary response text.
			snapshot, err := c.Fetch(context.Background(), cfg)
			if err != nil || snapshot == nil || len(snapshot.Sections) == 0 || !strings.Contains(snapshot.Sections[0].Content, "{malformed") {
				t.Fatalf("Fetch opaque body: snapshot = %+v, error = %v", snapshot, err)
			}
			return
		}
		assertFailedFetch(context.Background(), t, c, cfg, "malformed")
	})
}

func assertFailedFetch(ctx context.Context, t *testing.T, c connector.Connector, cfg map[string]any, failure string) {
	t.Helper()
	snapshot, err := c.Fetch(ctx, cfg)
	if err != nil {
		var auth *connector.AuthError
		var timeout *connector.TimeoutError
		var malformed *connector.MalformedResponseError
		if (failure == "auth error" && errors.As(err, &auth)) ||
			(failure == "timeout" && errors.As(err, &timeout)) ||
			(failure == "malformed" && errors.As(err, &malformed)) {
			return
		}
		t.Fatalf("Fetch error = %v, want %s", err, failure)
	}
	if snapshot == nil {
		t.Fatal("Fetch returned nil snapshot and nil error")
	}
	for _, section := range snapshot.Sections {
		content := strings.ToLower(section.Content)
		if strings.Contains(content, failure) {
			return
		}
	}
	t.Errorf("Fetch returned a healthy snapshot despite upstream failure: %+v", snapshot.Sections)
}
