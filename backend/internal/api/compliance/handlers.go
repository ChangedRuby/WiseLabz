// Package compliance provides API handlers for the entity attribute schema
// that PR3's user-defined compliance rule engine builds on.
package compliance

import (
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/httputil"
)

// Handler holds dependencies for compliance endpoints.
type Handler struct{}

// NewHandler creates a compliance handler.
func NewHandler() *Handler { return &Handler{} }

// Schema handles GET /api/compliance/schema, returning connectorType ->
// entity kind -> attribute specs for every connector that registered an
// attribute catalog (see connector.RegisterAttributeCatalog).
func (h *Handler) Schema(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, connector.AttributeCatalog())
}
