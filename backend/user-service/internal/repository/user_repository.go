package repository

import (
	"context"

	"bus-booking/shared/ginext"
	"bus-booking/user-service/internal/model"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*model.User, int64, error)
	ListByRole(ctx context.Context, role model.UserRole, limit, offset int) ([]*model.User, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Err(err).Msg("Failed to create user in database")
		return ginext.NewInternalServerError("Failed to create user")
	}
	return nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		log.Err(err).Msg("Failed to get user by ID")
		return nil, ginext.NewInternalServerError("Failed to get user by ID")
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		log.Err(err).Msg("Failed to get user by email")
		return nil, ginext.NewInternalServerError("Failed to get user by email")
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		log.Err(err).Msg("Failed to get user by username")
		return nil, ginext.NewInternalServerError("Failed to get user by username")
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("firebase_uid = ?", firebaseUID).First(&user).Error
	if err != nil {
		log.Err(err).Msg("Failed to get user by Firebase UID")
		return nil, ginext.NewInternalServerError("Failed to get user by Firebase UID")
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		log.Err(err).Msg("Failed to update user")
		return ginext.NewInternalServerError("Failed to update user")
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		log.Err(err).Msg("Failed to delete user")
		return ginext.NewInternalServerError("Failed to delete user")
	}
	return nil
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		log.Err(err).Msg("Failed to count users")
		return nil, 0, ginext.NewInternalServerError("Failed to count users")
	}

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		log.Err(err).Msg("Failed to list users")
		return nil, 0, ginext.NewInternalServerError("Failed to list users")
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) ListByRole(ctx context.Context, role model.UserRole, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("role = ?", role).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("role = ?", role).Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		log.Err(err).Msg("Failed to list users by role")
		return nil, 0, ginext.NewInternalServerError("Failed to list users by role")
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	if err := r.db.WithContext(ctx).
		Model(&model.User{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		log.Err(err).Msg("Failed to update user status in database")
		return ginext.NewInternalServerError("Failed to update user status")
	}
	return nil
}

func (r *UserRepositoryImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email = ?", email).Count(&count).Error; err != nil {
		log.Err(err).Msg("Failed to check email existence")
		return false, ginext.NewInternalServerError("Failed to check email existence")
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("username = ?", username).Count(&count).Error; err != nil {
		log.Err(err).Msg("Failed to check username existence")
		return false, ginext.NewInternalServerError("Failed to check username existence")
	}
	return count > 0, nil
}
