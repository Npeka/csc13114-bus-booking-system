package service

import (
	"context"
	"mime/multipart"

	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
	"bus-booking/shared/storage"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserService interface {
	CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.UserResponse, error)
	ListUsers(ctx context.Context, req model.UserListQuery) ([]*model.UserResponse, int64, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *model.UserUpdateRequest) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsersByRole(ctx context.Context, role constants.UserRole, limit, offset int) ([]*model.UserResponse, int64, error)
	UploadAvatar(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*model.UserResponse, error)
	DeleteAvatar(ctx context.Context, userID uuid.UUID) error
}

type UserServiceImpl struct {
	userRepo       repository.UserRepository
	storageService storage.StorageService
}

func NewUserService(
	userRepo repository.UserRepository,
	storageService storage.StorageService,
) UserService {
	return &UserServiceImpl{
		userRepo:       userRepo,
		storageService: storageService,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.UserResponse, error) {
	// Validate email if provided
	if req.Email != "" {
		if emailExists, err := s.userRepo.EmailExists(ctx, req.Email); err != nil {
			log.Error().Err(err).Msg("Failed to check email existence")
			return nil, ginext.NewInternalServerError("Không thể xác thực email")
		} else if emailExists {
			return nil, ginext.NewConflictError("email đã tồn tại")
		}
	}

	// Check if Firebase UID already exists
	if existingUser, err := s.userRepo.GetByFirebaseUID(ctx, req.FirebaseUID); err == nil && existingUser != nil {
		return nil, ginext.NewConflictError("người dùng với Firebase UID này đã tồn tại")
	}

	user := &model.User{
		Email:         req.Email,
		Phone:         req.Phone,
		FullName:      req.FullName,
		Avatar:        req.Avatar,
		Role:          req.Role,
		Status:        constants.UserStatusActive,
		FirebaseUID:   &req.FirebaseUID,
		EmailVerified: false,
		PhoneVerified: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to create user in database")
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ginext.NewInternalServerError("Không thể lấy người dùng theo ID")
	}
	return user.ToResponse(), nil
}

// ListUsers gets a paginated list of users with filtering
func (s *UserServiceImpl) ListUsers(ctx context.Context, req model.UserListQuery) ([]*model.UserResponse, int64, error) {
	users, total, err := s.userRepo.List(ctx, req)
	if err != nil {
		return nil, 0, ginext.NewInternalServerError("Không thể lấy danh sách người dùng")
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id uuid.UUID, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate and update email if provided
	if req.Email != nil && *req.Email != user.Email {
		if emailExists, err := s.userRepo.EmailExists(ctx, *req.Email); err != nil {
			log.Error().Err(err).Msg("Failed to check email existence during update")
			return nil, ginext.NewInternalServerError("Không thể xác thực email")
		} else if emailExists {
			return nil, ginext.NewConflictError("email đã tồn tại")
		}
		user.Email = *req.Email
	}

	// Update other fields
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user in database")
		return nil, ginext.NewInternalServerError("Không thể cập nhật người dùng")
	}

	return user.ToResponse(), nil
}

// DeleteUser soft deletes a user
func (s *UserServiceImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Check if user exists
	if _, err := s.userRepo.GetByID(ctx, id); err != nil {
		return err // Repository already returns proper ginext errors
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		log.Error().Err(err).Msg("Failed to delete user from database")
		return ginext.NewInternalServerError("Không thể xóa người dùng")
	}

	return nil
}

// ListUsersByRole gets a paginated list of users by role
func (s *UserServiceImpl) ListUsersByRole(ctx context.Context, role constants.UserRole, limit, offset int) ([]*model.UserResponse, int64, error) {
	users, total, err := s.userRepo.ListByRole(ctx, role, limit, offset)
	if err != nil {
		return nil, 0, ginext.NewInternalServerError("Không thể lấy danh sách người dùng theo vai trò")
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

// UploadAvatar uploads a new avatar for the user
func (s *UserServiceImpl) UploadAvatar(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*model.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" && contentType != "image/webp" {
		return nil, ginext.NewBadRequestError("Chỉ chấp nhận file ảnh (JPEG, PNG, WebP)")
	}

	// Validate file size (max 5MB)
	if header.Size > 5*1024*1024 {
		return nil, ginext.NewBadRequestError("Kích thước file không được vượt quá 5MB")
	}

	// Delete old avatar if exists
	if user.Avatar != "" {
		if err := s.storageService.DeleteFile(ctx, user.Avatar); err != nil {
			log.Warn().Err(err).Msg("Failed to delete old avatar, continuing with upload")
		}
	}

	// Upload new avatar
	avatarURL, err := s.storageService.UploadFile(ctx, file, header, "avatars")
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload avatar to storage")
		return nil, ginext.NewInternalServerError("Không thể tải ảnh lên")
	}

	// Update user's avatar URL
	user.Avatar = avatarURL
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user avatar in database")
		// Try to clean up uploaded file
		if err := s.storageService.DeleteFile(ctx, avatarURL); err != nil {
			log.Warn().Err(err).Msg("Failed to delete uploaded avatar, continuing with error")
		}
		return nil, ginext.NewInternalServerError("Không thể cập nhật avatar")
	}

	return user.ToResponse(), nil
}

// DeleteAvatar removes the user's avatar
func (s *UserServiceImpl) DeleteAvatar(ctx context.Context, userID uuid.UUID) error {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if user has an avatar
	if user.Avatar == "" {
		return ginext.NewBadRequestError("Người dùng không có avatar")
	}

	// Delete from storage
	if err := s.storageService.DeleteFile(ctx, user.Avatar); err != nil {
		log.Error().Err(err).Msg("Failed to delete avatar from storage")
		return ginext.NewInternalServerError("Không thể xóa avatar")
	}

	// Update user's avatar to empty string
	user.Avatar = ""
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user in database after avatar deletion")
		return ginext.NewInternalServerError("Không thể cập nhật người dùng")
	}

	return nil
}
