package connector

import "testing"

func TestValidateRefSegment(t *testing.T) {
	good := []string{"abc123", "a1b2c3d4e5f6", "my-service_1", "host.example.com", "3fa85f64-5717-4562-b3fc-2c963f66afa6"}
	bad := []string{"", ".", "..", "x/../../networks/prune?", "a/b", `a\b`, "a?b", "a#b", "a%2fb", "a b", "a\nb", "a\x00b", "a;b", "é"}
	for _, s := range good {
		if err := ValidateRefSegment(s); err != nil {
			t.Errorf("%q should be valid: %v", s, err)
		}
	}
	for _, s := range bad {
		if err := ValidateRefSegment(s); err == nil {
			t.Errorf("%q should be rejected", s)
		}
	}
}

func TestValidateCompositeRef(t *testing.T) {
	for _, s := range []string{"", "zone/record", "abc"} {
		if err := ValidateCompositeRef(s); err != nil {
			t.Errorf("%q should be valid: %v", s, err)
		}
	}
	for _, s := range []string{"a//b", "a/../b", "/a", "a/", "a/b?x", "a/b#"} {
		if err := ValidateCompositeRef(s); err == nil {
			t.Errorf("%q should be rejected", s)
		}
	}
}

func TestValidateUnixSocketPath(t *testing.T) {
	for _, s := range []string{"/var/run/docker.sock", "/run/user/1000/docker.sock"} {
		if err := ValidateUnixSocketPath(s); err != nil {
			t.Errorf("%q should be valid: %v", s, err)
		}
	}
	for _, s := range []string{"", "docker.sock", "./docker.sock", "/var/run/../docker.sock", "/var//run/docker.sock", "/var/run/docker.sock/", "/var/run/\ndocker.sock"} {
		if err := ValidateUnixSocketPath(s); err == nil {
			t.Errorf("%q should be rejected", s)
		}
	}
}
