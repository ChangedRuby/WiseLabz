package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func TestConnectorsCreateRotationValidation(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	tests := []struct {
		name string
		body map[string]any
	}{
		{"negative rotationMaxAgeDays", map[string]any{
			"name": "svc", "category": "virtualization", "type": "proxmox", "url": "https://example.com",
			"rotationMaxAgeDays": -1,
		}},
		{"zero rotationMaxAgeDays", map[string]any{
			"name": "svc", "category": "virtualization", "type": "proxmox", "url": "https://example.com",
			"rotationMaxAgeDays": 0,
		}},
		{"malformed userExpiresAt", map[string]any{
			"name": "svc", "category": "virtualization", "type": "proxmox", "url": "https://example.com",
			"userExpiresAt": "not-a-date",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := app.req(t, http.MethodPost, "/api/connectors", tt.body, opToken)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestConnectorsUpdateRotationValidation(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	conn := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("seed connector: %v", err)
	}
	app.connectorGrant(t, opUserID, conn.ID, "operator")

	tests := []struct {
		name string
		body map[string]any
	}{
		{"negative rotationMaxAgeDays", map[string]any{"rotationMaxAgeDays": -5}},
		{"malformed userExpiresAt", map[string]any{"userExpiresAt": "banana"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := app.req(t, http.MethodPut, "/api/connectors/"+conn.ID, tt.body, opToken)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
			}
		})
	}
}

// TestConnectorsUpdateRotationFieldsRoleBoundary reuses the #240 per-connector
// role model: a viewer grant must not be able to change rotation settings.
func TestConnectorsUpdateRotationFieldsRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	viewerUserID, viewerToken := app.user(t, "viewer")
	conn := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("seed connector: %v", err)
	}
	app.connectorGrant(t, viewerUserID, conn.ID, "viewer")

	rec := app.req(t, http.MethodPut, "/api/connectors/"+conn.ID, map[string]any{"rotationMaxAgeDays": 30}, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

// TestConnectorsRotationFieldsRoundTrip checks userExpiresAt/rotationMaxAgeDays
// travel through create, GET, update (including clearing via explicit null),
// and that secretRotatedAt is present and read-only (never accepted on write).
func TestConnectorsRotationFieldsRoundTrip(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")

	createRec := app.req(t, http.MethodPost, "/api/connectors", map[string]any{
		"name": "svc", "category": "virtualization", "type": "proxmox", "url": "https://example.com",
		"userExpiresAt": "2027-01-01T00:00:00Z", "rotationMaxAgeDays": 45,
	}, opToken)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", createRec.Code, createRec.Body)
	}
	var created store.ConnectorRecord
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create body: %v", err)
	}
	if created.UserExpiresAt != "2027-01-01T00:00:00Z" {
		t.Errorf("created.UserExpiresAt = %q, want 2027-01-01T00:00:00Z", created.UserExpiresAt)
	}
	if created.RotationMaxAgeDays == nil || *created.RotationMaxAgeDays != 45 {
		t.Errorf("created.RotationMaxAgeDays = %v, want 45", created.RotationMaxAgeDays)
	}
	if created.SecretRotatedAt == "" {
		t.Error("created.SecretRotatedAt is empty, want set on create")
	}
	app.connectorGrant(t, opUserID, created.ID, "operator")

	getRec := app.req(t, http.MethodGet, "/api/connectors/"+created.ID, nil, opToken)
	if getRec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200; body = %s", getRec.Code, getRec.Body)
	}
	var got store.ConnectorRecord
	if err := json.Unmarshal(getRec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode get body: %v", err)
	}
	if got.UserExpiresAt != created.UserExpiresAt || got.RotationMaxAgeDays == nil || *got.RotationMaxAgeDays != 45 {
		t.Errorf("GET rotation fields = %+v, want to match create response", got)
	}

	updateRec := app.req(t, http.MethodPut, "/api/connectors/"+created.ID, map[string]any{
		"userExpiresAt": nil, "rotationMaxAgeDays": nil,
	}, opToken)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body = %s", updateRec.Code, updateRec.Body)
	}
	cleared, err := app.Store.GetConnector(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetConnector: %v", err)
	}
	if cleared.UserExpiresAt != "" {
		t.Errorf("UserExpiresAt after clear = %q, want empty", cleared.UserExpiresAt)
	}
	if cleared.RotationMaxAgeDays != nil {
		t.Errorf("RotationMaxAgeDays after clear = %v, want nil", cleared.RotationMaxAgeDays)
	}

	// secretRotatedAt is read-only: an update request carrying it alongside a
	// real field (owner) must not let it through — the request struct simply
	// has no such field, so it's silently ignored rather than applied.
	before := cleared.SecretRotatedAt
	rec := app.req(t, http.MethodPut, "/api/connectors/"+created.ID, map[string]any{
		"owner":           "ops-team",
		"secretRotatedAt": "2020-01-01T00:00:00Z",
	}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	after, err := app.Store.GetConnector(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetConnector: %v", err)
	}
	if after.SecretRotatedAt != before {
		t.Errorf("SecretRotatedAt = %q after submitting a write-attempt, want unchanged %q (read-only field)", after.SecretRotatedAt, before)
	}
}

