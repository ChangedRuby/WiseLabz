package netbird

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConnector_Fetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/peers":
			_, _ = w.Write([]byte(`[{"id":"p1","name":"laptop","ip":"100.64.0.1","os":"linux","approval_required":false,"connected":true}]`))
		case "/api/routes":
			_, _ = w.Write([]byte(`[{"id":"r1","network":"10.0.0.0/24","enabled":true}]`))
		case "/api/policies":
			_, _ = w.Write([]byte(`[{"id":"pol1","name":"allow-eng","enabled":true}]`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := &Connector{url: srv.URL, apiToken: "fake-token", client: srv.Client()}

	snap, err := c.Fetch(context.Background(), map[string]any{})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if snap.ServiceName != "Netbird" {
		t.Errorf("ServiceName = %q, want %q", snap.ServiceName, "Netbird")
	}
	if len(snap.Sections) != 3 {
		t.Errorf("len(Sections) = %d, want 3", len(snap.Sections))
	}
	if len(snap.Entities) != 2 {
		t.Errorf("len(Entities) = %d, want 2 (1 peer + 1 policy)", len(snap.Entities))
	}
}

func TestConnector_Validate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := &Connector{url: srv.URL, apiToken: "bad-token", client: srv.Client()}
	if err := c.Validate(context.Background(), map[string]any{}); err == nil {
		t.Fatal("Validate() error = nil, want auth error")
	}
}

func TestConnector_ConfigPush(t *testing.T) {
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && r.URL.Path == "/api/peers/p1" {
			buf := make([]byte, r.ContentLength)
			_, _ = r.Body.Read(buf)
			gotBody = string(buf)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := &Connector{url: srv.URL, apiToken: "fake-token", client: srv.Client()}
	if err := c.ConfigPush(context.Background(), map[string]any{}, "p1", "approved", true); err != nil {
		t.Fatalf("ConfigPush() error = %v", err)
	}
	if gotBody == "" {
		t.Fatal("ConfigPush() did not send a request body")
	}
}
