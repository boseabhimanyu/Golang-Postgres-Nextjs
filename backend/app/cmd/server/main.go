package main

import (
	"basic-app/config"
	"basic-app/database"
	"basic-app/repository"
	"basic-app/repository/postgres"
	"basic-app/router"
	"basic-app/services"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Config Error: %v", err)
		return
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Printf("DB Error: %v", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get SQL connection: %v", err)
		return
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("PostgreSQL disconnect error: %v", err)
		}
	}()

	// Initialize repositories
	refreshTokenRepository := postgres.NewRefreshTokenRepository(db)

	settingsRepository := postgres.NewApplicationSettingsRepository(db)

	settingsService := services.NewSettingsService(
		settingsRepository,
		cfg,
	)

	// Ensure application settings exist
	settingsCtx, settingsCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer settingsCancel()

	if err := settingsRepository.EnsureDefaults(settingsCtx); err != nil {
		log.Printf("Failed to ensure application settings: %v", err)
		return
	}

	if err := settingsService.ValidateTwoFactorStartupConfig(
		settingsCtx,
	); err != nil {
		log.Printf("Two-factor startup validation failed: %v", err)
		return
	}

	// Refresh token cleanup
	cleanupCtx, cleanupCancel := context.WithCancel(
		context.Background(),
	)
	defer cleanupCancel()

	log.Printf(
		"Refresh token revoked cleanup retention: %d days",
		cfg.RefreshTokenRevokedRetentionDays,
	)

	log.Printf(
		"Refresh token expired cleanup retention: %d days",
		cfg.RefreshTokenExpiredRetentionDays,
	)

	go startRefreshTokenCleanup(
		cleanupCtx,
		refreshTokenRepository,
		cfg.RefreshTokenRevokedRetentionDays,
		cfg.RefreshTokenExpiredRetentionDays,
	)

	gin.SetMode(cfg.GinMode)

	engine := router.NewRouter(db, cfg)

	addr := ":" + cfg.ServerPort

	server := &http.Server{
		Addr:              addr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Server listening on http://localhost%s", addr)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("Server Failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

func startRefreshTokenCleanup(
	ctx context.Context,
	repo repository.RefreshTokenRepository,
	revokedRetentionDays int,
	expiredRetentionDays int,
) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	cleanup := func() {
		now := time.Now().UTC()

		revokedBefore := now.AddDate(
			0, 0, -revokedRetentionDays,
		)

		expiredBefore := now.AddDate(
			0, 0, -expiredRetentionDays,
		)

		cleanupCtx, cancel := context.WithTimeout(
			ctx,
			30*time.Second,
		)
		defer cancel()

		if err := repo.DeleteRevokedBefore(
			cleanupCtx,
			revokedBefore,
		); err != nil {
			log.Printf("Revoked token cleanup failed: %v", err)
			return
		}

		if err := repo.DeleteExpiredBefore(
			cleanupCtx,
			expiredBefore,
		); err != nil {
			log.Printf("Expired token cleanup failed: %v", err)
			return
		}

		log.Println("Refresh token cleanup completed")
	}

	cleanup()

	for {
		select {
		case <-ticker.C:
			cleanup()
		case <-ctx.Done():
			return
		}
	}
}
