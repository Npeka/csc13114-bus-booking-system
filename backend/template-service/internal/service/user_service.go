package service

import (
	"context"
	"errors"
	"time"

	"bus-booking/shared/constants"
	"bus-booking/shared/utils"
	"bus-booking/template-service/internal/model"
	"bus-booking/template-service/internal/repository"

	"github.com/rs/zerolog/log"
)

// UserServiceInterface defines the interface for user service
type UserServiceInterface interface {
	CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, req *model.UserUpdateRequest) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, limit, offset int) ([]*model.UserResponse, int64, error)
	ListUsersByRole(ctx context.Context, role string, limit, offset int) ([]*model.UserResponse, int64, error)
	UpdateUserStatus(ctx context.Context, id uint, status string) error
}

// UserService implements UserServiceInterface
type UserService struct {
	userRepo repository.UserRepositoryInterface
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.UserResponse, error) {
	// Check if email already exists
	emailExists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check email existence")
		return nil, errors.New(constants.ErrInternalServer)
	}
	if emailExists {
		return nil, errors.New("email already exists")
	}

	// Check if username already exists
	usernameExists, err := s.userRepo.UsernameExists(ctx, req.Username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check username existence")
		return nil, errors.New(constants.ErrInternalServer)
	}
	if usernameExists {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, errors.New(constants.ErrInternalServer)
	}

	// Create user model
	user := &model.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      req.Role,
		Status:    "active",
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// Create user
	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, errors.New(constants.ErrInternalServer)
	}

	return user.ToResponse(), nil
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by ID")
		return nil, errors.New(constants.ErrInternalServer)
	}
	if user == nil {
		return nil, errors.New(constants.ErrNotFound)
	}

	return user.ToResponse(), nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, id uint, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by ID")
		return nil, errors.New(constants.ErrInternalServer)
	}
	if user == nil {
		return nil, errors.New(constants.ErrNotFound)
	}

	// Update fields if provided
	if req.Email != nil && *req.Email != user.Email {
		emailExists, err := s.userRepo.EmailExists(ctx, *req.Email)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check email existence")
			return nil, errors.New(constants.ErrInternalServer)
		}
		if emailExists {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}

	if req.Username != nil && *req.Username != user.Username {
		usernameExists, err := s.userRepo.UsernameExists(ctx, *req.Username)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check username existence")
			return nil, errors.New(constants.ErrInternalServer)
		}
		if usernameExists {
			return nil, errors.New("username already exists")
		}
		user.Username = *req.Username
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
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

	user.UpdatedAt = time.Now()

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		return nil, errors.New(constants.ErrInternalServer)
	}

	return user.ToResponse(), nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by ID")
		return errors.New(constants.ErrInternalServer)
	}
	if user == nil {
		return errors.New(constants.ErrNotFound)
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		log.Error().Err(err).Msg("Failed to delete user")
		return errors.New(constants.ErrInternalServer)
	}

	return nil
}

// ListUsers gets a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*model.UserResponse, int64, error) {
	users, total, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list users")
		return nil, 0, errors.New(constants.ErrInternalServer)
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

// ListUsersByRole gets a paginated list of users by role
func (s *UserService) ListUsersByRole(ctx context.Context, role string, limit, offset int) ([]*model.UserResponse, int64, error) {
	users, total, err := s.userRepo.ListByRole(ctx, role, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list users by role")
		return nil, 0, errors.New(constants.ErrInternalServer)
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

// UpdateUserStatus updates user status
func (s *UserService) UpdateUserStatus(ctx context.Context, id uint, status string) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by ID")
		return errors.New(constants.ErrInternalServer)
	}
	if user == nil {
		return errors.New(constants.ErrNotFound)
	}

	// Update status
	if err := s.userRepo.UpdateStatus(ctx, id, status); err != nil {
		log.Error().Err(err).Msg("Failed to update user status")
		return errors.New(constants.ErrInternalServer)
	}

	return nil
}
