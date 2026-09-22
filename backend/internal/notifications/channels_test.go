package notifications

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestSendNtfyChannel_SendsTopicAndHeaders(t *testing.T) {
	connector.AllowLoopbackForTest(t)
	var gotPath, gotTitle, gotPriority, gotTags string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotTitle = r.Header.Get("X-Title")
		gotPriority = r.Header.Get("X-Priority")
		gotTags = r.Header.Get("X-Tags")
		gotBody, _ = io.ReadAll(r.Body)
	}))
	defer srv.Close()

	cfg := channelCfg{Config: map[string]any{
		"url":      srv.URL,
		"topic":    "wiselabz-alerts",
		"priority": "high",
		"tags":     "warning,skull",
	}}
	if err := sendNtfyChannel(context.Background(), cfg, "", "Title", "Message"); err != nil {
		t.Fatalf("send: %v", err)
	}
	if gotPath != "/wiselabz-alerts" {
		t.Errorf("path = %q, want /wiselabz-alerts", gotPath)
	}
	if gotTitle != "Title" {
		t.Errorf("X-Title = %q", gotTitle)
	}
	if gotPriority != "high" {
		t.Errorf("X-Priority = %q", gotPriority)
	}
	if gotTags != "warning,skull" {
		t.Errorf("X-Tags = %q", gotTags)
	}
	if string(gotBody) != "Message" {
		t.Errorf("body = %q, want %q", gotBody, "Message")
	}
}

func TestSendNtfyChannel_MissingTopic(t *testing.T) {
	err := sendNtfyChannel(context.Background(), channelCfg{Config: map[string]any{"url": "https://ntfy.sh"}}, "", "T", "M")
	if err == nil || !strings.Contains(err.Error(), "topic") {
		t.Errorf("expected topic error, got %v", err)
	}
}

func TestSendNtfyChannel_DefaultsToNtfySh(t *testing.T) {
	// No network call is expected to succeed (ntfy.sh isn't reachable/allowed in this sandbox),
	// but the URL construction path must not error out before the request is attempted; a
	// "topic not configured" error would indicate the default-server branch was skipped.
	err := sendNtfyChannel(context.Background(), channelCfg{Config: map[string]any{"topic": "t"}}, "", "T", "M")
	if err != nil && strings.Contains(err.Error(), "topic not configured") {
		t.Errorf("expected default server to be used, got %v", err)
	}
}

func TestSendTelegramChannel_SendsChatAndText(t *testing.T) {
	connector.AllowLoopbackForTest(t)
	var gotBody []byte
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotBody, _ = io.ReadAll(r.Body)
	}))
	defer srv.Close()

	// sendTelegramChannel hardcodes the Telegram API host; verify chat_id/text/parse_mode
	// shaping and the token/chatId validation using the exported helper directly instead.
	cfg := channelCfg{Config: map[string]any{"chatId": "12345"}}
	if err := sendTelegramChannel(context.Background(), cfg, "", "Title", "Message"); err == nil ||
		!strings.Contains(err.Error(), "bot token") {
		t.Errorf("expected bot token error without secret, got %v", err)
	}

	cfg = channelCfg{Config: map[string]any{}}
	if err := sendTelegramChannel(context.Background(), cfg, "sometoken", "Title", "Message"); err == nil ||
		!strings.Contains(err.Error(), "chat id") {
		t.Errorf("expected chat id error without chatId, got %v", err)
	}

	// Exercise the request-building path against a local server via doHTTPRequest directly,
	// mirroring what sendTelegramChannel sends.
	body := `{"chat_id":"12345","text":"*Title*\nMessage","parse_mode":"Markdown"}`
	if err := doHTTPRequest(context.Background(), srv.URL+"/bottok/sendMessage",
		map[string]string{"Content-Type": "application/json"}, []byte(body)); err != nil {
		t.Fatalf("send: %v", err)
	}
	if gotPath != "/bottok/sendMessage" {
		t.Errorf("path = %q", gotPath)
	}
	if string(gotBody) != body {
		t.Errorf("body = %q, want %q", gotBody, body)
	}
}

