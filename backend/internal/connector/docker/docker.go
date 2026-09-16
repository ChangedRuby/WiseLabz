// Package docker implements a real Docker Engine API connector.
package docker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "docker"

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "containers_paas",
		Name:     "Docker",
		Fields: []connector.SchemaField{
			{Key: "host", Label: "Docker Host", Type: "text", Required: true, Placeholder: "unix:///var/run/docker.sock, tcp://host:2375, or ssh://user@host"},
			{Key: "tls_cert", Label: "TLS Client Certificate (PEM)", Type: "secret", Description: "For tcp:// hosts using mutual TLS"},
			{Key: "tls_key", Label: "TLS Client Key (PEM)", Type: "secret", Description: "For tcp:// hosts using mutual TLS"},
			{Key: "tls_ca", Label: "TLS CA Certificate (PEM)", Type: "secret", Description: "Optional; verifies the server against this CA instead of the system pool"},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Default: "true", Description: "For tcp:// hosts with a client certificate configured"},
			{Key: "ssh_user", Label: "SSH Username", Type: "text", Description: "For ssh:// hosts; overrides any user@ in the host URL"},
			{Key: "ssh_password", Label: "SSH Password", Type: "password", Description: "For ssh:// hosts; ignored if an SSH private key is set"},
			{Key: "ssh_private_key", Label: "SSH Private Key (PEM)", Type: "secret", Description: "For ssh:// hosts"},
			{Key: "ssh_private_key_passphrase", Label: "SSH Private Key Passphrase", Type: "password", Description: "Optional; only used with an encrypted SSH private key"},
			{Key: "ssh_host_key", Label: "SSH Host Public Key", Type: "secret", Description: "For ssh:// hosts; pinned host key in authorized_keys format (e.g. output of ssh-keyscan), required to verify the server's identity"},
		},
	}, func(config map[string]any) (connector.Connector, error) {
		host, _ := config["host"].(string)
		client, baseURL, err := newDockerClient(host, config)
		// Construction never fails here: an invalid/empty host (e.g. an
		// unconfigured connector instance) surfaces as an error from
		// Validate/Fetch instead, matching how other connectors treat
		// missing config as a runtime rather than a registration error.
		return &Connector{host: host, baseURL: baseURL, client: client, configErr: err}, nil
	})
}

// Connector fetches data from the Docker Engine API.
type Connector struct {
	host      string
	baseURL   string
	client    *http.Client
	configErr error
}

// Name returns the connector display name.
func (d *Connector) Name() string { return "Docker" }

// Type returns the connector type identifier.
func (d *Connector) Type() string { return typeName }

// Category returns the connector category.
func (d *Connector) Category() string { return "containers_paas" }

// Validate tests the connection to the Docker Engine API.
func (d *Connector) Validate(ctx context.Context, _ map[string]any) error {
	_, err := d.doRequest(ctx, "/version")
	return err
}

