package service

import (
	"context"

	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
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
}

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(
	userRepo repository.UserRepository,
) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.UserResponse, error) {
	// Validate email if provided
	if req.Email != "" {
		if emailExists, err := s.userRepo.EmailExists(ctx, req.Email); err != nil {
			log.Error().Err(err).Msg("Failed to check email existence")
			return nil, ginext.NewInternalServerError("Failed to validate email")
		} else if emailExists {
			return nil, ginext.NewConflictError("email already exists")
		}
	}

	// Check if Firebase UID already exists
	if existingUser, err := s.userRepo.GetByFirebaseUID(ctx, req.FirebaseUID); err == nil && existingUser != nil {
		return nil, ginext.NewConflictError("user with this Firebase UID already exists")
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
		return nil, ginext.NewInternalServerError("Failed to get user by ID")
	}
	return user.ToResponse(), nil
}

// ListUsers gets a paginated list of users with filtering
func (s *UserServiceImpl) ListUsers(ctx context.Context, req model.UserListQuery) ([]*model.UserResponse, int64, error) {
	users, total, err := s.userRepo.List(ctx, req)
	if err != nil {
		return nil, 0, ginext.NewInternalServerError("Failed to list users")
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
			return nil, ginext.NewInternalServerError("Failed to validate email")
		} else if emailExists {
			return nil, ginext.NewConflictError("email already exists")
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
		return nil, ginext.NewInternalServerError("Failed to update user")
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
		return ginext.NewInternalServerError("Failed to delete user")
	}

	return nil
}

// ListUsersByRole gets a paginated list of users by role
func (s *UserServiceImpl) ListUsersByRole(ctx context.Context, role constants.UserRole, limit, offset int) ([]*model.UserResponse, int64, error) {
	users, total, err := s.userRepo.ListByRole(ctx, role, limit, offset)
	if err != nil {
		return nil, 0, ginext.NewInternalServerError("Failed to list users by role")
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}
