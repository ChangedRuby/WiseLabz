package store

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// maxCachedQueries bounds the rewrite cache so dynamically built queries
// (e.g. variable-length IN lists) cannot grow it without limit.
const maxCachedQueries = 2048

var (
	placeholderCache     sync.Map // query string -> rewritten string
	placeholderCacheSize atomic.Int64
)

// rewritePlaceholders converts SQLite-style `?` positional placeholders to
// PostgreSQL-style `$1`, `$2`, ... placeholders. `?` characters inside
// single-quoted string literals, including SQL's doubled-quote escape for a
// literal quote character, are left untouched.
//
// Results are cached (up to maxCachedQueries distinct queries).
func rewritePlaceholders(query string) string {
	if !strings.Contains(query, "?") {
		return query
	}
	if v, ok := placeholderCache.Load(query); ok {
		return v.(string)
	}
	out := doRewritePlaceholders(query)
	if placeholderCacheSize.Load() < maxCachedQueries {
		if _, loaded := placeholderCache.LoadOrStore(query, out); !loaded {
			placeholderCacheSize.Add(1)
		}
	}
	return out
}

func doRewritePlaceholders(query string) string {

	var b strings.Builder
	b.Grow(len(query) + 8)

	inString := false
	n := 0
	for i := 0; i < len(query); i++ {
		c := query[i]
		switch {
		case c == '\'':
			inString = !inString
			b.WriteByte(c)
		case c == '?' && !inString:
			n++
			b.WriteByte('$')
			b.WriteString(strconv.Itoa(n))
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}
