package quality

import (
	"context"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

const rotationRefresherConnType = "quality_test_rotation_refresher"

type fakeRefresherConnector struct{}

func (f *fakeRefresherConnector) Name() string     { return "fake refresher" }
func (f *fakeRefresherConnector) Type() string     { return rotationRefresherConnType }
func (f *fakeRefresherConnector) Category() string { return "test" }
func (f *fakeRefresherConnector) Fetch(_ context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	return nil, nil
}
func (f *fakeRefresherConnector) Validate(_ context.Context, _ map[string]any) error { return nil }
func (f *fakeRefresherConnector) RefreshCredentials(_ context.Context, config map[string]any) (map[string]any, time.Time, error) {
	return config, time.Time{}, nil
}

func init() {
	connector.Register(connector.TypeSchema{
		Type:     rotationRefresherConnType,
		Category: "virtualization",
		Name:     "Rotation test refresher connector",
	}, func(_ map[string]any) (connector.Connector, error) { return &fakeRefresherConnector{}, nil })
}

// createConnectorWithRotation creates a connector with an explicit
// secretRotatedAt (backdating CreateConnector's created_at default), and
// optional userExpiresAt/rotationMaxAgeDays overrides.
func createConnectorWithRotation(t *testing.T, s *store.Store, typ string, secretRotatedAt time.Time, userExpiresAt string, maxAgeDays *int) *store.ConnectorRecord {
	t.Helper()
	c := &store.ConnectorRecord{
		// Owner set so these tests don't also trip the unrelated
		// ownership_incomplete check.
		Name: "Rotation test connector", Category: "virtualization", Type: typ, URL: "https://example.test", Owner: "platform-team",
		SecretRotatedAt: secretRotatedAt.UTC().Format(time.RFC3339), UserExpiresAt: userExpiresAt, RotationMaxAgeDays: maxAgeDays,
	}
	if err := s.CreateConnector(context.Background(), c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	return c
}

func TestCheckCredentialRotationTableDriven(t *testing.T) {
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	fixedNow := func() time.Time { return now }

	days := func(d int) time.Time { return now.AddDate(0, 0, -d) }
	intPtr := func(v int) *int { return &v }

	tests := []struct {
		name            string
		connType        string
		secretRotatedAt time.Time
		userExpiresAt   string
		maxAgeDays      *int
		wantSeverity    string // "" means no open finding
	}{
		{
			name: "no expiry and young", connType: "proxmox",
			secretRotatedAt: days(1), wantSeverity: "",
		},
		{
			name: "near age is warning", connType: "proxmox",
			// global max age 90, warn 14: 80 days old -> due in 10 days, within warn window.
			secretRotatedAt: days(80), wantSeverity: "warning",
		},
		{
			name: "past age is critical", connType: "proxmox",
			secretRotatedAt: days(91), wantSeverity: "critical",
		},
		{
			name: "user expiry sooner than age forces critical", connType: "proxmox",
			secretRotatedAt: days(1), userExpiresAt: now.Add(-time.Hour).Format(time.RFC3339), wantSeverity: "critical",
		},
		{
			name: "user expiry sooner than age but still in warn window", connType: "proxmox",
			secretRotatedAt: days(1), userExpiresAt: now.AddDate(0, 0, 10).Format(time.RFC3339), wantSeverity: "warning",
		},
		{
			name: "per-connector override longer than global avoids finding", connType: "proxmox",
			secretRotatedAt: days(91), maxAgeDays: intPtr(180), wantSeverity: "",
		},
		{
			name: "per-connector override shorter than global forces critical", connType: "proxmox",
			secretRotatedAt: days(10), maxAgeDays: intPtr(5), wantSeverity: "critical",
		},
		{
			name: "refresher connector is always skipped", connType: rotationRefresherConnType,
			secretRotatedAt: days(9999), wantSeverity: "",
		},
		{
			name: "exact boundary instant at due is critical", connType: "proxmox",
			secretRotatedAt: now.AddDate(0, 0, -90), wantSeverity: "critical",
		},
		{
			name: "exact boundary instant at warn window start is warning", connType: "proxmox",
			secretRotatedAt: now.AddDate(0, 0, -90).AddDate(0, 0, 14), wantSeverity: "warning",
		},
		{
			name: "one second before warn window is not yet due", connType: "proxmox",
			secretRotatedAt: now.AddDate(0, 0, -90).AddDate(0, 0, 14).Add(time.Second), wantSeverity: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStore(t)
			c := createConnectorWithRotation(t, s, tt.connType, tt.secretRotatedAt, tt.userExpiresAt, tt.maxAgeDays)
			checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
			checker.now = fixedNow

			if err := checker.RunForConnector(context.Background(), c.ID); err != nil {
				t.Fatalf("RunForConnector() error: %v", err)
			}
			open := findings(t, s, c.ID, "credential_rotation", "open")
			if tt.wantSeverity == "" {
				if len(open) != 0 {
					t.Fatalf("open credential_rotation findings = %d, want 0 (got %#v)", len(open), open)
				}
				return
			}
			if len(open) != 1 {
				t.Fatalf("open credential_rotation findings = %d, want 1", len(open))
			}
			if open[0].Severity != tt.wantSeverity {
				t.Fatalf("severity = %q, want %q", open[0].Severity, tt.wantSeverity)
			}
		})
	}
}

func TestCheckCredentialRotationResolvesAfterRotationThenReopens(t *testing.T) {
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	s := newTestStore(t)
	c := createConnectorWithRotation(t, s, "proxmox", now.AddDate(0, 0, -91), "", nil)
	checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	checker.now = func() time.Time { return now }

	ctx := context.Background()
	if err := checker.RunForConnector(ctx, c.ID); err != nil {
		t.Fatalf("RunForConnector() detect error: %v", err)
	}
	if got := findings(t, s, c.ID, "credential_rotation", "open"); len(got) != 1 {
		t.Fatalf("open findings after detect = %d, want 1", len(got))
	}

	// Rotate the secret: bump secret_rotated_at to "now".
	if err := s.UpdateConnector(ctx, c.ID, map[string]any{"secret_rotated_at": now.Format(time.RFC3339)}); err != nil {
		t.Fatalf("UpdateConnector() error: %v", err)
	}
	if err := checker.RunForConnector(ctx, c.ID); err != nil {
		t.Fatalf("RunForConnector() resolve error: %v", err)
	}
	if got := findings(t, s, c.ID, "credential_rotation", "open"); len(got) != 0 {
		t.Fatalf("open findings after rotation = %d, want 0", len(got))
	}
	if got := findings(t, s, c.ID, "credential_rotation", "resolved"); len(got) != 1 {
		t.Fatalf("resolved findings after rotation = %d, want 1", len(got))
	}

	// Advance time past due again: must re-open.
	checker.now = func() time.Time { return now.AddDate(0, 0, 91) }
	if err := checker.RunForConnector(ctx, c.ID); err != nil {
		t.Fatalf("RunForConnector() reopen error: %v", err)
	}
	if got := findings(t, s, c.ID, "credential_rotation", "open"); len(got) != 1 {
		t.Fatalf("open findings after reopen = %d, want 1", len(got))
	}
}
