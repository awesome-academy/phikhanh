# Public Service Management System

A hybrid system with REST API for users and SSR for admin dashboard.

## Tech Stack
- Go 1.21
- Gin Framework
- GORM
- PostgreSQL

## Project Structure
```
phikhanh/
├── config/           # Database & Environment configuration
├── controllers/
│   ├── user/         # API handlers (JSON)
│   └── admin/        # SSR handlers (HTML)
├── dto/
│   ├── user/         # API request/response structs
│   └── admin/        # Admin form structs
├── middlewares/      # CORS, Auth middlewares
├── models/           # GORM models (shared)
├── repositories/
│   ├── user/         # User data access
│   └── admin/        # Admin data access
├── routes/           # Route definitions
├── services/
│   ├── user/         # User business logic
│   └── admin/        # Admin business logic
├── utils/            # Helper functions
├── templates/        # HTML templates for admin
│   └── admin/
├── assets/           # CSS/JS/Images
├── .env
├── go.mod
└── main.go
```

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Configure `.env`:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. Run migrations:
```bash
go run main.go
```

## API Endpoints

### User API (JSON)
- `GET /api/v1/health` - Health check

### Admin (SSR)
- `GET /admin/dashboard` - Admin dashboard

## Run
```bash
go run main.go
```

Server runs on http://localhost:8080
