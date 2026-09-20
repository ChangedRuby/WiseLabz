package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// Exercise the production router: handlers alone cannot enforce route middleware.
func TestTemplateMutationRoleMatrix(t *testing.T) {
	app := newTestApp(t)
	_, admin := app.user(t, "operator")
	_, viewer := app.user(t, "viewer")
	tmpl := seedTemplate(t, app, admin, nil)
	for _, role := range []struct {
		name, token string
		denied      int
	}{{"anonymous", "", 401}, {"viewer", viewer, 403}, {"admin", admin, 0}} {
		for _, route := range []struct {
			method, path string
			body         any
			success      int
		}{
			{"POST", "/api/templates", map[string]any{"name": "Created"}, 201},
			{"PUT", "/api/templates/" + tmpl.ID, map[string]any{"name": "Updated"}, 200},
			{"POST", "/api/templates/" + tmpl.ID + "/versions/1/restore", nil, 200},
		} {
			t.Run(role.name+"/"+route.method+route.path, func(t *testing.T) {
				want := role.denied
				if want == 0 {
					want = route.success
				}
				rr := app.req(t, route.method, route.path, route.body, role.token)
				if rr.Code != want {
					t.Fatalf("status=%d want=%d: %s", rr.Code, want, rr.Body.String())
				}
			})
		}
	}
}

func TestConnectorGrantRouteMatrix(t *testing.T) {
	app := newTestApp(t)
	c := &store.ConnectorRecord{Name: "Private", Category: "networking", Type: "custom", URL: "https://example.com", ConfigData: "{}"}
	if err := app.Store.CreateConnector(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	for _, role := range []struct {
		name, grant          string
		admin, authenticated bool
	}{
		{"anonymous", "", false, false}, {"ungranted user", "", false, true}, {"ungranted admin", "", true, true}, {"viewer", "viewer", false, true}, {"operator", "operator", false, true}, {"admin viewer", "viewer", true, true}, {"admin operator", "operator", true, true},
	} {
		t.Run(role.name, func(t *testing.T) {
			token := ""
			if role.authenticated {
				instanceRole := "viewer"
				if role.admin {
					instanceRole = "operator"
				}
				var user string
				user, token = app.user(t, instanceRole)
				if role.grant != "" {
					app.connectorGrant(t, user, c.ID, role.grant)
				}
			}
			for _, route := range []struct {
				method, suffix string
				write          bool
				body           any
			}{
				{"GET", "", false, nil}, {"GET", "/data", false, nil}, {"GET", "/syncs", false, nil},
				{"PUT", "", true, map[string]any{"name": "Private"}}, {"PUT", "/enabled", true, map[string]any{"enabled": false}},
			} {
				want := http.StatusOK
				if !role.authenticated {
					want = 401
				} else if route.write && role.grant != "operator" {
					want = 403
				} else if role.grant == "" {
					want = 403
					if route.suffix == "" {
						want = 404
					}
				}
				rr := app.req(t, route.method, "/api/connectors/"+c.ID+route.suffix, route.body, token)
				if rr.Code != want {
					t.Errorf("%s %s: status=%d want=%d: %s", route.method, route.suffix, rr.Code, want, rr.Body.String())
				}
			}
		})
	}
}
