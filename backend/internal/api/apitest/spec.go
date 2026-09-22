package apitest

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

// specRouter is the parsed docs/openapi.yaml, built once per test binary:
// loading and validating the spec is expensive relative to a single
// assertion, and every caller wants the same document.
var (
	specOnce   sync.Once
	specRoutes routers.Router
	specErr    error
)

// specPath locates docs/openapi.yaml relative to this source file, so the
// helper works from whichever internal/api/* package directory `go test`
// happens to run in.
func specPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "docs", "openapi.yaml")
}

func loadSpec() (routers.Router, error) {
	specOnce.Do(func() {
		loader := &openapi3.Loader{IsExternalRefsAllowed: true}
		doc, err := loader.LoadFromFile(specPath())
		if err != nil {
			specErr = err
			return
		}
		if err := doc.Validate(loader.Context); err != nil {
			specErr = err
			return
		}
		specRoutes, specErr = gorillamux.NewRouter(doc)
	})
	return specRoutes, specErr
}

// AssertMatchesSpec checks resp against the response schema docs/openapi.yaml
// declares for req's operation, failing the test on any divergence: an
// undocumented path or method, an undocumented status code, or a body that
// does not satisfy the declared schema.
//
// It is deliberately opt-in. internal/api/* tests adopt it one at a time; the
// shared harness does not impose it, so adding an endpoint never turns the
// whole suite red at once. That matters while success payloads that predate
// this helper still diverge from the spec; the connectors package has since
// adopted it on its success paths too.
//
// req must carry the path as the router serves it, including the /api or
// /api/v1 prefix. resp.Body is fully read and replaced with an equivalent
// reader, so callers can still inspect it afterwards.
func AssertMatchesSpec(t *testing.T, req *http.Request, resp *http.Response) {
	t.Helper()

	router, err := loadSpec()
	if err != nil {
		t.Fatalf("load %s: %v", specPath(), err)
	}
	route, pathParams, err := router.FindRoute(req)
	if err != nil {
		t.Fatalf("no spec operation for %s %s: %v", req.Method, req.URL.Path, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))

	input := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParams,
			Route:      route,
			Options: &openapi3filter.Options{
				// Handler tests exercise auth through the real middleware;
				// re-checking securitySchemes here would only duplicate that.
				AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
			},
		},
		Status: resp.StatusCode,
		Header: resp.Header,
		Body:   io.NopCloser(bytes.NewReader(body)),
		Options: &openapi3filter.Options{
			// Without this, kin-openapi silently passes any status the
			// operation does not document -- which is precisely the drift
			// this helper exists to catch.
			IncludeResponseStatus: true,
			MultiError:            true,
		},
	}
	if err := openapi3filter.ValidateResponse(req.Context(), input); err != nil {
		t.Errorf("%s %s -> %d does not match docs/openapi.yaml: %v\nbody: %s",
			req.Method, req.URL.Path, resp.StatusCode, err, body)
	}
}
