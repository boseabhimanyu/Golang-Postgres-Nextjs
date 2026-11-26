package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	TokenHash string    `gorm:"type:text;not null" json:"-"`

	CreatedAt  time.Time  `gorm:"type:timestamptz;not null" json:"createdAt"`
	LastUsedAt *time.Time `gorm:"type:timestamptz" json:"lastUsedAt"`
	ExpiresAt  time.Time  `gorm:"type:timestamptz;not null;index" json:"expiresAt"`
	RevokedAt  *time.Time `gorm:"type:timestamptz;index" json:"revokedAt"`
}