// Fetch retrieves engine info, containers, images, volumes, and networks.
// config may carry a "fields" selective-fetch hint naming a subset of
// {"containers","images","volumes","networks"} to skip the other calls.
// The System section always runs. Each section fetch failure is tolerated
// as a placeholder rather than failing the whole Fetch.
func (d *Connector) Fetch(ctx context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	fields := connector.RequestedFields(config)

	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity
	metadata := map[string]string{"docker_host": d.host}

	if raw, err := d.doRequest(ctx, "/info"); err != nil {
		sections = append(sections, connector.SnapshotSection{Title: "System", Content: "_System info unavailable: " + err.Error() + "_"})
	} else {
		var info struct {
			Name              string `json:"Name"`
			ServerVersion     string `json:"ServerVersion"`
			Containers        int    `json:"Containers"`
			ContainersRunning int    `json:"ContainersRunning"`
			Images            int    `json:"Images"`
			NCPU              int    `json:"NCPU"`
			MemTotal          int64  `json:"MemTotal"`
		}
		if err := json.Unmarshal(raw, &info); err != nil {
			sections = append(sections, connector.SnapshotSection{
				Title:   "System",
				Content: "_System info unavailable: " + connector.NewMalformedResponseError(err).Error() + "_",
			})
		} else {
			content := fmt.Sprintf("**Host**: %s\n**Engine Version**: %s\n**Containers**: %d (%d running)\n**Images**: %d\n**CPUs**: %d\n**Memory**: %d bytes\n",
				info.Name, info.ServerVersion, info.Containers, info.ContainersRunning, info.Images, info.NCPU, info.MemTotal)
			sections = append(sections, connector.SnapshotSection{Title: "System", Content: content})
			metadata["engine_version"] = info.ServerVersion
		}
	}

	if connector.WantsField(fields, "containers") {
		if raw, err := d.doRequest(ctx, "/containers/json?all=true"); err != nil {
			sections = append(sections, connector.SnapshotSection{Title: "Containers", Content: "_Containers unavailable: " + err.Error() + "_"})
		} else {
			content, ents := buildContainerTable(raw)
			sections = append(sections, connector.SnapshotSection{Title: "Containers", Content: content})
			entities = append(entities, ents...)
		}
	}
	if connector.WantsField(fields, "images") {
		sections = append(sections, d.fetchSection(ctx, "Images", "/images/json", buildImageTable))
	}
	if connector.WantsField(fields, "volumes") {
		sections = append(sections, d.fetchSection(ctx, "Volumes", "/volumes", buildVolumeTable))
	}
	if connector.WantsField(fields, "networks") {
		sections = append(sections, d.fetchSection(ctx, "Networks", "/networks", buildNetworkTable))
	}

	return &connector.ServiceSnapshot{
		ServiceName: "Docker",
		Type:        typeName,
		Sections:    sections,
		Dependencies: []connector.ServiceDependency{
			{Kind: "host", Name: d.host},
		},
		Entities:  entities,
		Metadata:  metadata,
		FetchedAt: start,
	}, nil
}

// Restart restarts the container identified by entityRef (a container ID).
func (d *Connector) Restart(ctx context.Context, _ map[string]any, entityRef string) error {
	if entityRef == "" {
		return fmt.Errorf("docker restart requires a target container ID")
	}
	return d.doPost(ctx, "/containers/"+entityRef+"/restart")
}

// doPost issues a POST with no body, tolerating a 204 No Content response.
func (d *Connector) doPost(ctx context.Context, path string) error {
	if d.configErr != nil {
		return d.configErr
	}
	req, err := http.NewRequestWithContext(ctx, "POST", d.baseURL+path, nil)
	if err != nil {
		return err
	}
	resp, err := d.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return connector.NewTimeoutError(fmt.Errorf("request failed: %w", err))
		}
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return fmt.Errorf("not found: %s", string(data))
	case resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout:
		return connector.NewServiceUnavailableError(fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data)))
	case resp.StatusCode >= 400:
		return fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data))
	}
	return nil
}

// fetchSection runs a single endpoint fetch and renders it with build,
// tolerating failure as a placeholder section rather than failing Fetch.
func (d *Connector) fetchSection(ctx context.Context, title, path string, build func([]byte) string) connector.SnapshotSection {
	raw, err := d.doRequest(ctx, path)
	if err != nil {
		return connector.SnapshotSection{Title: title, Content: "_" + title + " unavailable: " + err.Error() + "_"}
	}
	return connector.SnapshotSection{Title: title, Content: build(raw)}
}

