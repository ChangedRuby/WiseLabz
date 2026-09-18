package connector

import (
	"bytes"
	"testing"
)

func TestReadBodyLimit(t *testing.T) {
	if _, err := ReadBody(bytes.NewReader(make([]byte, MaxResponseBytes))); err != nil {
		t.Fatalf("at limit: %v", err)
	}
	if _, err := ReadBody(bytes.NewReader(make([]byte, MaxResponseBytes+1))); err == nil {
		t.Fatal("expected error over limit")
	}
}
