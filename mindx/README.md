# Student Dropout Risk Evaluation Service

This is a Go backend service that evaluates student dropout risk based on attendance, assignment completion, and communication response data.

## Features

- Reads student data from a JSON file
- Evaluates dropout risk using configurable thresholds
- Stores data in PostgreSQL database
- Provides RESTful API endpoints
- Containerized with Docker

## Architecture

The service follows idiomatic Go practices with a clear package structure:

- `config`: Configuration management
- `database`: Database connection and migrations
- `models`: Data models
- `services`: Business logic
- `handlers`: HTTP handlers
- `router`: API routing

## Detailed Component Documentation

### Config Package

The configuration package manages application settings through environment variables with sensible defaults. It uses a structured approach to configuration management:

```go
// Example configuration structure
type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    Risk     RiskConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
}
```

Configuration is loaded at application startup and validated to ensure all required settings are present.

### Database Package

The database package handles:

1. Connection establishment to PostgreSQL
2. Database migrations
3. Connection pooling for optimal performance
4. Error handling and reconnection logic

It provides a clean API for other packages to interact with the database without exposing implementation details.

### Models Package

The models package defines the core data structures used throughout the application:

- `Student`: Represents a student with their personal information and risk assessment
- `AttendanceRecord`: Tracks student attendance data
- `AssignmentRecord`: Tracks assignment submission data
- `ContactRecord`: Tracks communication attempts and responses

Each model includes validation logic and database mapping.

### Services Package

The services package contains the core business logic:

#### Risk Service

The risk evaluation service implements the algorithm for calculating dropout risk:

1. Processes raw student data
2. Applies configurable thresholds to determine risk factors
3. Calculates an overall risk score and risk level
4. Generates detailed notes explaining the risk assessment

#### Student Service

The student service manages student data operations:

1. Retrieval of student records from the database
2. Filtering and sorting of student records
3. Creation and updating of student records
4. Data validation and error handling

### Handlers Package

The handlers package implements the HTTP request handlers for the API endpoints:

- `EvaluateHandler`: Processes evaluation requests
- `StudentsHandler`: Manages student data retrieval with filtering and sorting
- Error handling middleware for consistent API responses

### Router Package

The router package defines the API routes and connects them to the appropriate handlers:

```go
// Example router setup
func SetupRouter(handlers *handlers.Handlers) *gin.Engine {
    router := gin.Default()
    
    // CORS middleware
    router.Use(cors.Default())
    
    // API endpoints
    router.POST("/evaluate", handlers.EvaluateHandler)
    router.GET("/students", handlers.GetStudentsHandler)
    
    return router
}
```

## API Endpoints

### POST /evaluate

Parses the JSON file, evaluates student dropout risk, and stores results in the database.

**Request**: No request body required

**Response**:
```json
{
  "message": "Evaluation completed successfully",
  "students_processed": 42,
  "timestamp": "2023-05-15T14:30:45Z"
}
```

**Status Codes**:
- `200 OK`: Successful evaluation
- `500 Internal Server Error`: Server error during evaluation

### GET /students

Lists all students with their dropout risk evaluations.

**Query Parameters**:
- `risk_level` (optional): Filter by risk level (LOW, MEDIUM, HIGH)
- `sort_by` (optional): Sort by different criteria
  - `risk_level`: Sort by risk level (HIGH to LOW)
  - `risk_level_asc`: Sort by risk level (LOW to HIGH)
  - `score`: Sort by risk score (HIGH to LOW)
  - `score_asc`: Sort by risk score (LOW to HIGH)

**Response**:
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "student_id": "ST12345",
    "student_name": "John Doe",
    "attendance": [
      {
        "date": "2023-05-01",
        "status": "PRESENT"
      }
    ],
    "assignments": [
      {
        "date": "2023-05-02",
        "name": "Math Assignment 1",
        "submitted": true
      }
    ],
    "contacts": [
      {
        "date": "2023-05-03",
        "status": "RESPONDED"
      }
    ],
    "dropout_score": 1,
    "dropout_risk_level": "LOW",
    "dropout_note": "Low risk due to good attendance and assignment completion.",
    "created_at": 1683648000,
    "updated_at": 1683648000
  }
]
```

**Status Codes**:
- `200 OK`: Successful retrieval
- `400 Bad Request`: Invalid query parameters
- `500 Internal Server Error`: Server error during retrieval

## Risk Evaluation Logic

For each student:
1. Calculate attendance rate = attended_days / total_days
2. Flag attendance as risky if attendance_rate < 75%
3. Calculate assignment completion rate = submitted_assignments / total_assignments
4. Flag assignment as risky if completion_rate < 50%
5. Count contact failures and flag as risky if failed_attempts >= 2
6. Score = total number of risky flags (0 to 3)
7. Risk level mapping:
   - 0-1: LOW
   - 2: MEDIUM
   - 3: HIGH

## Running the Service

### Prerequisites

- Docker
- Docker Compose

### Starting the Service

```bash
docker-compose up -d
```

This will start both the PostgreSQL database and the application server.

### API Usage

#### Evaluate Student Risk

```bash
curl -X POST http://localhost:8080/evaluate
```

This will read the data.json file, evaluate risk, and store results in the database.

#### List Students with Risk Evaluations

```bash
curl -X GET http://localhost:8080/students
```

## Configuration

Configuration is managed through environment variables:

- Database settings:
  - `DB_HOST`: PostgreSQL host (default: localhost)
  - `DB_PORT`: PostgreSQL port (default: 5432)
  - `DB_USER`: PostgreSQL user (default: postgres)
  - `DB_PASSWORD`: PostgreSQL password (default: postgres)
  - `DB_NAME`: PostgreSQL database name (default: studentrisk)
  - `DB_SSLMODE`: PostgreSQL SSL mode (default: disable)

- Server settings:
  - `SERVER_ADDRESS`: Server address and port (default: :8080)

- Risk evaluation settings:
  - `RISK_ATTENDANCE_THRESHOLD`: Attendance threshold percentage (default: 75.0)
  - `RISK_ASSIGNMENT_THRESHOLD`: Assignment completion threshold percentage (default: 50.0)
  - `RISK_CONTACT_THRESHOLD`: Contact failure threshold count (default: 2)
  - `RISK_LOW_THRESHOLD`: Score threshold for low risk level (default: 0)
  - `RISK_MEDIUM_THRESHOLD`: Score threshold for medium risk level (default: 2)
  - `RISK_HIGH_THRESHOLD`: Score threshold for high risk level (default: 3)

## Development

### Building from Source

```bash
go build -o app ./main.go
```

### Running Tests

```bash
go test ./...
```

### Adding New Features

When adding new features:

1. Create appropriate models in the `models` package
2. Implement business logic in the `services` package
3. Add HTTP handlers in the `handlers` package
4. Update routes in the `router` package
5. Update configuration in the `config` package if needed

Follow Go best practices for error handling, documentation, and testing.