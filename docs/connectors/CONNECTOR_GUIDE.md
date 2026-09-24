# WiseLabz Connector Guide

A connector is a Go package that implements the `Connector` interface. It tells
WiseLabz how to talk to a specific service, what data to fetch, and how to shape
that data into a documentation-ready snapshot.

This guide walks you through building one from scratch.

---

## The Connector interface

Every connector implements this interface, defined in `backend/internal/connector/connector.go`:

```go
package connector

import "context"

type Connector interface {
    Name() string
    Type() string
    Category() string
    Fetch(ctx context.Context, config map[string]any) (*ServiceSnapshot, error)
    Validate(ctx context.Context, config map[string]any) error
}
```

| Method                  | Purpose                                                                                                                                                                 |
|-------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Name()`                | Display name, e.g. `"Proxmox VE"`.                                                                                                                                      |
| `Type()`                | Stable type identifier, e.g. `"proxmox"`. Must match the string used during registration.                                                                               |
| `Category()`            | Grouping label, e.g. `"virtualization"`, `"networking"`.                                                                                                                |
| `Fetch(ctx, config)`    | Connect to the service, pull the data you care about, and return a `ServiceSnapshot`. The context carries a deadline — don't ignore it.                                 |
| `Validate(ctx, config)` | Test that the connector's configuration is usable before we attempt a `Fetch`. Return `nil` if everything looks good, or an error describing what's missing or invalid. |

Registration happens via `init()` using `connector.Register()` — see [Step 4: Register the connector](#4-register-the-connector) below.

## The ServiceSnapshot

`Fetch` returns a `ServiceSnapshot`:

```go
type ServiceSnapshot struct {
    ServiceName  string               // matches connector Name()
    Type         string               // e.g. "proxmox", "portainer", "pfsense"
    Sections     []SnapshotSection    // the documentation sections this service contributes
    Metadata     map[string]string    // arbitrary key-value pairs (version, endpoint, etc.)
    Dependencies []ServiceDependency  // operational dependencies this service relies on
    FetchedAt    time.Time            // set by the caller, but you can set it if needed
}

type SnapshotSection struct {
    Title   string // section heading in the generated doc
    Content string // markdown body — tables, lists, code fences are all fine
}

type ServiceDependency struct {
    Kind string // host | network | storage | upstream_service
    Name string
    Ref  string // optional connector/service ID if known
}
```

Guidelines for a good snapshot:

- **One section per logical unit.** For Proxmox, that might be VMs, storage,
  and network — three sections. Don't dump everything into one blob.
- **Content is Markdown.** You control the formatting. Use GFM tables for
  structured data, fenced code blocks for config dumps, and lists for
  enumerations.
- **Metadata is optional but useful.** Stick the service version, API endpoint,
  or connection mode in `Metadata`. It surfaces in the dashboard without
  bloating the generated doc.
- **Keep it read-only by design.** WiseLabz connectors never mutate the
  target service. If your `Fetch` changes state, it's wrong.

### Dependencies

`Dependencies` is for the operational dependencies a service relies on to run
— not its own configuration or state, which belong in `Sections`/`Metadata`.
Populate it when your connector knows (or can infer) what the service sits
on top of or calls out to. Each entry has a `Kind`:

- **`host`** — the physical or virtual host the service runs on.
- **`network`** — a network or VLAN it's reachable through.
- **`storage`** — a volume, datastore, or share it consumes.
- **`upstream_service`** — another service it calls (e.g. a database, an API).

Set `Ref` when you know the connector/service ID for the dependency — it lets
the dashboard link the two together instead of showing a plain label.

```go
Dependencies: []connector.ServiceDependency{
    {Kind: "host", Name: "pve-node-01"},
    {Kind: "network", Name: "vmbr0"},
    {Kind: "storage", Name: "local-zfs"},
    {Kind: "upstream_service", Name: "Postgres", Ref: "postgres-primary"},
},
```

## Step by step

### 1. Create the package

```bash
mkdir -p backend/internal/connector/mynewservice
touch backend/internal/connector/mynewservice/mynewservice.go
```

The package name should be the service name, lowercased, no hyphens
(`truenas`, not `true-nas`).

### 2. Define your config schema

Connectors self-describe their config via `connector.TypeSchema`. This is what
the frontend renders as a form — no separate config struct is required unless
your connector needs one internally:

```go
package mynewservice

