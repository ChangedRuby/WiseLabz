// Package templates provides API handlers for template CRUD and preview.
package templates

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/doc"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for template endpoints.
type Handler struct {
	Store     *store.Store
	DocEngine *doc.Engine
	// ponytail: one process-wide lock is enough at homelab scale; use per-template locks if save throughput matters.
	versionMu sync.Mutex
}

// NewHandler creates a new template handler.
func NewHandler(s *store.Store, eng *doc.Engine) *Handler {
	return &Handler{Store: s, DocEngine: eng}
}

// List handles GET /api/templates.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	_, pageSize, offset := httputil.Paginate(r)
	templates, _, err := h.Store.ListTemplates(r.Context(), offset, pageSize)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if templates == nil {
		templates = []store.TemplateRecord{}
	}

	// Spec: GET /templates returns a bare Template[] (see openapi.yaml).
	httputil.JSON(w, http.StatusOK, templates)
}

// templateResponse builds the flat Template JSON shape the spec expects,
// decoding the opaque appliesTo text column into an object when present.
func templateResponse(t *store.TemplateRecord, sections []store.TemplateSectionRecord) map[string]any {
	var appliesTo any
	if t.AppliesTo != "" {
		var m map[string]any
		if err := json.Unmarshal([]byte(t.AppliesTo), &m); err == nil {
			appliesTo = m
		}
	}
	return map[string]any{
		"id":             t.ID,
		"name":           t.Name,
		"description":    t.Description,
		"appliesTo":      appliesTo,
		"sections":       sections,
		"currentVersion": t.CurrentVersion,
	}
}

func versionSections(sections []store.TemplateSectionRecord) []store.TemplateVersionSection {
	result := make([]store.TemplateVersionSection, 0, len(sections))
	for _, section := range sections {
		result = append(result, store.TemplateVersionSection{
			Title: section.Title,
			Order: section.Ord,
			Body:  section.Body,
		})
	}
	return result
}

func createVersion(
	ctx context.Context,
	s *store.Store,
	t *store.TemplateRecord,
	sections []store.TemplateSectionRecord,
	trigger string,
) error {
	snapshot, err := json.Marshal(versionSections(sections))
	if err != nil {
		return err
	}
	return s.CreateTemplateVersion(ctx, &store.TemplateVersionRecord{
		TemplateID:  t.ID,
		Rev:         t.CurrentVersion,
		Name:        t.Name,
		Description: t.Description,
		AppliesTo:   t.AppliesTo,
		Sections:    string(snapshot),
		Author:      auth.UserIDFromContext(ctx),
		Trigger:     trigger,
	})
}

func templateVersionResponse(v store.TemplateVersionRecord) (map[string]any, error) {
	var appliesTo any
	if v.AppliesTo != "" {
		if err := json.Unmarshal([]byte(v.AppliesTo), &appliesTo); err != nil {
			return nil, err
		}
	}
	var sections []store.TemplateVersionSection
	if err := json.Unmarshal([]byte(v.Sections), &sections); err != nil {
		return nil, err
	}
	if sections == nil {
		sections = []store.TemplateVersionSection{}
	}
	var author any
	if v.Author != "" {
		author = v.Author
	}
	return map[string]any{
		"rev":         v.Rev,
		"name":        v.Name,
		"description": v.Description,
		"appliesTo":   appliesTo,
		"sections":    sections,
		"author":      author,
		"trigger":     v.Trigger,
		"createdAt":   v.CreatedAt,
	}, nil
}

// Get handles GET /api/templates/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := h.Store.GetTemplate(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Template not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	sections, _ := h.Store.GetTemplateSections(r.Context(), id)

	httputil.JSON(w, http.StatusOK, templateResponse(t, sections))
}

