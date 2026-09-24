package config

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func validConfig() *Config {
	return &Config{
		DB:         Database{Driver: "sqlite"},
		Server:     Server{Port: 8080, Origin: "http://localhost"},
		Encryption: EncryptionSettings{Key: base64.StdEncoding.EncodeToString(make([]byte, 32))},
		Auth:       AuthSettings{Secret: strings.Repeat("s", 32)},
	}
}

func TestValidate(t *testing.T) {
	if err := validConfig().Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	c := &Config{DB: Database{Driver: "mysql"}}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error")
	}
	for _, want := range []string{"AUTH_SECRET", "ENCRYPTION_KEY", "SERVER_ORIGIN", "db.driver", "server.port"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error missing %q: %v", want, err)
		}
	}
}

func TestRedacted(t *testing.T) {
	c := validConfig()
	c.DB.DSN = "postgres://wl:hunter2@db:5432/wl"
	c.AI.APIKey = "sk-ai"
	c.AI.EmbedAPIKey = "sk-embed"
	c.DocExport.Git.Token = "ghp-doc-token"
	c.DocExport.Git.Remote = "https://bot:remote-pass@git.example.com/o/r.git"
	c.Auth.OIDC = []OIDCProvider{{ID: "x", ClientSecret: "oidc-secret"}}

	r := c.Redacted()
	out, _ := json.Marshal(r)
	for _, secret := range []string{"hunter2", "sk-ai", "sk-embed", "oidc-secret", "ghp-doc-token", "remote-pass", c.Auth.Secret, c.Encryption.Key} {
		if strings.Contains(string(out), secret) {
			t.Errorf("redacted output leaks %q", secret)
		}
	}
	if r.AI.BaseURL != "" || r.AI.EmbedAPIKey == "" {
		t.Errorf("empty stays empty, set becomes masked: %+v", r.AI)
	}
	if c.Auth.OIDC[0].ClientSecret != "oidc-secret" || c.AI.APIKey != "sk-ai" {
		t.Error("Redacted mutated the original")
	}
}

func TestRedactDSN(t *testing.T) {
	cases := map[string]string{
		"file:/data/wl.db?cache=shared":              "file:/data/wl.db?cache=shared",
		"host=db user=wl password=hunter2 dbname=wl": "host=db user=wl password=" + redactedValue + " dbname=wl",
		"postgres://wl@db/wl":                        "postgres://wl@db/wl",
	}
	for in, want := range cases {
		if got := redactDSN(in); got != want {
			t.Errorf("redactDSN(%q) = %q, want %q", in, got, want)
		}
	}
	if got := redactDSN("postgres://wl:hunter2@db/wl"); strings.Contains(got, "hunter2") {
		t.Errorf("leaked: %s", got)
	}
}

// TestEveryKeyEnvOverridable ensures each scalar Config field (including
// fields of nested sections such as doc_export.git) can be set via
// WISELABZ_<SECTION>_<KEY>, so a new field can't silently miss its BindEnv.
func TestEveryKeyEnvOverridable(t *testing.T) {
	var setEnv func(prefix string, rt reflect.Type)
	setEnv = func(prefix string, rt reflect.Type) {
		for j := 0; j < rt.NumField(); j++ {
			f := rt.Field(j)
			name := prefix + "_" + f.Tag.Get("mapstructure")
			switch {
			case f.Type.Kind() == reflect.Struct:
				setEnv(name, f.Type)
			case f.Type.Kind() == reflect.Slice: // OIDC providers: file-only
			case f.Type.Kind() == reflect.Int:
				t.Setenv("WISELABZ_"+strings.ToUpper(name), "7")
			case f.Type.Kind() == reflect.Bool:
				t.Setenv("WISELABZ_"+strings.ToUpper(name), "true")
			case strings.HasSuffix(name, "cron_expr") || name == "sync_schedule":
				t.Setenv("WISELABZ_"+strings.ToUpper(name), "*/5 * * * *")
			case name == "server_public_url":
				t.Setenv("WISELABZ_"+strings.ToUpper(name), "https://example.test")
			default:
				t.Setenv("WISELABZ_"+strings.ToUpper(name), "v")
			}
		}
	}
	rt := reflect.TypeOf(Config{})
	for i := 0; i < rt.NumField(); i++ {
		sec := rt.Field(i)
		setEnv(sec.Tag.Get("mapstructure"), sec.Type)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	var check func(path string, v reflect.Value)
	check = func(path string, v reflect.Value) {
		for j := 0; j < v.NumField(); j++ {
			f := v.Field(j)
			name := path + "." + v.Type().Field(j).Name
			switch {
			case f.Kind() == reflect.Struct:
				check(name, f)
			case f.Kind() == reflect.Slice:
			case f.IsZero():
				t.Errorf("%s not set by its WISELABZ_ env var", name)
			}
		}
	}
	v := reflect.ValueOf(*cfg)
	for i := 0; i < v.NumField(); i++ {
		check(v.Type().Field(i).Name, v.Field(i))
	}
}

func TestValidateRejectsInvalidDocExportGit(t *testing.T) {
	c := validConfig()
	c.DocExport.Git = DocExportGitSettings{Remote: "http://git.example.com/o/r.git", Branch: "main", Path: "docs"}
	err := c.Validate()
	if err == nil || !strings.Contains(err.Error(), "doc_export.git.remote") {
		t.Fatalf("Validate() error = %v, want doc_export.git.remote error", err)
	}
}
