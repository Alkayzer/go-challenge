package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateEventRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateEventRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateEventRequest{
				Title:     "Test Event",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
			wantErr: false,
		},
		{
			name: "empty title",
			req: CreateEventRequest{
				Title:     "",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
			wantErr: true,
		},
		{
			name: "title too long",
			req: CreateEventRequest{
				Title:     "This is a very long title that exceeds the maximum allowed length of 100 characters and should cause a validation error",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
			wantErr: true,
		},
		{
			name: "start time after end time",
			req: CreateEventRequest{
				Title:     "Test Event",
				StartTime: time.Now().Add(time.Hour),
				EndTime:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "zero start time",
			req: CreateEventRequest{
				Title:     "Test Event",
				StartTime: time.Time{},
				EndTime:   time.Now().Add(time.Hour),
			},
			wantErr: true,
		},
		{
			name: "zero end time",
			req: CreateEventRequest{
				Title:     "Test Event",
				StartTime: time.Now(),
				EndTime:   time.Time{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEventRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEvent(t *testing.T) {
	req := &CreateEventRequest{
		Title:       "Test Event",
		Description: stringPtr("Test Description"),
		StartTime:   time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC),
	}

	event := NewEvent(req)

	if event.Title != req.Title {
		t.Errorf("Expected title %s, got %s", req.Title, event.Title)
	}

	if *event.Description != *req.Description {
		t.Errorf("Expected description %s, got %s", *req.Description, *event.Description)
	}

	if !event.StartTime.Equal(req.StartTime) {
		t.Errorf("Expected start time %v, got %v", req.StartTime, event.StartTime)
	}

	if !event.EndTime.Equal(req.EndTime) {
		t.Errorf("Expected end time %v, got %v", req.EndTime, event.EndTime)
	}

	if event.ID == uuid.Nil {
		t.Error("Expected non-nil UUID")
	}

	if event.CreatedAt.IsZero() {
		t.Error("Expected non-zero created at time")
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
