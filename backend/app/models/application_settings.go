package models

import (
	"time"
)

const ApplicationSettingsID = "application"

type ApplicationSettings struct {
	ID                    string               `gorm:"type:varchar(50);primaryKey" json:"-"`
	RegistrationEnabled   bool                 `gorm:"default:true" json:"registrationEnabled"`
	TwoFactorEnabled      bool                 `gorm:"default:false" json:"twoFactorEnabled"`
	LoginWithPrimaryEmail bool                 `gorm:"default:true" json:"loginWithPrimaryEmail"`
	LoginWithUsername     bool                 `gorm:"default:true" json:"loginWithUsername"`
	LoginWithPhone        bool                 `gorm:"default:false" json:"loginWithPhone"`
	LoginWithAltEmail     bool                 `gorm:"default:false" json:"loginWithAltEmail"`
	MenuMaxDepth          MenuMaxDepthSettings `gorm:"type:jsonb;serializer:json" json:"menuMaxDepth"`
	UpdatedAt             time.Time            `gorm:"type:timestamptz" json:"-"`
	UpdatedBy             string               `gorm:"type:varchar(255)" json:"-"`
}

type MenuMaxDepthSettings struct {
	Top    int `json:"top"`
	Left   int `json:"left"`
	Bottom int `json:"bottom"`
}
