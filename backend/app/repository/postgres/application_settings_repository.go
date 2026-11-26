package postgres

import (
	"context"
	"errors"
	"time"

	"basic-app/apperrors"
	"basic-app/models"
	"basic-app/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ApplicationSettingsRepository struct {
	db *gorm.DB
}

func NewApplicationSettingsRepository(
	db *gorm.DB,
) repository.ApplicationSettingsRepository {
	return &ApplicationSettingsRepository{
		db: db,
	}
}

// Get retrieves application settings.
func (r *ApplicationSettingsRepository) Get(
	ctx context.Context,
) (*models.ApplicationSettings, error) {
	var settings models.ApplicationSettings

	err := r.db.WithContext(ctx).
		First(&settings, "id = ?", models.ApplicationSettingsID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrApplicationSettingsNotFound
	}

	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// EnsureDefaults creates application settings if they do not exist.
func (r *ApplicationSettingsRepository) EnsureDefaults(
	ctx context.Context,
) error {
	now := time.Now().UTC()

	settings := models.ApplicationSettings{
		ID:                    models.ApplicationSettingsID,
		RegistrationEnabled:   true,
		TwoFactorEnabled:      false,
		LoginWithPrimaryEmail: true,
		LoginWithUsername:     true,
		LoginWithPhone:        false,
		LoginWithAltEmail:     false,
		MenuMaxDepth: models.MenuMaxDepthSettings{
			Top:    3,
			Left:   1,
			Bottom: 2,
		},
		UpdatedAt: now,
		UpdatedBy: "system",
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoNothing: true,
		}).
		Create(&settings).
		Error
}

// UpdateRegistrationEnabled updates registration availability.
func (r *ApplicationSettingsRepository) UpdateRegistrationEnabled(
	ctx context.Context,
	enabled bool,
	updatedBy string,
) error {
	return r.update(ctx, map[string]interface{}{
		"registration_enabled": enabled,
		"updated_by":           updatedBy,
	})
}

// UpdateTwoFactorEnabled updates global 2FA availability.
func (r *ApplicationSettingsRepository) UpdateTwoFactorEnabled(
	ctx context.Context,
	enabled bool,
	updatedBy string,
) error {
	return r.update(ctx, map[string]interface{}{
		"two_factor_enabled": enabled,
		"updated_by":         updatedBy,
	})
}

// UpdateLoginIdentifiers updates available login methods.
func (r *ApplicationSettingsRepository) UpdateLoginIdentifiers(
	ctx context.Context,
	primaryEmail bool,
	username bool,
	phone bool,
	altEmail bool,
	updatedBy string,
) error {
	return r.update(ctx, map[string]interface{}{
		"login_with_primary_email": primaryEmail,
		"login_with_username":      username,
		"login_with_phone":         phone,
		"login_with_alt_email":     altEmail,
		"updated_by":               updatedBy,
	})
}

// UpdateMenuMaxDepth updates menu nesting limits.
func (r *ApplicationSettingsRepository) UpdateMenuMaxDepth(
	ctx context.Context,
	top int,
	left int,
	bottom int,
	updatedBy string,
) error {
	return r.update(ctx, map[string]interface{}{
		"menu_max_depth": models.MenuMaxDepthSettings{
			Top:    top,
			Left:   left,
			Bottom: bottom,
		},
		"updated_by": updatedBy,
	})
}

// update is a shared helper for application settings updates.
func (r *ApplicationSettingsRepository) update(
	ctx context.Context,
	values map[string]interface{},
) error {
	values["updated_at"] = time.Now().UTC()

	result := r.db.WithContext(ctx).
		Model(&models.ApplicationSettings{}).
		Where("id = ?", models.ApplicationSettingsID).
		Updates(values)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrApplicationSettingsNotFound
	}

	return nil
}
