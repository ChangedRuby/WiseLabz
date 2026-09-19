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
	c.Auth.OIDC = []OIDCProvider{{ID: "x", ClientSecret: "oidc-secret"}}

	r := c.Redacted()
	out, _ := json.Marshal(r)
	for _, secret := range []string{"hunter2", "sk-ai", "sk-embed", "oidc-secret", c.Auth.Secret, c.Encryption.Key} {
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

// TestEveryKeyEnvOverridable ensures each scalar Config field can be set via
// WISELABZ_<SECTION>_<KEY>, so a new field can't silently miss its BindEnv.
func TestEveryKeyEnvOverridable(t *testing.T) {
	rt := reflect.TypeOf(Config{})
	for i := 0; i < rt.NumField(); i++ {
		sec := rt.Field(i)
		for j := 0; j < sec.Type.NumField(); j++ {
			f := sec.Type.Field(j)
			name := sec.Tag.Get("mapstructure") + "_" + f.Tag.Get("mapstructure")
			switch {
			case f.Type.Kind() == reflect.Slice: // OIDC providers: file-only
			case f.Type.Kind() == reflect.Int:
				t.Setenv("WISELABZ_"+strings.ToUpper(name), "7")
			case f.Type.Kind() == reflect.Bool:
				t.Setenv("WISELABZ_"+strings.ToUpper(name), "true")
			case strings.HasSuffix(name, "cron_expr") || name == "sync_schedule":
				t.Setenv("WISELABZ_"+strings.ToUpper(name), "*/5 * * * *")
			default:
				t.Setenv("WISELABZ_"+strings.ToUpper(name), "v")
			}
		}
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	v := reflect.ValueOf(*cfg)
	for i := 0; i < v.NumField(); i++ {
		sec := v.Field(i)
		for j := 0; j < sec.NumField(); j++ {
			f := sec.Field(j)
			if f.Kind() == reflect.Slice {
				continue
			}
			if f.IsZero() {
				t.Errorf("%s.%s not set by its WISELABZ_ env var", v.Type().Field(i).Name, sec.Type().Field(j).Name)
			}
		}
	}
}
