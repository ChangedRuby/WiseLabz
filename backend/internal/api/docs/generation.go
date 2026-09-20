package docs

import (
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/httputil"
)

// Generate handles POST /api/docs/generate.
// Generates a doc from a template and connector snapshot.
func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		TemplateID  string `json:"templateId"`
		ConnectorID string `json:"connectorId"`
	}](w, r)
	if !ok {
		return
	}
	if req.TemplateID == "" || req.ConnectorID == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "templateId and connectorId are required")
		return
	}
	if !h.requireDocOperator(w, r, req.ConnectorID) {
		return
	}

	result, err := h.DocEngine.GenerateFromTemplate(r.Context(), req.TemplateID, req.ConnectorID)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	h.syncDocEmbeddings(r.Context(), result.DocID, result.Content)
	httputil.JSON(w, http.StatusCreated, result)
}

// GenerateTopology handles POST /api/docs/topology: (re)generates the
// single lab-wide "Lab Topology" doc from every connector's latest
// snapshot. Idempotent — safe to call repeatedly (e.g. on every
// TopologyPage load) since it updates the existing doc in place.
func (h *Handler) GenerateTopology(w http.ResponseWriter, r *http.Request) {
	result, err := h.DocEngine.GenerateLabTopology(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	h.syncDocEmbeddings(r.Context(), result.DocID, result.Content)
	httputil.JSON(w, http.StatusOK, result)
}

// TemplateSchema handles GET /api/docs/template-schema.
// Returns the JSON schema of templateData (available fields/types) and the list
// of available template functions for use in template autocomplete/reference.
func (h *Handler) TemplateSchema(w http.ResponseWriter, _ *http.Request) {
	schema := map[string]any{
		"fields": map[string]map[string]any{
			"ServiceName": {
				"type":        "string",
				"description": "The name of the service",
			},
			"Type": {
				"type":        "string",
				"description": "The type of service",
			},
			"Sections": {
				"type":        "array",
				"description": "Array of snapshot sections",
				"items": map[string]any{
					"type": "object",
					"fields": map[string]map[string]any{
						"Title": {
							"type":        "string",
							"description": "Section title",
						},
						"Content": {
							"type":        "string",
							"description": "Section content",
						},
					},
				},
			},
			"Dependencies": {
				"type":        "array",
				"description": "Array of service dependencies",
				"items": map[string]any{
					"type": "object",
					"fields": map[string]map[string]any{
						"Kind": {
							"type":        "string",
							"description": "Dependency kind (host|network|storage|upstream_service)",
						},
						"Name": {
							"type":        "string",
							"description": "Dependency name",
						},
						"Ref": {
							"type":        "string",
							"description": "Optional connector/service ID reference",
						},
					},
				},
			},
			"Metadata": {
				"type":        "object",
				"description": "Key-value metadata map (map[string]string)",
			},
			"GeneratedAt": {
				"type":        "string",
				"description": "ISO-8601 timestamp when the template was rendered",
			},
		},
		"functions": []map[string]string{
			{
				"name":        "dateFormat",
				"description": "Format a time value or RFC3339 string using a Go layout string. Usage: .GeneratedAt | dateFormat \"2006-01-02\"",
			},
			{
				"name":        "truncate",
				"description": "Truncate a string to N characters, appending '...' if truncated. Usage: .SomeText | truncate 50",
			},
			{
				"name":        "toJSON",
				"description": "Marshal a value to a JSON string for embedding structured data. Usage: .Metadata | toJSON",
			},
			{
				"name":        "filterByTitle",
				"description": "Filter sections by exact title match. Usage: .Sections | filterByTitle \"Configuration\"",
			},
			{
				"name":        "join",
				"description": "Join an array of strings with a separator. Usage: .SomeStrings | join \", \"",
			},
		},
	}
	httputil.JSON(w, http.StatusOK, schema)
}
