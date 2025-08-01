# Student Dropout Risk Evaluation System

A full-stack application for evaluating and monitoring student dropout risk based on attendance, assignment completion, and communication data.

## Project Overview

This system helps educational institutions identify students at risk of dropping out by analyzing key behavioral indicators:
- Attendance patterns
- Assignment completion rates
- Communication responsiveness

The application provides a visual dashboard for educators to monitor student risk levels and take proactive intervention measures.

## Tech Stack

### Backend (Go)
- RESTful API built with Go
- PostgreSQL database for data storage
- Docker containerization

### Frontend (React)
- React with TypeScript
- Material-UI component library
- Responsive design for desktop and mobile

## Features

- **Risk Evaluation Algorithm**: Analyzes student data using configurable thresholds
- **Filtering and Sorting**: View students by risk level or sort by various metrics
- **Real-time Evaluation**: Trigger new evaluations with updated data
- **Visual Dashboard**: Clear presentation of student risk levels
- **Containerized Deployment**: Easy setup with Docker Compose

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Node.js and npm (for frontend development)
- Go (for backend development)

### Running the Application

1. Clone this repository:
   ```bash
   git clone <repository-url>
   cd challendmindx
   ```

2. Start the application using Docker Compose:
   ```bash
   cd mindx
   docker-compose up -d
   ```

3. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## API Endpoints

- `POST /evaluate`: Evaluate student dropout risk using data from data.json
- `GET /students`: List all students with their risk evaluations
  - Query parameters:
    - `risk_level`: Filter by risk level (LOW, MEDIUM, HIGH)
    - `sort_by`: Sort by different criteria (risk_level, score)

## Project Structure

```
mindx/
├── config/            # Backend configuration
├── database/          # Database connection and migrations
├── frontend/          # React frontend application
│   ├── public/        # Static assets
│   └── src/
│       ├── components/  # React components
│       ├── services/    # API services
│       └── types/       # TypeScript type definitions
├── handlers/          # HTTP request handlers
├── models/            # Data models
├── router/            # API routing
└── services/          # Business logic services
```

## Development

For detailed information about the backend and frontend components, please refer to their respective README files:
- [Backend README](mindx/README.md)
- [Frontend README](mindx/frontend/README.md)

## License

[Your license information here]