import "github.com/WiseLabz/wiselabz/internal/connector"

const typeName = "mynewservice"

func init() {
    connector.Register(connector.TypeSchema{
        Type:     typeName,
        Category: "infrastructure",   // groups connectors in the UI
        Name:     "My New Service",
        Fields: []connector.SchemaField{
            {Key: "url", Label: "API URL", Type: "text", Required: true, Placeholder: "https://..."},
            {Key: "api_key", Label: "API Key", Type: "password", Required: true},
            {Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Default: "true"},
        },
    }, func(config map[string]any) (connector.Connector, error) {
        // The config map contains all field values keyed by their Key above,
        // plus "url" and "verify_tls" injected automatically from the connector record.
        return New(config)
    })
}
```

Supported field types: `"text"`, `"password"`, `"number"`, `"select"`, `"toggle"`.

Secret fields should use `"password"` — the UI renders a password input and the
value is stored alongside other config.

### 3. Implement the interface

```go
type Connector struct {
    url    string
    apiKey string
    client *http.Client
}

func New(config map[string]any) (*Connector, error) {
    url, _ := config["url"].(string)
    apiKey, _ := config["api_key"].(string)
    if apiKey == "" {
        return nil, fmt.Errorf("mynewservice: API key is required")
    }
    verifyTLS := true
    if v, ok := config["verify_tls"]; ok {
        if b, ok := v.(bool); ok { verifyTLS = b }
    }
    return &Connector{
        url:    strings.TrimSuffix(url, "/"),
        apiKey: apiKey,
        // Always use the shared client: it blocks loopback/link-local
        // targets, refuses redirects, enforces TLS 1.2+ and a 30s timeout,
        // and retries idempotent requests on 429/502/503/504.
        client: connector.NewHTTPClient(connector.HTTPClientOptions{SkipTLSVerify: !verifyTLS}),
    }, nil
}

func (c *Connector) Name() string     { return "My New Service" }
func (c *Connector) Type() string     { return typeName }
func (c *Connector) Category() string { return "infrastructure" }

func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
    if c.url == "" {
        return fmt.Errorf("mynewservice: URL is required")
    }
    if c.apiKey == "" {
        return fmt.Errorf("mynewservice: APIKey is required")
    }
    return nil
}

func (c *Connector) Fetch(ctx context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
    // 1. Build the request
    req, err := http.NewRequestWithContext(ctx, "GET", c.url+"/api/v1/status", nil)
    if err != nil {
        return nil, fmt.Errorf("mynewservice: %w", err)
    }
    req.Header.Set("Authorization", "Bearer "+c.apiKey)

    // 2. Execute. MapTransportError turns a timeout into a TimeoutError and
    // CheckStatus maps 401/403 to AuthError and 502/503/504 to
    // ServiceUnavailableError, so the UI can tell the failures apart.
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, connector.MapTransportError(err)
    }
    defer resp.Body.Close()
    data, err := connector.ReadBody(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("mynewservice: read response: %w", err)
    }
    if err := connector.CheckStatus(resp.StatusCode, data); err != nil {
        return nil, err
    }

    // 3. Parse into your domain types
    var status MyServiceStatus
    if err := json.Unmarshal(data, &status); err != nil {
        return nil, connector.NewMalformedResponseError(fmt.Errorf("mynewservice: decode: %w", err))
    }

    // 4. Shape into a ServiceSnapshot
    return &connector.ServiceSnapshot{
        ServiceName: c.Name(),
        Type:        typeName,
        Sections: []connector.SnapshotSection{
            {
                Title:   "Status",
                Content: formatStatusSection(status),
            },
        },
        Metadata: map[string]string{
            "version":  status.Version,
            "endpoint": c.url,
        },
        FetchedAt: time.Now(),
    }, nil
}
```

### 4. Register the connector

Registration is handled in the same file via `init()` — see step 2 above.
The `connector.Register(schema, factory)` call registers both the config
schema (used by the frontend form) and the factory (used to construct instances).

### 5. Add the barrel import

In `backend/internal/connector/all/all.go`, add one blank import line:

```go
package all

