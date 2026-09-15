// Package doc provides documentation generation from templates and snapshots.
package doc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Engine generates documentation from templates and connector snapshots.
type Engine struct {
	store *store.Store
}

// NewEngine creates a new doc engine.
func NewEngine(s *store.Store) *Engine {
	return &Engine{store: s}
}

// GenerateResult holds the output of document generation.
type GenerateResult struct {
	DocID   string `json:"docId"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type renderResult struct {
	Title   string
	Content string
}

// render executes a template against a connector's latest snapshot without persisting it.
func (e *Engine) render(ctx context.Context, templateID, connectorID string) (*renderResult, error) {
	tmpl, err := e.store.GetTemplate(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("get template: %w", err)
	}
	sections, err := e.store.GetTemplateSections(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("get template sections: %w", err)
	}

	sn, err := e.store.GetLatestSnapshot(ctx, connectorID)
	if err != nil {
		return nil, fmt.Errorf("get snapshot: %w", err)
	}

	var snap connector.ServiceSnapshot
	if err := json.Unmarshal([]byte(sn.Data), &snap); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}

	links, err := matchEntities(ctx, e.store, connectorID, snap.Entities)
	if err != nil {
		return nil, fmt.Errorf("match entities: %w", err)
	}

	var buf bytes.Buffer
	data := templateData{
		ServiceName:  snap.ServiceName,
		Type:         snap.Type,
		Sections:     snap.Sections,
		Dependencies: snap.Dependencies,
		Metadata:     snap.Metadata,
		Links:        links,
		GeneratedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	fmt.Fprintf(&buf, "# %s\n\n", snap.ServiceName)
	if tmpl.Description != "" {
		fmt.Fprintf(&buf, "> %s\n\n", tmpl.Description)
	}

	for _, sec := range sections {
		sectionTemplate, err := template.New("section").Funcs(TemplateFuncs()).Parse(sec.Body)
		if err != nil {
			fmt.Fprintf(&buf, "## %s\n\n_Template error: %v_\n\n", sec.Title, err)
			continue
		}
		fmt.Fprintf(&buf, "## %s\n\n", sec.Title)
		if err := sectionTemplate.Execute(&buf, data); err != nil {
			fmt.Fprintf(&buf, "\n_Template error: %v_\n", err)
		}
		buf.WriteString("\n")
	}

	return &renderResult{
		Title:   snap.ServiceName,
		Content: buf.String(),
	}, nil
}

// PreviewFromTemplate renders a document without persisting it.
func (e *Engine) PreviewFromTemplate(ctx context.Context, templateID, connectorID string) (*GenerateResult, error) {
	rendered, err := e.render(ctx, templateID, connectorID)
	if err != nil {
		return nil, err
	}

	return &GenerateResult{
		Title:   rendered.Title,
		Content: rendered.Content,
	}, nil
}

// GenerateFromTemplate generates and persists a document using a template and a snapshot.
func (e *Engine) GenerateFromTemplate(ctx context.Context, templateID, connectorID string) (*GenerateResult, error) {
	rendered, err := e.render(ctx, templateID, connectorID)
	if err != nil {
		return nil, err
	}

	var docID string
	if err := e.store.WithinTransaction(ctx, func(tx *store.Store) error {
		existingDocs, err := tx.ListDocsByService(ctx, connectorID)
		if err != nil {
			return fmt.Errorf("list existing docs: %w", err)
		}

		if len(existingDocs) > 0 {
			docID = existingDocs[0].ID
			if err := tx.UpdateDoc(ctx, docID, rendered.Content, nil); err != nil {
				return fmt.Errorf("update doc: %w", err)
			}
			doc, err := tx.GetDoc(ctx, docID)
			if err != nil {
				return fmt.Errorf("get updated doc: %w", err)
			}
			if err := tx.CreateDocVersion(ctx, &store.DocVersionRecord{
				DocID: docID, Rev: doc.CurrentVersion, Content: rendered.Content, Trigger: "template",
			}); err != nil {
				return fmt.Errorf("create doc version: %w", err)
			}
			return nil
		}

		doc := &store.DocRecord{
			Title: rendered.Title, Kind: "service", ServiceID: connectorID, Content: rendered.Content,
		}
		if err := tx.CreateDoc(ctx, doc); err != nil {
			return fmt.Errorf("create doc: %w", err)
		}
		docID = doc.ID
		if err := tx.CreateDocVersion(ctx, &store.DocVersionRecord{
			DocID: docID, Rev: 1, Content: rendered.Content, Trigger: "template",
		}); err != nil {
			return fmt.Errorf("create doc version: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &GenerateResult{
		DocID:   docID,
		Title:   rendered.Title,
		Content: rendered.Content,
	}, nil
}

// MatchingConnectors returns connectors covered by a template's applicability scope.
func (e *Engine) MatchingConnectors(ctx context.Context, templateID string) ([]store.ConnectorRecord, error) {
	tmpl, err := e.store.GetTemplate(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("get template: %w", err)
	}

	var appliesTo struct {
		Category string `json:"category"`
		Type     string `json:"type"`
	}
	if tmpl.AppliesTo != "" {
		if err := json.Unmarshal([]byte(tmpl.AppliesTo), &appliesTo); err != nil {
			return nil, fmt.Errorf("unmarshal template applies_to: %w", err)
		}
	}

	connectors, err := e.store.ListAllConnectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("list connectors: %w", err)
	}

	matches := make([]store.ConnectorRecord, 0, len(connectors))
	for _, candidate := range connectors {
		if appliesTo.Category != "" && candidate.Category != appliesTo.Category {
			continue
		}
		if appliesTo.Type != "" && candidate.Type != appliesTo.Type {
			continue
		}
		matches = append(matches, candidate)
	}

	return matches, nil
}

// renderSnapshot renders a document from a connector's latest snapshot
// without a template, without persisting it. Shared by GenerateFromSnapshot
// (first-time generation) and RegenerateForConnector (sync-triggered refresh
// of existing template-less docs).
func (e *Engine) renderSnapshot(ctx context.Context, connectorID string) (*renderResult, error) {
	sn, err := e.store.GetLatestSnapshot(ctx, connectorID)
	if err != nil {
		return nil, fmt.Errorf("get snapshot: %w", err)
	}

	var snap connector.ServiceSnapshot
	if err := json.Unmarshal([]byte(sn.Data), &snap); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "# %s\n\n", snap.ServiceName)
	fmt.Fprintf(&buf, "**Type:** %s\n", snap.Type)
	fmt.Fprintf(&buf, "**Fetched:** %s\n\n", snap.FetchedAt.Format(time.RFC3339))

	for _, sec := range snap.Sections {
		buf.WriteString(sec.Content)
		buf.WriteString("\n")
	}

	if len(snap.Dependencies) > 0 {
		buf.WriteString("## Dependencies\n\n")
		for _, dep := range snap.Dependencies {
			fmt.Fprintf(&buf, "- **%s**: %s", dep.Kind, dep.Name)
			if dep.Ref != "" {
				fmt.Fprintf(&buf, " (%s)", dep.Ref)
			}
			buf.WriteString("\n")
		}
		buf.WriteString("\n")
	}

	if links, err := matchEntities(ctx, e.store, connectorID, snap.Entities); err == nil && len(links) > 0 {
		buf.WriteString("## Related Entities\n\n")
		buf.WriteString(relatedEntities(snap.ServiceName, links))
		buf.WriteString("\n")
	}

	return &renderResult{Title: snap.ServiceName, Content: buf.String()}, nil
}

// GenerateFromSnapshot generates a raw document from a snapshot without a template.
func (e *Engine) GenerateFromSnapshot(ctx context.Context, connectorID string) (*GenerateResult, error) {
	rendered, err := e.renderSnapshot(ctx, connectorID)
	if err != nil {
		return nil, err
	}

	doc := &store.DocRecord{
		Title:     rendered.Title,
		Kind:      "service",
		ServiceID: connectorID,
		Content:   rendered.Content,
	}
	if err := e.store.CreateDoc(ctx, doc); err != nil {
		return nil, fmt.Errorf("create doc: %w", err)
	}

	docID := doc.ID
	_ = e.store.CreateDocVersion(ctx, &store.DocVersionRecord{
		DocID:   docID,
		Rev:     1,
		Content: rendered.Content,
		Trigger: "manual",
	})

	return &GenerateResult{
		DocID:   docID,
		Title:   rendered.Title,
		Content: rendered.Content,
	}, nil
}

// RegenerateForConnector re-renders every existing doc for a connector via
// the template-less snapshot path and, for any whose content actually
// changed, updates it and records a new "sync" doc version. Docs generated
// from a template are re-rendered the same way (DocRecord has no
// TemplateID to look up), so a template-derived doc's next manual
// regeneration will re-apply its template; sync only refreshes snapshot
// content in between. A connector with no existing docs is a no-op.
func (e *Engine) RegenerateForConnector(ctx context.Context, connectorID string) error {
	docs, err := e.store.ListDocsByService(ctx, connectorID)
	if err != nil {
		return fmt.Errorf("list docs by service: %w", err)
	}
	if len(docs) == 0 {
		return nil
	}

	rendered, err := e.renderSnapshot(ctx, connectorID)
	if err != nil {
		return err
	}

	for _, d := range docs {
		if d.Content == rendered.Content {
			continue
		}
		if err := e.store.WithinTransaction(ctx, func(tx *store.Store) error {
			if err := tx.UpdateDoc(ctx, d.ID, rendered.Content, nil); err != nil {
				return fmt.Errorf("update doc: %w", err)
			}
			updated, err := tx.GetDoc(ctx, d.ID)
			if err != nil {
				return fmt.Errorf("get updated doc: %w", err)
			}
			if err := tx.CreateDocVersion(ctx, &store.DocVersionRecord{
				DocID: d.ID, Rev: updated.CurrentVersion, Content: rendered.Content, Trigger: "sync",
			}); err != nil {
				return fmt.Errorf("create doc version: %w", err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("regenerate doc %s: %w", d.ID, err)
		}
	}

	return nil
}

type templateData struct {
	ServiceName  string
	Type         string
	Sections     []connector.SnapshotSection
	Dependencies []connector.ServiceDependency
	Metadata     map[string]string
	Links        []EntityLink
	GeneratedAt  string
}
