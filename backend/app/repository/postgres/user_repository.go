package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"basic-app/apperrors"
	"basic-app/models"
	"basic-app/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const databaseTimeout = 5 * time.Second

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user.
func (r *UserRepository) Create(
	ctx context.Context,
	user *models.User,
) error {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	now := time.Now().UTC()

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}

	user.UpdatedAt = now

	user.Email = strings.ToLower(
		strings.TrimSpace(user.Email),
	)

	user.Username = strings.TrimSpace(
		user.Username,
	)

	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID finds a user by UUID.
func (r *UserRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.ErrInvalidUserID
	}

	var user models.User

	err = r.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByUsername finds a user by username.
func (r *UserRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	username = strings.TrimSpace(username)

	var user models.User

	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates general user information.
func (r *UserRepository) Update(
	ctx context.Context,
	user *models.User,
) error {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	user.UpdatedAt = time.Now().UTC()

	backupCodesJSON, err := json.Marshal(user.BackupCodeHashes)
	if err != nil {
		return fmt.Errorf("marshal backup_code_hashes: %w", err)
	}
	if user.BackupCodeHashes == nil {
		backupCodesJSON = []byte("[]")
	}

	updates := map[string]interface{}{
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"username":      strings.TrimSpace(user.Username),
		"email":         strings.ToLower(strings.TrimSpace(user.Email)),
		"alt_email":     user.AltEmail,
		"phone":         user.Phone,
		"date_of_birth": user.DateOfBirth,

		"address_line_1": user.AddressLine1,
		"address_line_2": user.AddressLine2,
		"city":           user.City,
		"state":          user.State,
		"pin_code":       user.PinCode,

		"two_factor_enabled":        user.TwoFactorEnabled,
		"two_factor_secret":         user.TwoFactorSecret,
		"two_factor_pending_secret": user.TwoFactorPendingSecret,
		"backup_code_hashes":        string(backupCodesJSON),

		"updated_at": user.UpdatedAt,
	}

	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

// UpdateProfilePic updates only the user's profile picture.
func (r *UserRepository) UpdateProfilePic(
	ctx context.Context,
	userID string,
	profilePic string,
) error {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	id, err := uuid.Parse(userID)
	if err != nil {
		return apperrors.ErrInvalidUserID
	}

	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"profile_pic": profilePic,
			"updated_at":  time.Now().UTC(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

// UpdatePassword updates the password hash.
func (r *UserRepository) UpdatePassword(
	ctx context.Context,
	userID string,
	passwordHash string,
) error {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	id, err := uuid.Parse(userID)
	if err != nil {
		return apperrors.ErrInvalidUserID
	}

	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password_hash": passwordHash,
			"updated_at":    time.Now().UTC(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

// Change user account status without deleting it.
func (r *UserRepository) UserStatus(
	ctx context.Context,
	userID string,
	status bool,
) error {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	id, err := uuid.Parse(userID)
	if err != nil {
		return apperrors.ErrInvalidUserID
	}

	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now().UTC(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

// FindByPhone finds a user by phone number.
func (r *UserRepository) FindByPhone(
	ctx context.Context,
	phone string,
) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	var user models.User

	err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByAnyEmail finds a user by primary or alternate email.
func (r *UserRepository) FindByAnyEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	email = strings.ToLower(strings.TrimSpace(email))

	var user models.User

	err := r.db.WithContext(ctx).
		Where("email = ? OR alt_email = ?", email, email).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByLoginIdentifier finds a user using any supported login identifier.
func (r *UserRepository) FindByLoginIdentifier(
	ctx context.Context,
	identifier string,
) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	identifier = strings.TrimSpace(identifier)

	var user models.User

	err := r.db.WithContext(ctx).
		Where(
			"username = ? OR email = ? OR alt_email = ? OR phone = ?",
			identifier,
			strings.ToLower(identifier),
			strings.ToLower(identifier),
			identifier,
		).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ListCustomers retrieves customers with optional filters and pagination.

func (r *UserRepository) ListCustomers(
	ctx context.Context,
	listFilter repository.CustomerListFilter,
) ([]models.User, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	query := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("role = ?", models.RoleCustomer)

	// Filter by account status.
	if listFilter.Status != nil {
		query = query.Where("status = ?", *listFilter.Status)
	}

	// Search across customer fields.
	if listFilter.Search != "" {
		search := strings.NewReplacer(
			`\`, `\\`,
			"%", `\%`,
			"_", `\_`,
		).Replace(listFilter.Search)

		searchPattern := "%" + search + "%"

		query = query.Where(
			`(
				first_name ILIKE ? ESCAPE '\'
				OR last_name ILIKE ? ESCAPE '\'
				OR username ILIKE ? ESCAPE '\'
				OR email ILIKE ? ESCAPE '\'
				OR alt_email ILIKE ? ESCAPE '\'
				OR phone ILIKE ? ESCAPE '\'
			)`,
			searchPattern,
			searchPattern,
			searchPattern,
			searchPattern,
			searchPattern,
			searchPattern,
		)
	}

	// Get total count before pagination.
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch customers.
	var customers []models.User

	query = query.Order("created_at DESC")

	if listFilter.Skip > 0 {
		query = query.Offset(int(listFilter.Skip))
	}

	if listFilter.Limit > 0 {
		query = query.Limit(int(listFilter.Limit))
	}

	if err := query.Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

// UpdateLastLoginAt updates the user's last login timestamp.
func (r *UserRepository) UpdateLastLoginAt(
	ctx context.Context,
	userID string,
	lastLoginAt time.Time,
) error {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	id, err := uuid.Parse(userID)
	if err != nil {
		return apperrors.ErrInvalidUserID
	}

	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("last_login_at", lastLoginAt)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}
