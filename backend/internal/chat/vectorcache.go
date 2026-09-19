package chat

import (
	"container/list"
	"sync"
)

// maxCachedVectors bounds the number of decoded section vectors held in
// memory (~6 KiB each at 1536 dims, so roughly 25 MiB at the cap).
const maxCachedVectors = 4096

type vectorKey struct{ docID, sectionKey string }

type vectorEntry struct {
	key    vectorKey
	vector []float32
}

// vectorCache is a bounded, concurrency-safe LRU of decoded embedding
// vectors. It is per-process: other server replicas keep their own copy and
// only see changes made elsewhere when the entry is evicted.
type vectorCache struct {
	mu      sync.Mutex
	limit   int
	ll      *list.List
	entries map[vectorKey]*list.Element
	// gen increments on every invalidation so a reader that decoded rows
	// loaded before a concurrent sync cannot store stale vectors afterwards.
	gen uint64
}

func newVectorCache(limit int) *vectorCache {
	return &vectorCache{limit: limit, ll: list.New(), entries: make(map[vectorKey]*list.Element)}
}

// sharedVectorCache is used by Retrieve and invalidated by SyncDocEmbeddings.
var sharedVectorCache = newVectorCache(maxCachedVectors)

// generation returns the current invalidation generation; pass it to put.
func (c *vectorCache) generation() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.gen
}

func (c *vectorCache) get(k vectorKey) ([]float32, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.entries[k]
	if !ok {
		return nil, false
	}
	c.ll.MoveToFront(el)
	return el.Value.(*vectorEntry).vector, true
}

// put stores v unless an invalidation happened since gen was read.
func (c *vectorCache) put(k vectorKey, v []float32, gen uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if gen != c.gen {
		return
	}
	if el, ok := c.entries[k]; ok {
		el.Value.(*vectorEntry).vector = v
		c.ll.MoveToFront(el)
		return
	}
	c.entries[k] = c.ll.PushFront(&vectorEntry{key: k, vector: v})
	for c.ll.Len() > c.limit {
		oldest := c.ll.Back()
		c.ll.Remove(oldest)
		delete(c.entries, oldest.Value.(*vectorEntry).key)
	}
}

// invalidateDoc drops every cached vector for docID.
func (c *vectorCache) invalidateDoc(docID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.gen++
	for k, el := range c.entries {
		if k.docID == docID {
			c.ll.Remove(el)
			delete(c.entries, k)
		}
	}
}