func (d *Connector) doRequest(ctx context.Context, path string) ([]byte, error) {
	if d.configErr != nil {
		return nil, d.configErr
	}
	req, err := http.NewRequestWithContext(ctx, "GET", d.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return nil, connector.NewTimeoutError(fmt.Errorf("request failed: %w", err))
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	switch {
	case resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout:
		return nil, connector.NewServiceUnavailableError(fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data)))
	case resp.StatusCode >= 400:
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// newDockerClient builds an HTTP client and base URL for the given Docker
// host address. unix:// sockets are dialed directly; tcp:// hosts are dialed
// through a guarded dialer that rejects loopback/link-local targets (mirrors
// custom.newGuardedClient — to be shared via connector.go in the pfSense PR)
// and use mutual TLS when tls_cert/tls_key are configured; ssh:// hosts
// tunnel the Engine API over an SSH connection. Any other scheme is
// rejected.
func newDockerClient(host string, config map[string]any) (*http.Client, string, error) {
	switch {
	case strings.HasPrefix(host, "unix://"):
		socketPath := strings.TrimPrefix(host, "unix://")
		transport := &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, "unix", socketPath)
			},
		}
		return &http.Client{Timeout: 30 * time.Second, Transport: transport}, "http://unix", nil

	case strings.HasPrefix(host, "tcp://"):
		return newTCPDockerClient(strings.TrimPrefix(host, "tcp://"), config)

	case strings.HasPrefix(host, "ssh://"):
		return newSSHDockerClient(host, config)

	default:
		return nil, "", fmt.Errorf("unsupported docker host scheme in %q (only unix://, tcp://, and ssh:// are supported)", host)
	}
}

func newTCPDockerClient(addr string, config map[string]any) (*http.Client, string, error) {
	dialer := &net.Dialer{
		Timeout: 30 * time.Second,
		Control: func(_, address string, _ syscall.RawConn) error {
			h, _, err := net.SplitHostPort(address)
			if err != nil {
				return fmt.Errorf("split address %q: %w", address, err)
			}
			ip := net.ParseIP(h)
			if ip == nil {
				return fmt.Errorf("unresolvable address %q", h)
			}
			if ip.IsLoopback() || ip.IsLinkLocalUnicast() {
				return fmt.Errorf("connection to blocked address %s denied", ip)
			}
			return nil
		},
	}

	tlsConfig, err := buildDockerTLSConfig(config)
	if err != nil {
		return nil, "", err
	}
	if tlsConfig == nil {
		transport := &http.Transport{DialContext: dialer.DialContext}
		return &http.Client{Timeout: 30 * time.Second, Transport: transport}, "http://" + addr, nil
	}
	transport := &http.Transport{DialContext: dialer.DialContext, TLSClientConfig: tlsConfig}
	return &http.Client{Timeout: 30 * time.Second, Transport: transport}, "https://" + addr, nil
}

// buildDockerTLSConfig builds the mutual-TLS config for a tcp:// Docker host
// from its tls_cert/tls_key/tls_ca/verify_tls fields, or returns a nil
// config (no error) when no client certificate is configured.
func buildDockerTLSConfig(config map[string]any) (*tls.Config, error) {
	certPEM, _ := config["tls_cert"].(string)
	keyPEM, _ := config["tls_key"].(string)
	if certPEM == "" || keyPEM == "" {
		return nil, nil
	}

	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return nil, fmt.Errorf("parse TLS client certificate/key: %w", err)
	}
	verifyTLS := true
	if v, ok := config["verify_tls"]; ok {
		if b, ok := v.(bool); ok {
			verifyTLS = b
		}
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: !verifyTLS,
	}
	if caPEM, _ := config["tls_ca"].(string); caPEM != "" {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM([]byte(caPEM)) {
			return nil, errors.New("parse TLS CA certificate: no valid certificates found")
		}
		tlsConfig.RootCAs = pool
	}
	return tlsConfig, nil
}

