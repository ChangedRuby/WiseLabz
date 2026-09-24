package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DocSectionEmbeddingRecord represents a row in the doc_section_embeddings
// table: one embedded chunk (a doc section) used for chat retrieval.
type DocSectionEmbeddingRecord struct {
	DocID      string `json:"docId"`
	SectionKey string `json:"sectionKey"`
	Content    string `json:"content"`
	Vector     []byte `json:"-"`
	Model      string `json:"model"`
	UpdatedAt  string `json:"updatedAt"`
}

// UpsertDocSectionEmbedding inserts or replaces the embedding for one doc section.
func (s *Store) UpsertDocSectionEmbedding(ctx context.Context, rec *DocSectionEmbeddingRecord) error {
	if rec.UpdatedAt == "" {
		rec.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO doc_section_embeddings (doc_id, section_key, content, vector, model, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(doc_id, section_key) DO UPDATE SET
			content = excluded.content,
			vector = excluded.vector,
			model = excluded.model,
			updated_at = excluded.updated_at
	`, rec.DocID, rec.SectionKey, rec.Content, rec.Vector, rec.Model, rec.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert doc section embedding: %w", err)
	}
	return nil
}

// DeleteDocSectionEmbeddings removes all embedded sections for a doc, so a
// resync can recreate them from scratch (section keys/counts may change
// between regenerations).
func (s *Store) DeleteDocSectionEmbeddings(ctx context.Context, docID string) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM doc_section_embeddings WHERE doc_id = ?`, docID); err != nil {
		return fmt.Errorf("delete doc section embeddings: %w", err)
	}
	return nil
}

// ListDocSectionEmbeddings returns embedded sections for retrieval, limited to
// docs userID may view: lab-wide docs (no service) plus docs on connectors the
// user holds a grant on. An empty docID lists across all visible docs (lab
// scope); a non-empty docID filters to that doc (doc scope).
func (s *Store) ListDocSectionEmbeddings(ctx context.Context, userID, docID string) ([]DocSectionEmbeddingRecord, error) {
	query := `SELECT e.doc_id, e.section_key, e.content, e.vector, e.model, e.updated_at
		FROM doc_section_embeddings e
		JOIN docs d ON d.id = e.doc_id
		WHERE (d.service_id IS NULL OR d.service_id = ''
			OR (d.service_id IN (SELECT connector_id FROM user_connector_roles WHERE user_id = ?)%s))`
	keyFilter, keyArgs := apiKeyConnectorFilter(ctx, "d.service_id")
	query = fmt.Sprintf(query, keyFilter)
	args := append([]any{userID}, keyArgs...)
	if docID != "" {
		query += ` AND e.doc_id = ?`
		args = append(args, docID)
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list doc section embeddings: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	sections := []DocSectionEmbeddingRecord{}
	for rows.Next() {
		var rec DocSectionEmbeddingRecord
		if err := rows.Scan(&rec.DocID, &rec.SectionKey, &rec.Content, &rec.Vector, &rec.Model, &rec.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		sections = append(sections, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate doc section embeddings: %w", err)
	}
	return sections, nil
}

// ChatConversationRecord represents a row in the chat_conversations table.
type ChatConversationRecord struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	ScopeType string `json:"scopeType"`
	ScopeID   string `json:"scopeId"`
	CreatedAt string `json:"createdAt"`
}

// CreateChatConversation inserts a new chat conversation.
func (s *Store) CreateChatConversation(ctx context.Context, c *ChatConversationRecord) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.CreatedAt == "" {
		c.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO chat_conversations (id, user_id, scope_type, scope_id, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, c.ID, c.UserID, c.ScopeType, nilToStr(c.ScopeID), c.CreatedAt)
	if err != nil {
		return fmt.Errorf("create chat conversation: %w", err)
	}
	return nil
}

// GetChatConversation retrieves a single conversation by ID.
func (s *Store) GetChatConversation(ctx context.Context, id string) (*ChatConversationRecord, error) {
	c := &ChatConversationRecord{}
	var scopeID sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, scope_type, scope_id, created_at FROM chat_conversations WHERE id = ?
	`, id).Scan(&c.ID, &c.UserID, &c.ScopeType, &scopeID, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get chat conversation: %w", err)
	}
	c.ScopeID = scopeID.String
	return c, nil
}

// ListChatConversationsByUser returns a user's conversations, newest first.
func (s *Store) ListChatConversationsByUser(ctx context.Context, userID string) ([]ChatConversationRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, scope_type, scope_id, created_at FROM chat_conversations
		WHERE user_id = ? ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list chat conversations: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	conversations := []ChatConversationRecord{}
	for rows.Next() {
		var c ChatConversationRecord
		var scopeID sql.NullString
		if err := rows.Scan(&c.ID, &c.UserID, &c.ScopeType, &scopeID, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		c.ScopeID = scopeID.String
		conversations = append(conversations, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate chat conversations: %w", err)
	}
	return conversations, nil
}

// ChatMessageRecord represents a row in the chat_messages table.
type ChatMessageRecord struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversationId"`
	Role           string `json:"role"`
	Content        string `json:"content"`
	Provider       string `json:"provider"`
	CreatedAt      string `json:"createdAt"`
}

// CreateChatMessage inserts a new chat message.
func (s *Store) CreateChatMessage(ctx context.Context, m *ChatMessageRecord) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.CreatedAt == "" {
		m.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO chat_messages (id, conversation_id, role, content, provider, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, m.ID, m.ConversationID, m.Role, m.Content, m.Provider, m.CreatedAt)
	if err != nil {
		return fmt.Errorf("create chat message: %w", err)
	}
	return nil
}

// ListChatMessages returns a conversation's messages, oldest first.
func (s *Store) ListChatMessages(ctx context.Context, conversationID string) ([]ChatMessageRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, conversation_id, role, content, provider, created_at FROM chat_messages
		WHERE conversation_id = ? ORDER BY created_at ASC
	`, conversationID)
	if err != nil {
		return nil, fmt.Errorf("list chat messages: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	messages := []ChatMessageRecord{}
	for rows.Next() {
		var m ChatMessageRecord
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.Provider, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate chat messages: %w", err)
	}
	return messages, nil
}