// TestConnectorsSecretChangeBumpsSecretRotatedAtViaAPI exercises the full
// create -> update HTTP path: only an actual secret change bumps
// secret_rotated_at, not a rename or a resubmitted-identical secret.
func TestConnectorsSecretChangeBumpsSecretRotatedAtViaAPI(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")

	createRec := app.req(t, http.MethodPost, "/api/connectors", map[string]any{
		"name": "svc", "category": "virtualization", "type": "proxmox", "url": "https://pve.example",
		"config": map[string]any{"token_id": "root@pam!mon", "token_secret": "s3cr3t-v1"},
	}, opToken)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", createRec.Code, createRec.Body)
	}
	var created store.ConnectorRecord
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create body: %v", err)
	}
	app.connectorGrant(t, opUserID, created.ID, "operator")

	// Backdate secret_rotated_at so the RFC3339 (second-precision) comparisons
	// below can't collide with "now" purely from test speed.
	const initialRotatedAt = "2020-01-01T00:00:00Z"
	if err := app.Store.UpdateConnector(context.Background(), created.ID, map[string]any{"secret_rotated_at": initialRotatedAt}); err != nil {
		t.Fatalf("backdate secret_rotated_at: %v", err)
	}

	// Rename only: resubmit the exact same secret value plus a changed name.
	renameRec := app.req(t, http.MethodPut, "/api/connectors/"+created.ID, map[string]any{
		"name":   "svc-renamed",
		"config": map[string]any{"token_id": "root@pam!mon", "token_secret": "s3cr3t-v1"},
	}, opToken)
	if renameRec.Code != http.StatusOK {
		t.Fatalf("rename status = %d, want 200; body = %s", renameRec.Code, renameRec.Body)
	}
	afterRename, err := app.Store.GetConnector(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetConnector: %v", err)
	}
	if afterRename.SecretRotatedAt != initialRotatedAt {
		t.Errorf("SecretRotatedAt after rename-only update = %q, want unchanged %q", afterRename.SecretRotatedAt, initialRotatedAt)
	}

	// Actual secret rotation.
	rotateRec := app.req(t, http.MethodPut, "/api/connectors/"+created.ID, map[string]any{
		"config": map[string]any{"token_id": "root@pam!mon", "token_secret": "s3cr3t-v2"},
	}, opToken)
	if rotateRec.Code != http.StatusOK {
		t.Fatalf("rotate status = %d, want 200; body = %s", rotateRec.Code, rotateRec.Body)
	}
	afterRotate, err := app.Store.GetConnector(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetConnector: %v", err)
	}
	if afterRotate.SecretRotatedAt == initialRotatedAt {
		t.Error("SecretRotatedAt unchanged after an actual secret rotation, want bumped")
	}
}

// TestConnectorUpdateAuditNeverLeaksSecretValue confirms the audit trail
// records which fields changed by name, never the config value itself — the
// existing connector.update discipline (docs/AUDIT.md) applies unchanged to
// rotation fields.
func TestConnectorUpdateAuditNeverLeaksSecretValue(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	conn := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("seed connector: %v", err)
	}
	app.connectorGrant(t, opUserID, conn.ID, "operator")

	const secretValue = "super-secret-token-xyz"
	rec := app.req(t, http.MethodPut, "/api/connectors/"+conn.ID, map[string]any{
		"userExpiresAt":      "2027-06-01T00:00:00Z",
		"rotationMaxAgeDays": 30,
		"config":             map[string]any{"token_id": "root@pam!mon", "token_secret": secretValue},
	}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	listRec := app.req(t, http.MethodGet, "/api/system/audit", nil, opToken)
	if listRec.Code != http.StatusOK {
		t.Fatalf("audit list status = %d, want 200; body = %s", listRec.Code, listRec.Body)
	}
	if strings.Contains(listRec.Body.String(), secretValue) {
		t.Fatal("audit log response contains the plaintext secret value")
	}

	var page struct {
		Items []store.AuditRecord `json:"items"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode audit body: %v", err)
	}
	var found *store.AuditRecord
	for i := range page.Items {
		if page.Items[i].Action == "connector.update" && page.Items[i].TargetID == conn.ID {
			found = &page.Items[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("no connector.update audit record for %s in %+v", conn.ID, page.Items)
	}

	var detail struct {
		Fields []string `json:"fields"`
	}
	if err := json.Unmarshal([]byte(found.Detail), &detail); err != nil {
		t.Fatalf("decode audit detail: %v", err)
	}
	fieldSet := make(map[string]bool, len(detail.Fields))
	for _, f := range detail.Fields {
		fieldSet[f] = true
	}
	for _, want := range []string{"user_expires_at", "rotation_max_age_days", "config_data", "secret_rotated_at"} {
		if want == "secret_rotated_at" {
			if !fieldSet[want] {
				t.Errorf("audit fields = %v, want %q present (secret actually changed)", detail.Fields, want)
			}
			continue
		}
		if !fieldSet[want] {
			t.Errorf("audit fields = %v, want %q present", detail.Fields, want)
		}
	}
	if fieldSet["token_secret"] {
		t.Error("audit fields list names the raw config key, not just config_data — check field-name discipline")
	}
}
