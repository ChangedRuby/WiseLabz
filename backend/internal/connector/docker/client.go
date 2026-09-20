package docker

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// doPost issues a POST with no body, tolerating a 204 No Content response.
func (d *Connector) doPost(ctx context.Context, path string) error {
	return d.doPostBody(ctx, path, nil)
}

// doPostBody issues a POST with an optional JSON body, tolerating a 204 No
// Content response.
func (d *Connector) doPostBody(ctx context.Context, path string, body []byte) error {
	if d.configErr != nil {
		return d.configErr
	}
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", d.baseURL+path, reqBody)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := d.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return connector.NewTimeoutError(fmt.Errorf("request failed: %w", err))
		}
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck
	data, err := connector.ReadBody(resp.Body)
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

	data, err := connector.ReadBody(resp.Body)
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
		if err := connector.ValidateUnixSocketPath(socketPath); err != nil {
			return nil, "", err
		}
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
			if connector.IsDangerousIP(ip) {
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
		MinVersion:         tls.VersionTLS12,
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
