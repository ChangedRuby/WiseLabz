package main

import (
	"io/fs"
	"net/http/httptest"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api"
	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/web"
)

func TestEmbeddedSPAWithoutFrontendBuild(t *testing.T) {
	// Use the same embedded subtree and router wiring as main. The committed
	// placeholder must serve both the root and a client-side route on checkout.
	files, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		t.Fatal(err)
	}
	index, err := fs.ReadFile(files, "index.html")
	if err != nil {
		t.Fatal(err)
	}
	if len(index) == 0 {
		t.Fatal("empty embedded entry point")
	}
	cfg := &config.Config{}
	cfg.Server.Embed = true
	router := api.NewRouter(api.Config{Store: apitest.NewStore(t), JWT: apitest.JWTService(), Config: cfg, SPAFiles: files})
	for _, path := range []string{"/", "/services/fixture"} {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, httptest.NewRequest("GET", path, nil))
		if rr.Code != 200 || rr.Body.String() != string(index) {
			t.Fatalf("%s: status=%d body=%s", path, rr.Code, rr.Body.String())
		}
	}
}
