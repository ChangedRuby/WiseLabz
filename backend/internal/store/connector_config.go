package store

import (
	"encoding/json"
	"fmt"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/crypto"
)

// IsSecretFieldType reports whether a SchemaField.Type holds a value that
// must be encrypted at rest in config_data. "password" is the only kind
// today; "secret" (a future multi-line paste field for PEM certs/SSH keys)
// must be treated the same way, so check this helper rather than comparing
// against "password" directly.
func IsSecretFieldType(t string) bool {
	return t == "password" || t == "secret"
}

// ParseConnectorConfig parses the config_data JSON string into a map,
// decrypting the value of any field the connector type's schema marks as
// secret-bearing (see IsSecretFieldType). A stored value that fails to
// decrypt is assumed to be legacy plaintext (saved before encryption was
// added) and is returned as-is; MarshalConnectorConfig will encrypt it on
// the connector's next save.
func ParseConnectorConfig(connType, data, encKeyB64 string) (map[string]any, error) {
	var cfg map[string]any
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, fmt.Errorf("parse connector config: %w", err)
	}
	if cfg == nil {
		cfg = map[string]any{}
	}

	schema, err := connector.GetTypeSchema(connType)
	if err != nil {
		// Unknown/unregistered connector type: nothing known to decrypt.
		return cfg, nil //nolint:nilerr
	}

	var key []byte
	for _, f := range schema.Fields {
		if !IsSecretFieldType(f.Type) {
			continue
		}
		raw, ok := cfg[f.Key].(string)
		if !ok || raw == "" {
			continue
		}
		if key == nil {
			key, err = crypto.DecodeKey(encKeyB64)
			if err != nil {
				return nil, fmt.Errorf("decode encryption key: %w", err)
			}
		}
		if plaintext, err := crypto.Decrypt(raw, key); err == nil {
			cfg[f.Key] = plaintext
		}
		// else: not valid ciphertext, treat as legacy plaintext and leave as-is.
	}
	return cfg, nil
}

// SecretFieldsChanged reports whether any secret-typed config field's
// plaintext value differs between the connector's currently stored config
// (oldConfigData, as read from connectors.config_data) and newConfig (the
// plaintext config a request is about to save). A rename-only edit or a
// resubmitted-but-unchanged secret must report false; a newly set, changed,
// or cleared secret field reports true.
func SecretFieldsChanged(connType, oldConfigData string, newConfig map[string]any, encKeyB64 string) (bool, error) {
	oldConfig, err := ParseConnectorConfig(connType, oldConfigData, encKeyB64)
	if err != nil {
		return false, fmt.Errorf("parse existing connector config: %w", err)
	}
	schema, err := connector.GetTypeSchema(connType)
	if err != nil {
		// Unknown/unregistered connector type: no secret fields are known,
		// so nothing can have changed.
		return false, nil //nolint:nilerr
	}
	for _, f := range schema.Fields {
		if !IsSecretFieldType(f.Type) {
			continue
		}
		oldVal, _ := oldConfig[f.Key].(string)
		newVal, _ := newConfig[f.Key].(string)
		if oldVal != newVal {
			return true, nil
		}
	}
	return false, nil
}

// MarshalConnectorConfig marshals a config map to a JSON string, encrypting
// the value of any field the connector type's schema marks as secret-bearing
// (see IsSecretFieldType) before marshaling.
func MarshalConnectorConfig(connType string, cfg map[string]any, encKeyB64 string) (string, error) {
	schema, err := connector.GetTypeSchema(connType)
	if err == nil {
		var key []byte
		toEncrypt := make(map[string]any, len(cfg))
		for k, v := range cfg {
			toEncrypt[k] = v
		}
		for _, f := range schema.Fields {
			if !IsSecretFieldType(f.Type) {
				continue
			}
			raw, ok := toEncrypt[f.Key].(string)
			if !ok || raw == "" {
				continue
			}
			if key == nil {
				key, err = crypto.DecodeKey(encKeyB64)
				if err != nil {
					return "", fmt.Errorf("decode encryption key: %w", err)
				}
			}
			encrypted, err := crypto.Encrypt(raw, key)
			if err != nil {
				return "", fmt.Errorf("encrypt connector config field %q: %w", f.Key, err)
			}
			toEncrypt[f.Key] = encrypted
		}
		cfg = toEncrypt
	}
	// Unknown/unregistered connector type: nothing known to encrypt, marshal as-is.

	b, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal connector config: %w", err)
	}
	return string(b), nil
}
