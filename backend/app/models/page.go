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
	Title      string         `gorm:"type:varchar(255)" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	AuthorID   uuid.UUID      `gorm:"type:uuid;index" json:"authorId"`
	Slug       string         `gorm:"type:varchar(255);index" json:"slug"`
	Visibility PageVisibility `gorm:"type:varchar(50)" json:"visibility"`
	IsHomepage bool           `gorm:"default:false" json:"isHomepage"`

	CreatedAt time.Time `gorm:"type:timestamptz" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz" json:"updatedAt"`
	CSS       string    `gorm:"type:text" json:"css"`
}
