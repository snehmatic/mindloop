# Mindloop API v1 Documentation

## Overview

The Mindloop API provides a RESTful interface for managing productivity tracking features including habits, intents, focus sessions, journal entries, and summaries.

**Base URL**: `http://localhost:8080/api/v1`

## Authentication

Currently, the API does not require authentication. All endpoints are publicly accessible.

## Response Format

All API responses follow a standard format:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... }
}
```

Error responses:

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

## Endpoints

### Health & Info

#### GET `/`
Get API information and available features.

**Response:**
```json
{
  "success": true,
  "message": "Welcome to Mindloop API!",
  "data": {
    "version": "1.0.0",
    "features": ["habits", "intents", "focus", "journal", "summary"]
  }
}
```

#### GET `/healthz`
Check API health status.

**Response:**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "healthy"
  }
}
```

### Habits

#### POST `/habits`
Create a new habit.

**Request Body:**
```json
{
  "title": "Exercise",
  "description": "Daily workout routine",
  "target_count": 1,
  "interval": "daily"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Habit created successfully",
  "data": {
    "id": 1,
    "title": "Exercise",
    "description": "Daily workout routine",
    "target_count": 1,
    "interval": "daily",
    "created_at": "2025-08-16T22:04:49Z"
  }
}
```

#### GET `/habits`
List all habits.

**Response:**
```json
{
  "success": true,
  "message": "Habits retrieved successfully",
  "data": [
    {
      "id": 1,
      "title": "Exercise",
      "description": "Daily workout routine",
      "target_count": 1,
      "interval": "daily",
      "created_at": "2025-08-16T22:04:49Z"
    }
  ]
}
```

#### GET `/habits/{id}`
Get a specific habit by ID.

**Response:**
```json
{
  "success": true,
  "message": "Habit retrieved successfully",
  "data": {
    "id": 1,
    "title": "Exercise",
    "description": "Daily workout routine",
    "target_count": 1,
    "interval": "daily",
    "created_at": "2025-08-16T22:04:49Z"
  }
}
```

#### DELETE `/habits/{id}`
Delete a habit.

**Response:**
```json
{
  "success": true,
  "message": "Habit deleted successfully"
}
```

#### POST `/habits/{id}/log`
Log habit completion.

**Request Body:**
```json
{
  "actual_count": 1
}
```

**Response:**
```json
{
  "success": true,
  "message": "Habit logged successfully",
  "data": {
    "habit_id": 1,
    "actual_count": 1,
    "logged_at": "2025-08-16T22:04:49Z"
  }
}
```

### Intents

#### POST `/intents`
Create a new intent.

**Request Body:**
```json
{
  "name": "Complete project documentation"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Intent created successfully",
  "data": {
    "id": 1,
    "name": "Complete project documentation",
    "status": "active",
    "created_at": "2025-08-16T22:04:49Z"
  }
}
```

#### GET `/intents`
List all intents.

**Query Parameters:**
- `active=true` - Filter to show only active intents

**Response:**
```json
{
  "success": true,
  "message": "Intents retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Complete project documentation",
      "status": "active",
      "created_at": "2025-08-16T22:04:49Z"
    }
  ]
}
```

#### POST `/intents/{id}/end`
End an intent.

**Response:**
```json
{
  "success": true,
  "message": "Intent ended successfully",
  "data": {
    "id": 1,
    "name": "Complete project documentation",
    "status": "done",
    "ended_at": "2025-08-16T22:04:49Z"
  }
}
```

#### DELETE `/intents/{id}`
Delete an intent.

**Response:**
```json
{
  "success": true,
  "message": "Intent deleted successfully"
}
```

### Focus Sessions

#### POST `/focus`
Create a new focus session.

**Request Body:**
```json
{
  "title": "Complete project documentation"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Focus session created successfully",
  "data": {
    "id": 1,
    "title": "Complete project documentation",
    "status": "active",
    "created_at": "2025-08-16T22:04:49Z"
  }
}
```

#### GET `/focus`
List all focus sessions.

**Query Parameters:**
- `active=true` - Filter to show only active sessions

**Response:**
```json
{
  "success": true,
  "message": "Focus sessions retrieved successfully",
  "data": [
    {
      "id": 1,
      "title": "Complete project documentation",
      "status": "active",
      "duration": 30.5,
      "rating": -1,
      "created_at": "2025-08-16T22:04:49Z"
    }
  ]
}
```

#### POST `/focus/{id}/end`
End a focus session.

**Response:**
```json
{
  "success": true,
  "message": "Focus session ended successfully",
  "data": {
    "id": 1,
    "title": "Complete project documentation",
    "status": "ended",
    "duration": 30.5,
    "end_time": "2025-08-16T22:35:00Z"
  }
}
```

#### POST `/focus/{id}/pause`
Pause a focus session.

**Response:**
```json
{
  "success": true,
  "message": "Focus session paused successfully",
  "data": {
    "id": 1,
    "title": "Complete project documentation",
    "status": "paused"
  }
}
```

#### POST `/focus/{id}/resume`
Resume a paused focus session.

**Response:**
```json
{
  "success": true,
  "message": "Focus session resumed successfully",
  "data": {
    "id": 1,
    "title": "Complete project documentation",
    "status": "active"
  }
}
```

