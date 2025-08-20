package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"events-service/internal/models"
	"events-service/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	event, err := h.eventService.CreateEvent(ctx, &req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, models.EventResponse{Event: event})
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	event, err := h.eventService.GetEventByID(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Event not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.EventResponse{Event: event})
}

func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	events, err := h.eventService.GetAllEvents(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}

	respondWithJSON(w, http.StatusOK, models.EventsResponse{
		Events: events,
		Count:  len(events),
	})
}

func (h *EventHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{
		"error": message,
	})
}
