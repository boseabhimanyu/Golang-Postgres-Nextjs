package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
)

func EnsureIndexes(db *gorm.DB) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		indexes := []string{

			// -------------------------------------------------
			// Menu indexes
			// -------------------------------------------------

			`CREATE UNIQUE INDEX IF NOT EXISTS unique_menu_location
			ON menus (location)`,

			// -------------------------------------------------
			// User indexes
			// -------------------------------------------------

			`CREATE UNIQUE INDEX IF NOT EXISTS unique_email
			ON users (email)`,

			`CREATE UNIQUE INDEX IF NOT EXISTS unique_username
			ON users (username)`,

			`CREATE UNIQUE INDEX IF NOT EXISTS unique_phone
			ON users (phone)`,

			`CREATE UNIQUE INDEX IF NOT EXISTS unique_alt_email
			ON users (alt_email)
			WHERE alt_email IS NOT NULL AND alt_email <> ''`,

			// -------------------------------------------------
			// Refresh token indexes
			// -------------------------------------------------

			`CREATE UNIQUE INDEX IF NOT EXISTS unique_token_hash
			ON refresh_tokens (token_hash)`,

			`CREATE INDEX IF NOT EXISTS refresh_token_user_id
			ON refresh_tokens (user_id)`,

			`CREATE INDEX IF NOT EXISTS refresh_token_expiry
			ON refresh_tokens (expires_at)`,

			// -------------------------------------------------
			// Page indexes
			// -------------------------------------------------

			`CREATE UNIQUE INDEX IF NOT EXISTS unique_page_slug
			ON pages (slug)`,

			`CREATE INDEX IF NOT EXISTS page_visibility_created_at
			ON pages (visibility, created_at DESC)`,

			`CREATE INDEX IF NOT EXISTS page_author_created_at
			ON pages (author_id, created_at DESC)`,

			`CREATE UNIQUE INDEX IF NOT EXISTS unique_homepage
			ON pages (is_homepage)
			WHERE is_homepage = TRUE`,
		}

		for _, index := range indexes {
			if err := tx.Exec(index).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
