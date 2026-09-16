// Package pihole implements a Pi-hole local DNS records connector,
// targeting the Pi-hole v6 REST API.
package pihole

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "pihole"

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "dns",
		Name:     "Pi-hole",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "Pi-hole URL", Type: "text", Required: true, Placeholder: "https://pihole.example.com"},
			{Key: "password", Label: "Password / App Password", Type: "password", Required: true},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Required: false, Default: "true"},
		},
	}, func(config map[string]any) (connector.Connector, error) {
		url, _ := config["url"].(string)
		password, _ := config["password"].(string)
		verifyTLS := true
		if v, ok := config["verify_tls"]; ok {
			if b, ok := v.(bool); ok {
				verifyTLS = b
			}
		}
		dialer := connector.GuardedDialer(30 * time.Second)
		client := &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				DialContext:     dialer.DialContext,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifyTLS},
			},
		}
		return &Connector{
			url:      strings.TrimSuffix(url, "/"),
			password: password,
			client:   client,
		}, nil
	})
}

// Connector fetches local DNS records from a Pi-hole v6 API.
type Connector struct {
	url      string
	password string
	client   *http.Client
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "Pi-hole" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "dns" }

// Validate tests the connection to the Pi-hole API.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	sid, err := c.authenticate(ctx)
	if err != nil {
		return err
	}
	_, err = c.doRequest(ctx, sid, "/api/config/dns/hosts")
	return err
}

// Fetch retrieves Pi-hole local DNS records.
func (c *Connector) Fetch(ctx context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity
	metadata := map[string]string{"pihole_url": c.url}

	sid, err := c.authenticate(ctx)
	if err != nil {
		sections = append(sections, connector.SnapshotSection{Title: "Local DNS Records", Content: "_Local DNS records unavailable: " + err.Error() + "_"})
		return &connector.ServiceSnapshot{
			ServiceName: "Pi-hole",
			Type:        typeName,
			Sections:    sections,
			Metadata:    metadata,
			FetchedAt:   start,
		}, nil
	}

	if raw, err := c.doRequest(ctx, sid, "/api/config/dns/hosts"); err != nil {
		sections = append(sections, connector.SnapshotSection{Title: "Local DNS Records", Content: "_Local DNS records unavailable: " + err.Error() + "_"})
	} else {
		content, ents := buildHostsTable(raw)
		sections = append(sections, connector.SnapshotSection{Title: "Local DNS Records", Content: content})
		entities = ents
	}

	return &connector.ServiceSnapshot{
		ServiceName: "Pi-hole",
		Type:        typeName,
		Sections:    sections,
		Entities:    entities,
		Metadata:    metadata,
		FetchedAt:   start,
	}, nil
}

// Restart restarts the Pi-hole DNS resolver (FTL). Pi-hole manages a single
// implicit DNS service, so entityRef is ignored.
func (c *Connector) Restart(ctx context.Context, _ map[string]any, _ string) error {
	sid, err := c.authenticate(ctx)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/api/action/restartdns", nil)
	if err != nil {
		return err
	}
	req.Header.Set("sid", sid)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
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
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data))
	}
	return nil
}

// Start re-enables Pi-hole DNS blocking. Pi-hole's FTL service has no
// start/stop lifecycle exposed via the API (only restart), so Start/Stop
// here map to the closest on/off concept the API actually exposes: the
// blocking toggle.
func (c *Connector) Start(ctx context.Context, _ map[string]any, _ string) error {
	return c.setBlocking(ctx, true)
}

// Stop disables Pi-hole DNS blocking. See Start for why this maps to the
// blocking toggle rather than a literal service stop.
func (c *Connector) Stop(ctx context.Context, _ map[string]any, _ string) error {
	return c.setBlocking(ctx, false)
}

func (c *Connector) setBlocking(ctx context.Context, blocking bool) error {
	sid, err := c.authenticate(ctx)
	if err != nil {
		return err
	}
	body, err := json.Marshal(map[string]any{"blocking": blocking, "timer": nil})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/api/dns/blocking", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("sid", sid)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
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
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data))
	}
	return nil
}

// WritableFields lists the config-push-eligible local DNS record field.
func (c *Connector) WritableFields() []connector.ConfigField {
	return []connector.ConfigField{
		{Key: "ip", Label: "Record IP", Type: "text", EntityScope: true},
	}
}