// newSSHDockerClient connects to an ssh:// Docker host by opening an SSH
// session and running "docker system dial-stdio" on the remote end, the
// same mechanism the Docker CLI itself uses for SSH contexts. The session's
// stdin/stdout pipe is wrapped as a net.Conn and reused as the single
// underlying connection for all Engine API requests.
//
// ponytail: one dial-stdio process per connector instance, no connection
// pooling — fine since Fetch only issues serial requests; add pooling if
// concurrent Docker connector requests are ever needed.
func newSSHDockerClient(host string, config map[string]any) (*http.Client, string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, "", fmt.Errorf("parse ssh host %q: %w", host, err)
	}
	addr := u.Host
	if u.Port() == "" {
		addr = net.JoinHostPort(u.Hostname(), "22")
	}
	user, _ := config["ssh_user"].(string)
	if user == "" && u.User != nil {
		user = u.User.Username()
	}
	if user == "" {
		return nil, "", errors.New("ssh docker host requires a username (set ssh_user or user@ in the host URL)")
	}

	auth, err := sshAuthMethods(config)
	if err != nil {
		return nil, "", err
	}

	hostKeyText, _ := config["ssh_host_key"].(string)
	if strings.TrimSpace(hostKeyText) == "" {
		return nil, "", errors.New("ssh docker host requires a pinned host public key (set ssh_host_key)")
	}
	hostPublicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(hostKeyText))
	if err != nil {
		return nil, "", fmt.Errorf("parse ssh_host_key: %w", err)
	}

	sshClient, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.FixedHostKey(hostPublicKey),
		Timeout:         30 * time.Second,
	})
	if err != nil {
		return nil, "", fmt.Errorf("ssh dial %q: %w", addr, err)
	}

	session, err := sshClient.NewSession()
	if err != nil {
		_ = sshClient.Close()
		return nil, "", fmt.Errorf("open ssh session: %w", err)
	}
	stdin, err := session.StdinPipe()
	if err != nil {
		_ = session.Close()
		_ = sshClient.Close()
		return nil, "", fmt.Errorf("open ssh stdin pipe: %w", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		_ = session.Close()
		_ = sshClient.Close()
		return nil, "", fmt.Errorf("open ssh stdout pipe: %w", err)
	}
	session.Stderr = io.Discard
	if err := session.Start("docker system dial-stdio"); err != nil {
		_ = session.Close()
		_ = sshClient.Close()
		return nil, "", fmt.Errorf("start docker system dial-stdio: %w", err)
	}

	conn := &sshStdioConn{stdin: stdin, stdout: stdout, session: session, client: sshClient}
	transport := &http.Transport{
		DialContext: func(context.Context, string, string) (net.Conn, error) {
			if !conn.claim() {
				return nil, errors.New("ssh docker connection already in use")
			}
			return conn, nil
		},
	}
	return &http.Client{Timeout: 30 * time.Second, Transport: transport}, "http://docker", nil
}

func sshAuthMethods(config map[string]any) ([]ssh.AuthMethod, error) {
	if keyPEM, _ := config["ssh_private_key"].(string); keyPEM != "" {
		passphrase, _ := config["ssh_private_key_passphrase"].(string)
		var signer ssh.Signer
		var err error
		if passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(keyPEM), []byte(passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(keyPEM))
		}
		if err != nil {
			return nil, fmt.Errorf("parse ssh private key: %w", err)
		}
		return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
	}
	if password, _ := config["ssh_password"].(string); password != "" {
		return []ssh.AuthMethod{ssh.Password(password)}, nil
	}
	return nil, errors.New("ssh docker host requires ssh_private_key or ssh_password")
}

// sshStdioConn adapts an SSH session's stdin/stdout pipes to a net.Conn for
// use as an http.Transport connection. claim() lets the transport detect
// reuse beyond the single supported connection instead of silently
// corrupting the stream.
type sshStdioConn struct {
	stdin   io.WriteCloser
	stdout  io.Reader
	session *ssh.Session
	client  *ssh.Client
	claimed bool
}

func (c *sshStdioConn) claim() bool {
	if c.claimed {
		return false
	}
	c.claimed = true
	return true
}

func (c *sshStdioConn) Read(b []byte) (int, error)  { return c.stdout.Read(b) }
func (c *sshStdioConn) Write(b []byte) (int, error) { return c.stdin.Write(b) }

func (c *sshStdioConn) Close() error {
	_ = c.stdin.Close()
	sessErr := c.session.Close()
	cliErr := c.client.Close()
	if sessErr != nil {
		return sessErr
	}
	return cliErr
}

