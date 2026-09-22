package system

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// ListAudit handles GET /api/system/audit. Operator-only.
// Returns the audit trail (docs/AUDIT.md), newest first, optionally
// filtered by action, targetType, and/or a createdAt range.
//
// Passing ?cursor= (empty for the first page, then the previous response's
// nextCursor) switches to keyset pagination, which stays cheap as the table
// grows. Without the parameter the historical offset behaviour is unchanged.
func (h *Handler) ListAudit(w http.ResponseWriter, r *http.Request) {
	page, pageSize, offset := httputil.Paginate(r)
	action := r.URL.Query().Get("action")
	targetType := r.URL.Query().Get("targetType")
	createdAfter := r.URL.Query().Get("createdAfter")
	createdBefore := r.URL.Query().Get("createdBefore")

	keyset, curSort, curID, ok := httputil.Cursor(w, r)
	if !ok {
		return
	}
	if keyset {
		records, total, err := h.Store.ListAuditRecordsKeyset(r.Context(), action, targetType, createdAfter, createdBefore,
			store.Keyset{Sort: curSort, ID: curID}, pageSize)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		next := httputil.NextCursor(records, pageSize, func(a store.AuditRecord) (string, string) {
			return a.CreatedAt, a.ID
		})
		httputil.WritePaginatedCursor(w, records, page, pageSize, total, next)
		return
	}

	records, total, err := h.Store.ListAuditRecords(r.Context(), action, targetType, createdAfter, createdBefore, offset, pageSize)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.WritePaginated(w, records, page, pageSize, total)
}

// auditCSVHeader is the column order for CSV export.
var auditCSVHeader = []string{"id", "actorUserId", "actorRole", "action", "targetType", "targetId", "detail", "createdAt"}

// ExportAudit handles GET /api/system/audit/export. Operator-only.
// Exports the full (unpaginated) set of audit records matching the given
// filters as CSV or JSON.
func (h *Handler) ExportAudit(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}
	if format != "json" && format != "csv" {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_format", "format must be \"json\" or \"csv\"", []httputil.FieldError{{Field: "format", Msg: "must be \"json\" or \"csv\""}})
		return
	}

	action := r.URL.Query().Get("action")
	targetType := r.URL.Query().Get("targetType")
	createdAfter := r.URL.Query().Get("createdAfter")
	createdBefore := r.URL.Query().Get("createdBefore")

	records, err := h.Store.ListAllAuditRecords(r.Context(), action, targetType, createdAfter, createdBefore)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	timestamp := time.Now().UTC().Format("20060102-150405")

	if format == "csv" {
		filename := fmt.Sprintf("wiselabz-audit-%s.csv", timestamp)
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)

		cw := csv.NewWriter(w)
		if err := cw.Write(auditCSVHeader); err != nil {
			return // headers already sent; nothing more we can do
		}
		for _, a := range records {
			if err := cw.Write([]string{a.ID, a.ActorUserID, a.ActorRole, a.Action, a.TargetType, a.TargetID, a.Detail, a.CreatedAt}); err != nil {
				slog.Error("audit csv export: write row failed, aborting (client likely disconnected)", "error", err)
				return
			}
		}
		cw.Flush()
		return
	}

	filename := fmt.Sprintf("wiselabz-audit-%s.json", timestamp)
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	httputil.JSON(w, http.StatusOK, records)
}