import (
    _ "github.com/WiseLabz/wiselabz/internal/connector/adguardhome"
    _ "github.com/WiseLabz/wiselabz/internal/connector/custom"
    _ "github.com/WiseLabz/wiselabz/internal/connector/docker"
    _ "github.com/WiseLabz/wiselabz/internal/connector/home_assistant"
    _ "github.com/WiseLabz/wiselabz/internal/connector/mynewservice" // <-- add this
    _ "github.com/WiseLabz/wiselabz/internal/connector/opnsense"
    _ "github.com/WiseLabz/wiselabz/internal/connector/pfsense"
    _ "github.com/WiseLabz/wiselabz/internal/connector/portainer"
    _ "github.com/WiseLabz/wiselabz/internal/connector/proxmox"
    _ "github.com/WiseLabz/wiselabz/internal/connector/traefik"
    _ "github.com/WiseLabz/wiselabz/internal/connector/truenas"
    _ "github.com/WiseLabz/wiselabz/internal/connector/unifi"
)
```

The barrel import triggers `init()` in your package, which calls
`connector.Register`. The router imports `connector/all` once — no other
files need to change.

### 6. Write tests

```go
func TestConnector_Fetch(t *testing.T) {
    // Start a test HTTP server that mimics your service's API
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"version":"1.2.3","status":"ok"}`))
    }))
    defer srv.Close()

    c, err := New(map[string]any{"url": srv.URL, "api_key": "fake-key"})
    require.NoError(t, err)

    snap, err := c.Fetch(context.Background(), map[string]any{})
    require.NoError(t, err)
    assert.Equal(t, "My New Service", snap.ServiceName)
    assert.Len(t, snap.Sections, 1)
}

func TestConnector_Validate(t *testing.T) {
    _, err := New(map[string]any{})
    require.Error(t, err)
}
```

### 7. Document config fields

Add a commented block to `config.example.yaml` so users know what to set:

```yaml
# Example: MyNewService connector
#   - name: my-newservice
#     type: mynewservice
#     url: https://192.168.1.50:8443
#     tls_verify: false
#     # WISELABZ_SERVICES_N_APIKEY must be set in the environment
```

## Sync flow

When a user creates a connector and triggers sync:

1. **Frontend** calls `POST /api/connectors/{id}/sync` → the handler returns
   `202` with `{jobId, serviceId}` and spawns a background goroutine.
2. **Sync engine** runs `Fetch` → diffs against the previous snapshot →
   creates changes/alerts. Throughout the run it broadcasts WebSocket events:
   - `sync.progress` at each phase (`queued` → `fetching` → `diffing` →
     `generating` → `done`) so the UI shows a live progress bar.
   - `change.detected` for each diff result.
   - `alert.created` for each alert raised.
   - `sync.complete` with final counts on completion.
3. **Frontend** listens via `WebSocketProvider` and dispatches events to the
   `useLive` Zustand store, which drives the sync progress UI and the
   dashboard activity feed.

Your connector doesn't need to do anything special for sync — just implement
`Fetch` and `Validate`. The engine handles the rest.

## Health checks vs. sync

`POST /api/connectors/{id}/health` checks availability without the cost of a
full sync: it calls only `Validate`, never `Fetch`, so no snapshot, diff,
change, alert, or doc is produced. The result is classified into one of the
connector's `online` / `degraded` / `offline` statuses
(`internal/connector.ClassifyHealth`) and persisted on the connector record —
`degraded` means `Validate` succeeded but took longer than
`connector.DegradedLatencyThreshold`; `offline` means it returned an error.

