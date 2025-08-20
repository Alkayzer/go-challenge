# Events Service - Go RESTful API

A simple backend service in Go that manages a collection of "Events" with PostgreSQL integration.

## Features

- **RESTful API** with proper HTTP status codes
- **PostgreSQL Integration** with connection pooling
- **Input Validation** with comprehensive error messages
- **UUID Generation** for event IDs
- **JSON Handling** with proper struct tags
- **Context Usage** for request handling

## API Endpoints

### Create Event
```http
POST /events
Content-Type: application/json

{
  "title": "Team Meeting",
  "description": "Weekly team sync",
  "start_time": "2024-01-15T14:00:00Z",
  "end_time": "2024-01-15T15:00:00Z"
}
```

**Response (201 Created):**
```json
{
  "event": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "description": "Weekly team sync",
    "start_time": "2024-01-15T14:00:00Z",
    "end_time": "2024-01-15T15:00:00Z",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### List Events
```http
GET /events
```

**Response (200 OK):**
```json
{
  "events": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Team Meeting",
      "description": "Weekly team sync",
      "start_time": "2024-01-15T14:00:00Z",
      "end_time": "2024-01-15T15:00:00Z",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

### Get Event by ID
```http
GET /events/{id}
```

**Response (200 OK):**
```json
{
  "event": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "description": "Weekly team sync",
    "start_time": "2024-01-15T14:00:00Z",
    "end_time": "2024-01-15T15:00:00Z",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response (404 Not Found):**
```json
{
  "error": "Event not found"
}
```

## Quick Start

### 1. Start the service with Docker
```bash
docker-compose up -d
```

### 2. Test the API with curl commands

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Create an Event:**
```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "description": "Weekly team sync",
    "start_time": "2024-01-15T14:00:00Z",
    "end_time": "2024-01-15T15:00:00Z"
  }'
```

**Get All Events:**
```bash
curl http://localhost:8080/events
```

**Get Event by ID (use the ID from the create response):**
```bash
curl http://localhost:8080/events/b7c774bf-f0ce-4225-8739-51facbfc1ce7
```

**Test Validation - Empty Title:**
```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "",
    "start_time": "2024-01-15T14:00:00Z",
    "end_time": "2024-01-15T15:00:00Z"
  }'
```

**Test Validation - Invalid Time Range:**
```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Invalid Event",
    "start_time": "2024-01-15T15:00:00Z",
    "end_time": "2024-01-15T14:00:00Z"
  }'
```

## Database Schema

```sql
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(100) NOT NULL CHECK (length(title) > 0),
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_time_range CHECK (start_time < end_time)
);
```

## Validation Rules

- **Title**: Required, maximum 100 characters
- **Description**: Optional
- **Start Time**: Required, must be before end time
- **End Time**: Required, must be after start time

## Error Handling

The service returns appropriate HTTP status codes:

- `400 Bad Request` - Invalid input data
- `404 Not Found` - Event not found
- `500 Internal Server Error` - Server errors

## Configuration

The service uses these environment variables (already configured in docker-compose.yml):

.env:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=events_db
DB_SSL_MODE=disable

SERVER_PORT=8080
SERVER_HOST=0.0.0.0

LOG_LEVEL=info
```