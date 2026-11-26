package main

import (
	"basic-app/config"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHostName,
		cfg.DBPort,
		cfg.DBName,
		cfg.SSLMode,
	)

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("Migration path error: %v", err)
	}

	m, err := migrate.New(
		"file://"+filepath.ToSlash(migrationsPath),
		dsn,
	)
	if err != nil {
		log.Fatalf("Migration initialization error: %v", err)
	}

	defer func() {
		sourceErr, databaseErr := m.Close()

		if sourceErr != nil {
			log.Printf("Migration source close error: %v", sourceErr)
		}

		if databaseErr != nil {
			log.Printf("Migration database close error: %v", databaseErr)
		}
	}()

	command := "up"

	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		runUp(m)

	case "down":
		runDown(m)

	case "version":
		showVersion(m)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run ./cmd/migrate force <version>")
		}

		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid migration version: %v", err)
		}

		if err := m.Force(version); err != nil {
			log.Fatalf("Migration force failed: %v", err)
		}

		log.Printf("Migration version forced to %d", version)
	default:
		log.Fatalf(
			"Unknown migration command: %s. Use: up, down, version",
			command,
		)
	}
}

func runUp(m *migrate.Migrate) {
	err := m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migration failed: %v", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("No new migrations to apply")
		return
	}

	log.Println("Database migrations completed successfully")
}

func runDown(m *migrate.Migrate) {
	err := m.Steps(-1)

	if err != nil {
		log.Fatalf("Migration rollback failed: %v", err)
	}

	log.Println("Last migration rolled back successfully")
}

func showVersion(m *migrate.Migrate) {
	version, dirty, err := m.Version()

	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Println("No migrations have been applied")
			return
		}

		log.Fatalf("Could not get migration version: %v", err)
	}

	log.Printf("Migration version: %d, dirty: %t", version, dirty)
}
