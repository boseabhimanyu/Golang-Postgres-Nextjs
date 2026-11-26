package postgres

import (
	"basic-app/apperrors"
	"basic-app/models"
	"basic-app/repository"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PageRepository struct {
	db *gorm.DB
}

func NewPageRepository(
	db *gorm.DB,
) repository.PageRepository {
	return &PageRepository{
		db: db,
	}
}

func (r *PageRepository) Create(
	ctx context.Context,
	page *models.Page,
) error {
	now := time.Now().UTC()

	if page.ID == uuid.Nil {
		page.ID = uuid.New()
	}

	if page.CreatedAt.IsZero() {
		page.CreatedAt = now
	}

	page.UpdatedAt = now

	return r.db.WithContext(ctx).Create(page).Error
}

func (r *PageRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.Page, error) {
	pageID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.ErrPageNotFound
	}

	var page models.Page

	err = r.db.WithContext(ctx).
		First(&page, "id = ?", pageID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrPageNotFound
		}

		return nil, err
	}

	return &page, nil
}

func (r *PageRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*models.Page, error) {
	var page models.Page

	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&page).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrPageNotFound
		}

		return nil, err
	}

	return &page, nil
}

func (r *PageRepository) List(
	ctx context.Context,
	filter repository.PageListFilter,
) ([]*models.Page, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Page{})

	if filter.Visibility != nil {
		query = query.Where(
			"visibility = ?",
			*filter.Visibility,
		)
	}

	if filter.AuthorID != "" {
		authorID, err := uuid.Parse(filter.AuthorID)
		if err != nil {
			return nil, 0, err
		}

		query = query.Where("author_id = ?", authorID)
	}

	// Text search.
	if filter.Search != "" {
		search := strings.NewReplacer(
			`\`, `\\`,
			"%", `\%`,
			"_", `\_`,
		).Replace(filter.Search)

		searchPattern := "%" + search + "%"

		query = query.Where(
			`(
				title ILIKE ? ESCAPE '\'
				OR slug ILIKE ? ESCAPE '\'
			)`,
			searchPattern,
			searchPattern,
		)
	}

	// Count matching pages before pagination.
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch matching pages.
	var pages []*models.Page

	query = query.
		Select(
			"id",
			"title",
			"content",
			"author_id",
			"slug",
			"visibility",
			"is_homepage",
			"created_at",
			"updated_at",
		).
		Order("created_at DESC")

	if filter.Skip > 0 {
		query = query.Offset(int(filter.Skip))
	}

	if filter.Limit > 0 {
		query = query.Limit(int(filter.Limit))
	}

	if err := query.Find(&pages).Error; err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

func (r *PageRepository) Update(
	ctx context.Context,
	id string,
	page *models.Page,
) error {
	pageID, err := uuid.Parse(id)
	if err != nil {
		return apperrors.ErrPageNotFound
	}

	result := r.db.WithContext(ctx).
		Model(&models.Page{}).
		Where("id = ?", pageID).
		Updates(map[string]interface{}{
			"title":      page.Title,
			"content":    page.Content,
			"css":        page.CSS,
			"slug":       page.Slug,
			"visibility": page.Visibility,
			"updated_at": time.Now().UTC(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrPageNotFound
	}

	return nil
}

func (r *PageRepository) Delete(
	ctx context.Context,
	id string,
) error {
	pageID, err := uuid.Parse(id)
	if err != nil {
		return apperrors.ErrPageNotFound
	}

	result := r.db.WithContext(ctx).
		Where("id = ? AND is_homepage = ?", pageID, false).
		Delete(&models.Page{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		var page models.Page

		err := r.db.WithContext(ctx).
			Select("id", "is_homepage").
			First(&page, "id = ?", pageID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrPageNotFound
		}

		if err != nil {
			return err
		}

		if page.IsHomepage {
			return apperrors.ErrHomepageCannotBeDeleted
		}

		return apperrors.ErrPageNotFound
	}

	return nil
}

func (r *PageRepository) FindByIDs(
	ctx context.Context,
	ids []string,
) (map[string]*models.Page, error) {
	result := make(map[string]*models.Page, len(ids))

	if len(ids) == 0 {
		return result, nil
	}

	pageIDs := make([]uuid.UUID, 0, len(ids))

	for _, id := range ids {
		pageID, err := uuid.Parse(id)
		if err != nil {
			return nil, apperrors.ErrInvalidMenuItem
		}

		pageIDs = append(pageIDs, pageID)
	}

	var pages []*models.Page

	err := r.db.WithContext(ctx).
		Select("id", "slug").
		Where("id IN ?", pageIDs).
		Find(&pages).Error

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		result[page.ID.String()] = page
	}

	return result, nil
}

func (r *PageRepository) SetHomepage(
	ctx context.Context,
	id string,
	isHomepage bool,
) error {
	pageID, err := uuid.Parse(id)
	if err != nil {
		return apperrors.ErrPageNotFound
	}

	now := time.Now().UTC()

	if !isHomepage {
		result := r.db.WithContext(ctx).
			Model(&models.Page{}).
			Where("id = ?", pageID).
			Updates(map[string]interface{}{
				"is_homepage": false,
				"updated_at":  now,
			})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return apperrors.ErrPageNotFound
		}

		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Verify that the target page exists and is public.
		var page models.Page

		err := tx.
			Where("id = ?", pageID).
			First(&page).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrPageNotFound
		}

		if err != nil {
			return err
		}

		if page.Visibility != models.PageVisibilityPublic {
			return apperrors.ErrHomepageRequiresPublicPage
		}

		// Remove homepage designation from existing pages.
		if err := tx.
			Model(&models.Page{}).
			Where("is_homepage = ?", true).
			Updates(map[string]interface{}{
				"is_homepage": false,
				"updated_at":  now,
			}).Error; err != nil {
			return err
		}

		// Assign the new homepage.
		result := tx.
			Model(&models.Page{}).
			Where(
				"id = ? AND visibility = ?",
				pageID,
				models.PageVisibilityPublic,
			).
			Updates(map[string]interface{}{
				"is_homepage": true,
				"updated_at":  now,
			})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return apperrors.ErrHomepageRequiresPublicPage
		}

		return nil
	})
}

func (r *PageRepository) FindHomepage(
	ctx context.Context,
) (*models.Page, error) {
	var page models.Page

	err := r.db.WithContext(ctx).
		Where("is_homepage = ?", true).
		First(&page).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrPageNotFound
		}

		return nil, err
	}

	return &page, nil
}
