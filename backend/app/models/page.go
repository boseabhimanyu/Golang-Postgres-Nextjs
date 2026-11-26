package models

import (
	"time"

	"github.com/google/uuid"
)

type PageVisibility string

const (
	PageVisibilityPublic     PageVisibility = "public"
	PageVisibilityRegistered PageVisibility = "registered"
)

type Page struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	AuthorID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"authorId"`
	Slug       string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"slug"`
	Visibility PageVisibility `gorm:"type:varchar(50);not null" json:"visibility"`
	IsHomepage bool           `gorm:"not null;default:false" json:"isHomepage"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updatedAt"`
	CSS       string    `gorm:"type:text" json:"css"`
}
