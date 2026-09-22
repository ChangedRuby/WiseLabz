package truenas

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const (
	systemInfoJSON = `{
		"version":"TrueNAS-SCALE-24.04.2","hostname":"nas.lan","model":"Intel(R) Xeon(R) E-2246G",
		"cores":12,"physical_cores":6,"physmem":68719476736,"system_product":"PowerEdge T140",
		"system_manufacturer":"Dell Inc.","system_serial":"ABC1234","ecc_memory":true,"license":null,
		"uptime":"10 days","uptime_seconds":864000.5,"loadavg":[0.4,0.3,0.2],"datetime":{"$date":1700000000000}
	}`
	poolsJSON = `[
		{"name":"tank","status":"ONLINE","healthy":true,"path":"/mnt/tank","encrypt":0,
		 "autotrim":{"value":"off","rawvalue":"off","source":"DEFAULT"},
		 "topology":{"data":[{"type":"MIRROR","children":[{"type":"DISK"},{"type":"DISK"}]}],
		             "cache":[],"log":[{"type":"DISK","children":[]}],"spare":[],"special":[]}},
		{"name":"vault","status":"DEGRADED","healthy":false,"path":"/mnt/vault","encrypt":2,
		 "autotrim":"on",
		 "topology":{"data":[{"type":"RAIDZ2","children":[{"type":"DISK"},{"type":"DISK"},{"type":"DISK"},{"type":"DISK"}]}]}}
	]`
	datasetsJSON = `[
		{"id":"tank","name":"tank","type":"FILESYSTEM","pool":"tank","encrypted":false,"mountpoint":"/mnt/tank",
		 "compression":{"value":"LZ4","rawvalue":"lz4"},"deduplication":{"value":"OFF"},"atime":{"value":"ON"},
		 "readonly":{"value":"OFF"},"sync":{"value":"STANDARD"},"quota":{"value":"0","rawvalue":"0"},
		 "used":{"value":"1.2T"},"available":{"value":"3.1T"},
		 "children":[
		   {"id":"tank/media","name":"tank/media","type":"FILESYSTEM","pool":"tank","encrypted":true,
		    "encryption_algorithm":"AES-256-GCM","mountpoint":"/mnt/tank/media",
		    "compression":{"value":"ZSTD"},"deduplication":{"value":"OFF"},"atime":{"value":"OFF"},
		    "readonly":{"value":"ON"},"sync":{"value":"DISABLED"},"quota":{"value":"500G","rawvalue":"536870912000"}}
		 ]},
		{"id":"tank/vm-disk","name":"tank/vm-disk","type":"VOLUME","pool":"tank","encrypted":false,
		 "compression":{"value":"LZ4"},"deduplication":{"value":"OFF"},"atime":null,
		 "readonly":{"value":"OFF"},"sync":{"value":"ALWAYS"},"quota":{"value":"0"}}
	]`
	disksJSON = `[
		{"name":"sda","serial":"WD-A1","model":"WDC WD40EFRX","type":"HDD","size":4000787030016,"pool":"tank","togglesmart":true,"description":"bay 1"},
		{"name":"nvme0n1","serial":"NV-77","model":"Samsung 980","type":"SSD","size":512110190592,"pool":"","togglesmart":false,"description":""}
	]`
	smbSharesJSON = `[
		{"id":1,"name":"media","path":"/mnt/tank/media","enabled":true,"ro":false,"guestok":false,"browsable":true,"purpose":"NO_PRESET","comment":"family | media"},
		{"id":2,"name":"public","path":"/mnt/tank/public","enabled":false,"ro":true,"guestok":true,"browsable":false,"purpose":"","comment":""}
	]`
	nfsSharesJSON = `[
		{"id":3,"paths":["/mnt/tank/backups"],"enabled":true,"ro":false,"maproot_user":"root","networks":["10.0.0.0/24"],"hosts":["backup.lan"],"comment":"veeam"},
		{"id":4,"path":"/mnt/vault/archive","enabled":true,"ro":true,"maproot_user":"","networks":[],"hosts":[],"comment":""}
	]`
	servicesJSON = `[
		{"id":2,"service":"ssh","enable":true,"state":"RUNNING"},
		{"id":1,"service":"cifs","enable":true,"state":"RUNNING"},
		{"id":3,"service":"nfs","enable":false,"state":"STOPPED"}
	]`
	interfacesJSON = `[
		{"id":"eno1","name":"eno1","type":"PHYSICAL","description":"lan uplink","ipv4_dhcp":false,"mtu":9000,
		 "aliases":[{"type":"INET","address":"10.0.0.10","netmask":24}]},
		{"id":"br0","name":"br0","type":"BRIDGE","description":"","ipv4_dhcp":true,"mtu":0,"aliases":[]}
	]`
	snapshotTasksJSON = `[
		{"id":1,"dataset":"tank/media","recursive":true,"enabled":true,"lifetime_value":2,"lifetime_unit":"WEEK",
		 "naming_schema":"auto-%Y%m%d.%H%M","exclude":["tank/media/tmp"],
		 "schedule":{"minute":"0","hour":"2","dom":"*","month":"*","dow":"*"}},
		{"id":2,"dataset":"tank","recursive":false,"enabled":false,"lifetime_value":0,"lifetime_unit":"",
		 "naming_schema":"","exclude":[],"schedule":{"minute":"*/15","hour":"*","dom":"*","month":"*","dow":"*"}}
	]`
	replicationTasksJSON = `[
		{"id":1,"name":"offsite","direction":"PUSH","transport":"SSH","source_datasets":["tank/media"],
		 "target_dataset":"backup/media","recursive":true,"enabled":true,"auto":true,"retention_policy":"SOURCE"}
	]`
)

