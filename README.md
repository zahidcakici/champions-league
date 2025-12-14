# ⚽ Champions League Simulation

A full-stack football league simulation application that simulates a round-robin tournament with realistic match outcomes and championship predictions.

## 📋 Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running with Docker](#running-with-docker)
  - [Running Locally](#running-locally)
- [Application URLs](#application-urls)
- [API Documentation](#api-documentation)
- [Mathematical Models](#mathematical-models)
  - [Match Simulation Algorithm](#match-simulation-algorithm)
  - [Championship Prediction Algorithm](#championship-prediction-algorithm)
- [Project Structure](#project-structure)

## Overview

This application simulates a football league tournament where:

- Teams compete in a **round-robin format** (each team plays every other team twice - home and away)
- Match results are **simulated based on team power ratings** with home advantage
- **Championship predictions** are calculated dynamically as the league progresses
- Users can **manually edit match results** to explore different scenarios
- Full **CRUD operations** for teams before the tournament starts

## Tech Stack

| Layer      | Technology                                    |
| ---------- | --------------------------------------------- |
| Frontend   | Vue 3.5, Vite 6, Vue Router 4, Axios          |
| Backend    | Go 1.25, Fiber v2, GORM                       |
| Database   | PostgreSQL 18                                 |
| API Docs   | Swagger (swaggo/swag)                         |
| Testing    | Vitest (Frontend), Go testing (Backend)       |
| CI/CD      | GitHub Actions                                |
| Container  | Docker, Docker Compose                        |

## Getting Started

### Prerequisites

- Docker and Docker Compose (for containerized setup)
- Go 1.25+ (for local backend development)
- Node.js 24+ (for local frontend development)
- PostgreSQL 18+ (for local development without Docker)

### Running with Docker

The easiest way to run the application:

```bash
# Clone the repository
git clone https://github.com/zahidcakici/champions-league.git
cd champions-league

# Start all services
docker-compose up --build

# Or run in detached mode
docker-compose up -d --build
```

### Running Locally

**Backend:**

```bash
cd backend

# Set environment variables
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/champions_league?sslmode=disable"
export SERVER_PORT=8080

# Install dependencies and run
go mod download
go run cmd/server/main.go
```

**Frontend:**

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

## Application URLs

| Service           | URL                              | Description                        |
| ----------------- | -------------------------------- | ---------------------------------- |
| **Frontend**      | http://localhost:3000            | Vue.js web application             |
| **Backend API**   | http://localhost:8080/api        | REST API endpoints                 |
| **Swagger UI**    | http://localhost:8080/swagger/   | Interactive API documentation      |
| **Health Check**  | http://localhost:8080/health     | Backend health status              |

## API Documentation

Interactive API documentation is available via Swagger UI at:

```
http://localhost:8080/swagger/
```

### Main Endpoints

| Method | Endpoint                    | Description                          |
| ------ | --------------------------- | ------------------------------------ |
| GET    | `/api/teams`                | Get all teams                        |
| POST   | `/api/teams`                | Create a new team                    |
| DELETE | `/api/teams/:id`            | Delete a team                        |
| GET    | `/api/fixtures`             | Get all fixtures                     |
| GET    | `/api/fixtures/:week`       | Get fixtures for a specific week     |
| POST   | `/api/fixtures/generate`    | Generate fixtures for the tournament |
| GET    | `/api/simulation/state`     | Get current simulation state         |
| POST   | `/api/simulation/play-week` | Simulate next week's matches         |
| POST   | `/api/simulation/play-all`  | Simulate all remaining matches       |
| PUT    | `/api/simulation/match/:id` | Update a match result manually       |
| POST   | `/api/simulation/reset`     | Reset the entire simulation          |
| GET    | `/api/standings`            | Get current league standings         |
| GET    | `/api/predictions`          | Get championship predictions         |

## Mathematical Models

### Match Simulation Algorithm

The match simulation uses a **power-based probabilistic model** with home advantage:

#### 1. Effective Power Calculation

```
HomePower_effective = HomePower × HomeAdvantageFactor
AwayPower_effective = AwayPower
```

Where `HomeAdvantageFactor = 1.1` (10% boost for home team)

#### 2. Relative Power Calculation

```
TotalPower = HomePower_effective + AwayPower_effective

HomeRelativePower = HomePower_effective / TotalPower
AwayRelativePower = AwayPower_effective / TotalPower
```

#### 3. Expected Goals Calculation

```
HomeExpectedGoals = BaseExpectedGoals × 2 × HomeRelativePower
AwayExpectedGoals = BaseExpectedGoals × 2 × AwayRelativePower
```

Where `BaseExpectedGoals = 1.5`

**Key insight:** A team's expected goals depends on their power *relative* to their opponent. Playing against a stronger opponent reduces your expected goals.

#### 4. Goal Generation (Poisson Distribution)

Goals are generated using the **Poisson distribution** via inverse transform sampling. This models real football where:
- 0-1 goals are most common
- 2-3 goals are fairly common  
- 4+ goals are rare but possible

**Algorithm:**

```
L = e^(-ExpectedGoals)
k = 0
p = 1.0

while p > L:
    k++
    p = p × random()

Goals = min(k - 1, MaxGoalsPerTeam)
```

Where `MaxGoalsPerTeam = 7`

**Why Poisson?** Football goals are discrete, relatively rare events (~1.5 per team per match) that occur independently - exactly what Poisson distribution models.

**Probability Distribution for xG = 1.5:**

| Goals | Probability |
|-------|-------------|
| 0     | 22.3%       |
| 1     | 33.5%       |
| 2     | 25.1%       |
| 3     | 12.6%       |
| 4     | 4.7%        |
| 5+    | 1.8%        |

#### Example Calculation

**Team A (Power: 90) vs Team B (Power: 60) at Team A's home:**

```
HomePower_effective = 90 × 1.1 = 99
AwayPower_effective = 60
TotalPower = 99 + 60 = 159

HomeRelativePower = 99 / 159 = 0.623
AwayRelativePower = 60 / 159 = 0.377

HomeExpectedGoals = 1.5 × 2 × 0.623 = 1.87 goals
AwayExpectedGoals = 1.5 × 2 × 0.377 = 1.13 goals
```

---

### Championship Prediction Algorithm

The championship prediction uses a **weighted probability model** based on current standings and mathematical possibilities:

#### 1. Activation Condition

Predictions are only calculated when **3 or fewer weeks remain**:

```
RemainingWeeks = TotalWeeks - CurrentWeek

if RemainingWeeks > 3:
    All predictions = 0% (too early to predict)
```

#### 2. Mathematical Elimination Check

```
MaxRemainingPoints = RemainingWeeks × 3
PointsGap = LeaderPoints - TeamPoints

if PointsGap > MaxRemainingPoints:
    Team is mathematically eliminated (0% chance)
```

#### 3. Weight Calculation

For teams still in contention:

```
Weight = (CurrentPoints + 1) × 0.7^PointsGap

if GoalDifference > 0:
    Weight = Weight × 1.1  (10% bonus for positive GD)
```

The exponential decay factor (`0.7^PointsGap`) ensures:
- Teams at the top have significantly higher probabilities
- Each point gap roughly halves the championship probability
- Goal difference serves as a tiebreaker consideration

#### 4. Normalization

```
TotalWeight = Sum of all weights

For each team:
    Percentage = (TeamWeight / TotalWeight) × 100

// Adjust to ensure sum = 100%
```

#### Example Calculation

**With 2 weeks remaining (6 max points available):**

| Team | Points | Gap | Weight Calculation | Raw Weight |
|------|--------|-----|-------------------|------------|
| A    | 15     | 0   | (15+1) × 0.7^0 × 1.1 = 17.6 | 17.6 |
| B    | 13     | 2   | (13+1) × 0.7^2 × 1.1 = 7.55 | 7.55 |
| C    | 12     | 3   | (12+1) × 0.7^3 = 4.46 | 4.46 |
| D    | 8      | 7   | Eliminated (7 > 6) | 0 |

```
Total Weight = 17.6 + 7.55 + 4.46 = 29.61

Team A: 17.6 / 29.61 × 100 = 59%
Team B: 7.55 / 29.61 × 100 = 26%
Team C: 4.46 / 29.61 × 100 = 15%
Team D: 0%
```

## Project Structure

```
champions-league-case/
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── config/          # Configuration
│   │   ├── database/        # Database connection
│   │   ├── handlers/        # HTTP handlers & DTOs
│   │   │   └── docs/        # Swagger documentation
│   │   ├── models/          # Domain models
│   │   ├── repository/      # Data access layer
│   │   ├── routes/          # Route definitions
│   │   └── services/        # Business logic
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── views/           # Vue components
│   │   ├── api.js           # API client
│   │   ├── router.js        # Vue Router config
│   │   └── main.js          # App entry point
│   ├── __tests__/           # Vitest tests
│   ├── Dockerfile
│   └── package.json
├── .github/workflows/       # CI/CD pipelines
├── docker-compose.yml
└── README.md
```

## License

MIT License
