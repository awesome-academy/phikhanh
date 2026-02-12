# Public Service Management System

Hệ thống quản lý dịch vụ công với kiến trúc hybrid: REST API cho user và SSR cho admin.

## Công nghệ sử dụng
- **Go** 1.21+
- **Gin Framework** - Web framework
- **GORM** - ORM
- **PostgreSQL** - Database
- **Swagger** - API Documentation

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

## API Documentation

### Cho Frontend Developer:

**Xem API Documentation (không cần run app):**
- File Swagger: `/docs/swagger.yaml`
- Import vào Swagger Editor: https://editor.swagger.io/
- Hoặc dùng Postman: Import file `docs/swagger.yaml`

**Online Swagger UI (khi app đang chạy):**
- URL: http://localhost:8080/docs/index.html
- Swagger JSON: http://localhost:8080/docs/doc.json

### API Endpoints

#### Auth
- `POST /api/v1/auth/register` - Đăng ký tài khoản
- `POST /api/v1/auth/login` - Đăng nhập
- `POST /api/v1/auth/logout` - Đăng xuất (cần token)

#### Profile
- `GET /api/v1/profile` - Lấy thông tin profile (cần token)
- `PUT /api/v1/profile` - Cập nhật profile (cần token)

#### Health
- `GET /api/v1/health` - Health check

### Response Format

**Success Response:**
```json
{
  "status": 200,
  "message": "Success message",
  "data": { }
}
```

**Error Response:**
```json
{
  "status": 400,
  "message": "Error message"
}
```

### Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request successful |
| 201 | Created - Resource created |
| 400 | Bad Request - Validation error |
| 401 | Unauthorized - Authentication failed |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error |

### Authentication

Sử dụng JWT Bearer Token trong header:
```
Authorization: Bearer <your_token_here>
```

**Example:**
```bash
curl -H "Authorization: Bearer eyJhbGc..." http://localhost:8080/api/v1/profile
```

## Run
```bash
go run main.go
```

Server runs on http://localhost:8080
