package repository

import (
	"context"
	"errors"
	"fmt"

	"bus-booking/shared/constants"
	"bus-booking/user-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*model.User, int64, error)
	ListByRole(ctx context.Context, role constants.UserRole, limit, offset int) ([]*model.User, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	EmailExists(ctx context.Context, email string) (bool, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("User already exists: %w", err)
		}
		return fmt.Errorf("Failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("Failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("Failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("Failed to get user by Firebase UID: %w", err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("Failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return fmt.Errorf("Failed to delete user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("Failed to count users: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("Failed to list users: %w", err)
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) ListByRole(ctx context.Context, role constants.UserRole, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("role = ?", role).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("role = ?", role).Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("Failed to list users by role: %w", err)
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	if err := r.db.WithContext(ctx).
		Model(&model.User{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("Failed to update user status: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("Failed to check email existence: %w", err)
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("username = ?", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("Failed to check username existence: %w", err)
	}
	return count > 0, nil
}