// Create handles POST /api/templates.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		AppliesTo   json.RawMessage `json:"appliesTo"`
		Sections    []struct {
			Title string `json:"title"`
			Order int    `json:"order"`
			Body  string `json:"body"`
		} `json:"sections"`
	}](w, r)
	if !ok {
		return
	}
	if req.Name == "" {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "name is required", []httputil.FieldError{{Field: "name", Msg: "is required"}})
		return
	}

	t := &store.TemplateRecord{
		Name:        req.Name,
		Description: req.Description,
		AppliesTo:   string(req.AppliesTo),
	}
	sections := make([]store.TemplateSectionRecord, 0, len(req.Sections))
	if err := h.Store.WithinTransaction(r.Context(), func(tx *store.Store) error {
		if err := tx.CreateTemplate(r.Context(), t); err != nil {
			return err
		}
		for _, sec := range req.Sections {
			rec := store.TemplateSectionRecord{
				TemplateID: t.ID, Title: sec.Title, Ord: sec.Order, Body: sec.Body,
			}
			if err := tx.CreateTemplateSection(r.Context(), &rec); err != nil {
				return err
			}
			sections = append(sections, rec)
		}
		return createVersion(r.Context(), tx, t, sections, "save")
	}); err != nil {
		httputil.Errorf(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, templateResponse(t, sections))
}

// Update handles PUT /api/templates/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	h.versionMu.Lock()
	defer h.versionMu.Unlock()

	id := r.PathValue("id")
	req, ok := httputil.DecodeJSON[struct {
		Name        *string          `json:"name"`
		Description *string          `json:"description"`
		AppliesTo   *json.RawMessage `json:"appliesTo"`
		Sections    *[]struct {
			ID    *string `json:"id"`
			Title string  `json:"title"`
			Order int     `json:"order"`
			Body  string  `json:"body"`
		} `json:"sections"`
	}](w, r)
	if !ok {
		return
	}

	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.AppliesTo != nil {
		updates["applies_to"] = string(*req.AppliesTo)
	}
	var (
		t        *store.TemplateRecord
		sections []store.TemplateSectionRecord
	)
	err := h.Store.WithinTransaction(r.Context(), func(tx *store.Store) error {
		if err := tx.UpdateTemplate(r.Context(), id, updates); err != nil {
			return err
		}

		if req.Sections != nil {
			if err := tx.DeleteTemplateSections(r.Context(), id); err != nil {
				return err
			}
			for _, sec := range *req.Sections {
				if err := tx.CreateTemplateSection(r.Context(), &store.TemplateSectionRecord{
					TemplateID: id, Title: sec.Title, Ord: sec.Order, Body: sec.Body,
				}); err != nil {
					return err
				}
			}
		}

		var err error
		t, err = tx.GetTemplate(r.Context(), id)
		if err != nil {
			return err
		}
		sections, err = tx.GetTemplateSections(r.Context(), id)
		if err != nil {
			return err
		}
		return createVersion(r.Context(), tx, t, sections, "save")
	})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Template not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, templateResponse(t, sections))
}

// Delete handles DELETE /api/templates/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.Store.DeleteTemplate(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Template not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}
	httputil.NoContent(w)
}

// Versions handles GET /api/templates/{id}/versions.
func (h *Handler) Versions(w http.ResponseWriter, r *http.Request) {
	versions, err := h.Store.GetTemplateVersions(r.Context(), r.PathValue("id"))
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	result := make([]map[string]any, 0, len(versions))
	for _, version := range versions {
		item, err := templateVersionResponse(version)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		delete(item, "sections")
		delete(item, "name")
		delete(item, "description")
		delete(item, "appliesTo")
		result = append(result, item)
	}
	if result == nil {
		result = []map[string]any{}
	}
	httputil.JSON(w, http.StatusOK, result)
}

// Version handles GET /api/templates/{id}/versions/{rev}.
func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	rev, err := strconv.Atoi(r.PathValue("rev"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid rev")
		return
	}
	versions, err := h.Store.GetTemplateVersions(r.Context(), r.PathValue("id"))
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	for _, version := range versions {
		if version.Rev != rev {
			continue
		}
		result, err := templateVersionResponse(version)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		httputil.JSON(w, http.StatusOK, result)
		return
	}
	httputil.Error(w, http.StatusNotFound, "not_found", "Version not found")
}

