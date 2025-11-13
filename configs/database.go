package configs

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql" // MySQL driver untuk database/sql
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver untuk database/sql
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ConnectDatabase membuat koneksi database dinamis berdasarkan DB_TYPE
// Support: mysql, postgres, sqlite
// Auto-create database jika belum ada (untuk MySQL dan PostgreSQL)
func ConnectDatabase() (*gorm.DB, error) {
	// Ambil database type dari env (default: sqlite)
	dbType := GetEnvOrDefault("DB_DRIVER", "sqlite")
	dbType = strings.ToLower(dbType)

	var db *gorm.DB
	var err error

	switch dbType {
	case "mysql":
		// Auto-create database jika belum ada
		if err := ensureMySQLDatabase(); err != nil {
			return nil, err
		}
		db, err = connectMySQL()
	case "postgres", "pgql", "postgresql":
		// Auto-create database jika belum ada
		if err := ensurePostgresDatabase(); err != nil {
			return nil, err
		}
		db, err = connectPostgres()
	case "sqlite":
		db, err = connectSQLite()
	default:
		return nil, fmt.Errorf("unsupported database type: %s. Supported types: mysql, postgres, sqlite", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s database: %w", dbType, err)
	}

	return db, nil
}

// connectMySQL membuat koneksi ke MySQL
func connectMySQL() (*gorm.DB, error) {
	// Cek apakah ada DB_URL atau DB_DSN
	dbURL, err := GetEnv("DB_URL")
	if err != nil {
		// Jika DB_URL tidak ada, coba build dari komponen terpisah
		dbURL, err = buildMySQLDSN()
		if err != nil {
			return nil, err
		}
	}

	return gorm.Open(mysql.Open(dbURL), &gorm.Config{})
}

// connectPostgres membuat koneksi ke PostgreSQL
func connectPostgres() (*gorm.DB, error) {
	// Cek apakah ada DB_URL atau DB_DSN
	dbURL, err := GetEnv("DB_URL")
	if err != nil {
		// Jika DB_URL tidak ada, coba build dari komponen terpisah
		dbURL, err = buildPostgresDSN()
		if err != nil {
			return nil, err
		}
	}

	return gorm.Open(postgres.Open(dbURL), &gorm.Config{})
}

// connectSQLite membuat koneksi ke SQLite
func connectSQLite() (*gorm.DB, error) {
	// Untuk SQLite, bisa pakai DB_PATH
	dbPath := GetEnvOrDefault("DB_PATH", "database/database.sqlite")
	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}

// buildMySQLDSN membangun MySQL DSN dari komponen terpisah
func buildMySQLDSN() (string, error) {
	host := GetEnvOrDefault("DB_HOST", "localhost")
	port := GetEnvOrDefault("DB_PORT", "3306")
	user, err := GetEnv("DB_USERNAME")
	if err != nil {
		return "", fmt.Errorf("DB_USERNAME is required for MySQL connection")
	}
	password, err := GetEnv("DB_PASSWORD")
	if err != nil {
		return "", fmt.Errorf("DB_PASSWORD is required for MySQL connection")
	}
	dbname, err := GetEnv("DB_DATABASE")
	if err != nil {
		return "", fmt.Errorf("DB_DATABASE is required for MySQL connection")
	}

	// Format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	return dsn, nil
}

// buildPostgresDSN membangun PostgreSQL DSN dari komponen terpisah
func buildPostgresDSN() (string, error) {
	host := GetEnvOrDefault("DB_HOST", "localhost")
	port := GetEnvOrDefault("DB_PORT", "5432")
	user, err := GetEnv("DB_USERNAME")
	if err != nil {
		return "", fmt.Errorf("DB_USERNAME is required for PostgreSQL connection")
	}
	password, err := GetEnv("DB_PASSWORD")
	if err != nil {
		return "", fmt.Errorf("DB_PASSWORD is required for PostgreSQL connection")
	}
	dbname, err := GetEnv("DB_DATABASE")
	if err != nil {
		return "", fmt.Errorf("DB_DATABASE is required for PostgreSQL connection")
	}
	sslmode := GetEnvOrDefault("DB_SSLMODE", "disable")

	// Format: host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Jakarta
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, dbname, port, sslmode)

	return dsn, nil
}

// ensureMySQLDatabase membuat database MySQL jika belum ada
func ensureMySQLDatabase() error {
	host := GetEnvOrDefault("DB_HOST", "localhost")
	port := GetEnvOrDefault("DB_PORT", "3306")
	user, err := GetEnv("DB_USERNAME")
	if err != nil {
		return fmt.Errorf("DB_USERNAME is required for MySQL connection")
	}
	password, err := GetEnv("DB_PASSWORD")
	if err != nil {
		return fmt.Errorf("DB_PASSWORD is required for MySQL connection")
	}
	dbname, err := GetEnv("DB_DATABASE")
	if err != nil {
		return fmt.Errorf("DB_DATABASE is required for MySQL connection")
	}

	// Connect ke MySQL server tanpa specify database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port)

	// Open connection untuk create database
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w", err)
	}
	defer sqlDB.Close()

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL server: %w", err)
	}

	// Check apakah database sudah ada
	var exists int
	err = sqlDB.QueryRow("SELECT 1 FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", dbname).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	// Jika database belum ada, buat database
	if err == sql.ErrNoRows {
		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbname))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		fmt.Printf("Database '%s' created successfully\n", dbname)
	}

	return nil
}

// ensurePostgresDatabase membuat database PostgreSQL jika belum ada
func ensurePostgresDatabase() error {
	host := GetEnvOrDefault("DB_HOST", "localhost")
	port := GetEnvOrDefault("DB_PORT", "5432")
	user, err := GetEnv("DB_USERNAME")
	if err != nil {
		return fmt.Errorf("DB_USERNAME is required for PostgreSQL connection")
	}
	password, err := GetEnv("DB_PASSWORD")
	if err != nil {
		return fmt.Errorf("DB_PASSWORD is required for PostgreSQL connection")
	}
	dbname, err := GetEnv("DB_DATABASE")
	if err != nil {
		return fmt.Errorf("DB_DATABASE is required for PostgreSQL connection")
	}
	sslmode := GetEnvOrDefault("DB_SSLMODE", "disable")

	// Connect ke PostgreSQL default database (postgres) untuk create database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=%s",
		host, user, password, port, sslmode)

	// Open connection untuk create database
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL server: %w", err)
	}
	defer sqlDB.Close()

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL server: %w", err)
	}

	// Check apakah database sudah ada
	var exists int
	err = sqlDB.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", dbname).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	// Jika database belum ada, buat database
	if err == sql.ErrNoRows {
		// Set connection limit untuk database baru
		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		fmt.Printf("Database '%s' created successfully\n", dbname)
	}

	return nil
}
