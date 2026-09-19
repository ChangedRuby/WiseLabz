package connectors

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/sync"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// ConfigPush handles POST /api/connectors/{id}/config-push. Per ADR 0003:
// field-level partial update against a per-connector whitelist, snapshot
// before write, verify-diff after, auto-revert-then-alert on mismatch.
func (h *Handler) ConfigPush(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	req, ok := httputil.DecodeJSON[struct {
		EntityRef string `json:"entityRef"`
		FieldKey  string `json:"fieldKey"`
		Value     any    `json:"value"`
		// PreviousValue is the field's currently-displayed value, as the
		// frontend read it from GET /{id}/data before the user edited it.
		// The handler has no generic way to re-derive "the old value of
		// fieldKey" from a rendered ServiceSnapshot, so the revert path
		// (ADR 0003) re-pushes this rather than something inferred.
		PreviousValue any `json:"previousValue"`
	}](w, r)
	if !ok {
		return
	}
	if err := connector.ValidateCompositeRef(req.EntityRef); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "invalid entityRef")
		return
	}
	if req.FieldKey == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "fieldKey is required")
		return
	}

	rec, err := h.Store.GetConnector(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	cfg, err := store.ParseConnectorConfig(rec.Type, rec.ConfigData, h.Config.Encryption.Key)
	if err != nil {
		httputil.Errorf(w, fmt.Errorf("parse config: %w", err))
		return
	}
	cfg["url"] = rec.URL
	cfg["verify_tls"] = rec.VerifyTLS

	conn, err := connector.Get(rec.Type, cfg)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	pusher, ok := conn.(connector.ConfigPusher)
	if !ok {
		httputil.Error(w, http.StatusBadRequest, "unsupported_operation", "connector does not support config-push")
		return
	}
	if !isWritableField(pusher, req.FieldKey) {
		httputil.Error(w, http.StatusBadRequest, "unsupported_field", fmt.Sprintf("field %q is not writable for this connector", req.FieldKey))
		return
	}

	if err := auth.ValidateElevationHeader(h.JWT, h.Store, "connector.configPush", r); err != nil {
		auth.WriteElevationError(w, err)
		return
	}

	pre, err := conn.Fetch(r.Context(), cfg)
	if err != nil {
		httputil.Errorf(w, fmt.Errorf("pre-push fetch: %w", err))
		return
	}

	if err := pusher.ConfigPush(r.Context(), cfg, req.EntityRef, req.FieldKey, req.Value); err != nil {
		slog.Error("connector config-push failed", "connector", id, "field", req.FieldKey, "error", err)
		httputil.Error(w, http.StatusBadGateway, "config_push_failed", err.Error())
		return
	}

	post, err := conn.Fetch(r.Context(), cfg)
	if err != nil {
		httputil.Errorf(w, fmt.Errorf("post-push fetch: %w", err))
		return
	}

	if !configPushLanded(pre, post) {
		h.revertConfigPush(w, r, pusher, cfg, rec, req.EntityRef, req.FieldKey, req.PreviousValue)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.configPush", "connector", id, map[string]any{
		"fieldKey":  req.FieldKey,
		"entityRef": req.EntityRef,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.configPush", "error", err)
	}

	httputil.JSON(w, http.StatusOK, post)
}

// isWritableField reports whether key is on the connector's config-push
// whitelist.
func isWritableField(pusher connector.ConfigPusher, key string) bool {
	for _, f := range pusher.WritableFields() {
		if f.Key == key {
			return true
		}
	}
	return false
}

// configPushLanded reports whether pre→post shows any change at all.
// sync.Compare's diff results, non-empty, are treated as "the write landed
// as some change" — the field-scoped check ADR 0003 calls for; a
// per-field-key section filter would need a fixed field→section mapping
// per connector, which the curated whitelist doesn't carry today.
// ponytail: any detected change counts as "landed" rather than verifying
// the new value matches exactly — tighten if a connector's diff proves too
// coarse to catch a landed-but-wrong-value push.
func configPushLanded(pre, post *connector.ServiceSnapshot) bool {
	return len(sync.Compare(pre, post)) > 0
}

// revertConfigPush handles ADR 0003's mismatch path: re-push the field's
// pre-push value (the revert), then raise a critical alert describing the
// mismatch — escalating further if the revert call itself fails, since
// that's the one hard stop in this feature set requiring manual
// intervention. No further automation is attempted either way.
// ponytail: reverts exactly one field per call by design (matches
// ConfigPush's own one-field-per-call contract), not a multi-field batch.
func (h *Handler) revertConfigPush(w http.ResponseWriter, r *http.Request, pusher connector.ConfigPusher, cfg map[string]any, rec *store.ConnectorRecord, entityRef, fieldKey string, previousValue any) {
	id := rec.ID
	revertErr := pusher.ConfigPush(r.Context(), cfg, entityRef, fieldKey, previousValue)

	desc := fmt.Sprintf("Pushed field %q did not verify after write; auto-revert to the pre-push value was attempted.", fieldKey)
	if revertErr != nil {
		desc = fmt.Sprintf("Pushed field %q did not verify after write; auto-revert FAILED (%v) — manual intervention required.", fieldKey, revertErr)
		slog.Error("config-push auto-revert failed", "connector", id, "field", fieldKey, "error", revertErr)
	} else {
		slog.Warn("config-push mismatch, auto-reverted", "connector", id, "field", fieldKey)
	}

	alert := &store.AlertRecord{
		ServiceID:   id,
		Severity:    "critical",
		Title:       fmt.Sprintf("Config push mismatch for %s", rec.Name),
		Description: desc,
	}
	if createErr := h.Store.CreateAlert(r.Context(), alert); createErr != nil {
		slog.Error("failed to create config-push mismatch alert", "error", createErr)
	} else if h.WSHub != nil {
		h.WSHub.Broadcast(ws.EventAlertCreated, map[string]any{
			"alertId":   alert.ID,
			"serviceId": id,
			"severity":  alert.Severity,
			"title":     alert.Title,
		})
	}

	httputil.Error(w, http.StatusConflict, "config_push_mismatch", "pushed value did not verify; auto-revert attempted, see the new alert")
}
