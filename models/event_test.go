package models

import (
	"testing"
	"time"
)

func TestEventValidate(t *testing.T) {
	tests := []struct {
		name    string
		event   *Event
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid event",
			event: &Event{
				Id:            "123",
				Date:          time.Now(),
				StatusCode:    200,
				StatusMessage: "OK",
				Expiration:    0,
				Metadata: map[string]string{
					"key": "value",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid date - zero time",
			event: &Event{
				Id:            "123",
				Date:          time.Time{},
				StatusCode:    200,
				StatusMessage: "OK",
			},
			wantErr: true,
			errMsg:  "invalid date",
		},
		{
			name: "Invalid status code - negative",
			event: &Event{
				Id:            "123",
				Date:          time.Now(),
				StatusCode:    -1,
				StatusMessage: "Error",
			},
			wantErr: true,
			errMsg:  "invalid status code",
		},
		{
			name: "Valid event with zero status code",
			event: &Event{
				Id:            "123",
				Date:          time.Now(),
				StatusCode:    0,
				StatusMessage: "No Status",
			},
			wantErr: false,
		},
		{
			name: "Valid event with large status code",
			event: &Event{
				Id:            "123",
				Date:          time.Now(),
				StatusCode:    599,
				StatusMessage: "Error",
			},
			wantErr: false,
		},
		{
			name: "Valid event without metadata",
			event: &Event{
				Id:            "123",
				Date:          time.Now(),
				StatusCode:    200,
				StatusMessage: "OK",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestEventFields(t *testing.T) {
	now := time.Now()
	metadata := map[string]string{
		"env":    "test",
		"region": "us-east-1",
	}

	event := &Event{
		Id:            "test-id",
		Date:          now,
		StatusCode:    201,
		StatusMessage: "Created",
		Expiration:    now.Add(1 * time.Hour).Unix(),
		Metadata:      metadata,
	}

	if event.Id != "test-id" {
		t.Errorf("Id field: expected 'test-id', got '%s'", event.Id)
	}
	if event.Date != now {
		t.Errorf("Date field: expected %v, got %v", now, event.Date)
	}
	if event.StatusCode != 201 {
		t.Errorf("StatusCode field: expected 201, got %d", event.StatusCode)
	}
	if event.StatusMessage != "Created" {
		t.Errorf("StatusMessage field: expected 'Created', got '%s'", event.StatusMessage)
	}
	if len(event.Metadata) != 2 {
		t.Errorf("Metadata length: expected 2, got %d", len(event.Metadata))
	}
}