// truenasAPI returns an httptest server serving the full happy-path API. If
// authCheck is non-nil it gates every request; returning false makes the
// handler reply 401.
func truenasAPI(t *testing.T, authCheck func(*http.Request) bool) *httptest.Server {
	t.Helper()
	body := map[string]string{
		pathSystemInfo:      systemInfoJSON,
		pathPools:           poolsJSON,
		pathDatasets:        datasetsJSON,
		pathDisks:           disksJSON,
		pathSMBShares:       smbSharesJSON,
		pathNFSShares:       nfsSharesJSON,
		pathServices:        servicesJSON,
		pathInterfaces:      interfacesJSON,
		pathSnapshotTasks:   snapshotTasksJSON,
		pathReplicationTask: replicationTasksJSON,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authCheck != nil && !authCheck(r) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"Not authenticated"}`))
			return
		}
		payload, ok := body[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(payload))
	}))
	t.Cleanup(server.Close)
	return server
}

// newTestConnector builds the connector through the registered factory so
// the config plumbing (api_key, verify_tls) is exercised too.
func newTestConnector(t *testing.T, config map[string]any) *Connector {
	t.Helper()
	connector.AllowLoopbackForTest(t)
	impl, err := connector.Get(typeName, config)
	if err != nil {
		t.Fatalf("connector.Get: %v", err)
	}
	c, ok := impl.(*Connector)
	if !ok {
		t.Fatalf("connector.Get returned %T, want *Connector", impl)
	}
	return c
}

func TestRegisteredSchema(t *testing.T) {
	schema, err := connector.GetTypeSchema(typeName)
	if err != nil {
		t.Fatalf("GetTypeSchema: %v", err)
	}
	if schema.Category != category || schema.Name != "TrueNAS" || schema.Stub {
		t.Fatalf("schema = %+v", schema)
	}
	want := map[string]bool{"url": true, "api_key": true, "verify_tls": false}
	got := make(map[string]bool, len(schema.Fields))
	for _, f := range schema.Fields {
		got[f.Key] = f.Required
	}
	if len(got) != len(want) {
		t.Errorf("schema fields = %+v, want exactly %v", got, want)
	}
	for key, required := range want {
		gotRequired, ok := got[key]
		if !ok {
			t.Errorf("schema missing field %q", key)
			continue
		}
		if gotRequired != required {
			t.Errorf("field %q Required = %t, want %t", key, gotRequired, required)
		}
	}
	for _, f := range schema.Fields {
		if f.Key == "api_key" && f.Type != "password" {
			t.Errorf("api_key field kind = %q, want password", f.Type)
		}
	}
	if _, ok := connector.AttributeCatalog()[typeName]; !ok {
		t.Errorf("attribute catalog not registered for %q", typeName)
	}
}

func TestIdentity(t *testing.T) {
	c := &Connector{}
	if c.Name() != "TrueNAS" || c.Type() != "truenas" || c.Category() != "virtualization" {
		t.Fatalf("identity = %s/%s/%s", c.Name(), c.Type(), c.Category())
	}
}

func TestValidate(t *testing.T) {
	server := truenasAPI(t, func(r *http.Request) bool {
		return r.Header.Get("Authorization") == "Bearer key123"
	})

	tests := []struct {
		name    string
		config  map[string]any
		wantErr string
	}{
		{name: "valid", config: map[string]any{"url": server.URL, "api_key": "key123"}},
		{name: "trailing slash trimmed", config: map[string]any{"url": server.URL + "/", "api_key": "key123"}},
		{name: "missing url", config: map[string]any{"api_key": "key123"}, wantErr: "url is required"},
		{name: "missing api key", config: map[string]any{"url": server.URL}, wantErr: "api_key is required"},
		{name: "wrong api key", config: map[string]any{"url": server.URL, "api_key": "nope"}, wantErr: "auth error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestConnector(t, tt.config)
			err := c.Validate(context.Background(), nil)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Validate() error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

// TestBearerHeaderIsSent asserts the exact wire format of the API key.
func TestBearerHeaderIsSent(t *testing.T) {
	var gotAuth, gotAccept string
	server := truenasAPI(t, func(r *http.Request) bool {
		gotAuth = r.Header.Get("Authorization")
		gotAccept = r.Header.Get("Accept")
		return true
	})

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "1-abcdef"})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if gotAuth != "Bearer 1-abcdef" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer 1-abcdef")
	}
	if gotAccept != "application/json" {
		t.Errorf("Accept = %q, want application/json", gotAccept)
	}
}

func TestFetchHappyPath(t *testing.T) {
	server := truenasAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123"})

	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if snapshot.ServiceName != "TrueNAS" || snapshot.Type != typeName {
		t.Fatalf("snapshot = %s/%s", snapshot.ServiceName, snapshot.Type)
	}

	wantSections := []string{
		"System", "Pools", "Datasets", "Disks", "SMB Shares", "NFS Shares",
		"Services", "Network Interfaces", "Snapshot Tasks", "Replication Tasks",
	}
	if len(snapshot.Sections) != len(wantSections) {
		t.Fatalf("sections = %d, want %d", len(snapshot.Sections), len(wantSections))
	}
	for i, title := range wantSections {
		if snapshot.Sections[i].Title != title {
			t.Errorf("sections[%d].Title = %q, want %q", i, snapshot.Sections[i].Title, title)
		}
		if strings.Contains(snapshot.Sections[i].Content, "unavailable") {
			t.Errorf("sections[%d] unexpectedly unavailable: %q", i, snapshot.Sections[i].Content)
		}
	}

	wantMetadata := map[string]string{
		"truenas_url":            server.URL,
		"truenas_version":        "TrueNAS-SCALE-24.04.2",
		"truenas_flavor":         "SCALE",
		"truenas_hostname":       "nas.lan",
		"pool_count":             "2",
		"dataset_count":          "3",
		"disk_count":             "2",
		"smb_share_count":        "2",
		"nfs_share_count":        "2",
		"service_count":          "3",
		"services_running":       "2",
		"interface_count":        "2",
		"snapshot_task_count":    "2",
		"replication_task_count": "1",
	}
	for k, want := range wantMetadata {
		if got := snapshot.Metadata[k]; got != want {
			t.Errorf("metadata[%q] = %q, want %q", k, got, want)
		}
	}

	byKind := map[string]int{}
	for _, e := range snapshot.Entities {
		byKind[e.Kind]++
	}
	want := map[string]int{
		"pool": 2, "dataset": 3, "disk": 2, "share": 4,
		"service": 3, "interface": 2, "snapshot_task": 2, "replication_task": 1,
	}
	for kind, n := range want {
		if byKind[kind] != n {
			t.Errorf("entities of kind %q = %d, want %d", kind, byKind[kind], n)
		}
	}

	wantDeps := []connector.ServiceDependency{
		{Kind: "storage", Name: "tank"},
		{Kind: "storage", Name: "vault"},
	}
	if len(snapshot.Dependencies) != len(wantDeps) {
		t.Fatalf("dependencies = %+v, want %+v", snapshot.Dependencies, wantDeps)
	}
	for i, d := range wantDeps {
		if snapshot.Dependencies[i] != d {
			t.Errorf("dependencies[%d] = %+v, want %+v", i, snapshot.Dependencies[i], d)
		}
	}
}

// TestFetchIsStableAcrossCalls is the guarantee the connector exists for: an
// unchanged appliance must produce byte-identical sections and entities, so
// the diff engine reports no change.
func TestFetchIsStableAcrossCalls(t *testing.T) {
	server := truenasAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123"})

	first, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	second, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(first.Sections) != len(second.Sections) {
		t.Fatalf("section counts differ: %d vs %d", len(first.Sections), len(second.Sections))
	}
	for i := range first.Sections {
		if first.Sections[i] != second.Sections[i] {
			t.Errorf("section %q is not stable:\n%q\n%q", first.Sections[i].Title,
				first.Sections[i].Content, second.Sections[i].Content)
		}
	}
	for i := range first.Entities {
		if first.Entities[i].ExternalID != second.Entities[i].ExternalID {
			t.Errorf("entity %d order is not stable: %q vs %q", i,
				first.Entities[i].ExternalID, second.Entities[i].ExternalID)
		}
	}
}

// TestSnapshotOmitsVolatileValues guards the values that would otherwise
// churn the diff on every sync.
func TestSnapshotOmitsVolatileValues(t *testing.T) {
	server := truenasAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123"})

	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	for _, s := range snapshot.Sections {
		for _, banned := range []string{"10 days", "864000", "0.4", "1.2T", "3.1T"} {
			if strings.Contains(s.Content, banned) {
				t.Errorf("section %q contains volatile value %q:\n%s", s.Title, banned, s.Content)
			}
		}
	}
	for _, e := range snapshot.Entities {
		for _, banned := range []string{"used", "available", "uptime", "temperature", "allocated"} {
			if _, ok := e.Attributes[banned]; ok {
				t.Errorf("entity %s/%s carries volatile attribute %q", e.Kind, e.Name, banned)
			}
		}
	}
}

// TestFetchSelectiveFields checks the "fields" hint limits which endpoints
// are hit and which sections come back.
func TestFetchSelectiveFields(t *testing.T) {
	var hit []string
	server := truenasAPI(t, func(r *http.Request) bool {
		hit = append(hit, r.URL.Path)
		return true
	})
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123"})

	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"pools"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 1 || snapshot.Sections[0].Title != "Pools" {
		t.Fatalf("sections = %+v, want only Pools", snapshot.Sections)
	}
	if len(hit) != 1 || hit[0] != pathPools {
		t.Errorf("requests = %v, want only %q", hit, pathPools)
	}

	// "shares" gates both share endpoints.
	hit = nil
	snapshot, err = c.Fetch(context.Background(), map[string]any{"fields": []any{"shares"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 2 {
		t.Fatalf("sections = %+v, want SMB and NFS shares", snapshot.Sections)
	}
	if len(hit) != 2 || hit[0] != pathSMBShares || hit[1] != pathNFSShares {
		t.Errorf("requests = %v, want the two share endpoints", hit)
	}
}

// TestFetchDegradesPerSection checks a failing endpoint only degrades its
// own section instead of failing the whole snapshot.
func TestFetchDegradesPerSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case pathSystemInfo:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(systemInfoJSON))
		case pathPools:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("boom"))
		case pathDatasets:
			_, _ = w.Write([]byte("not json"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123"})
	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != len(sections) {
		t.Fatalf("sections = %d, want %d", len(snapshot.Sections), len(sections))
	}
	byTitle := map[string]string{}
	for _, s := range snapshot.Sections {
		byTitle[s.Title] = s.Content
	}
	if strings.Contains(byTitle["System"], "unavailable") {
		t.Errorf("System should have succeeded: %q", byTitle["System"])
	}
	if !strings.Contains(byTitle["Pools"], "API returned 500") {
		t.Errorf("Pools = %q, want the upstream 500 reported", byTitle["Pools"])
	}
	if !strings.Contains(byTitle["Datasets"], "malformed response") {
		t.Errorf("Datasets = %q, want a malformed-response placeholder", byTitle["Datasets"])
	}
	if !strings.Contains(byTitle["Replication Tasks"], "API returned 404") {
		t.Errorf("Replication Tasks = %q, want the upstream 404 reported", byTitle["Replication Tasks"])
	}
	if snapshot.Metadata["truenas_hostname"] != "nas.lan" {
		t.Errorf("metadata = %+v, want the System section's metadata kept", snapshot.Metadata)
	}
}

func TestErrorMapping(t *testing.T) {
	tests := []struct {
		name   string
		status int
		assert func(*testing.T, error)
	}{
		{
			name: "401 maps to AuthError", status: http.StatusUnauthorized,
			assert: func(t *testing.T, err error) {
				var target *connector.AuthError
				if !errors.As(err, &target) {
					t.Fatalf("error = %v (%T), want *connector.AuthError", err, err)
				}
			},
		},
		{
			name: "403 maps to AuthError", status: http.StatusForbidden,
			assert: func(t *testing.T, err error) {
				var target *connector.AuthError
				if !errors.As(err, &target) {
					t.Fatalf("error = %v (%T), want *connector.AuthError", err, err)
				}
			},
		},
		{
			name: "502 maps to ServiceUnavailableError", status: http.StatusBadGateway,
			assert: func(t *testing.T, err error) {
				var target *connector.ServiceUnavailableError
				if !errors.As(err, &target) {
					t.Fatalf("error = %v (%T), want *connector.ServiceUnavailableError", err, err)
				}
			},
		},
		{
			name: "503 maps to ServiceUnavailableError", status: http.StatusServiceUnavailable,
			assert: func(t *testing.T, err error) {
				var target *connector.ServiceUnavailableError
				if !errors.As(err, &target) {
					t.Fatalf("error = %v (%T), want *connector.ServiceUnavailableError", err, err)
				}
			},
		},
		{
			name: "504 maps to ServiceUnavailableError", status: http.StatusGatewayTimeout,
			assert: func(t *testing.T, err error) {
				var target *connector.ServiceUnavailableError
				if !errors.As(err, &target) {
					t.Fatalf("error = %v (%T), want *connector.ServiceUnavailableError", err, err)
				}
			},
		},
		{
			name: "500 stays a plain error", status: http.StatusInternalServerError,
			assert: func(t *testing.T, err error) {
				var authErr *connector.AuthError
				var unavailable *connector.ServiceUnavailableError
				if errors.As(err, &authErr) || errors.As(err, &unavailable) {
					t.Fatalf("error = %v (%T), want a plain error", err, err)
				}
				if !strings.Contains(err.Error(), "API returned 500") {
					t.Fatalf("error = %v, want it to mention the status", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte("upstream said no"))
			}))
			defer server.Close()

			c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123"})
			err := c.Validate(context.Background(), nil)
			if err == nil {
				t.Fatal("Validate() error = nil, want an error")
			}
			tt.assert(t, err)
		})
	}
}

func TestExpiredDeadlineMapsToTimeout(t *testing.T) {
	server := truenasAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123"})

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	err := c.Validate(ctx, nil)
	var timeoutErr *connector.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("Validate() error = %v (%T), want *connector.TimeoutError", err, err)
	}
}

// TestOversizedResponseIsRejected checks the shared body cap is enforced.
func TestOversizedResponseIsRejected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		chunk := strings.Repeat("a", 1<<20)
		for written := 0; written <= connector.MaxResponseBytes; written += len(chunk) {
			if _, err := w.Write([]byte(chunk)); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123"})
	err := c.Validate(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("Validate() error = %v, want the response-size cap to trip", err)
	}
}

// TestGuardedDialerBlocksLoopback confirms the SSRF guard is wired in when
// the test override is not active.
func TestGuardedDialerBlocksLoopback(t *testing.T) {
	server := truenasAPI(t, nil)
	impl, err := connector.Get(typeName, map[string]any{"url": server.URL, "api_key": "key123"})
	if err != nil {
		t.Fatalf("connector.Get: %v", err)
	}
	if err := impl.Validate(context.Background(), nil); err == nil {
		t.Fatal("Validate() error = nil, want the guarded dialer to block loopback")
	}
}

func TestRedirectsAreNotFollowed(t *testing.T) {
	var hits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if r.URL.Path == pathSystemInfo {
			http.Redirect(w, r, "/elsewhere", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123"})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate() error = %v, want the 302 to be returned as-is", err)
	}
	if hits != 1 {
		t.Fatalf("hits = %d, want 1 (redirect must not be followed)", hits)
	}
}

func TestVerifyTLSDisabledAllowsSelfSignedCert(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(systemInfoJSON))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123", "verify_tls": false})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate() with verify_tls=false error = %v, want nil", err)
	}

	strict := newTestConnector(t, map[string]any{"url": server.URL, "api_key": "key123", "verify_tls": true})
	if err := strict.Validate(context.Background(), nil); err == nil {
		t.Fatal("Validate() with verify_tls=true error = nil, want a certificate error")
	}
}
