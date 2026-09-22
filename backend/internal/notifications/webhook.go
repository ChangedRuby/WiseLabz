package notifications

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/crypto"
)

// maxWebhookResponseBytes bounds how much of a webhook response is read before the body is
// discarded. Response bodies are never stored or surfaced.
const maxWebhookResponseBytes = 64 << 10

// Signature headers sent when a channel has a signing secret. The signature is
// hex(HMAC-SHA256(secret, timestamp + "." + body)) prefixed with "sha256=".
const (
	signatureHeader = "X-WiseLabz-Signature"
	timestampHeader = "X-WiseLabz-Timestamp"
)

// webhookClient is the single client used by both first-attempt delivery and retries.
var webhookClient = newWebhookClient()

// newWebhookClient builds an HTTP client whose dialer rejects loopback, link-local,
// unspecified and multicast targets (RFC 1918 stays reachable for self-hosted receivers)
// and which never follows redirects.
func newWebhookClient() *http.Client {
	dialer := connector.GuardedDialer(10 * time.Second)
	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{DialContext: dialer.DialContext},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// signWebhook returns the signature header value for the given timestamp and body.
func signWebhook(secret, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "."))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// signingSecret returns the decrypted signing secret for a channel, or "" when none is
// configured or it cannot be decrypted (delivery then proceeds unsigned).
func (d *Dispatcher) signingSecret(cfg channelCfg) string {
	enc, _ := cfg.Config["secretEncrypted"].(string)
	if enc == "" || d.encKey == nil {
		return ""
	}
	secret, err := crypto.Decrypt(enc, d.encKey)
	if err != nil {
		slog.Error("decrypt webhook signing secret", "error", err, "channel", cfg.Type)
		return ""
	}
	return secret
}

// sendWebhook POSTs a JSON payload to rawURL and treats any transport error or non-2xx response
// as failure. When secret is non-empty the request is signed. Errors never embed the URL, since it
// may carry credentials or tokens and errors are persisted with the delivery record.
func sendWebhook(ctx context.Context, rawURL, secret string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	headers := map[string]string{"Content-Type": "application/json"}
	if secret != "" {
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		headers[timestampHeader] = ts
		headers[signatureHeader] = signWebhook(secret, ts, body)
	}
	return doHTTPRequest(ctx, rawURL, headers, body)
}

// doHTTPRequest POSTs body to rawURL with the given headers over the shared guarded webhook
// transport, treating any transport error or non-2xx response as failure. Shared by every HTTP-based
// channel sender (generic webhook, Discord, Slack, ntfy, Telegram) so they all inherit the same
// SSRF-guarded dialer, redirect policy, and bounded response read. Errors never embed the URL, since
// it may carry credentials or tokens and errors are persisted with the delivery record.
func doHTTPRequest(ctx context.Context, rawURL string, headers map[string]string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewReader(body))
	if err != nil {
		return redactURLError(err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := webhookClient.Do(req)
	if err != nil {
		return redactURLError(err)
	}
	defer resp.Body.Close() //nolint:errcheck
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxWebhookResponseBytes))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request returned status %d", resp.StatusCode)
	}
	return nil
}

// redactURLError drops the request URL that net/http embeds in *url.Error.
func redactURLError(err error) error {
	var ue *url.Error
	if !errors.As(err, &ue) {
		return err
	}
	if errors.Is(ue.Err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: request timed out", ue.Op)
	}
	var ne net.Error
	if errors.As(ue.Err, &ne) && ne.Timeout() {
		return fmt.Errorf("%s: request timed out", ue.Op)
	}
	return fmt.Errorf("%s: %w", ue.Op, ue.Err)
}
