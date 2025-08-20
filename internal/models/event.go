package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description,omitempty" db:"description"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required,max=100"`
	Description *string   `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
}

type EventResponse struct {
	Event *Event `json:"event,omitempty"`
	Error string `json:"error,omitempty"`
}

type EventsResponse struct {
	Events []Event `json:"events"`
	Count  int     `json:"count"`
}

func (req *CreateEventRequest) Validate() error {
	if req.Title == "" {
		return errors.New("title is required")
	}

	if len(req.Title) > 100 {
		return errors.New("title must be 100 characters or less")
	}

	if req.StartTime.IsZero() {
		return errors.New("start_time is required")
	}

	if req.EndTime.IsZero() {
		return errors.New("end_time is required")
	}

	if !req.StartTime.Before(req.EndTime) {
		return errors.New("start_time must be before end_time")
	}

	return nil
}

func NewEvent(req *CreateEventRequest) *Event {
	return &Event{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		CreatedAt:   time.Now().UTC(),
	}
}
