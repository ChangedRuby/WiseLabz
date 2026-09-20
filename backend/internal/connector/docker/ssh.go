package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// newSSHDockerClient connects to an ssh:// Docker host by opening an SSH
// session and running "docker system dial-stdio" on the remote end, the
// same mechanism the Docker CLI itself uses for SSH contexts. The session's
// stdin/stdout pipe is wrapped as a net.Conn. The SSH connection is dialed
// lazily per HTTP request and closed by the transport afterwards, so
// discarding the client never leaks a connection, session or remote process.
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

	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.FixedHostKey(hostPublicKey),
		Timeout:         30 * time.Second,
	}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialSSHStdio(ctx, addr, sshConfig)
		},
		// Each request gets its own SSH connection that the transport closes
		// as soon as the response is consumed, so nothing outlives the client.
		DisableKeepAlives: true,
	}
	return &http.Client{Timeout: 30 * time.Second, Transport: transport}, "http://docker", nil
}

// dialSSHStdio opens an SSH connection and starts "docker system dial-stdio"
// on it, returning the session's pipes as a net.Conn. Closing the conn tears
// down the remote process, session and SSH client. ctx cancels the TCP dial
// and the SSH handshake; cfg.Timeout still bounds the TCP dial.
func dialSSHStdio(ctx context.Context, addr string, cfg *ssh.ClientConfig) (net.Conn, error) {
	dialer := net.Dialer{Timeout: cfg.Timeout}
	rawConn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("ssh dial %q: %w", addr, err)
	}
	// ssh.NewClientConn takes no ctx; closing the conn aborts a blocked handshake.
	handshakeDone := make(chan struct{})
	stop := context.AfterFunc(ctx, func() {
		closeQuietly("rawConn", rawConn)
		close(handshakeDone)
	})
	if cfg.Timeout > 0 {
		_ = rawConn.SetDeadline(time.Now().Add(cfg.Timeout))
	}
	sshConn, chans, reqs, err := ssh.NewClientConn(rawConn, addr, cfg)
	if !stop() {
		// ctx fired; the AfterFunc closed the conn.
		<-handshakeDone
		if err == nil {
			closeQuietly("sshConn", sshConn)
		}
		return nil, fmt.Errorf("ssh dial %q: %w", addr, ctx.Err())
	}
	if err != nil {
		closeQuietly("rawConn", rawConn)
		return nil, fmt.Errorf("ssh dial %q: %w", addr, err)
	}
	_ = rawConn.SetDeadline(time.Time{})
	sshClient := ssh.NewClient(sshConn, chans, reqs)

	session, err := sshClient.NewSession()
	if err != nil {
		closeQuietly("sshClient", sshClient)
		return nil, fmt.Errorf("open ssh session: %w", err)
	}
	stdin, err := session.StdinPipe()
	if err != nil {
		closeQuietly("session", session)
		closeQuietly("sshClient", sshClient)
		return nil, fmt.Errorf("open ssh stdin pipe: %w", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		closeQuietly("session", session)
		closeQuietly("sshClient", sshClient)
		return nil, fmt.Errorf("open ssh stdout pipe: %w", err)
	}
	session.Stderr = io.Discard
	if err := session.Start("docker system dial-stdio"); err != nil {
		closeQuietly("session", session)
		closeQuietly("sshClient", sshClient)
		return nil, fmt.Errorf("start docker system dial-stdio: %w", err)
	}

	return &sshStdioConn{stdin: stdin, stdout: stdout, session: session, client: sshClient}, nil
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
// use as an http.Transport connection. Deadlines are not supported (the
// SetDeadline methods are no-ops); http.Client cancels a request by closing
// the conn, which tears down the session.
type sshStdioConn struct {
	stdin   io.WriteCloser
	stdout  io.Reader
	session *ssh.Session
	client  *ssh.Client
}

func (c *sshStdioConn) Read(b []byte) (int, error) { return c.stdout.Read(b) }

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

func (c *sshStdioConn) LocalAddr() net.Addr { return dockerSSHAddr{} }

func (c *sshStdioConn) RemoteAddr() net.Addr { return dockerSSHAddr{} }

func (c *sshStdioConn) SetDeadline(time.Time) error { return nil }

func (c *sshStdioConn) SetReadDeadline(time.Time) error { return nil }

func (c *sshStdioConn) SetWriteDeadline(time.Time) error { return nil }

type dockerSSHAddr struct{}

func (dockerSSHAddr) Network() string { return "ssh" }

func (dockerSSHAddr) String() string { return "docker-ssh-dial-stdio" }

// closeQuietly closes c on a cleanup path where the caller is already
// returning a more relevant error; a Close failure is logged at debug only.
func closeQuietly(name string, c io.Closer) {
	if err := c.Close(); err != nil {
		slog.Debug("docker connector: close failed", "resource", name, "error", err)
	}
}
