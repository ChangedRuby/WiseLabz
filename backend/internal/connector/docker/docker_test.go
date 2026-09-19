package docker

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestValidateHitsVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/version" {
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"Version":"24.0.0"}`))
	}))
	defer server.Close()

	c := &Connector{host: "tcp://" + server.Listener.Addr().String(), baseURL: server.URL, client: server.Client()}
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestFetchBuildsSectionsFromEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/info":
			_, _ = w.Write([]byte(`{"Name":"docker-host","ServerVersion":"24.0.0","Containers":1,"ContainersRunning":1,"Images":2,"NCPU":4,"MemTotal":1024}`))
		case "/containers/json":
			_, _ = w.Write([]byte(`[{"Names":["/web"],"Image":"nginx","State":"running","Status":"Up 2 hours"}]`))
		case "/images/json":
			_, _ = w.Write([]byte(`[{"RepoTags":["nginx:latest"],"Size":100}]`))
		case "/volumes":
			_, _ = w.Write([]byte(`{"Volumes":[{"Name":"data","Driver":"local","Mountpoint":"/var/lib/docker/volumes/data"}]}`))
		case "/networks":
			_, _ = w.Write([]byte(`[{"Name":"bridge","Driver":"bridge","Scope":"local"}]`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
	snap, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if len(snap.Sections) != 5 {
		t.Fatalf("Sections = %d, want 5", len(snap.Sections))
	}
	if !strings.Contains(snap.Sections[0].Content, "docker-host") {
		t.Errorf("System section = %q", snap.Sections[0].Content)
	}
	wantDeps := []connector.ServiceDependency{{Kind: "host", Name: "tcp://example"}}
	if len(snap.Dependencies) != 1 || snap.Dependencies[0] != wantDeps[0] {
		t.Errorf("Dependencies = %+v, want %+v", snap.Dependencies, wantDeps)
	}
}

func TestFetchWithFieldsHintSkipsUnrequestedCalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/info":
			_, _ = w.Write([]byte(`{"Name":"docker-host"}`))
		case "/containers/json":
			_, _ = w.Write([]byte(`[{"Names":["/web"],"Image":"nginx","State":"running","Status":"Up"}]`))
		default:
			t.Fatalf("unexpected request path %s: selective fetch should only hit /info and /containers/json", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
	snap, err := c.Fetch(context.Background(), map[string]any{"fields": []any{"containers"}})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if len(snap.Sections) != 2 {
		t.Fatalf("Sections = %d, want 2 (System + Containers)", len(snap.Sections))
	}
}

func TestFetchToleratesEndpointFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/info":
			_, _ = w.Write([]byte(`{"Name":"docker-host"}`))
		case "/containers/json":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
	snap, err := c.Fetch(context.Background(), map[string]any{"fields": []any{"containers"}})
	if err != nil {
		t.Fatalf("Fetch() error = %v, want nil (failures should be tolerated as placeholders)", err)
	}
	if !strings.Contains(snap.Sections[1].Content, "unavailable") {
		t.Errorf("Containers section = %q, want unavailable placeholder", snap.Sections[1].Content)
	}
}

func TestNewDockerClientRejectsUnsupportedScheme(t *testing.T) {
	if _, _, err := newDockerClient("ftp://user@host", nil); err == nil {
		t.Fatal("newDockerClient(ftp://...) error = nil, want rejection")
	}
}

func TestNewDockerClientDialsUnixSocket(t *testing.T) {
	dir := t.TempDir()
	socketPath := filepath.Join(dir, "docker.sock")

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("net.Listen(unix): %v", err)
	}
	defer func() { _ = listener.Close() }()

	go func() {
		srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/version" {
				return
			}
			_, _ = w.Write([]byte(`{"Version":"24.0.0"}`))
		})}
		_ = srv.Serve(listener)
	}()

	client, baseURL, err := newDockerClient("unix://"+socketPath, nil)
	if err != nil {
		t.Fatalf("newDockerClient(unix://...) error = %v", err)
	}
	if baseURL != "http://unix" {
		t.Fatalf("baseURL = %q, want http://unix", baseURL)
	}

	c := &Connector{host: "unix://" + socketPath, baseURL: baseURL, client: client}
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate() over unix socket error = %v", err)
	}
}

func TestFetchSurfacesMalformedSystemResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/info" {
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
	snap, err := c.Fetch(context.Background(), map[string]any{"fields": []any{"none"}})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if !strings.Contains(snap.Sections[0].Content, "malformed response") {
		t.Fatalf("System section = %q, want malformed response placeholder", snap.Sections[0].Content)
	}
}

func generateSelfSignedCert(t *testing.T) (certPEM, keyPEM []byte) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "docker-test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	return certPEM, keyPEM
}

func TestNewTCPDockerClientMutualTLS(t *testing.T) {
	certPEM, keyPEM := generateSelfSignedCert(t)
	clientCAPool := x509.NewCertPool()
	if !clientCAPool.AppendCertsFromPEM(certPEM) {
		t.Fatal("failed to load client cert into pool")
	}

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.TLS.PeerCertificates) == 0 {
			t.Error("server saw no client certificate")
		}
		_, _ = w.Write([]byte(`{"Version":"24.0.0"}`))
	}))
	server.TLS = &tls.Config{ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientCAPool}
	server.StartTLS()
	defer server.Close()

	addr := strings.TrimPrefix(server.URL, "https://")

	// newTCPDockerClient's guarded dialer intentionally refuses loopback
	// addresses (see connector.GuardedDialer), which httptest always binds
	// to — so exercise the TLS config it builds directly against the
	// server via tls.Dial instead of going through the full client, still
	// proving the cert/key/verify_tls wiring produces a working handshake.
	tlsConfig, err := buildDockerTLSConfig(map[string]any{
		"tls_cert":   string(certPEM),
		"tls_key":    string(keyPEM),
		"verify_tls": false,
	})
	if err != nil {
		t.Fatalf("buildDockerTLSConfig() error = %v", err)
	}
	if tlsConfig == nil {
		t.Fatal("buildDockerTLSConfig() = nil, want a config")
	}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		t.Fatalf("mutual TLS handshake failed: %v", err)
	}
	defer conn.Close() //nolint:errcheck
	if _, err := conn.Write([]byte("GET /version HTTP/1.1\r\nHost: docker\r\nConnection: close\r\n\r\n")); err != nil {
		t.Fatalf("write request: %v", err)
	}
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
}

func TestNewTCPDockerClientNoTLSWhenNoCert(t *testing.T) {
	client, baseURL, err := newTCPDockerClient("example:2375", nil)
	if err != nil {
		t.Fatalf("newTCPDockerClient() error = %v", err)
	}
	if baseURL != "http://example:2375" {
		t.Fatalf("baseURL = %q, want http://example:2375", baseURL)
	}
	if client.Transport.(*http.Transport).TLSClientConfig != nil {
		t.Fatal("TLSClientConfig set with no tls_cert/tls_key configured")
	}
}

func TestNewTCPDockerClientRejectsInvalidCertPair(t *testing.T) {
	if _, _, err := newTCPDockerClient("example:2376", map[string]any{
		"tls_cert": "not a cert",
		"tls_key":  "not a key",
	}); err == nil {
		t.Fatal("newTCPDockerClient() error = nil, want rejection of invalid cert/key")
	}
}

// sshDockerServer runs a minimal SSH server accepting one exec request of
// "docker system dial-stdio" and speaking a single canned HTTP exchange over
// the resulting channel, enough to prove newSSHDockerClient's dial/auth/pipe
// wiring actually carries Engine API traffic end to end.
func startSSHDockerServer(t *testing.T, user, password string) (addr, hostPublicKeyText string) {
	t.Helper()
	hostKey, err := generateSSHHostKey()
	if err != nil {
		t.Fatalf("generate ssh host key: %v", err)
	}
	hostPublicKeyText = string(ssh.MarshalAuthorizedKey(hostKey.PublicKey()))

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == user && string(pass) == password {
				return nil, nil
			}
			return nil, fmt.Errorf("wrong credentials")
		},
	}
	config.AddHostKey(hostKey)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close() })

	go func() {
		for {
			nConn, err := listener.Accept()
			if err != nil {
				return
			}
			go serveSSHDockerConn(nConn, config)
		}
	}()

	return listener.Addr().String(), hostPublicKeyText
}

func serveSSHDockerConn(nConn net.Conn, config *ssh.ServerConfig) {
	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		return
	}
	defer sshConn.Close() //nolint:errcheck
	go ssh.DiscardRequests(reqs)
	for newChan := range chans {
		if newChan.ChannelType() != "session" {
			_ = newChan.Reject(ssh.UnknownChannelType, "unsupported")
			continue
		}
		channel, requests, err := newChan.Accept()
		if err != nil {
			return
		}
		go func() {
			for req := range requests {
				if req.Type == "exec" {
					_ = req.Reply(true, nil)
					go serveOneHTTPExchange(channel)
				} else {
					_ = req.Reply(false, nil)
				}
			}
		}()
	}
}

func serveOneHTTPExchange(channel ssh.Channel) {
	defer channel.Close() //nolint:errcheck
	reader := bufio.NewReader(channel)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if line == "\r\n" || line == "\n" {
			break
		}
	}
	_, _ = channel.Write([]byte("HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n{\"Version\":\"24.0.0\"}"))
}

func generateSSHHostKey() (ssh.Signer, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return ssh.NewSignerFromKey(key)
}

func TestNewSSHDockerClientDialsAndExecutesDialStdio(t *testing.T) {
	addr, hostKey := startSSHDockerServer(t, "testuser", "testpass")

	client, baseURL, err := newSSHDockerClient("ssh://testuser@"+addr, map[string]any{
		"ssh_password": "testpass",
		"ssh_host_key": hostKey,
	})
	if err != nil {
		t.Fatalf("newSSHDockerClient() error = %v", err)
	}
	if baseURL != "http://docker" {
		t.Fatalf("baseURL = %q, want http://docker", baseURL)
	}

	resp, err := client.Get(baseURL + "/version")
	if err != nil {
		t.Fatalf("request over ssh dial-stdio failed: %v", err)
	}
	defer resp.Body.Close() //nolint:errcheck
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if !strings.Contains(string(body), "24.0.0") {
		t.Fatalf("body = %q, want version 24.0.0", body)
	}
}

func TestNewSSHDockerClientSupportsSequentialRequests(t *testing.T) {
	addr, hostKey := startSSHDockerServer(t, "testuser", "testpass")
	client, baseURL, err := newSSHDockerClient("ssh://testuser@"+addr, map[string]any{
		"ssh_password": "testpass",
		"ssh_host_key": hostKey,
	})
	if err != nil {
		t.Fatalf("newSSHDockerClient() error = %v", err)
	}
	for i := 0; i < 3; i++ {
		resp, err := client.Get(baseURL + "/version")
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

func TestNewSSHDockerClientRejectsWrongCredentials(t *testing.T) {
	addr, hostKey := startSSHDockerServer(t, "testuser", "testpass")

	client, baseURL, err := newSSHDockerClient("ssh://testuser@"+addr, map[string]any{
		"ssh_password": "wrongpass",
		"ssh_host_key": hostKey,
	})
	if err != nil {
		t.Fatalf("newSSHDockerClient() error = %v", err)
	}
	if resp, err := client.Get(baseURL + "/version"); err == nil {
		_ = resp.Body.Close()
		t.Fatal("request error = nil, want auth failure")
	}
}

func TestNewSSHDockerClientRejectsMissingHostKey(t *testing.T) {
	addr, _ := startSSHDockerServer(t, "testuser", "testpass")

	if _, _, err := newSSHDockerClient("ssh://testuser@"+addr, map[string]any{
		"ssh_password": "testpass",
	}); err == nil {
		t.Fatal("newSSHDockerClient() error = nil, want rejection for missing ssh_host_key")
	}
}

func TestNewSSHDockerClientRejectsWrongHostKey(t *testing.T) {
	addr, _ := startSSHDockerServer(t, "testuser", "testpass")
	_, otherHostKey := startSSHDockerServer(t, "testuser", "testpass")

	client, baseURL, err := newSSHDockerClient("ssh://testuser@"+addr, map[string]any{
		"ssh_password": "testpass",
		"ssh_host_key": otherHostKey,
	})
	if err != nil {
		t.Fatalf("newSSHDockerClient() error = %v", err)
	}
	if resp, err := client.Get(baseURL + "/version"); err == nil {
		_ = resp.Body.Close()
		t.Fatal("request error = nil, want rejection for mismatched host key")
	}
}

func TestDoRequestErrorCases(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "502 BadGateway returns ServiceUnavailableError",
			statusCode: http.StatusBadGateway,
			wantErr:    true,
		},
		{
			name:       "503 ServiceUnavailable returns ServiceUnavailableError",
			statusCode: http.StatusServiceUnavailable,
			wantErr:    true,
		},
		{
			name:       "504 GatewayTimeout returns ServiceUnavailableError",
			statusCode: http.StatusGatewayTimeout,
			wantErr:    true,
		},
		{
			name:       "500 InternalServerError returns generic error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte("error response"))
			}))
			defer server.Close()

			c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
			_, err := c.doRequest(context.Background(), "/info")

			if !tt.wantErr && err != nil {
				t.Errorf("doRequest() error = %v, wantErr = false", err)
				return
			}
			if tt.wantErr && err == nil {
				t.Errorf("doRequest() error = nil, wantErr = true")
				return
			}
		})
	}
}

func TestDoRequestContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(1 * time.Second)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := c.doRequest(ctx, "/info")
	var timeoutErr *connector.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Errorf("doRequest() error = %v, want *connector.TimeoutError", err)
	}
}

func TestValidateAndFetchWithMalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}

	// Validate returns nil on successful connection even with malformed JSON
	// (Validate only checks that the endpoint is reachable, not that it returns valid JSON)
	err := c.Validate(context.Background(), nil)
	if err != nil {
		t.Errorf("Validate() error = %v, want nil (Validate only checks connectivity)", err)
	}

	// Fetch should handle malformed /info by creating a malformed response section
	snap, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Errorf("Fetch() error = %v, want nil (tolerates malformed JSON as placeholder)", err)
		return
	}
	if !strings.Contains(snap.Sections[0].Content, "malformed response") {
		t.Errorf("System section = %q, want malformed response placeholder", snap.Sections[0].Content)
	}
}

func TestRestart(t *testing.T) {
	tests := []struct {
		name       string
		entityRef  string
		statusCode int
		wantErr    bool
	}{
		{name: "success", entityRef: "abc123", statusCode: http.StatusNoContent},
		{name: "empty entityRef errors", entityRef: "", wantErr: true},
		{name: "not found", entityRef: "missing", statusCode: http.StatusNotFound, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath, gotMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath, gotMethod = r.URL.Path, r.Method
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
			err := c.Restart(context.Background(), nil, tt.entityRef)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Restart() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Restart() error = %v", err)
			}
			if gotMethod != "POST" || gotPath != "/containers/"+tt.entityRef+"/restart" {
				t.Errorf("request = %s %s, want POST /containers/%s/restart", gotMethod, gotPath, tt.entityRef)
			}
		})
	}
}

func TestStartStop(t *testing.T) {
	tests := []struct {
		name       string
		action     string
		entityRef  string
		statusCode int
		wantErr    bool
	}{
		{name: "start success", action: "start", entityRef: "abc123", statusCode: http.StatusNoContent},
		{name: "stop success", action: "stop", entityRef: "abc123", statusCode: http.StatusNoContent},
		{name: "start empty entityRef errors", action: "start", entityRef: "", wantErr: true},
		{name: "stop not found", action: "stop", entityRef: "missing", statusCode: http.StatusNotFound, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath, gotMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath, gotMethod = r.URL.Path, r.Method
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
			var err error
			if tt.action == "start" {
				err = c.Start(context.Background(), nil, tt.entityRef)
			} else {
				err = c.Stop(context.Background(), nil, tt.entityRef)
			}
			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s() error = nil, want error", tt.action)
				}
				return
			}
			if err != nil {
				t.Fatalf("%s() error = %v", tt.action, err)
			}
			wantPath := "/containers/" + tt.entityRef + "/" + tt.action
			if gotMethod != "POST" || gotPath != wantPath {
				t.Errorf("request = %s %s, want POST %s", gotMethod, gotPath, wantPath)
			}
		})
	}
}

func TestConfigPush(t *testing.T) {
	tests := []struct {
		name       string
		entityRef  string
		fieldKey   string
		value      any
		statusCode int
		wantErr    bool
	}{
		{name: "success", entityRef: "abc123", fieldKey: "restartPolicy", value: "always", statusCode: http.StatusOK},
		{name: "empty entityRef errors", entityRef: "", fieldKey: "restartPolicy", value: "always", wantErr: true},
		{name: "unsupported field errors", entityRef: "abc123", fieldKey: "image", value: "nginx:latest", wantErr: true},
		{name: "upstream error", entityRef: "abc123", fieldKey: "restartPolicy", value: "always", statusCode: http.StatusInternalServerError, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath, gotMethod, gotBody string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath, gotMethod = r.URL.Path, r.Method
				b, _ := io.ReadAll(r.Body)
				gotBody = string(b)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
			err := c.ConfigPush(context.Background(), nil, tt.entityRef, tt.fieldKey, tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ConfigPush() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ConfigPush() error = %v", err)
			}
			wantPath := "/containers/" + tt.entityRef + "/update"
			if gotMethod != "POST" || gotPath != wantPath {
				t.Errorf("request = %s %s, want POST %s", gotMethod, gotPath, wantPath)
			}
			if !strings.Contains(gotBody, "always") {
				t.Errorf("body = %q, want it to reference the pushed value", gotBody)
			}
		})
	}
}

func TestDockerWritableFields(t *testing.T) {
	c := &Connector{}
	fields := c.WritableFields()
	if len(fields) != 1 || fields[0].Key != "restartPolicy" {
		t.Errorf("WritableFields() = %+v, want one field \"restartPolicy\"", fields)
	}
}
