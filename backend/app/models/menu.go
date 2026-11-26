package models

import (
	"time"

	"github.com/google/uuid"
)

type MenuLocation string

const (
	MenuLocationTop    MenuLocation = "top"
	MenuLocationLeft   MenuLocation = "left"
	MenuLocationBottom MenuLocation = "bottom"
)

type MenuItemType string

const (
	MenuItemTypePage  MenuItemType = "page"
	MenuItemTypeURL   MenuItemType = "url"
	MenuItemTypeGroup MenuItemType = "group"
)

type Menu struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string       `gorm:"type:varchar(255)" json:"name"`
	Location  MenuLocation `gorm:"type:varchar(50)" json:"location"`
	Status    bool         `gorm:"default:true" json:"status"`
	Items     []MenuItem   `gorm:"type:jsonb;serializer:json" json:"items"`
	CreatedAt time.Time    `gorm:"type:timestamptz" json:"createdAt"`
	UpdatedAt time.Time    `gorm:"type:timestamptz" json:"updatedAt"`
}

type MenuItem struct {
	ID            uuid.UUID    `json:"id"`
	Label         string       `json:"label"`
	Type          MenuItemType `json:"type"`
	PageID        *uuid.UUID   `json:"pageId"`
	URL           string       `json:"url"`
	Order         int          `json:"order"`
	DisplayStatus bool         `json:"displayStatus"`
	Children      []MenuItem   `json:"children"`
}
