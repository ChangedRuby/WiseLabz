package ttlcache

import (
	"testing"
	"time"
)

func TestCacheExpires(t *testing.T) {
	c := New[int](20 * time.Millisecond)
	c.Set("k", 1)
	if v, ok := c.Get("k"); !ok || v != 1 {
		t.Fatalf("Get = %d,%v; want 1,true", v, ok)
	}
	time.Sleep(30 * time.Millisecond)
	if _, ok := c.Get("k"); ok {
		t.Fatal("entry should have expired")
	}
}
