package store

import (
	"context"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const rotationTestConnType = "store_test_rotation_conn"

func init() {
	// Two secret-bearing fields plus one plain field, so SecretFieldsChanged
	// has more than one secret to compare (and can't just check "any field
	// changed").
	connector.Register(connector.TypeSchema{
		Type:     rotationTestConnType,
		Category: "test",
		Name:     "Rotation test connector",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "URL", Type: "text"},
			{Key: "api_key", Label: "API Key", Type: "password"},
			{Key: "api_secret", Label: "API Secret", Type: "password"},
		},
	}, func(_ map[string]any) (connector.Connector, error) { return nil, nil })
}

func TestCreateConnectorDefaultsSecretRotatedAtToCreatedAt(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	if c.SecretRotatedAt != c.CreatedAt {
		t.Fatalf("SecretRotatedAt = %q, want CreatedAt %q", c.SecretRotatedAt, c.CreatedAt)
	}

	got, err := s.GetConnector(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetConnector() error: %v", err)
	}
	if got.SecretRotatedAt != c.CreatedAt {
		t.Fatalf("stored SecretRotatedAt = %q, want %q", got.SecretRotatedAt, c.CreatedAt)
	}
	if got.UserExpiresAt != "" {
		t.Fatalf("UserExpiresAt = %q, want empty", got.UserExpiresAt)
	}
	if got.RotationMaxAgeDays != nil {
		t.Fatalf("RotationMaxAgeDays = %v, want nil", got.RotationMaxAgeDays)
	}
}

func TestConnectorRotationFieldsRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	maxAge := 30
	c := &ConnectorRecord{
		Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com",
		UserExpiresAt: "2027-01-01T00:00:00Z", RotationMaxAgeDays: &maxAge,
	}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	got, err := s.GetConnector(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetConnector() error: %v", err)
	}
	if got.UserExpiresAt != "2027-01-01T00:00:00Z" {
		t.Fatalf("UserExpiresAt = %q, want 2027-01-01T00:00:00Z", got.UserExpiresAt)
	}
	if got.RotationMaxAgeDays == nil || *got.RotationMaxAgeDays != 30 {
		t.Fatalf("RotationMaxAgeDays = %v, want 30", got.RotationMaxAgeDays)
	}

	newMaxAge := 60
	if err := s.UpdateConnector(ctx, c.ID, map[string]any{
		"user_expires_at":       "2028-01-01T00:00:00Z",
		"rotation_max_age_days": &newMaxAge,
	}); err != nil {
		t.Fatalf("UpdateConnector() error: %v", err)
	}
	got, err = s.GetConnector(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetConnector() after update error: %v", err)
	}
	if got.UserExpiresAt != "2028-01-01T00:00:00Z" {
		t.Fatalf("UserExpiresAt after update = %q, want 2028-01-01T00:00:00Z", got.UserExpiresAt)
	}
	if got.RotationMaxAgeDays == nil || *got.RotationMaxAgeDays != 60 {
		t.Fatalf("RotationMaxAgeDays after update = %v, want 60", got.RotationMaxAgeDays)
	}

	// Clearing (nil pointer / empty string) must null the columns out.
	if err := s.UpdateConnector(ctx, c.ID, map[string]any{
		"user_expires_at":       nil,
		"rotation_max_age_days": (*int)(nil),
	}); err != nil {
		t.Fatalf("UpdateConnector() clear error: %v", err)
	}
	got, err = s.GetConnector(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetConnector() after clear error: %v", err)
	}
	if got.UserExpiresAt != "" {
		t.Fatalf("UserExpiresAt after clear = %q, want empty", got.UserExpiresAt)
	}
	if got.RotationMaxAgeDays != nil {
		t.Fatalf("RotationMaxAgeDays after clear = %v, want nil", got.RotationMaxAgeDays)
	}
}

func TestSecretFieldsChangedOnActualSecretChange(t *testing.T) {
	old, err := MarshalConnectorConfig(rotationTestConnType, map[string]any{
		"url": "https://example.com", "api_key": "old-key", "api_secret": "old-secret",
	}, testEncKey)
	if err != nil {
		t.Fatalf("MarshalConnectorConfig() error: %v", err)
	}

	changed, err := SecretFieldsChanged(rotationTestConnType, old, map[string]any{
		"url": "https://example.com", "api_key": "new-key", "api_secret": "old-secret",
	}, testEncKey)
	if err != nil {
		t.Fatalf("SecretFieldsChanged() error: %v", err)
	}
	if !changed {
		t.Fatal("SecretFieldsChanged() = false, want true when api_key actually changes")
	}
}

