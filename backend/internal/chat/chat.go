// Package chat implements retrieval-augmented "ask your lab" chat: splitting
// generated docs into embeddable sections, syncing their embeddings, and
// ranking sections by similarity to a question.
package chat

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Section is a chunk of a doc's rendered Markdown, split on "## " headings —
// the same granularity the doc engine's templates render at.
type Section struct {
	Key     string // heading text; "" for content before the first heading
	Content string
}

// SplitSections splits rendered doc Markdown into per-heading chunks.
func SplitSections(content string) []Section {
	var sections []Section
	var key string
	var body strings.Builder

	flush := func() {
		text := strings.TrimSpace(body.String())
		if text != "" {
			sections = append(sections, Section{Key: key, Content: text})
		}
		body.Reset()
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if after, ok := strings.CutPrefix(line, "## "); ok {
			flush()
			key = strings.TrimSpace(after)
			continue
		}
		body.WriteString(line)
		body.WriteByte('\n')
	}
	flush()

	return sections
}

// packVector encodes a float32 vector as a little-endian byte slice for BLOB/BYTEA storage.
func packVector(v []float32) []byte {
	buf := make([]byte, 4*len(v))
	for i, f := range v {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(f))
	}
	return buf
}

// unpackVector decodes a byte slice produced by packVector back into a float32 vector.
func unpackVector(b []byte) []float32 {
	v := make([]float32, len(b)/4)
	for i := range v {
		v[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[i*4:]))
	}
	return v
}

// cosineSimilarity returns the cosine similarity of two equal-length vectors, or 0 if either is zero-length/zero-norm.
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// SyncDocEmbeddings re-splits a doc's rendered content into sections and
// (re)computes their embeddings, replacing whatever was previously stored for
// this doc. Called eagerly whenever a doc is generated or saved so retrieval
// is always current.
func SyncDocEmbeddings(ctx context.Context, s *store.Store, embedder ai.Embedder, model, docID, content string) error {
	// Invalidate after the write (and on failure, harmlessly) so retrieval
	// never keeps serving vectors decoded from the replaced rows.
	defer sharedVectorCache.invalidateDoc(docID)

	sections := SplitSections(content)
	if len(sections) == 0 {
		return s.DeleteDocSectionEmbeddings(ctx, docID)
	}

	texts := make([]string, len(sections))
	for i, sec := range sections {
		texts[i] = sec.Content
	}
	// Embed before touching stored rows so a failed embed keeps the old ones.
	vectors, err := embedder.Embed(ctx, texts)
	if err != nil {
		return fmt.Errorf("embed doc sections: %w", err)
	}
	if len(vectors) != len(sections) {
		return fmt.Errorf("embed doc sections: got %d vectors for %d sections", len(vectors), len(sections))
	}

	return s.WithinTransaction(ctx, func(tx *store.Store) error {
		if err := tx.DeleteDocSectionEmbeddings(ctx, docID); err != nil {
			return err
		}
		for i, sec := range sections {
			key := sec.Key
			if key == "" {
				key = fmt.Sprintf("section-%d", i)
			}
			if err := tx.UpsertDocSectionEmbedding(ctx, &store.DocSectionEmbeddingRecord{
				DocID:      docID,
				SectionKey: key,
				Content:    sec.Content,
				Vector:     packVector(vectors[i]),
				Model:      model,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

// Match is a retrieved doc section ranked by relevance to a question.
type Match struct {
	DocID      string
	SectionKey string
	Content    string
	Score      float64
}

// Retrieve embeds question and returns the topN most similar doc sections.
// When scopeDocID is non-empty, retrieval is limited to that doc (doc scope);
// otherwise it ranks across every doc userID may view (lab scope).
func Retrieve(ctx context.Context, s *store.Store, embedder ai.Embedder, question, userID, scopeDocID string, topN int) ([]Match, error) {
	vectors, err := embedder.Embed(ctx, []string{question})
	if err != nil {
		return nil, fmt.Errorf("embed question: %w", err)
	}
	queryVector := vectors[0]

	gen := sharedVectorCache.generation()
	rows, err := s.ListDocSectionEmbeddings(ctx, userID, scopeDocID)
	if err != nil {
		return nil, err
	}

	matches := make([]Match, 0, len(rows))
	for _, row := range rows {
		key := vectorKey{row.DocID, row.SectionKey}
		vec, ok := sharedVectorCache.get(key)
		if !ok {
			vec = unpackVector(row.Vector)
			sharedVectorCache.put(key, vec, gen)
		}
		matches = append(matches, Match{
			DocID:      row.DocID,
			SectionKey: row.SectionKey,
			Content:    row.Content,
			Score:      cosineSimilarity(queryVector, vec),
		})
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i].Score > matches[j].Score })
	if len(matches) > topN {
		matches = matches[:topN]
	}
	return matches, nil
}
