package chat

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func TestVectorCacheBoundedLRU(t *testing.T) {
	c := newVectorCache(2)
	g := c.generation()
	c.put(vectorKey{"d", "a"}, []float32{1}, g)
	c.put(vectorKey{"d", "b"}, []float32{2}, g)
	c.get(vectorKey{"d", "a"}) // a is now most recent
	c.put(vectorKey{"d", "c"}, []float32{3}, g)
	if _, ok := c.get(vectorKey{"d", "b"}); ok {
		t.Fatal("expected b evicted")
	}
	if _, ok := c.get(vectorKey{"d", "a"}); !ok {
		t.Fatal("expected a retained")
	}
	if c.ll.Len() != 2 || len(c.entries) != 2 {
		t.Fatalf("cache exceeded bound: %d/%d", c.ll.Len(), len(c.entries))
	}
}

func TestVectorCacheInvalidateDocAndStalePut(t *testing.T) {
	c := newVectorCache(10)
	g := c.generation()
	c.put(vectorKey{"d1", "a"}, []float32{1}, g)
	c.put(vectorKey{"d2", "a"}, []float32{2}, g)
	c.invalidateDoc("d1")
	if _, ok := c.get(vectorKey{"d1", "a"}); ok {
		t.Fatal("d1 should be invalidated")
	}
	if _, ok := c.get(vectorKey{"d2", "a"}); !ok {
		t.Fatal("d2 should be untouched")
	}
	c.put(vectorKey{"d1", "a"}, []float32{9}, g) // stale generation
	if _, ok := c.get(vectorKey{"d1", "a"}); ok {
		t.Fatal("stale put must be dropped")
	}
}

func TestVectorCacheConcurrent(t *testing.T) {
	c := newVectorCache(8)
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				k := vectorKey{fmt.Sprintf("d%d", i%5), fmt.Sprintf("s%d", i%7)}
				if _, ok := c.get(k); !ok {
					c.put(k, []float32{float32(w)}, c.generation())
				}
				if i%50 == 0 {
					c.invalidateDoc(k.docID)
				}
			}
		}(w)
	}
	wg.Wait()
	if len(c.entries) > 8 {
		t.Fatalf("cache exceeded bound: %d", len(c.entries))
	}
}

func TestRetrieveUsesCacheAndSyncInvalidates(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", "file:"+t.TempDir()+"/test.db?cache=shared")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	if err := store.RunMigrations(db, "sqlite", slog.New(slog.NewTextHandler(io.Discard, nil))); err != nil {
		t.Fatalf("migrations: %v", err)
	}
	s := store.New(db, "sqlite")
	if err := s.CreateDoc(ctx, &store.DocRecord{ID: "cache-d1", Title: "d1", Kind: "lab", Content: "x"}); err != nil {
		t.Fatalf("CreateDoc: %v", err)
	}
	if err := SyncDocEmbeddings(ctx, s, stubEmbedder{}, "m", "cache-d1", "## A\n\none\n"); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if _, err := Retrieve(ctx, s, stubEmbedder{}, "q", "", "cache-d1", 5); err != nil {
		t.Fatalf("retrieve: %v", err)
	}
	key := vectorKey{"cache-d1", "A"}
	if _, ok := sharedVectorCache.get(key); !ok {
		t.Fatal("expected vector cached after Retrieve")
	}
	if err := SyncDocEmbeddings(ctx, s, stubEmbedder{}, "m", "cache-d1", "## A\n\nchanged\n"); err != nil {
		t.Fatalf("resync: %v", err)
	}
	if _, ok := sharedVectorCache.get(key); ok {
		t.Fatal("expected cache invalidated by SyncDocEmbeddings")
	}
}
