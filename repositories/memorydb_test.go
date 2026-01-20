package repositories

import (
	"api/models"
	"context"
	"testing"
	"time"
)

func TestNewMemoryDB(t *testing.T) {
	ttl := 1 * time.Hour
	repo := NewMemoryDB(ttl)

	if repo == nil {
		t.Fatalf("NewMemoryDB returned nil")
	}
	if repo.ttl != ttl {
		t.Errorf("TTL: expected %v, got %v", ttl, repo.ttl)
	}
	if len(repo.db) != 0 {
		t.Errorf("Initial db length: expected 0, got %d", len(repo.db))
	}
}

func TestMemoryDBCreate(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	err := repo.Create(ctx)
	if err != nil {
		t.Errorf("Create() returned error: %v", err)
	}
}

func TestMemoryDBSave(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	tests := []struct {
		name  string
		event *models.Event
	}{
		{
			name: "Save new event",
			event: &models.Event{
				Id:            "123",
				Date:          time.Now(),
				StatusCode:    200,
				StatusMessage: "OK",
				Metadata: map[string]string{
					"key": "value",
				},
			},
		},
		{
			name: "Update existing event",
			event: &models.Event{
				Id:            "456",
				Date:          time.Now(),
				StatusCode:    404,
				StatusMessage: "Not Found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(ctx, tt.event)
			if err != nil {
				t.Errorf("Save() returned error: %v", err)
			}

			// Verify it's in the db
			found, _ := repo.Get(ctx, tt.event.Id)
			if found == nil {
				t.Errorf("Event not found after Save()")
			}
			if found.Id != tt.event.Id {
				t.Errorf("Event ID mismatch: expected %s, got %s", tt.event.Id, found.Id)
			}
		})
	}
}

func TestMemoryDBSaveExpiration(t *testing.T) {
	ttl := 1 * time.Hour
	repo := NewMemoryDB(ttl)
	ctx := context.Background()

	event := &models.Event{
		Id:            "123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
		Expiration:    0,
	}

	repo.Save(ctx, event)

	if event.Expiration == 0 {
		t.Errorf("Expiration should be set after Save()")
	}

	expected := time.Now().Add(ttl).Unix()
	if event.Expiration < expected-2 || event.Expiration > expected+2 {
		t.Errorf("Expiration not set correctly: got %d, expected around %d", event.Expiration, expected)
	}
}

