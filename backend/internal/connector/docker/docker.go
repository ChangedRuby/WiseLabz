// Package docker implements a real Docker Engine API connector.
package docker

import (
	"context"
	"net/http"

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
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on "container" entities (see buildContainerTable/enrichContainerAttributes),
// exposed via GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"container": {
		{Name: "privileged", Type: "boolean", Description: "Whether the container runs in privileged mode"},
		{Name: "network_mode", Type: "string", Description: "Docker network mode (bridge, host, none, container:<id>, ...)"},
		{Name: "restart_policy", Type: "string", Description: "Restart policy name (no, always, unless-stopped, on-failure)"},
		{Name: "published_ports", Type: "string_array", Description: "Published host:container/protocol port mappings"},
		{Name: "user", Type: "string", Description: "User the container's main process runs as"},
		{Name: "read_only_rootfs", Type: "boolean", Description: "Whether the container's root filesystem is read-only"},
		{Name: "image", Type: "string", Description: "Image reference the container was created from"},
	},
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
