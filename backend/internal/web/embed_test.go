package web

import (
	"io/fs"
	"testing"
)

func TestEmbeddedFrontendEntryPoint(t *testing.T) {
	// Both a source checkout's placeholder and a production build must work
	// with the same fs.Sub path used by cmd/server, without building web first.
	dist, err := fs.Sub(DistFS, "dist")
	if err != nil {
		t.Fatal(err)
	}
	index, err := fs.ReadFile(dist, "index.html")
	if err != nil {
		t.Fatal(err)
	}
	if len(index) == 0 {
		t.Fatal("embedded frontend entry point is empty")
	}
}
