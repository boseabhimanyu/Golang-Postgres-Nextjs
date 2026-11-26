package database

import (
	"basic-app/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	// Prevents your app from freezing during startup.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHostName,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.SSLMode,
		cfg.DBTimezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres connect failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql database failed: %w", err)
	}

	// Ping check.
	if err := pingDatabase(ctx, sqlDB); err != nil {
		if closeErr := sqlDB.Close(); closeErr != nil {
			log.Printf("postgres disconnect after failed ping: %v", closeErr)
		}

		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	log.Printf("PostgreSQL connected (database: %s)", cfg.DBName)

	return db, nil
}

func pingDatabase(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}

func Disconnect(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get underlying sql database failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = ctx // database/sql Close does not accept a context.
	return sqlDB.Close()
}
