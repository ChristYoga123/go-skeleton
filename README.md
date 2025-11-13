# Golang Skeleton

Skeleton project untuk aplikasi Go dengan Fiber framework, GORM, dan support multiple database (MySQL, PostgreSQL, SQLite).

## 📋 Daftar Isi

- [Fitur](#-fitur)
- [Struktur Project](#-struktur-project)
- [Memulai](#-memulai)
- [Environment Variables](#-environment-variables)
- [Setup Database](#-setup-database)
- [Migration](#-migration)
- [Best Practices](#-best-practices)
- [Standar Response API](#-standar-response-api)
- [Autentikasi](#-autentikasi)
- [Password Hashing](#-password-hashing)
- [Development](#️-development)
- [Dependencies](#-dependencies)

## ✨ Fitur

- 🚀 **Fiber Framework** - Fast HTTP web framework
- 🗄️ **Multi-Database Support** - MySQL, PostgreSQL, SQLite
- 🔄 **Auto Database Creation** - Otomatis membuat database jika belum ada (seperti Laravel)
- 🔐 **JWT Authentication** - Generate dan verifikasi JWT token
- 🔒 **Password Hashing** - Bcrypt password hashing
- 📦 **Clean Architecture** - Pattern Repository, Service, Controller
- 📝 **Auto Migration** - GORM auto migration
- 🎯 **Standardized Response** - Base response untuk konsistensi API

## 📁 Struktur Project

```
golang-skeleton/
├── app/
│   ├── commands/             # CLI commands
│   ├── controllers/           # HTTP handlers/controllers
│   ├── entities/
│   │   ├── dtos/
│   │   │   ├── requests/      # Request DTOs
│   │   │   └── responses/     # Response DTOs (BaseResponse)
│   │   └── models/            # GORM models
│   ├── repositories/          # Data access layer
│   └── services/              # Business logic layer
├── cmd/
│   ├── api/
│   │   └── main.go           # Entry point aplikasi
│   └── cli/
│       └── main.go           # Entry point CLI
├── configs/                   # File konfigurasi
│   ├── bcrypt.go             # Password hashing
│   ├── database.go            # Koneksi database
│   ├── env_config.go          # Environment variables
│   └── jwt.go                 # JWT utilities
├── database/
│   ├── auto_migrate.go        # Setup auto migration
│   └── database.sqlite        # File database SQLite (jika menggunakan SQLite)
├── routes/
│   └── api.go                 # API routes
└── tests/                     # File test
```

## 🚀 Memulai

### Prerequisites

- Go 1.25.4 atau lebih tinggi
- MySQL/PostgreSQL (opsional, jika tidak menggunakan SQLite)

### Instalasi

1. Clone repository
```bash
git clone <repository-url>
cd golang-skeleton
```

2. Install dependencies
```bash
go mod download
```

3. Copy file environment
```bash
cp .env.example .env
```

4. Konfigurasi environment variables (lihat [Environment Variables](#-environment-variables))

5. Jalankan aplikasi
```bash
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8000` (atau port yang ditentukan di `APP_PORT`)

## 🔧 Environment Variables

Buat file `.env` di root directory dengan variabel berikut:

### Konfigurasi Aplikasi

```env
APP_PORT=8000
```

### Konfigurasi Database

#### Opsi 1: Menggunakan SQLite (Default)

```env
DB_DRIVER=sqlite
DB_PATH=database/database.sqlite
```

#### Opsi 2: Menggunakan MySQL

```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_DATABASE=myapp_db
```

**Atau menggunakan DB_URL:**
```env
DB_DRIVER=mysql
DB_URL=user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

#### Opsi 3: Menggunakan PostgreSQL

```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_DATABASE=myapp_db
DB_SSLMODE=disable
```

**Atau menggunakan DB_URL:**
```env
DB_DRIVER=postgres
DB_URL=host=localhost user=postgres password=password dbname=myapp_db port=5432 sslmode=disable
```

### Konfigurasi JWT

```env
JWT_SECRET="your-secret-key-change-this-in-production"
```

**Gunakan CLI untuk generate JWT secret:**
```bash
go run cmd/cli/main.go generate-jwt-secret
```

## 🗄️ Setup Database

### Auto Database Creation

Aplikasi secara otomatis membuat database jika belum ada (untuk MySQL dan PostgreSQL), mirip dengan behavior Laravel.

**Cara kerjanya:**
1. Aplikasi connect ke database server (tanpa specify nama database)
2. Cek apakah database sudah ada
3. Buat database jika belum ada
4. Connect ke database yang sudah dibuat

**Catatan:** Pastikan database user memiliki privilege `CREATE DATABASE`.

### Manual Database Creation

Jika ingin membuat database secara manual:

**MySQL:**
```sql
CREATE DATABASE myapp_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

**PostgreSQL:**
```sql
CREATE DATABASE myapp_db;
```

## 🔄 Migration

### Menambah Model Baru

Untuk menambah model baru dan memasukkannya ke auto-migration:

1. **Buat file model baru** di `app/entities/models/`

Contoh: `app/entities/models/product.go`
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

2. **Tambahkan model ke auto-migration** di `database/auto_migrate.go`

```go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Product{},  // Tambahkan model baru di sini
        // ... tambahkan model lainnya di sini
    )
}
```

3. **Restart aplikasi** - Migration akan berjalan otomatis saat startup

### Best Practices Migration

- Selalu tambahkan model baru ke function `AutoMigrate()`
- GORM akan otomatis:
  - Membuat tabel baru
  - Menambah kolom baru
  - Membuat index
  - **Catatan:** GORM tidak akan menghapus kolom atau tabel secara otomatis (untuk keamanan)

## 📝 Best Practices

### 1. Gunakan BaseResponse untuk Response API

Selalu gunakan `BaseResponse` untuk standarisasi response API:

```go
import "golang-skeleton/app/entities/dtos/responses"
import "net/http"

// Success response
return c.JSON(responses.NewSuccessResponse(
    "Data berhasil diambil",
    data,
    http.StatusOK,
))

// Error response
return c.JSON(responses.NewErrorResponse(
    "Validasi gagal",
    nil,
    http.StatusBadRequest,
))
```

**Format Response:**
```json
{
  "code": 200,
  "message": "Pesan sukses",
  "data": { ... }
}
```

### 2. Pattern Arsitektur Project

Ikuti pattern clean architecture:

```
Controller → Service → Repository → Database
```

- **Controller**: Handle HTTP request/response
- **Service**: Business logic
- **Repository**: Data access layer
- **Model**: Database models

### 3. Error Handling

- Selalu return HTTP status code yang tepat
- Gunakan `BaseResponse` untuk error response
- Log error dengan tepat

### 4. Environment Variables

- Jangan commit file `.env` (sudah ada di `.gitignore`)
- Gunakan `GetEnv()` untuk variabel yang wajib (return error jika tidak diset)
- Gunakan `GetEnvOrDefault()` untuk variabel opsional dengan default value

## 📊 Standar Response API

Semua response API mengikuti format standar ini:

### Success Response

```json
{
  "code": 200,
  "message": "Operasi berhasil",
  "data": {
    // Data response
  }
}
```

### Error Response

```json
{
  "code": 400,
  "message": "Pesan error",
  "data": null
}
```

### Status Code yang Tersedia

- `200` - OK (Success)
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

Gunakan HTTP status constants:
```go
import "net/http"

responses.NewSuccessResponse("Success", data, http.StatusOK)
responses.NewErrorResponse("Not found", nil, http.StatusNotFound)
```

## 🔐 Autentikasi

### Generate JWT Token

```go
import "golang-skeleton/configs"

token, err := configs.GenerateJWT(userID, email)
```

### Verifikasi JWT Token

```go
claims, err := configs.VerifyJWT(tokenString)
if err != nil {
    // Handle error
}
userID := claims.UserID
email := claims.Email
```

### Ambil User ID dari Token

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
    // Password tidak cocok
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
- Menambahkan `JWT_SECRET` ke file `.env` jika belum ada
- Update value `JWT_SECRET` jika sudah ada

### Menjalankan Aplikasi

```bash
go run cmd/api/main.go
```

### Build Aplikasi

```bash
go build -o bin/api cmd/api/main.go
```

### Menjalankan Test

```bash
go test ./...
```

## 📦 Dependencies

Dependencies utama:
- **Fiber v2** - Web framework
- **GORM v1** - ORM
- **JWT v5** - JWT authentication
- **Bcrypt** - Password hashing
- **Godotenv** - Environment variables
- **Cobra** - CLI framework
