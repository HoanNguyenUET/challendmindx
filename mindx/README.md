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

## API Endpoints

- `POST /evaluate`: Parse JSON file, evaluate student dropout risk, and store results in database
- `GET /students`: List all students with their dropout risk evaluations

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