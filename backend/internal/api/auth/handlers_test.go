package auth

import "testing"

func TestEmailDomainAllowed(t *testing.T) {
	if !emailDomainAllowed("Alice@Example.COM", []string{"  @example.com  "}) {
		t.Error("emailDomainAllowed() rejected case-insensitive matching domain")
	}
	for _, email := range []string{"", "alice", "alice@", "@example.com", "alice@other.com"} {
		if emailDomainAllowed(email, []string{"example.com"}) {
			t.Errorf("emailDomainAllowed(%q) = true, want false", email)
		}
	}
}

func TestOIDCRoleForGroups(t *testing.T) {
	mapping := map[string]string{"readers": "viewer", "admins": "operator", "invalid": "unknown"}
	for _, tc := range []struct {
		name   string
		groups []string
		want   string
	}{
		{name: "default", want: "user"},
		{name: "viewer", groups: []string{"readers"}, want: "user"},
		{name: "operator takes precedence", groups: []string{"readers", "admins"}, want: "admin"},
		{name: "unknown ignored", groups: []string{"invalid"}, want: "user"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := oidcRoleForGroups(tc.groups, mapping); got != tc.want {
				t.Fatalf("oidcRoleForGroups(%v) = %q, want %q", tc.groups, got, tc.want)
			}
		})
	}
}