This is distinct from `POST /api/connectors/{id}/test`, which also calls only
`Validate` but is a preview — its result is never persisted — used when
adding/editing a connector before it's saved. `/health` is for an
already-configured connector, checked on demand or on a lighter cadence than a
full sync, so the dashboard's status reflects current availability even
between sync runs.

Your connector doesn't need to do anything special for health checks either —
they call the same `Validate` you already implemented for sync.

## Testing without a real instance

If you don't have the service running, you have three options:

1. **Mock HTTP server (recommended).** Use `net/http/httptest` as shown above.
   Capture a real API response once (even from a browser devtools copy-paste),
   save it as a test fixture in `testdata/`, and serve it from the test server.

2. **Interface-based mocking.** If your connector depends on an interface
   (e.g. a `Client` that wraps the HTTP calls), generate a mock with
   `go.uber.org/mock` or write one by hand. Test your business logic without
   touching the network at all.

3. **Integration test with Docker.** If the service has an official Docker
   image, add a `docker compose` service in `deploy/test-fixtures/` that
   starts a throwaway instance. Use `t.Setenv` to point your connector at
   `localhost:<mapped-port>`, then run `Fetch` against a real (but ephemeral)
   service. These tests are tagged `//go:build integration` so `make test`
   skips them by default.

## Conventions

- **Package location:** `backend/internal/connector/<servicename>/`
- **One connector per package.** Don't bundle multiple connectors in one package.
- **Wrapping errors:** Use `fmt.Errorf("servicename: %w", err)` so callers can
  trace which connector failed.
- **Context deadlines:** Don't create a new context inside `Fetch`. Use the one
  passed to you. The caller sets the timeout.
- **Secrets:** Use `"password"` field type in the schema. The frontend renders a
  password input and the value is stored alongside other config.
- **Go version and dependencies:** Match the `go.mod` at the repository root.
  Avoid pulling in large dependency trees for simple HTTP calls.

## Version-aware connectors

Some services change their API between major releases: different paths,
different auth, sometimes a different payload shape for the same concept. The
Pi-hole connector (`backend/internal/connector/pihole/`) is the reference for
handling that without splitting into two connectors:

- **Offer the choice, default to detection.** A `select` field
  (`api_version`, options `auto` / `v6` / `v5`, default `auto`) lets an
  operator pin a version when detection would be ambiguous. The factory
  normalises anything unrecognised back to `auto`, so an existing stored
  config without the field keeps working.
- **Resolve the version once per call.** A small `session` struct carries the
  resolved version plus whatever the auth handshake produced (Pi-hole v6
  returns a session id from `POST /api/auth`; v5 sends its API token on every
  request). On `auto`, probe the newer API first and fall back; report the
  newer API's error when neither answers, since that is the one an operator
  most likely misconfigured.
- **Normalise before rendering.** Each version gets its own parser producing
  the same internal row structs (`groupRow`, `adlistRow`, ...). The table
  builders and the `SnapshotEntity` values they emit are version-agnostic, so
  the same instance diffs identically whichever API served it — and entity
  kinds, IDs and attributes stay stable across an upgrade.
- **Be explicit about what a version cannot do.** Pi-hole v5 exposes no
  restart action, so `Restart` returns a plain error naming the limitation
  rather than silently calling a v6-only path.
- **Test both.** Stand up one `httptest` fake per version, each answering 404
  for any path it does not serve, and run the same assertions against both.

## Session-based and multi-flavour APIs

Some vendors ship the same API behind two different front doors, and
authenticate with a cookie rather than a header. The UniFi connector
(`backend/internal/connector/unifi/`) is the reference for that shape:

- **Give the client a cookie jar** (`HTTPClientOptions.Jar`). `http.Client.Jar` is what replays the
  session cookie the login endpoint hands out; nothing else in the connector
  has to know the cookie's name.