func TestSecretFieldsChangedFalseOnRenameOnly(t *testing.T) {
	old, err := MarshalConnectorConfig(rotationTestConnType, map[string]any{
		"url": "https://example.com", "api_key": "same-key", "api_secret": "same-secret",
	}, testEncKey)
	if err != nil {
		t.Fatalf("MarshalConnectorConfig() error: %v", err)
	}

	// Non-secret field ("url") changes; both secrets stay identical.
	changed, err := SecretFieldsChanged(rotationTestConnType, old, map[string]any{
		"url": "https://renamed.example.com", "api_key": "same-key", "api_secret": "same-secret",
	}, testEncKey)
	if err != nil {
		t.Fatalf("SecretFieldsChanged() error: %v", err)
	}
	if changed {
		t.Fatal("SecretFieldsChanged() = true, want false for a non-secret-only edit")
	}
}

func TestSecretFieldsChangedFalseOnResubmittedUnchangedSecret(t *testing.T) {
	old, err := MarshalConnectorConfig(rotationTestConnType, map[string]any{
		"url": "https://example.com", "api_key": "same-key", "api_secret": "same-secret",
	}, testEncKey)
	if err != nil {
		t.Fatalf("MarshalConnectorConfig() error: %v", err)
	}

	changed, err := SecretFieldsChanged(rotationTestConnType, old, map[string]any{
		"url": "https://example.com", "api_key": "same-key", "api_secret": "same-secret",
	}, testEncKey)
	if err != nil {
		t.Fatalf("SecretFieldsChanged() error: %v", err)
	}
	if changed {
		t.Fatal("SecretFieldsChanged() = true, want false when the resubmitted secret is unchanged")
	}
}

func TestSecretFieldsChangedMultipleSecretFields(t *testing.T) {
	old, err := MarshalConnectorConfig(rotationTestConnType, map[string]any{
		"url": "https://example.com", "api_key": "key-1", "api_secret": "secret-1",
	}, testEncKey)
	if err != nil {
		t.Fatalf("MarshalConnectorConfig() error: %v", err)
	}

	// Only the second secret field changes; SecretFieldsChanged must still
	// detect it rather than short-circuiting on the first field checked.
	changed, err := SecretFieldsChanged(rotationTestConnType, old, map[string]any{
		"url": "https://example.com", "api_key": "key-1", "api_secret": "secret-2",
	}, testEncKey)
	if err != nil {
		t.Fatalf("SecretFieldsChanged() error: %v", err)
	}
	if !changed {
		t.Fatal("SecretFieldsChanged() = false, want true when the second secret field changes")
	}
}

func TestSecretFieldsChangedLegacyPlaintext(t *testing.T) {
	// Simulates a row saved before encryption was added: api_key stored as
	// plain JSON text rather than ciphertext.
	legacyJSON := `{"url":"https://example.com","api_key":"legacy-plain","api_secret":"legacy-plain-2"}`

	changed, err := SecretFieldsChanged(rotationTestConnType, legacyJSON, map[string]any{
		"url": "https://example.com", "api_key": "legacy-plain", "api_secret": "legacy-plain-2",
	}, testEncKey)
	if err != nil {
		t.Fatalf("SecretFieldsChanged() error: %v", err)
	}
	if changed {
		t.Fatal("SecretFieldsChanged() = true, want false when a legacy-plaintext secret is resubmitted unchanged")
	}

	changed, err = SecretFieldsChanged(rotationTestConnType, legacyJSON, map[string]any{
		"url": "https://example.com", "api_key": "rotated", "api_secret": "legacy-plain-2",
	}, testEncKey)
	if err != nil {
		t.Fatalf("SecretFieldsChanged() error: %v", err)
	}
	if !changed {
		t.Fatal("SecretFieldsChanged() = false, want true when a legacy-plaintext secret is rotated")
	}
}