func TestMemoryDBSaveUpdate(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	event := &models.Event{
		Id:            "123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}

	repo.Save(ctx, event)
	initialLen := len(repo.db)

	event.StatusCode = 404
	repo.Save(ctx, event)

	if len(repo.db) != initialLen {
		t.Errorf("DB length changed after update: expected %d, got %d", initialLen, len(repo.db))
	}

	found, _ := repo.Get(ctx, event.Id)
	if found.StatusCode != 404 {
		t.Errorf("Event not updated: expected StatusCode 404, got %d", found.StatusCode)
	}
}

func TestMemoryDBGet(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	event := &models.Event{
		Id:            "123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}

	repo.Save(ctx, event)

	found, err := repo.Get(ctx, "123")
	if err != nil {
		t.Errorf("Get() returned error: %v", err)
	}
	if found == nil {
		t.Errorf("Get() returned nil event")
	}
	if found.Id != "123" {
		t.Errorf("Get() returned wrong event: expected ID 123, got %s", found.Id)
	}
}

func TestMemoryDBGetNotFound(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	found, err := repo.Get(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Get() returned error: %v", err)
	}
	if found != nil {
		t.Errorf("Get() should return nil for nonexistent event")
	}
}

func TestMemoryDBGetExpired(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	expiredTime := time.Now().Add(-1 * time.Hour)
	event := &models.Event{
		Id:            "expired",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
		Expiration:    expiredTime.Unix(),
	}

	repo.Save(ctx, event)

	found, _ := repo.Get(ctx, "expired")
	if found != nil {
		t.Errorf("Get() should skip expired events")
	}
}

func TestMemoryDBDelete(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	event := &models.Event{
		Id:            "123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}

	repo.Save(ctx, event)
	initialLen := len(repo.db)

	deleted, err := repo.Delete(ctx, "123")
	if err != nil {
		t.Errorf("Delete() returned error: %v", err)
	}
	if deleted == nil {
		t.Errorf("Delete() should return the deleted event")
	}
	if deleted.Id != "123" {
		t.Errorf("Delete() returned wrong event: expected ID 123, got %s", deleted.Id)
	}
	if len(repo.db) != initialLen-1 {
		t.Errorf("Delete() did not remove event from db")
	}

	found, _ := repo.Get(ctx, "123")
	if found != nil {
		t.Errorf("Get() should return nil after Delete()")
	}
}

func TestMemoryDBDeleteNotFound(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	deleted, err := repo.Delete(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Delete() returned error: %v", err)
	}
	if deleted != nil {
		t.Errorf("Delete() should return nil for nonexistent event")
	}
}

func TestMemoryDBFindByDateAndReturnCode(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	now := time.Now()
	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)

	events := []*models.Event{
		{
			Id:            "1",
			Date:          now,
			StatusCode:    200,
			StatusMessage: "OK",
		},
		{
			Id:            "2",
			Date:          now.Add(10 * time.Minute),
			StatusCode:    200,
			StatusMessage: "OK",
		},
		{
			Id:            "3",
			Date:          now.Add(20 * time.Minute),
			StatusCode:    404,
			StatusMessage: "Not Found",
		},
		{
			Id:            "4",
			Date:          now.Add(2 * time.Hour),
			StatusCode:    200,
			StatusMessage: "OK",
		},
	}

	for _, e := range events {
		repo.Save(ctx, e)
	}

	found, err := repo.FindByDateAndReturnCode(ctx, from, to, 200)
	if err != nil {
		t.Errorf("FindByDateAndReturnCode() returned error: %v", err)
	}
	if len(found) != 2 {
		t.Errorf("FindByDateAndReturnCode() returned %d events, expected 2", len(found))
	}

	for _, e := range found {
		if e.StatusCode != 200 {
			t.Errorf("Found event with wrong status code: %d", e.StatusCode)
		}
		if e.Date.Before(from) || e.Date.After(to) {
			t.Errorf("Found event outside date range")
		}
	}
}

func TestMemoryDBFindByDateAndReturnCodeExpired(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	now := time.Now()
	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)

	expiredEvent := &models.Event{
		Id:            "expired",
		Date:          now,
		StatusCode:    200,
		StatusMessage: "OK",
		Expiration:    time.Now().Add(-1 * time.Hour).Unix(),
	}

	validEvent := &models.Event{
		Id:            "valid",
		Date:          now,
		StatusCode:    200,
		StatusMessage: "OK",
		Expiration:    time.Now().Add(1 * time.Hour).Unix(),
	}

	repo.Save(ctx, expiredEvent)
	repo.Save(ctx, validEvent)

	found, _ := repo.FindByDateAndReturnCode(ctx, from, to, 200)
	if len(found) != 1 {
		t.Errorf("FindByDateAndReturnCode() should skip expired events, got %d", len(found))
	}
	if found[0].Id != "valid" {
		t.Errorf("FindByDateAndReturnCode() returned wrong event")
	}
}

func TestMemoryDBFindByDateAndReturnCodeEmpty(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	now := time.Now()
	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)

	found, err := repo.FindByDateAndReturnCode(ctx, from, to, 200)
	if err != nil {
		t.Errorf("FindByDateAndReturnCode() returned error: %v", err)
	}
	if len(found) != 0 {
		t.Errorf("FindByDateAndReturnCode() should return empty slice for no matches")
	}
}

func TestMemoryDBFindByDateAndReturnCodeBoundaries(t *testing.T) {
	repo := NewMemoryDB(1 * time.Hour)
	ctx := context.Background()

	now := time.Now()
	from := now
	to := now.Add(1 * time.Hour)

	event := &models.Event{
		Id:            "boundary",
		Date:          from,
		StatusCode:    200,
		StatusMessage: "OK",
	}

	repo.Save(ctx, event)

	found, _ := repo.FindByDateAndReturnCode(ctx, from, to, 200)
	if len(found) != 1 {
		t.Errorf("FindByDateAndReturnCode() should include boundary date")
	}
}
