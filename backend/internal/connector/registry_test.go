package connector

import (
	"context"
	"testing"
	"time"
)

type registryTestRefresher struct{}

func (registryTestRefresher) Name() string     { return "refresher" }
func (registryTestRefresher) Type() string     { return "registry_test_refresher" }
func (registryTestRefresher) Category() string { return "test" }
func (registryTestRefresher) Fetch(_ context.Context, _ map[string]any) (*ServiceSnapshot, error) {
	return nil, nil
}
func (registryTestRefresher) Validate(_ context.Context, _ map[string]any) error { return nil }
func (registryTestRefresher) RefreshCredentials(_ context.Context, config map[string]any) (map[string]any, time.Time, error) {
	return config, time.Time{}, nil
}

func TestRegisterStubRoundTrips(t *testing.T) {
	Register(TypeSchema{Type: "registry_test_stub", Category: "test", Name: "Stub", Stub: true},
		func(_ map[string]any) (Connector, error) { return nil, nil })

	got, err := GetTypeSchema("registry_test_stub")
	if err != nil {
		t.Fatalf("GetTypeSchema: %v", err)
	}
	if !got.Stub {
		t.Errorf("Stub = false, want true")
	}

	found := false
	for _, s := range ListSchemas() {
		if s.Type == "registry_test_stub" {
			found = true
			if !s.Stub {
				t.Errorf("ListSchemas: Stub = false, want true")
			}
		}
	}
	if !found {
		t.Fatal("registry_test_stub not found in ListSchemas()")
	}
}

func TestRegisterDefaultsToNonStub(t *testing.T) {
	Register(TypeSchema{Type: "registry_test_real", Category: "test", Name: "Real"},
		func(_ map[string]any) (Connector, error) { return nil, nil })

	got, err := GetTypeSchema("registry_test_real")
	if err != nil {
		t.Fatalf("GetTypeSchema: %v", err)
	}
	if got.Stub {
		t.Errorf("Stub = true, want false")
	}
}

func TestIsCredentialRefresherType(t *testing.T) {
	Register(TypeSchema{Type: "registry_test_refresher", Category: "test", Name: "Refresher"},
		func(_ map[string]any) (Connector, error) { return registryTestRefresher{}, nil })
	Register(TypeSchema{Type: "registry_test_nonrefresher", Category: "test", Name: "Non-refresher"},
		func(_ map[string]any) (Connector, error) { return nil, nil })

	if !IsCredentialRefresherType("registry_test_refresher") {
		t.Error("IsCredentialRefresherType(refresher) = false, want true")
	}
	if IsCredentialRefresherType("registry_test_nonrefresher") {
		t.Error("IsCredentialRefresherType(non-refresher) = true, want false")
	}
	if IsCredentialRefresherType("registry_test_unknown_type") {
		t.Error("IsCredentialRefresherType(unknown) = true, want false")
	}

	for _, s := range ListSchemas() {
		switch s.Type {
		case "registry_test_refresher":
			if !s.IsCredentialRefresher {
				t.Error("ListSchemas: refresher's IsCredentialRefresher = false, want true")
			}
		case "registry_test_nonrefresher":
			if s.IsCredentialRefresher {
				t.Error("ListSchemas: non-refresher's IsCredentialRefresher = true, want false")
			}
		}
	}
}
