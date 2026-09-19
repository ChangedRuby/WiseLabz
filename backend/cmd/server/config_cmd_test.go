package main

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func setValidEnv(t *testing.T) {
	t.Helper()
	t.Setenv("WISELABZ_AUTH_SECRET", strings.Repeat("a", 40))
	t.Setenv("WISELABZ_ENCRYPTION_KEY", base64.StdEncoding.EncodeToString(make([]byte, 32)))
	t.Setenv("WISELABZ_SERVER_ORIGIN", "http://localhost:5173")
}

func TestConfigValidate(t *testing.T) {
	t.Chdir(t.TempDir())
	var out, errOut bytes.Buffer

	t.Setenv("WISELABZ_AUTH_SECRET", "")
	t.Setenv("WISELABZ_ENCRYPTION_KEY", "")
	t.Setenv("WISELABZ_SERVER_ORIGIN", "")
	if code := runConfigCommand([]string{"validate"}, &out, &errOut); code != 1 {
		t.Fatalf("code = %d, want 1 (stderr: %s)", code, errOut.String())
	}

	setValidEnv(t)
	out.Reset()
	if code := runConfigCommand([]string{"validate"}, &out, &errOut); code != 0 {
		t.Fatalf("code = %d, want 0 (stderr: %s)", code, errOut.String())
	}
}

func TestConfigPrintRedacted(t *testing.T) {
	t.Chdir(t.TempDir())
	setValidEnv(t)
	t.Setenv("WISELABZ_AI_API_KEY", "sk-super-secret")
	var out, errOut bytes.Buffer

	if code := runConfigCommand([]string{"print"}, &out, &errOut); code != 2 {
		t.Fatalf("print without --redacted: code = %d, want 2", code)
	}
	if out.Len() != 0 {
		t.Fatal("print without --redacted wrote to stdout")
	}

	if code := runConfigCommand([]string{"print", "--redacted"}, &out, &errOut); code != 0 {
		t.Fatalf("code = %d (stderr: %s)", code, errOut.String())
	}
	for _, secret := range []string{"sk-super-secret", strings.Repeat("a", 40)} {
		if strings.Contains(out.String(), secret) {
			t.Errorf("output leaks secret %q", secret)
		}
	}
	if !strings.Contains(out.String(), "localhost:5173") {
		t.Errorf("non-secret value missing: %s", out.String())
	}
}

func TestConfigUnknown(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := runConfigCommand(nil, &out, &errOut); code != 2 {
		t.Errorf("no args: code = %d", code)
	}
	if code := runConfigCommand([]string{"bogus"}, &out, &errOut); code != 2 {
		t.Errorf("bogus: code = %d", code)
	}
}
