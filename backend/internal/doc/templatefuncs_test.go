package doc

import (
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestDateFormat(t *testing.T) {
	tests := []struct {
		name      string
		layout    string
		input     any
		want      string
		wantError bool
	}{
		{
			name:   "time.Time",
			layout: "2006-01-02",
			input:  time.Date(2026, time.September, 5, 12, 30, 45, 0, time.UTC),
			want:   "2026-09-05",
		},
		{
			name:   "RFC3339 string",
			layout: "2006-01-02",
			input:  "2026-09-05T12:30:45Z",
			want:   "2026-09-05",
		},
		{
			name:   "RFC3339 with custom layout",
			layout: "Jan 02, 2006",
			input:  "2026-09-05T12:30:45Z",
			want:   "Sep 05, 2026",
		},
		{
			name:      "invalid RFC3339 string",
			layout:    "2006-01-02",
			input:     "not-a-date",
			wantError: true,
		},
		{
			name:   "unsupported type",
			layout: "2006-01-02",
			input:  12345,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dateFormat(tt.layout, tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("dateFormat() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("dateFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		maxChars int
		input    string
		want     string
	}{
		{
			name:     "no truncation needed",
			maxChars: 50,
			input:    "hello world",
			want:     "hello world",
		},
		{
			name:     "truncation at exact limit",
			maxChars: 11,
			input:    "hello world",
			want:     "hello world",
		},
		{
			name:     "truncation needed",
			maxChars: 10,
			input:    "hello world",
			want:     "hello worl...",
		},
		{
			name:     "single character over limit",
			maxChars: 5,
			input:    "hello world",
			want:     "hello...",
		},
		{
			name:     "zero max chars",
			maxChars: 0,
			input:    "hello",
			want:     "...",
		},
		{
			name:     "empty string",
			maxChars: 50,
			input:    "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.maxChars, tt.input)
			if got != tt.want {
				t.Errorf("truncate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		want      string
		wantError bool
	}{
		{
			name:  "map",
			input: map[string]string{"key": "value"},
			want:  `{"key":"value"}`,
		},
		{
			name:  "slice",
			input: []string{"a", "b", "c"},
			want:  `["a","b","c"]`,
		},
		{
			name:  "struct",
			input: struct{ Name string }{Name: "test"},
			want:  `{"Name":"test"}`,
		},
		{
			name:  "string",
			input: "hello",
			want:  `"hello"`,
		},
		{
			name:  "number",
			input: 42,
			want:  `42`,
		},
		{
			name:  "boolean",
			input: true,
			want:  `true`,
		},
		{
			name:  "nil",
			input: nil,
			want:  `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toJSON(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("toJSON() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("toJSON() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFilterByTitle(t *testing.T) {
	sections := []connector.SnapshotSection{
		{Title: "Configuration", Content: "config content"},
		{Title: "Status", Content: "status content"},
		{Title: "Configuration", Content: "more config"},
		{Title: "Logs", Content: "log content"},
	}

	tests := []struct {
		name  string
		title string
		want  int
	}{
		{
			name:  "exact match single",
			title: "Status",
			want:  1,
		},
		{
			name:  "exact match multiple",
			title: "Configuration",
			want:  2,
		},
		{
			name:  "no match",
			title: "NonExistent",
			want:  0,
		},
		{
			name:  "case sensitive no match",
			title: "configuration",
			want:  0,
		},
		{
			name:  "empty title string",
			title: "",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterByTitle(tt.title, sections)
			if len(got) != tt.want {
				t.Errorf("filterByTitle() returned %d sections, want %d", len(got), tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name string
		sep  string
		strs []string
		want string
	}{
		{
			name: "basic join",
			sep:  ", ",
			strs: []string{"a", "b", "c"},
			want: "a, b, c",
		},
		{
			name: "single element",
			sep:  ", ",
			strs: []string{"a"},
			want: "a",
		},
		{
			name: "empty slice",
			sep:  ", ",
			strs: []string{},
			want: "",
		},
		{
			name: "empty separator",
			sep:  "",
			strs: []string{"a", "b", "c"},
			want: "abc",
		},
		{
			name: "multi-char separator",
			sep:  " | ",
			strs: []string{"x", "y", "z"},
			want: "x | y | z",
		},
		{
			name: "empty strings",
			sep:  ",",
			strs: []string{"", "b", ""},
			want: ",b,",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := join(tt.sep, tt.strs)
			if got != tt.want {
				t.Errorf("join() = %q, want %q", got, tt.want)
			}
		})
	}
}
