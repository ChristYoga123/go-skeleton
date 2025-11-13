# Golang Skeleton

Skeleton project untuk aplikasi Go dengan Fiber framework, GORM, dan support multiple database (MySQL, PostgreSQL, SQLite).

## 📋 Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Database Setup](#database-setup)
- [Migration](#migration)
- [Best Practices](#best-practices)
- [API Response Standard](#api-response-standard)

## ✨ Features

- 🚀 **Fiber Framework** - Fast HTTP web framework
- 🗄️ **Multi-Database Support** - MySQL, PostgreSQL, SQLite
- 🔄 **Auto Database Creation** - Automatically creates database if not exists (like Laravel)
- 🔐 **JWT Authentication** - JWT token generation and verification
- 🔒 **Password Hashing** - Bcrypt password hashing
- 📦 **Clean Architecture** - Repository, Service, Controller pattern
- 📝 **Auto Migration** - GORM auto migration
- 🎯 **Standardized Response** - Base response for API consistency

## 📁 Project Structure

```
golang-skeleton/
├── app/
│   ├── controllers/          # HTTP handlers/controllers
│   ├── entities/
│   │   ├── dtos/
│   │   │   ├── requests/     # Request DTOs
│   │   │   └── responses/    # Response DTOs (BaseResponse)
│   │   └── models/           # GORM models
│   ├── repositories/         # Data access layer
│   └── services/             # Business logic layer
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── configs/                  # Configuration files
│   ├── bcrypt.go            # Password hashing
│   ├── database.go          # Database connection
│   ├── env_config.go        # Environment variables
│   └── jwt.go               # JWT utilities
├── database/
│   ├── auto_migrate.go      # Auto migration setup
│   └── database.sqlite      # SQLite database file (if using SQLite)
├── routes/
│   └── api.go               # API routes
└── tests/                   # Test files
```

## 🚀 Getting Started

### Prerequisites

- Go 1.25.4 or higher
- MySQL/PostgreSQL (optional, if not using SQLite)

### Installation

1. Clone the repository
```bash
git clone <repository-url>
cd golang-skeleton
```

2. Install dependencies
```bash
go mod download
```

3. Copy environment file
```bash
cp .env.example .env
```

4. Configure environment variables (see [Environment Variables](#environment-variables))

5. Run the application
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8000` (or port specified in `APP_PORT`)

## 🔧 Environment Variables

Create a `.env` file in the root directory with the following variables:

### Application Configuration

```env
APP_PORT=8000
```

### Database Configuration

#### Option 1: Using SQLite (Default)

```env
DB_DRIVER=sqlite
DB_PATH=database/database.sqlite
```

#### Option 2: Using MySQL

```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_DATABASE=myapp_db
```

**Or using DB_URL:**
```env
DB_DRIVER=mysql
DB_URL=user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

#### Option 3: Using PostgreSQL

```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_DATABASE=myapp_db
DB_SSLMODE=disable
```

**Or using DB_URL:**
```env
DB_DRIVER=postgres
DB_URL=host=localhost user=postgres password=password dbname=myapp_db port=5432 sslmode=disable
```

### JWT Configuration

```env
JWT_SECRET=your-secret-key-change-this-in-production
```

## 🗄️ Database Setup

### Auto Database Creation

The application automatically creates the database if it doesn't exist (for MySQL and PostgreSQL), similar to Laravel's behavior.

**How it works:**
1. Application connects to the database server (without specifying database name)
2. Checks if the database exists
3. Creates the database if it doesn't exist
4. Connects to the created database

**Note:** Make sure the database user has `CREATE DATABASE` privileges.

### Manual Database Creation

If you prefer to create the database manually:

**MySQL:**
```sql
CREATE DATABASE myapp_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

**PostgreSQL:**
```sql
CREATE DATABASE myapp_db;
```

## 🔄 Migration

### Adding New Models

To add a new model and include it in auto-migration:

1. **Create a new model file** in `app/entities/models/`

Example: `app/entities/models/product.go`
```go
package models

import "time"

type Product struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `json:"name"`
    Price     float64   `json:"price"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

2. **Add the model to auto-migration** in `database/auto_migrate.go`

```go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Product{},  // Add your new model here
        // ... tambahkan model lainnya di sini
    )
}
```

3. **Restart the application** - Migration will run automatically on startup

### Migration Best Practices

- Always add new models to `AutoMigrate()` function
- GORM will automatically:
  - Create new tables
  - Add new columns
  - Create indexes
  - **Note:** GORM won't delete columns or tables automatically (for safety)

## 📝 Best Practices

### 1. Use BaseResponse for API Responses

Always use `BaseResponse` for standardized API responses:

```go
import "golang-skeleton/app/entities/dtos/responses"
import "net/http"

// Success response
return c.JSON(responses.NewSuccessResponse(
    "Data retrieved successfully",
    data,
    http.StatusOK,
))

// Error response
return c.JSON(responses.NewErrorResponse(
    "Validation failed",
    nil,
    http.StatusBadRequest,
))
```

**Response Format:**
```json
{
  "code": 200,
  "message": "Success message",
  "data": { ... }
}
```

### 2. Project Architecture Pattern

Follow the clean architecture pattern:

```
Controller → Service → Repository → Database
```

- **Controller**: Handle HTTP requests/responses
- **Service**: Business logic
- **Repository**: Data access layer
- **Model**: Database models

### 3. Error Handling

- Always return proper HTTP status codes
- Use `BaseResponse` for error responses
- Log errors appropriately

### 4. Environment Variables

- Never commit `.env` file (already in `.gitignore`)
- Use `GetEnv()` for required variables (returns error if not set)
- Use `GetEnvOrDefault()` for optional variables with defaults

## 📊 API Response Standard

All API responses follow this standard format:

### Success Response

```json
{
  "code": 200,
  "message": "Operation successful",
  "data": {
    // Response data
  }
}
```

### Error Response

```json
{
  "code": 400,
  "message": "Error message",
  "data": null
}
```

### Available Status Codes

- `200` - OK (Success)
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

Use HTTP status constants:
```go
import "net/http"

responses.NewSuccessResponse("Success", data, http.StatusOK)
responses.NewErrorResponse("Not found", nil, http.StatusNotFound)
```

## 🔐 Authentication

### JWT Token Generation

```go
import "golang-skeleton/configs"

token, err := configs.GenerateJWT(userID, email)
```

### JWT Token Verification

```go
claims, err := configs.VerifyJWT(tokenString)
if err != nil {
    // Handle error
}
userID := claims.UserID
email := claims.Email
```

### Get User ID from Token

```go
userID, err := configs.GetIDFromToken(tokenString)
```

## 🔒 Password Hashing

### Hash Password

```go
import "golang-skeleton/configs"

hashedPassword, err := configs.HashPassword("plain_password")
```

### Compare Password

```go
err := configs.ComparePassword("plain_password", hashedPassword)
if err != nil {
    // Password doesn't match
}
```

## 🛠️ Development

### CLI Commands

Build CLI tool:
```bash
go build -o bin/cli cmd/cli/main.go
```

Generate JWT Secret:
```bash
./bin/cli generate-jwt-secret
# atau
go run cmd/cli/main.go generate-jwt-secret
```

Command ini akan:
- Generate secure random JWT secret (32 bytes, base64 encoded)
- Add `JWT_SECRET` to `.env` file if not exists
- Update `JWT_SECRET` value if already exists

### Running the Application

```bash
go run cmd/api/main.go
```

### Building the Application

```bash
go build -o bin/api cmd/api/main.go
```

### Running Tests

```bash
go test ./...
```

## 📦 Dependencies

Main dependencies:
- **Fiber v2** - Web framework
- **GORM v1** - ORM
- **JWT v5** - JWT authentication
- **Bcrypt** - Password hashing
- **Godotenv** - Environment variables

## 📄 License

[Your License Here]

## 🤝 Contributing

[Your Contributing Guidelines Here]

## 📞 Support

[Your Support Information Here]

