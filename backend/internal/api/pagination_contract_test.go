package api_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// paginationWriters are the httputil helpers that own a pagination envelope.
// A handler that paginates must hand its page to one of these rather than
// hand-building the envelope, so the shape stays identical everywhere.
var paginationWriters = map[string]bool{
	"WritePaginated":       true,
	"WritePaginatedCursor": true,
	"WriteDataPaginated":   true,
}

// bareArrayListHandlers are the list handlers that legitimately do not write a
// pagination envelope, because their published response body is a bare array
// (see docs/openapi.yaml). They still page the query with httputil.Paginate;
// changing them to an envelope would break every existing client.
//
// Keep this list short and justified. A new list endpoint belongs in the
// envelope, not here.
var bareArrayListHandlers = map[string]string{
	// GET /connectors -> Connector[]; total rides on X-Total-Count.
	"connectors.Handler.List": "spec: bare Connector[]",
	// GET /users -> User[].
	"auth.Handler.ListUsers": "spec: bare User[]",
	// GET /templates -> Template[].
	"templates.Handler.List": "spec: bare Template[]",
}

// TestListHandlersUseSharedPaginationWriter is the CI guard for issue #267:
// every handler under internal/api that pages a query with httputil.Paginate
// must return it through one of the shared httputil pagination writers.
// Hand-building {items,total,page,pageSize} is how the envelope drifted in the
// first place, and a divergent envelope is invisible until a client breaks.
func TestListHandlersUseSharedPaginationWriter(t *testing.T) {
	root := "."
	fset := token.NewFileSet()

	var offenders []string
	checked := 0

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		pkg := file.Name.Name

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			calls := httputilCalls(fn.Body)
			if !calls["Paginate"] {
				continue
			}
			checked++

			name := pkg + "." + receiverName(fn) + "." + fn.Name.Name
			if _, allowed := bareArrayListHandlers[name]; allowed {
				continue
			}
			if !writesEnvelope(calls) {
				offenders = append(offenders, fmt.Sprintf("%s (%s:%d)", name, path, fset.Position(fn.Pos()).Line))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan internal/api: %v", err)
	}

	// A scan that suddenly finds nothing is a broken scan, not a clean repo.
	if checked < len(bareArrayListHandlers)+1 {
		t.Fatalf("only %d paginating handlers found; the AST scan is no longer matching httputil.Paginate calls", checked)
	}

	if len(offenders) > 0 {
		sort.Strings(offenders)
		t.Errorf("these handlers page a query but do not write the shared pagination envelope.\n"+
			"Use httputil.WritePaginated / WritePaginatedCursor (or WriteDataPaginated for the\n"+
			"legacy data-keyed envelope) instead of building the response by hand:\n  %s",
			strings.Join(offenders, "\n  "))
	}
}

// TestBareArrayAllowlistIsCurrent keeps the allowlist from outliving its
// entries: a handler listed there must still exist and still paginate.
func TestBareArrayAllowlistIsCurrent(t *testing.T) {
	found := map[string]bool{}
	fset := token.NewFileSet()

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			if httputilCalls(fn.Body)["Paginate"] {
				found[file.Name.Name+"."+receiverName(fn)+"."+fn.Name.Name] = true
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan internal/api: %v", err)
	}

	for name, why := range bareArrayListHandlers {
		if !found[name] {
			t.Errorf("bareArrayListHandlers entry %q (%s) no longer names a paginating handler; drop it", name, why)
		}
	}
}

// envelopeKeys are the pagination envelope's own keys. A handler that builds a
// map literal carrying these is rebuilding an envelope httputil already owns.
var envelopeKeys = []string{"items", "pageSize"}

// TestNoHandRolledPaginationEnvelopes catches the drift class that motivated
// issue #267 from the other direction: a handler assembling
// map[string]any{"items": ..., "pageSize": ...} by hand, whether or not it
// went through httputil.Paginate. WritePaginated is the only place the
// envelope's shape should be decided.
func TestNoHandRolledPaginationEnvelopes(t *testing.T) {
	fset := token.NewFileSet()
	var offenders []string

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.CompositeLit)
			if !ok || !hasAllStringKeys(lit, envelopeKeys) {
				return true
			}
			offenders = append(offenders, fmt.Sprintf("%s:%d", path, fset.Position(lit.Pos()).Line))
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("scan internal/api: %v", err)
	}

	if len(offenders) > 0 {
		sort.Strings(offenders)
		t.Errorf("pagination envelope built by hand; call httputil.WritePaginated instead:\n  %s",
			strings.Join(offenders, "\n  "))
	}
}

// hasAllStringKeys reports whether the composite literal has a string-literal
// key for every name in keys.
func hasAllStringKeys(lit *ast.CompositeLit, keys []string) bool {
	present := map[string]bool{}
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.BasicLit)
		if !ok || key.Kind != token.STRING {
			continue
		}
		present[strings.Trim(key.Value, `"`)] = true
	}
	for _, want := range keys {
		if !present[want] {
			return false
		}
	}
	return true
}

// httputilCalls returns the set of httputil.X function names called anywhere
// in body.
func httputilCalls(body *ast.BlockStmt) map[string]bool {
	calls := map[string]bool{}
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "httputil" {
			calls[sel.Sel.Name] = true
		}
		return true
	})
	return calls
}

func writesEnvelope(calls map[string]bool) bool {
	for name := range calls {
		if paginationWriters[name] {
			return true
		}
	}
	return false
}

// receiverName is the handler's receiver type name ("Handler"), or "" for a
// plain function.
func receiverName(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return ""
	}
	typ := fn.Recv.List[0].Type
	if star, ok := typ.(*ast.StarExpr); ok {
		typ = star.X
	}
	if ident, ok := typ.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}
