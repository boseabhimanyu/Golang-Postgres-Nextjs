package postgres

import (
	"context"
	"errors"
	"time"

	"basic-app/apperrors"
	"basic-app/models"
	"basic-app/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(
	db *gorm.DB,
) repository.MenuRepository {
	return &MenuRepository{
		db: db,
	}
}

// Create creates a new menu.
func (r *MenuRepository) Create(
	ctx context.Context,
	menu *models.Menu,
) error {
	if menu.ID == uuid.Nil {
		menu.ID = uuid.New()
	}

	now := time.Now().UTC()

	if menu.CreatedAt.IsZero() {
		menu.CreatedAt = now
	}

	menu.UpdatedAt = now

	return r.db.WithContext(ctx).Create(menu).Error
}

// FindByID finds a menu by UUID.
func (r *MenuRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.Menu, error) {
	menuID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.ErrMenuNotFound
	}

	var menu models.Menu

	err = r.db.WithContext(ctx).
		First(&menu, "id = ?", menuID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrMenuNotFound
	}

	if err != nil {
		return nil, err
	}

	return &menu, nil
}

// FindByLocation finds a menu by location.
func (r *MenuRepository) FindByLocation(
	ctx context.Context,
	location models.MenuLocation,
) (*models.Menu, error) {
	var menu models.Menu

	err := r.db.WithContext(ctx).
		Where("location = ?", location).
		First(&menu).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrMenuNotFound
	}

	if err != nil {
		return nil, err
	}

	return &menu, nil
}

// List returns all menus ordered by location.
func (r *MenuRepository) List(
	ctx context.Context,
) ([]*models.Menu, error) {
	menus := make([]*models.Menu, 0)

	err := r.db.WithContext(ctx).
		Order("location ASC").
		Find(&menus).
		Error

	if err != nil {
		return nil, err
	}

	return menus, nil
}

// Update updates an existing menu.
func (r *MenuRepository) Update(
	ctx context.Context,
	id string,
	menu *models.Menu,
) error {
	menuID, err := uuid.Parse(id)
	if err != nil {
		return apperrors.ErrMenuNotFound
	}

	menu.UpdatedAt = time.Now().UTC()

	result := r.db.WithContext(ctx).
		Model(&models.Menu{}).
		Where("id = ?", menuID).
		Updates(map[string]interface{}{
			"name":       menu.Name,
			"location":   menu.Location,
			"status":     menu.Status,
			"items":      menu.Items,
			"updated_at": menu.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrMenuNotFound
	}

	return nil
}

// Delete permanently removes a menu.
func (r *MenuRepository) Delete(
	ctx context.Context,
	id string,
) error {
	menuID, err := uuid.Parse(id)
	if err != nil {
		return apperrors.ErrMenuNotFound
	}

	result := r.db.WithContext(ctx).
		Delete(&models.Menu{}, "id = ?", menuID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrMenuNotFound
	}

	return nil
}

// FindAllByLocations finds menus matching the supplied locations.
func (r *MenuRepository) FindAllByLocations(
	ctx context.Context,
	locations []models.MenuLocation,
) ([]*models.Menu, error) {
	menus := make([]*models.Menu, 0)

	if len(locations) == 0 {
		return menus, nil
	}

	err := r.db.WithContext(ctx).
		Where("location IN ?", locations).
		Order("location ASC").
		Find(&menus).
		Error

	if err != nil {
		return nil, err
	}

	return menus, nil
}
