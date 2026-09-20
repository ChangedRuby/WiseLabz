package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// Restart restarts the container identified by entityRef (a container ID).
func (d *Connector) Restart(ctx context.Context, _ map[string]any, entityRef string) error {
	if entityRef == "" {
		return fmt.Errorf("docker restart requires a target container ID")
	}
	if err := connector.ValidateRefSegment(entityRef); err != nil {
		return fmt.Errorf("invalid entityRef: %w", err)
	}
	return d.doPost(ctx, "/containers/"+entityRef+"/restart")
}

// Start starts the container identified by entityRef (a container ID).
// Idempotent-safe: the Docker Engine API returns 304 Not Modified (treated
// as success) for an already-running container.
func (d *Connector) Start(ctx context.Context, _ map[string]any, entityRef string) error {
	if entityRef == "" {
		return fmt.Errorf("docker start requires a target container ID")
	}
	if err := connector.ValidateRefSegment(entityRef); err != nil {
		return fmt.Errorf("invalid entityRef: %w", err)
	}
	return d.doPost(ctx, "/containers/"+entityRef+"/start")
}

// Stop stops the container identified by entityRef (a container ID).
func (d *Connector) Stop(ctx context.Context, _ map[string]any, entityRef string) error {
	if entityRef == "" {
		return fmt.Errorf("docker stop requires a target container ID")
	}
	if err := connector.ValidateRefSegment(entityRef); err != nil {
		return fmt.Errorf("invalid entityRef: %w", err)
	}
	return d.doPost(ctx, "/containers/"+entityRef+"/stop")
}

// WritableFields lists the config-push-eligible container fields.
// ponytail: Docker's Engine API has no live image-swap for a running
// container (that needs a full stop/remove/recreate) so the whitelist
// targets what /containers/{id}/update can actually patch in place —
// restart policy — rather than the image tag; image-tag push is a future
// extension once recreate-with-rollback is designed.
func (d *Connector) WritableFields() []connector.ConfigField {
	return []connector.ConfigField{
		{Key: "restartPolicy", Label: "Restart Policy", Type: "select", EntityScope: true},
	}
}

// ConfigPush updates the container identified by entityRef's restart
// policy via Docker's /containers/{id}/update endpoint.
func (d *Connector) ConfigPush(ctx context.Context, _ map[string]any, entityRef, fieldKey string, value any) error {
	if entityRef == "" {
		return fmt.Errorf("docker config-push requires a target container ID")
	}
	if err := connector.ValidateRefSegment(entityRef); err != nil {
		return fmt.Errorf("invalid entityRef: %w", err)
	}
	if fieldKey != "restartPolicy" {
		return fmt.Errorf("unsupported field %q", fieldKey)
	}
	name, _ := value.(string)
	body, err := json.Marshal(map[string]any{"RestartPolicy": map[string]string{"Name": name}})
	if err != nil {
		return err
	}
	return d.doPostBody(ctx, "/containers/"+entityRef+"/update", body)
}
