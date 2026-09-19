package config

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/WiseLabz/wiselabz/internal/crypto"
)

// MinAuthSecretLen is the minimum length of the JWT signing secret.
const MinAuthSecretLen = 32

// redactedValue replaces secret values in printed configuration.
const redactedValue = "***REDACTED***"

// Validate checks the settings the server requires to start and returns all
// problems joined into one error. Cron expressions are already checked by Load.
func (c *Config) Validate() error {
	var errs []error
	if len(c.Auth.Secret) < MinAuthSecretLen {
		errs = append(errs, fmt.Errorf("WISELABZ_AUTH_SECRET is missing or too short (min %d chars)", MinAuthSecretLen))
	}
	if _, err := crypto.DecodeKey(c.Encryption.Key); err != nil {
		errs = append(errs, fmt.Errorf("WISELABZ_ENCRYPTION_KEY is missing or invalid (want base64-encoded 32 bytes, e.g. `openssl rand -base64 32`): %w", err))
	}
	if c.Server.Origin == "" {
		errs = append(errs, errors.New("WISELABZ_SERVER_ORIGIN is missing: an explicit allowed CORS origin is required"))
	}
	if c.DB.Driver != "sqlite" && c.DB.Driver != "sqlite3" && c.DB.Driver != "postgres" {
		errs = append(errs, fmt.Errorf("db.driver %q is not supported (use sqlite or postgres)", c.DB.Driver))
	}
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errs = append(errs, fmt.Errorf("server.port %d out of range 1-65535", c.Server.Port))
	}
	return errors.Join(errs...)
}

// Redacted returns a deep-enough copy of the config with secrets masked, safe
// to print or log. Empty secrets are left empty so "unset" stays visible.
func (c *Config) Redacted() Config {
	r := *c
	r.DB.DSN = redactDSN(r.DB.DSN)
	r.Encryption.Key = mask(r.Encryption.Key)
	r.Auth.Secret = mask(r.Auth.Secret)
	r.AI.APIKey = mask(r.AI.APIKey)
	r.AI.EmbedAPIKey = mask(r.AI.EmbedAPIKey)
	if c.Auth.OIDC != nil {
		r.Auth.OIDC = make([]OIDCProvider, len(c.Auth.OIDC))
		copy(r.Auth.OIDC, c.Auth.OIDC)
		for i := range r.Auth.OIDC {
			r.Auth.OIDC[i].ClientSecret = mask(r.Auth.OIDC[i].ClientSecret)
		}
	}
	return r
}

func mask(s string) string {
	if s == "" {
		return ""
	}
	return redactedValue
}

// redactDSN masks the password in URL-style DSNs (postgres://user:pass@host/db)
// and, for key=value DSNs, any password= field. SQLite file DSNs carry no
// secret and pass through unchanged.
func redactDSN(dsn string) string {
	if u, err := url.Parse(dsn); err == nil && u.User != nil {
		if _, ok := u.User.Password(); ok {
			u.User = url.UserPassword(u.User.Username(), redactedValue)
			return u.String()
		}
		return dsn
	}
	return redactKVPassword(dsn)
}

func redactKVPassword(dsn string) string {
	const key = "password="
	for i := 0; i+len(key) <= len(dsn); i++ {
		if dsn[i:i+len(key)] != key || (i > 0 && dsn[i-1] != ' ') {
			continue
		}
		end := i + len(key)
		for end < len(dsn) && dsn[end] != ' ' {
			end++
		}
		return dsn[:i+len(key)] + redactedValue + dsn[end:]
	}
	return dsn
}
