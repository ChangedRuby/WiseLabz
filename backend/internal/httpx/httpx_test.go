package httpx

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewClientTimeout(t *testing.T) {
	tests := []struct {
		name          string
		timeout       time.Duration
		wantClient    time.Duration
		wantRespHdrTO time.Duration
	}{
		{name: "zero uses default", timeout: 0, wantClient: DefaultTimeout, wantRespHdrTO: DefaultTimeout},
		{name: "custom timeout kept", timeout: 60 * time.Second, wantClient: 60 * time.Second, wantRespHdrTO: 60 * time.Second},
		{name: "negative disables client timeout", timeout: -1, wantClient: 0, wantRespHdrTO: defaultResponseHeaderTimeout},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(Options{Timeout: tt.timeout})
			if c.Timeout != tt.wantClient {
				t.Errorf("client timeout = %v, want %v", c.Timeout, tt.wantClient)
			}
			tr, ok := c.Transport.(*http.Transport)
			if !ok {
				t.Fatalf("transport type = %T, want *http.Transport", c.Transport)
			}
			if tr.ResponseHeaderTimeout != tt.wantRespHdrTO {
				t.Errorf("response header timeout = %v, want %v", tr.ResponseHeaderTimeout, tt.wantRespHdrTO)
			}
		})
	}
}

func TestNewTransportTLS(t *testing.T) {
	tests := []struct {
		name         string
		insecure     bool
		wantInsecure bool
	}{
		{name: "verify by default", insecure: false, wantInsecure: false},
		{name: "insecure propagates", insecure: true, wantInsecure: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTransport(Options{InsecureSkipVerify: tt.insecure})
			if tr.TLSClientConfig == nil {
				t.Fatal("TLSClientConfig is nil")
			}
			if tr.TLSClientConfig.MinVersion != tls.VersionTLS12 {
				t.Errorf("MinVersion = %x, want %x", tr.TLSClientConfig.MinVersion, tls.VersionTLS12)
			}
			if tr.TLSClientConfig.InsecureSkipVerify != tt.wantInsecure {
				t.Errorf("InsecureSkipVerify = %v, want %v", tr.TLSClientConfig.InsecureSkipVerify, tt.wantInsecure)
			}
		})
	}
}

func TestNewClientInsecureSkipVerifyConnects(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	tests := []struct {
		name     string
		insecure bool
		wantErr  bool
	}{
		{name: "self-signed rejected", insecure: false, wantErr: true},
		{name: "self-signed accepted when insecure", insecure: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(Options{InsecureSkipVerify: tt.insecure, Timeout: 5 * time.Second})
			resp, err := c.Get(srv.URL)
			if resp != nil {
				_ = resp.Body.Close()
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClientDoesNotFollowRedirects(t *testing.T) {
	var targetHits atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		targetHits.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()
	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer redirector.Close()

	resp, err := NewClient(Options{Timeout: 5 * time.Second}).Get(redirector.URL)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusFound)
	}
	if got := targetHits.Load(); got != 0 {
		t.Errorf("redirect target hit %d times, want 0", got)
	}
}

func TestNewClientUsesCustomDialContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	errBlocked := errors.New("blocked by test dialer")
	tests := []struct {
		name    string
		block   bool
		wantErr bool
	}{
		{name: "dialer allows", block: false, wantErr: false},
		{name: "dialer blocks", block: true, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls atomic.Int32
			dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
				calls.Add(1)
				if tt.block {
					return nil, errBlocked
				}
				var d net.Dialer
				return d.DialContext(ctx, network, addr)
			}
			c := NewClient(Options{Timeout: 5 * time.Second, DialContext: dial})
			if tr, ok := c.Transport.(*http.Transport); !ok || tr.Proxy != nil {
				t.Error("custom dialer must disable proxy so the dialer sees the real target")
			}
			resp, err := c.Get(srv.URL)
			if resp != nil {
				_ = resp.Body.Close()
			}
			if calls.Load() == 0 {
				t.Error("custom DialContext was not called")
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, errBlocked) {
				t.Errorf("err = %v, want wrapping %v", err, errBlocked)
			}
		})
	}
}
