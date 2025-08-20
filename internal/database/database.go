package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"events-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEventByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetAllEvents(ctx context.Context) ([]models.Event, error)
	Close() error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, databaseURL string) (*PostgresRepository, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	return &PostgresRepository{pool: pool}, nil
}

func (r *PostgresRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO events (id, title, description, start_time, end_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(ctx, query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetEventByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time, created_at
		FROM events
		WHERE id = $1
	`

	var event models.Event
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}

	return &event, nil
}

func (r *PostgresRepository) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time, created_at
		FROM events
		ORDER BY start_time ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

func (r *PostgresRepository) Close() error {
	r.pool.Close()
	return nil
}