func (c *sshStdioConn) LocalAddr() net.Addr              { return dockerSSHAddr{} }
func (c *sshStdioConn) RemoteAddr() net.Addr             { return dockerSSHAddr{} }
func (c *sshStdioConn) SetDeadline(time.Time) error      { return nil }
func (c *sshStdioConn) SetReadDeadline(time.Time) error  { return nil }
func (c *sshStdioConn) SetWriteDeadline(time.Time) error { return nil }

type dockerSSHAddr struct{}

func (dockerSSHAddr) Network() string { return "ssh" }
func (dockerSSHAddr) String() string  { return "docker-ssh-dial-stdio" }

func buildContainerTable(raw []byte) (string, []connector.SnapshotEntity) {
	var containers []struct {
		ID              string   `json:"Id"`
		Names           []string `json:"Names"`
		Image           string   `json:"Image"`
		State           string   `json:"State"`
		Status          string   `json:"Status"`
		NetworkSettings struct {
			Networks map[string]struct {
				IPAddress string `json:"IPAddress"`
			} `json:"Networks"`
		} `json:"NetworkSettings"`
	}
	if err := json.Unmarshal(raw, &containers); err != nil || len(containers) == 0 {
		return "_No containers returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Name | Image | State | Status |\n")
	b.WriteString("|------|-------|-------|--------|\n")
	var entities []connector.SnapshotEntity
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}
		if _, err := fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", name, c.Image, c.State, c.Status); err != nil {
			return "", nil
		}
		ent := connector.SnapshotEntity{Kind: "container", Name: name, ExternalID: c.ID}
		for _, net := range c.NetworkSettings.Networks {
			if net.IPAddress != "" {
				ent.IP = net.IPAddress
				break
			}
		}
		entities = append(entities, ent)
	}
	return b.String(), entities
}

func buildImageTable(raw []byte) string {
	var images []struct {
		RepoTags []string `json:"RepoTags"`
		Size     int64    `json:"Size"`
	}
	if err := json.Unmarshal(raw, &images); err != nil || len(images) == 0 {
		return "_No images returned_"
	}
	var b strings.Builder
	b.WriteString("| Tags | Size (bytes) |\n")
	b.WriteString("|------|---------------|\n")
	for _, img := range images {
		tags := strings.Join(img.RepoTags, ", ")
		if tags == "" {
			tags = "<none>"
		}
		if _, err := fmt.Fprintf(&b, "| %s | %d |\n", tags, img.Size); err != nil {
			return ""
		}
	}
	return b.String()
}

func buildVolumeTable(raw []byte) string {
	var resp struct {
		Volumes []struct {
			Name       string `json:"Name"`
			Driver     string `json:"Driver"`
			Mountpoint string `json:"Mountpoint"`
		} `json:"Volumes"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Volumes) == 0 {
		return "_No volumes returned_"
	}
	var b strings.Builder
	b.WriteString("| Name | Driver | Mountpoint |\n")
	b.WriteString("|------|--------|------------|\n")
	for _, v := range resp.Volumes {
		if _, err := fmt.Fprintf(&b, "| %s | %s | %s |\n", v.Name, v.Driver, v.Mountpoint); err != nil {
			return ""
		}
	}
	return b.String()
}

func buildNetworkTable(raw []byte) string {
	var networks []struct {
		Name   string `json:"Name"`
		Driver string `json:"Driver"`
		Scope  string `json:"Scope"`
	}
	if err := json.Unmarshal(raw, &networks); err != nil || len(networks) == 0 {
		return "_No networks returned_"
	}
	var b strings.Builder
	b.WriteString("| Name | Driver | Scope |\n")
	b.WriteString("|------|--------|-------|\n")
	for _, n := range networks {
		if _, err := fmt.Fprintf(&b, "| %s | %s | %s |\n", n.Name, n.Driver, n.Scope); err != nil {
			return ""
		}
	}
	return b.String()
}
