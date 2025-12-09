package repository

import (
	"context"
	"errors"
	"fmt"

	"bus-booking/shared/constants"
	"bus-booking/shared/utils/dbutils"
	"bus-booking/user-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*model.User, error)
	List(ctx context.Context, query model.UserListQuery) ([]*model.User, int64, error)
	ListByRole(ctx context.Context, role constants.UserRole, limit, offset int) ([]*model.User, int64, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, dbutils.WrapIfNotFound(err, "không tìm thấy người dùng theo ID")
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, dbutils.WrapIfNotFound(err, "không tìm thấy người dùng theo email")
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, dbutils.WrapIfNotFound(err, "không tìm thấy người dùng theo số điện thoại")
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return nil, dbutils.WrapIfNotFound(err, "không tìm thấy người dùng theo firebase uid")
	}
	return &user, nil
}

func (r *UserRepositoryImpl) List(ctx context.Context, query model.UserListQuery) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// Build query with filters
	db := r.db.WithContext(ctx).Model(&model.User{})

	// Apply search filter
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("full_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Apply role filter
	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}

	// Apply status filter
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// Count total with filters
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("không thể đếm người dùng: %w", err)
	}

	// Apply sorting
	sortBy := "created_at"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}
	sortOrder := "DESC"
	if !query.SortDesc {
		sortOrder = "ASC"
	}
	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)

	// Apply pagination and get results
	offset := (query.Page - 1) * query.PageSize
	if err := db.Limit(query.PageSize).Offset(offset).Order(orderClause).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("không thể lấy danh sách người dùng: %w", err)
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) ListByRole(ctx context.Context, role constants.UserRole, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("role = ?", role).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("không thể đếm người dùng theo vai trò: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("role = ?", role).Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("không thể lấy danh sách người dùng theo vai trò: %w", err)
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("không thể kiểm tra email: %w", err)
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("user already exists: %w", err)
		}
		return fmt.Errorf("không thể tạo người dùng: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("không thể cập nhật người dùng: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return fmt.Errorf("không thể xóa người dùng: %w", err)
	}
	return nil
}
