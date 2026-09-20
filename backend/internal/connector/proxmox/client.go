package proxmox

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func (p *Connector) doRequest(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	url := p.url + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("Authorization", "PVEAPIToken="+p.tokenID+"="+p.tokenSecret)
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
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
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return nil, connector.NewAuthError(fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data)))
	case resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout:
		return nil, connector.NewServiceUnavailableError(fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data)))
	case resp.StatusCode >= 400:
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// isTimeout reports whether err represents a request deadline being
// exceeded, covering both a canceled context and a net.Error timeout
// (e.g. a dial or read timing out on the underlying transport).
func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