// Restore handles POST /api/templates/{id}/versions/{rev}/restore.
func (h *Handler) Restore(w http.ResponseWriter, r *http.Request) {
	h.versionMu.Lock()
	defer h.versionMu.Unlock()

	id := r.PathValue("id")
	rev, err := strconv.Atoi(r.PathValue("rev"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid rev")
		return
	}
	versions, err := h.Store.GetTemplateVersions(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	var target *store.TemplateVersionRecord
	for i := range versions {
		if versions[i].Rev == rev {
			target = &versions[i]
			break
		}
	}
	if target == nil {
		httputil.Error(w, http.StatusNotFound, "not_found", "Version not found")
		return
	}
	var sections []store.TemplateVersionSection
	if err := json.Unmarshal([]byte(target.Sections), &sections); err != nil {
		httputil.Errorf(w, err)
		return
	}
	var (
		t            *store.TemplateRecord
		liveSections []store.TemplateSectionRecord
	)
	if err := h.Store.WithinTransaction(r.Context(), func(tx *store.Store) error {
		if err := tx.UpdateTemplate(r.Context(), id, map[string]any{
			"name": target.Name, "description": target.Description, "applies_to": target.AppliesTo,
		}); err != nil {
			return err
		}
		if err := tx.DeleteTemplateSections(r.Context(), id); err != nil {
			return err
		}
		for _, section := range sections {
			if err := tx.CreateTemplateSection(r.Context(), &store.TemplateSectionRecord{
				TemplateID: id, Title: section.Title, Ord: section.Order, Body: section.Body,
			}); err != nil {
				return err
			}
		}
		var err error
		t, err = tx.GetTemplate(r.Context(), id)
		if err != nil {
			return err
		}
		liveSections, err = tx.GetTemplateSections(r.Context(), id)
		if err != nil {
			return err
		}
		return createVersion(r.Context(), tx, t, liveSections, "restore")
	}); err != nil {
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "template.restore", "template", id, map[string]any{
		"rev": rev,
	}); err != nil {
		slog.Error("failed to record audit", "action", "template.restore", "error", err)
	}
	httputil.JSON(w, http.StatusOK, templateResponse(t, liveSections))
}

// Preview handles POST /api/templates/{id}/preview without persisting docs or versions.
func (h *Handler) Preview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		ConnectorID string `json:"connectorId"`
	}
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
			return
		}
	}
	connectors, err := h.DocEngine.MatchingConnectors(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	type affectedConnector struct {
		ConnectorID    string  `json:"connectorId"`
		ConnectorName  string  `json:"connectorName"`
		HasExistingDoc bool    `json:"hasExistingDoc"`
		WouldChange    bool    `json:"wouldChange"`
		RenderError    *string `json:"renderError"`
	}
	affected := make([]affectedConnector, 0, len(connectors))
	previews := make(map[string]*doc.GenerateResult, len(connectors))
	for _, connector := range connectors {
		item := affectedConnector{ConnectorID: connector.ID, ConnectorName: connector.Name}
		docs, err := h.Store.ListDocsByService(r.Context(), connector.ID)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		item.HasExistingDoc = len(docs) > 0
		preview, renderErr := h.DocEngine.PreviewFromTemplate(r.Context(), id, connector.ID)
		if renderErr != nil {
			message := renderErr.Error()
			item.RenderError = &message
			affected = append(affected, item)
			continue
		}
		previews[connector.ID] = preview
		item.WouldChange = len(docs) == 0 || docs[0].Content != preview.Content
		affected = append(affected, item)
	}
	var detail *doc.GenerateResult
	if req.ConnectorID != "" {
		detail = previews[req.ConnectorID]
		if detail == nil {
			detail, err = h.DocEngine.PreviewFromTemplate(r.Context(), id, req.ConnectorID)
			if err != nil {
				httputil.Errorf(w, err)
				return
			}
		}
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"affected": affected, "detail": detail})
}
