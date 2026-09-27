package logsafe

import (
	"errors"
	"testing"
)

func TestSanitize(t *testing.T) {
	got := Sanitize("evil\r\nINJECTED: fake line")
	want := "evilINJECTED: fake line"
	if got != want {
		t.Fatalf("Sanitize() = %q, want %q", got, want)
	}
}

func TestErr(t *testing.T) {
	if got := Err(nil); got != "" {
		t.Fatalf("Err(nil) = %q, want empty", got)
	}
	got := Err(errors.New("get job health \"x\r\nINJECTED\": boom"))
	want := "get job health \"xINJECTED\": boom"
	if got != want {
		t.Fatalf("Err() = %q, want %q", got, want)
	}
}