// ConfigPush repoints the local DNS record for the hostname identified by
// entityRef to a new IP. Pi-hole's hosts config is a flat set of "ip
// hostname" strings with no update verb, so this deletes the existing
// entry for entityRef (if any) and adds the new one.
// ponytail: one field (ip) per call, matching the handler's one-field
// revert contract — not a batch hosts-file replace.
func (c *Connector) ConfigPush(ctx context.Context, _ map[string]any, entityRef, fieldKey string, value any) error {
	if entityRef == "" {
		return fmt.Errorf("pihole config-push requires a target hostname")
	}
	if fieldKey != "ip" {
		return fmt.Errorf("unsupported field %q", fieldKey)
	}
	newIP, _ := value.(string)
	sid, err := c.authenticate(ctx)
	if err != nil {
		return err
	}
	if raw, err := c.doRequest(ctx, sid, "/api/config/dns/hosts"); err == nil {
		_, entities := buildHostsTable(raw)
		for _, e := range entities {
			if e.Hostname == entityRef && e.IP != "" {
				if err := c.hostsItem(ctx, sid, "DELETE", e.IP, e.Hostname); err != nil {
					return fmt.Errorf("remove old record: %w", err)
				}
				break
			}
		}
	}
	return c.hostsItem(ctx, sid, "PUT", newIP, entityRef)
}

func (c *Connector) hostsItem(ctx context.Context, sid, method, ip, hostname string) error {
	item := ip + " " + hostname
	req, err := http.NewRequestWithContext(ctx, method, c.url+"/api/config/dns/hosts/"+strings.TrimSpace(item), nil)
	if err != nil {
		return err
	}
	req.Header.Set("sid", sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
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
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data))
	}
	return nil
}

// authenticate exchanges the configured password for a session id (sid) via
// POST /api/auth.
func (c *Connector) authenticate(ctx context.Context) (sid string, err error) {
	body, err := json.Marshal(map[string]string{"password": c.password})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/api/auth", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return "", connector.NewTimeoutError(fmt.Errorf("auth request failed: %w", err))
		}
		return "", fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read auth response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", connector.NewAuthError(fmt.Errorf("auth returned %d: %s", resp.StatusCode, string(data)))
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("auth returned %d: %s", resp.StatusCode, string(data))
	}

	var authResp struct {
		Session struct {
			SID     string `json:"sid"`
			Valid   bool   `json:"valid"`
			Message string `json:"message"`
		} `json:"session"`
	}
	if err := json.Unmarshal(data, &authResp); err != nil {
		return "", connector.NewMalformedResponseError(err)
	}
	if authResp.Session.SID == "" {
		return "", connector.NewAuthError(fmt.Errorf("no session id returned: %s", authResp.Session.Message))
	}
	return authResp.Session.SID, nil
}

func (c *Connector) doRequest(ctx context.Context, sid, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.url+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("sid", sid)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
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
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return nil, connector.NewAuthError(fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data)))
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

// buildHostsTable renders the Pi-hole local DNS "IP hostname" entries as a
// markdown table and extracts one SnapshotEntity per record.
func buildHostsTable(raw []byte) (content string, entities []connector.SnapshotEntity) {
	var resp struct {
		Config struct {
			DNS struct {
				Hosts []string `json:"hosts"`
			} `json:"dns"`
		} `json:"config"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "_Local DNS records unavailable: " + connector.NewMalformedResponseError(err).Error() + "_", nil
	}
	if len(resp.Config.DNS.Hosts) == 0 {
		return "_No local DNS records returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Hostname | IP |\n")
	b.WriteString("|----------|----|\n")
	for _, entry := range resp.Config.DNS.Hosts {
		fields := strings.Fields(entry)
		if len(fields) < 2 {
			continue
		}
		ip, hostname := fields[0], fields[1]
		if _, err := fmt.Fprintf(&b, "| %s | %s |\n", hostname, ip); err != nil {
			return "", nil
		}
		entities = append(entities, connector.SnapshotEntity{
			Kind:     "dns_record",
			Hostname: hostname,
			IP:       ip,
		})
	}
	if len(entities) == 0 {
		return "_No local DNS records returned_", nil
	}
	return b.String(), entities
}
