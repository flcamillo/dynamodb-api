package apis

import (
	"api/repositories"
	"testing"
	"time"
)

func TestNewLambdaApi(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewLambdaApi(repo)

	if api == nil {
		t.Fatalf("NewLambdaApi returned nil")
	}
	if api.repository != repo {
		t.Errorf("Repository not set correctly")
	}
}

func TestLambdaApiRepositoryAssignment(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewLambdaApi(repo)

	if api.repository == nil {
		t.Errorf("Repository should not be nil")
	}
}

func TestLambdaApiMultipleInstances(t *testing.T) {
	repo1 := repositories.NewMemoryDB(1 * time.Hour)
	repo2 := repositories.NewMemoryDB(1 * time.Hour)

	api1 := NewLambdaApi(repo1)
	api2 := NewLambdaApi(repo2)

	if api1.repository == api2.repository {
		t.Errorf("Instances should have different repositories")
	}
}

func TestLambdaApiRepositoryNotNil(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewLambdaApi(repo)

	if api.repository == nil {
		t.Errorf("Repository should not be nil")
	}
}

func TestLambdaApiStructure(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	api := NewLambdaApi(repo)

	// Verify the API is properly structured
	if api == nil {
		t.Fatalf("API should not be nil")
	}

	// Verify we can create the API and it has the right repository
	if api.repository != repo {
		t.Errorf("Repository not assigned correctly")
	}
}

func TestLambdaApiDifferentRepositories(t *testing.T) {
	repo1 := repositories.NewMemoryDB(1 * time.Hour)
	repo2 := repositories.NewMemoryDB(2 * time.Hour)

	api1 := NewLambdaApi(repo1)
	api2 := NewLambdaApi(repo2)

	if api1.repository == api2.repository {
		t.Errorf("Different API instances should have different repositories")
	}

	// Verify that each API uses its own repository
	if api1.repository != repo1 {
		t.Errorf("API1 should use repo1")
	}
	if api2.repository != repo2 {
		t.Errorf("API2 should use repo2")
	}
}
