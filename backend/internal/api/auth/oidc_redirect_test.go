package auth

import (
	"net/http/httptest"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/config"
)

func TestOIDCRedirectURL(t *testing.T) {
	tests := []struct {
		name, origin, host, want string
	}{
		{"match among many", "https://a.example.com, https://b.example.com", "B.example.com", "https://b.example.com/auth/callback"},
		{"single origin overrides evil host", "https://a.example.com", "evil.com", "https://a.example.com/auth/callback"},
		{"multi no match valid host fallback", "https://a.example.com,https://b.example.com", "c.example.com:8080", "http://c.example.com:8080/auth/callback"},
		{"no origin valid host", "", "localhost:3000", "http://localhost:3000/auth/callback"},
		{"no origin invalid host", "", "evil.com/x@y", "http://localhost/auth/callback"},
		{"no origin userinfo", "", "user@evil.com", "http://localhost/auth/callback"},
		{"no origin space", "", "a b", "http://localhost/auth/callback"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{Config: &config.Config{}}
			h.Config.Server.Origin = tt.origin
			r := httptest.NewRequest("GET", "/", nil)
			r.Host = tt.host
			if got := h.oidcRedirectURL(r); got != tt.want {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}