#### POST `/focus/{id}/rate`
Rate a completed focus session.

**Request Body:**
```json
{
  "rating": 8
}
```

**Response:**
```json
{
  "success": true,
  "message": "Focus session rated successfully",
  "data": {
    "id": 1,
    "title": "Complete project documentation",
    "status": "ended",
    "rating": 8
  }
}
```

#### DELETE `/focus/{id}`
Delete a focus session.

**Response:**
```json
{
  "success": true,
  "message": "Focus session deleted successfully"
}
```

### Journal

#### POST `/journal`
Create a new journal entry.

**Request Body:**
```json
{
  "title": "Today's thoughts",
  "content": "Had a great day working on the project!",
  "mood": "happy"
}
```

**Valid moods:** `happy`, `sad`, `neutral`, `angry`, `excited`

**Response:**
```json
{
  "success": true,
  "message": "Journal entry created successfully",
  "data": {
    "id": 1,
    "title": "Today's thoughts",
    "content": "Had a great day working on the project!",
    "mood": "happy",
    "created_at": "2025-08-16T22:04:49Z"
  }
}
```

#### GET `/journal`
List all journal entries.

**Response:**
```json
{
  "success": true,
  "message": "Journal entries retrieved successfully",
  "data": [
    {
      "id": 1,
      "title": "Today's thoughts",
      "content": "Had a great day working on the project!",
      "mood": "happy",
      "created_at": "2025-08-16T22:04:49Z"
    }
  ]
}
```

#### GET `/journal/{id}`
Get a specific journal entry.

**Response:**
```json
{
  "success": true,
  "message": "Journal entry retrieved successfully",
  "data": {
    "id": 1,
    "title": "Today's thoughts",
    "content": "Had a great day working on the project!",
    "mood": "happy",
    "created_at": "2025-08-16T22:04:49Z"
  }
}
```

#### PUT `/journal/{id}`
Update a journal entry.

**Request Body:**
```json
{
  "title": "Updated title",
  "content": "Updated content",
  "mood": "excited"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Journal entry updated successfully",
  "data": {
    "id": 1,
    "title": "Updated title",
    "content": "Updated content",
    "mood": "excited",
    "created_at": "2025-08-16T22:04:49Z"
  }
}
```

#### DELETE `/journal/{id}`
Delete a journal entry.

**Response:**
```json
{
  "success": true,
  "message": "Journal entry deleted successfully"
}
```

### Summaries

#### GET `/summary/daily`
Get daily summary (last 24 hours).

**Response:**
```json
{
  "success": true,
  "message": "Daily summary generated successfully",
  "data": {
    "date_range": "15-Aug-2025 to 16-Aug-2025",
    "focus": {
      "total_sessions": 2,
      "total_duration": "45min",
      "longest_session": "30min"
    },
    "habits": [
      {
        "habit_name": "Exercise",
        "completion_rate": 0.8,
        "logs_tracked": 5,
        "logs_completed": 4
      }
    ],
    "intents": [
      {
        "intent_name": "Complete project",
        "status": "done"
      }
    ]
  }
}
```

#### GET `/summary/weekly`
Get weekly summary (last 7 days).

#### GET `/summary/monthly`
Get monthly summary (last 30 days).

#### GET `/summary/yearly`
Get yearly summary (last 365 days).

#### POST `/summary/custom`
Get custom date range summary.

**Request Body:**
```json
{
  "start_date": "2025-08-01",
  "end_date": "2025-08-31"
}
```

**Response:** Same format as other summary endpoints.

## Error Codes

- `400 Bad Request` - Invalid request data or parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## CORS

The API supports CORS and allows requests from any origin with the following methods:
- GET, POST, PUT, DELETE, OPTIONS

## Examples

### Using curl

```bash
# Create a habit
curl -X POST http://localhost:8080/api/v1/habits \
  -H "Content-Type: application/json" \
  -d '{"title": "Exercise", "description": "Daily workout", "target_count": 1, "interval": "daily"}'

# Start a focus session
curl -X POST http://localhost:8080/api/v1/focus \
  -H "Content-Type: application/json" \
  -d '{"title": "Complete project documentation"}'

# End a focus session
curl -X POST http://localhost:8080/api/v1/focus/1/end

# Rate a focus session
curl -X POST http://localhost:8080/api/v1/focus/1/rate \
  -H "Content-Type: application/json" \
  -d '{"rating": 8}'

# Create a journal entry
curl -X POST http://localhost:8080/api/v1/journal \
  -H "Content-Type: application/json" \
  -d '{"title": "Today", "content": "Great day!", "mood": "happy"}'

# Get daily summary
curl http://localhost:8080/api/v1/summary/daily
```

### Using JavaScript/Fetch

```javascript
// Create a habit
const response = await fetch('http://localhost:8080/api/v1/habits', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: 'Exercise',
    description: 'Daily workout',
    target_count: 1,
    interval: 'daily'
  })
});

const result = await response.json();
console.log(result);
```

## Rate Limiting

Currently, there are no rate limits implemented. Consider implementing rate limiting for production use.

## Versioning

This is API version 1. Future versions will be available at `/api/v2`, `/api/v3`, etc.