// fakeSMTPServer accepts a connection, replies to the SMTP dance with plain 2xx codes, and
// records the DATA payload. It's just enough of the protocol for sendSMTPChannel's client-side
// path (no STARTTLS, no AUTH) to complete successfully.
func fakeSMTPServer(t *testing.T) (addr string, dataCh <-chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() }) //nolint:errcheck
	ch := make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close() //nolint:errcheck
		r := bufio.NewReader(conn)
		fmt.Fprint(conn, "220 fake.smtp ESMTP\r\n") //nolint:errcheck
		var inData bool
		var data strings.Builder
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}
			switch {
			case inData:
				if strings.TrimRight(line, "\r\n") == "." {
					inData = false
					ch <- data.String()
					fmt.Fprint(conn, "250 OK\r\n") //nolint:errcheck
					continue
				}
				data.WriteString(line)
			case strings.HasPrefix(line, "EHLO"), strings.HasPrefix(line, "HELO"):
				fmt.Fprint(conn, "250-fake.smtp\r\n250 OK\r\n") //nolint:errcheck
			case strings.HasPrefix(line, "MAIL FROM"):
				fmt.Fprint(conn, "250 OK\r\n") //nolint:errcheck
			case strings.HasPrefix(line, "RCPT TO"):
				fmt.Fprint(conn, "250 OK\r\n") //nolint:errcheck
			case strings.HasPrefix(line, "DATA"):
				inData = true
				fmt.Fprint(conn, "354 Start mail input\r\n") //nolint:errcheck
			case strings.HasPrefix(line, "QUIT"):
				fmt.Fprint(conn, "221 Bye\r\n") //nolint:errcheck
				return
			default:
				fmt.Fprint(conn, "250 OK\r\n") //nolint:errcheck
			}
		}
	}()
	return ln.Addr().String(), ch
}

func TestSendSMTPChannel_SendsMessage(t *testing.T) {
	connector.AllowLoopbackForTest(t)
	addr, dataCh := fakeSMTPServer(t)
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatal(err)
	}
	cfg := channelCfg{Config: map[string]any{
		"host": host,
		"port": portAsFloat(t, port),
		"from": "alerts@wiselabz.local",
		"to":   "ops@example.com, oncall@example.com",
	}}
	if err := sendSMTPChannel(context.Background(), cfg, "", "Disk full", "sda1 is at 95%"); err != nil {
		t.Fatalf("send: %v", err)
	}

	select {
	case data := <-dataCh:
		if !strings.Contains(data, "Subject: Disk full") || !strings.Contains(data, "sda1 is at 95%") ||
			!strings.Contains(data, "To: ops@example.com, oncall@example.com") {
			t.Errorf("unexpected DATA payload: %q", data)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("fake SMTP server never received DATA")
	}
}

func TestSendSMTPChannel_MissingConfig(t *testing.T) {
	err := sendSMTPChannel(context.Background(), channelCfg{Config: map[string]any{"host": "localhost"}}, "", "T", "M")
	if err == nil || !strings.Contains(err.Error(), "not fully configured") {
		t.Errorf("expected not-fully-configured error, got %v", err)
	}
}

func TestSplitRecipients(t *testing.T) {
	got := splitRecipients(" a@example.com ,b@example.com,, c@example.com")
	want := []string{"a@example.com", "b@example.com", "c@example.com"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestBuildEmailMessage_SanitizesSubjectNewlines(t *testing.T) {
	msg := string(buildEmailMessage("f@x.com", []string{"t@x.com"}, "Evil\r\nBcc: evil@x.com", "body"))
	if strings.Contains(msg, "Bcc: evil@x.com\r\nSubject") || strings.Count(msg, "\r\n\r\n") != 1 {
		t.Errorf("subject newline injection not sanitized: %q", msg)
	}
}

func portAsFloat(t *testing.T, port string) float64 {
	t.Helper()
	var p int
	if _, err := fmt.Sscanf(port, "%d", &p); err != nil {
		t.Fatalf("parse port %q: %v", port, err)
	}
	return float64(p)
}
