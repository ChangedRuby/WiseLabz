package store

import "testing"

func TestRewritePlaceholders(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no placeholders",
			in:   "SELECT COUNT(*) FROM users",
			want: "SELECT COUNT(*) FROM users",
		},
		{
			name: "single placeholder",
			in:   "SELECT * FROM users WHERE id = ?",
			want: "SELECT * FROM users WHERE id = $1",
		},
		{
			name: "multiple placeholders",
			in:   "INSERT INTO users (id, username) VALUES (?, ?)",
			want: "INSERT INTO users (id, username) VALUES ($1, $2)",
		},
		{
			name: "placeholder adjacent to punctuation",
			in:   "UPDATE users SET role=? WHERE id=?",
			want: "UPDATE users SET role=$1 WHERE id=$2",
		},
		{
			name: "literal question mark in string is untouched",
			in:   "SELECT * FROM docs WHERE title = 'what?' AND id = ?",
			want: "SELECT * FROM docs WHERE title = 'what?' AND id = $1",
		},
		{
			name: "escaped quote inside string literal",
			in:   "SELECT * FROM docs WHERE title = 'it''s a ? test' AND id = ?",
			want: "SELECT * FROM docs WHERE title = 'it''s a ? test' AND id = $1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rewritePlaceholders(tc.in)
			if got != tc.want {
				t.Errorf("rewritePlaceholders(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestRewritePlaceholdersCached(t *testing.T) {
	q := "SELECT * FROM cache_probe WHERE a = ? AND b = ?"
	want := "SELECT * FROM cache_probe WHERE a = $1 AND b = $2"
	for i := 0; i < 2; i++ {
		if got := rewritePlaceholders(q); got != want {
			t.Fatalf("call %d: got %q, want %q", i, got, want)
		}
	}
	if v, ok := placeholderCache.Load(q); !ok || v.(string) != want {
		t.Errorf("expected rewritten query to be cached, got %v %v", v, ok)
	}
}
