package postgres

import (
	"basic-app/apperrors"
	"basic-app/models"
	"basic-app/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(
	db *gorm.DB,
) repository.RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) Create(
	ctx context.Context,
	token *models.RefreshToken,
) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}

	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now().UTC()
	}

	return r.db.WithContext(ctx).Create(token).Error
}

func (r *RefreshTokenRepository) FindByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*models.RefreshToken, error) {
	var token models.RefreshToken

	err := r.db.WithContext(ctx).
		Where(
			"token_hash = ? AND revoked_at IS NULL AND expires_at > ?",
			tokenHash,
			time.Now().UTC(),
		).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrRefreshTokenNotFound
		}

		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepository) Touch(
	ctx context.Context,
	tokenID string,
) error {
	id, err := uuid.Parse(tokenID)
	if err != nil {
		return apperrors.ErrRefreshTokenNotFound
	}

	result := r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("last_used_at", time.Now().UTC())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrRefreshTokenNotFound
	}

	return nil
}

func (r *RefreshTokenRepository) Revoke(
	ctx context.Context,
	tokenID string,
) error {
	id, err := uuid.Parse(tokenID)
	if err != nil {
		return apperrors.ErrRefreshTokenNotFound
	}

	result := r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", time.Now().UTC())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrRefreshTokenNotFound
	}

	return nil
}

func (r *RefreshTokenRepository) RevokeAllByUserID(
	ctx context.Context,
	userID string,
) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return apperrors.ErrRefreshTokenNotFound
	}

	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", time.Now().UTC()).Error
}

func (r *RefreshTokenRepository) DeleteRevokedBefore(
	ctx context.Context,
	before time.Time,
) error {
	return r.db.WithContext(ctx).
		Where("revoked_at IS NOT NULL AND revoked_at < ?", before).
		Delete(&models.RefreshToken{}).Error
}

func (r *RefreshTokenRepository) DeleteExpiredBefore(
	ctx context.Context,
	before time.Time,
) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", before).
		Delete(&models.RefreshToken{}).Error
}
