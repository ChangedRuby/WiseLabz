package proxmox

import (
	"context"
	"fmt"
	"io"
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
		return nil, connector.MapTransportError(err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := connector.ReadBody(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if statusErr := connector.CheckStatus(resp.StatusCode, data); statusErr != nil {
		return nil, statusErr
	}

	return data, nil
}
