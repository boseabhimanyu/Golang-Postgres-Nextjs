package main

import (
	"basic-app/config"
	"basic-app/models"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}

	seedCfg, err := config.LoadSeedConfig()
	if err != nil {
		log.Fatalf("Seed Config Error: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Kolkata",
		cfg.DBHostName,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	ctx := context.Background()

	var count int64

	if err := db.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", seedCfg.AdminUsername).
		Count(&count).Error; err != nil {
		log.Fatalf("Checking admin user failed: %v", err)
	}

	if count > 0 {
		log.Println("Admin user already exists")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(seedCfg.AdminPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatalf("Password hashing failed: %v", err)
	}

	admin := models.User{
		ID:           uuid.New(),
		FirstName:    seedCfg.AdminFirstName,
		LastName:     seedCfg.AdminLastName,
		Username:     seedCfg.AdminUsername,
		Email:        seedCfg.AdminEmail,
		Phone:        seedCfg.AdminPhone,
		PasswordHash: string(passwordHash),
		Role:         models.RoleAdmin,
		Status:       true,
	}

	if err := db.WithContext(ctx).Create(&admin).Error; err != nil {
		log.Fatalf("Admin seed failed: %v", err)
	}

	log.Println("Initial admin user created successfully")
	log.Printf("Username: %s", seedCfg.AdminUsername)
}
