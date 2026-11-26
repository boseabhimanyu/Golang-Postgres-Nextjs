package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	FirstName string `gorm:"type:varchar(255)" json:"firstName"`
	LastName  string `gorm:"type:varchar(255)" json:"lastName"`
	Username  string `gorm:"type:varchar(255)" json:"username"`
	Email     string `gorm:"type:varchar(255)" json:"email"`
	AltEmail  string `gorm:"type:varchar(255)" json:"altEmail"`
	Phone     string `gorm:"type:varchar(50)" json:"phone"`

	Role       UserRole `gorm:"type:varchar(50)" json:"role"`
	ProfilePic string   `gorm:"type:text" json:"profilePic"`
	Status     bool     `gorm:"default:true" json:"status"`

	PasswordHash string `gorm:"type:text" json:"-"`

	CreatedAt   time.Time  `gorm:"type:timestamptz" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"type:timestamptz" json:"updatedAt"`
	LastLoginAt *time.Time `gorm:"type:timestamptz" json:"lastLoginAt"`
	DateOfBirth *time.Time `gorm:"type:date" json:"dateOfBirth"`

	AddressLine1 string `gorm:"type:text" json:"addressLine1"`
	AddressLine2 string `gorm:"type:text" json:"addressLine2"`
	City         string `gorm:"type:varchar(255)" json:"city"`
	State        string `gorm:"type:varchar(255)" json:"state"`
	PinCode      string `gorm:"type:varchar(20)" json:"pinCode"`

	TwoFactorEnabled       bool     `gorm:"default:false" json:"twoFactorEnabled"`
	TwoFactorSecret        string   `gorm:"type:text" json:"-"`
	TwoFactorPendingSecret string   `gorm:"type:text" json:"-"`
	BackupCodeHashes       []string `gorm:"type:jsonb;serializer:json" json:"-"`
}

func NewCustomerUser() User {
	now := time.Now().UTC()

	return User{
		ID:        uuid.New(),
		Role:      RoleCustomer,
		Status:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
