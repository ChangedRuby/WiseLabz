package store

import (
	"context"
	"testing"
)

func seedScopeFixture(t *testing.T) (*Store, string, string, string) {
	t.Helper()
	ctx := context.Background()
	s := newDocTestStore(t)

	user := &User{Username: "viewer", InstanceAdminRole: "user", AuthSource: "local"}
	if err := s.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	granted := &ConnectorRecord{Name: "granted", Category: "networking", Type: "custom", Enabled: true}
	denied := &ConnectorRecord{Name: "denied", Category: "networking", Type: "custom", Enabled: true}
	for _, c := range []*ConnectorRecord{granted, denied} {
		if err := s.CreateConnector(ctx, c); err != nil {
			t.Fatalf("CreateConnector: %v", err)
		}
		if err := s.CreateAlert(ctx, &AlertRecord{ServiceID: c.ID, Severity: "critical", Title: c.Name, Status: "pending"}); err != nil {
			t.Fatalf("CreateAlert: %v", err)
		}
	}
	if _, err := s.UpsertConnectorGrant(ctx, user.ID, granted.ID, "viewer"); err != nil {
		t.Fatalf("UpsertConnectorGrant: %v", err)
	}
	return s, user.ID, granted.ID, denied.ID
}

func TestMergedAttentionItemsFiltersByGrant(t *testing.T) {
	s, userID, grantedID, _ := seedScopeFixture(t)

	items, total, err := s.MergedAttentionItems(context.Background(), userID, "", 0, 50)
	if err != nil {
		t.Fatalf("MergedAttentionItems: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].ConnectorID != grantedID {
		t.Fatalf("want only the granted connector's item, got total=%d items=%+v", total, items)
	}

	items, total, err = s.MergedAttentionItems(context.Background(), "nobody", "", 0, 50)
	if err != nil {
		t.Fatalf("MergedAttentionItems: %v", err)
	}
	if total != 0 || len(items) != 0 {
		t.Fatalf("user without grants must see nothing, got %+v", items)
	}
}

func TestListDocSectionEmbeddingsFiltersByGrant(t *testing.T) {
	ctx := context.Background()
	s, userID, grantedID, deniedID := seedScopeFixture(t)

	for docID, svc := range map[string]string{"doc-granted": grantedID, "doc-denied": deniedID, "doc-lab": ""} {
		d := &DocRecord{ID: docID, Title: docID, Kind: "service", ServiceID: svc, Content: "x"}
		if svc == "" {
			d.Kind = "lab"
		}
		if err := s.CreateDoc(ctx, d); err != nil {
			t.Fatalf("CreateDoc %s: %v", docID, err)
		}
		if err := s.UpsertDocSectionEmbedding(ctx, &DocSectionEmbeddingRecord{DocID: docID, SectionKey: "s", Content: docID, Vector: []byte{0}, Model: "m"}); err != nil {
			t.Fatalf("UpsertDocSectionEmbedding %s: %v", docID, err)
		}
	}

	got, err := s.ListDocSectionEmbeddings(ctx, userID, "")
	if err != nil {
		t.Fatalf("ListDocSectionEmbeddings: %v", err)
	}
	seen := map[string]bool{}
	for _, r := range got {
		seen[r.DocID] = true
	}
	if !seen["doc-granted"] || !seen["doc-lab"] || seen["doc-denied"] {
		t.Fatalf("unexpected visible docs: %v", seen)
	}

	got, err = s.ListDocSectionEmbeddings(ctx, userID, "doc-denied")
	if err != nil {
		t.Fatalf("ListDocSectionEmbeddings: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("doc scope must not bypass the grant filter, got %+v", got)
	}
}