- **Model the flavour as a path prefix.** A UniFi OS console serves the same
  Network API as a classic controller, just behind `/proxy/network`. A
  one-field `session` struct carrying that prefix keeps every call site
  identical, and a `controller_type` select (`auto` / `unifios` / `classic`)
  lets an operator pin it — same pattern as Pi-hole's `api_version`.
- **Prefer the actionable error when auto-detecting.** A 404 from probing the
  wrong flavour says nothing useful; an `AuthError` from the other probe means
  the credentials really were rejected, so that is the error to report.
- **Map the vendor's own error envelope.** UniFi answers failures with
  `{"meta":{"rc":"error","msg":"api.err.Invalid"}}` and uses 400, not 401, for
  bad credentials. Translate both into the shared connector error types so the
  sync engine can tell "retry later" from "fix your credentials".

## Keeping snapshots stable

A snapshot is diffed against the previous one, so anything that changes on
every poll turns into permanent noise. Decode the fields you need rather than
the whole payload, and leave out uptime, byte and packet counters, load
averages, signal strength and `last_seen`-style timestamps — the UniFi device
and client endpoints return all of those next to the inventory data, and the
connector simply does not decode them. The same rule is enforced for
`SnapshotEntity.Attributes` by `AttributeSpec` (see `RegisterAttributeCatalog`):
attributes are configuration, not telemetry. Where a resource is inherently
volatile — connected clients, say — prefer a summary (counts per SSID) over
one entity per instance.

When nearly everything an API returns is volatile, filtering has to go a step
further. The Home Assistant connector
(`backend/internal/connector/home_assistant/`) is the reference for that case:

- **Drop the timestamps outright.** Home Assistant stamps every entity with
  `last_changed`/`last_updated`. Neither reaches the snapshot — they say when
  something happened, not how the instance is configured.
- **Collapse measurements to availability.** A `sensor.*` state is a live
  reading. `stableState` reports such an entity as `available` (or
  `unavailable`/`unknown`, which *are* meaningful) instead of its value, using
  a volatile-domain list plus a numeric/timestamp check so an unrecognised
  integration is treated conservatively.
- **Report the counts, not the readings.** Per-domain totals and an
  unavailable-entity count carry the operational signal that the individual
  values would have, and only change when something actually changes.
- **Sort before you cap.** Large instances expose thousands of entities, so a
  `max_entities` field caps how many a snapshot carries. Entities are sorted by
  ID first, so the cap always keeps the same set instead of whatever order the
  API happened to return.

An appliance that reports far more than it should store needs the same
discipline applied section by section. The TrueNAS connector
(`backend/internal/connector/truenas/`) is the reference there:

- **Snapshot configuration, not telemetry.** Pool and dataset usage, uptime,
  load averages and disk temperatures are all available from `/api/v2.0/...`
  and all deliberately dropped. What stays is what an operator configured:
  topology, encryption, compression, share exports, which services start on
  boot.
- **Sort anything the API does not order.** Services are sorted by name and
  the dataset tree is flattened and sorted, so a reordered API response is
  not read as a change.
- **Derive stable summaries.** A vdev group becomes `"data: MIRROR (2
  disks)"` rather than a raw child list, and capacities are rendered from the
  static `size` field, never from live allocation.
- **Degrade per section.** Each endpoint is fetched independently and a
  failure leaves an `_<section> unavailable: ..._` placeholder, so an API key
  without one privilege (or a release that lacks an endpoint) still yields
  every other section.

## Getting your connector merged

Once your connector is implemented and tested:

1. Run `make lint` and `make test` to confirm everything passes.
2. Add a row to the Supported Services table in `README.md`.
3. Open a PR. Use the **New connector** type in the pull request template, and
   confirm you've tested against a real instance.
4. A maintainer will review. Expect feedback on error handling, logging, and
   the shape of your snapshot sections.

If you want to request a connector instead of building one, use the
[Connector Request form][connector-request].

[connector-request]: https://github.com/WiseLabz/WiseLabz/issues/new?template=connector_request.yml
