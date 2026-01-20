package apis

import (
	"api/repositories"
	"context"
	"testing"
	"time"
)

func TestNewHttpApi(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewHttpApi(8080, repo)

	if api == nil {
		t.Fatalf("NewHttpApi returned nil")
	}
	if api.repository != repo {
		t.Errorf("Repository not set correctly")
	}
	if api.router == nil {
		t.Errorf("Router is nil")
	}
	if api.server == nil {
		t.Errorf("Server is nil")
	}
	if api.server.Addr != ":8080" {
		t.Errorf("Server address: expected ':8080', got '%s'", api.server.Addr)
	}
}

func TestNewHttpApiDefaultPort(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewHttpApi(3000, repo)

	if api.server.Addr != ":3000" {
		t.Errorf("Server address: expected ':3000', got '%s'", api.server.Addr)
	}
}

func TestHttpApiServerConfig(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewHttpApi(8080, repo)

	if api.server.IdleTimeout != 60*time.Second {
		t.Errorf("IdleTimeout: expected 60s, got %v", api.server.IdleTimeout)
	}
	if api.server.ReadTimeout != 30*time.Second {
		t.Errorf("ReadTimeout: expected 30s, got %v", api.server.ReadTimeout)
	}
	if api.server.WriteTimeout != 30*time.Second {
		t.Errorf("WriteTimeout: expected 30s, got %v", api.server.WriteTimeout)
	}
	if api.server.MaxHeaderBytes != 1024*1024 {
		t.Errorf("MaxHeaderBytes: expected 1048576, got %d", api.server.MaxHeaderBytes)
	}
}

func TestHttpApiRouterNotNil(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewHttpApi(8080, repo)

	if api.router == nil {
		t.Errorf("Router should not be nil")
	}
}

func TestHttpApiRepositoryAssignment(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewHttpApi(8080, repo)

	if api.repository == nil {
		t.Errorf("Repository should not be nil")
	}

	// Test that we can use the repository
	err := api.repository.Create(context.Background())
	if err != nil {
		t.Errorf("Repository.Create() returned error: %v", err)
	}
}

func TestHttpApiMultipleInstances(t *testing.T) {
	repo1 := repositories.NewMemoryDB(1 * time.Hour)
	repo2 := repositories.NewMemoryDB(1 * time.Hour)

	api1 := NewHttpApi(8080, repo1)
	api2 := NewHttpApi(9080, repo2)

	if api1.server.Addr == api2.server.Addr {
		t.Errorf("Instances should have different addresses")
	}
}

func TestHttpApiHandler(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewHttpApi(8080, repo)

	if api.router == nil {
		t.Errorf("Router should not be nil")
	}
}

func TestHttpApiPortConfiguration(t *testing.T) {
	ports := []int{3000, 5000, 8000, 9000}

	for _, port := range ports {
		repo := repositories.NewMemoryDB(1 * time.Hour)
		api := NewHttpApi(port, repo)

		if api.server.Addr[0] != ':' {
			t.Errorf("Server address should start with ':'")
		}
	}
}

func TestHttpApiServerHandler(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewHttpApi(8080, repo)

	if api.server.Handler == nil {
		t.Errorf("Server handler should not be nil")
	}
}
