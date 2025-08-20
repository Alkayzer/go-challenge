package service

import (
	"context"
	"fmt"
	"log"

	"events-service/internal/database"
	"events-service/internal/models"

	"github.com/google/uuid"
)

type EventService struct {
	repo database.Repository
}

func NewEventService(repo database.Repository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) CreateEvent(ctx context.Context, req *models.CreateEventRequest) (*models.Event, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	event := models.NewEvent(req)

	if err := s.repo.CreateEvent(ctx, event); err != nil {
		log.Printf("Failed to create event: %v", err)
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	log.Printf("Created event with ID: %s", event.ID)
	return event, nil
}

func (s *EventService) GetEventByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	event, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		log.Printf("Failed to get event by ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return event, nil
}

func (s *EventService) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	events, err := s.repo.GetAllEvents(ctx)
	if err != nil {
		log.Printf("Failed to get all events: %v", err)
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}
