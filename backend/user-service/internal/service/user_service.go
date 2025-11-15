package service

import (
	"context"
	"strconv"

	"bus-booking/shared/ginext"
	contextlogger "bus-booking/shared/logger"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/repository"
	"bus-booking/user-service/internal/utils"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// UserServiceInterface defines the interface for user service
type UserServiceInterface interface {
	CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *model.UserUpdateRequest) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, limit, offset int) ([]*model.UserResponse, int64, error)
	ListUsersByRole(ctx context.Context, role string, limit, offset int) ([]*model.UserResponse, int64, error)
	UpdateUserStatus(ctx context.Context, id uuid.UUID, status string) error
}

// UserService implements UserServiceInterface
type UserService struct {
	userRepo repository.UserRepositoryInterface
	logger   *contextlogger.ContextLogger
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
		logger:   contextlogger.NewContextLogger("user-service", "UserService", ""),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.UserResponse, error) {
	if emailExists, err := s.userRepo.EmailExists(ctx, req.Email); err != nil {
		s.logger.Error(err, "Failed to check email existence")
		return nil, ginext.NewInternalServerError("Failed to validate email")
	} else if emailExists {
		return nil, ginext.NewConflictError("email already exists")
	}

	if usernameExists, err := s.userRepo.UsernameExists(ctx, req.Username); err != nil {
		s.logger.Error(err, "Failed to check username existence")
		return nil, ginext.NewInternalServerError("Failed to validate username")
	} else if usernameExists {
		return nil, ginext.NewConflictError("username already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Error(err, "Failed to hash password")
		return nil, ginext.NewInternalServerError("Failed to process password")
	}

	// Set default role if not provided
	role := req.Role
	if role == 0 {
		role = model.RolePassenger
	}

	user := &model.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      role,
		Status:    "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error(err, "Failed to create user in database")
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate and update email if provided
	if req.Email != nil && *req.Email != user.Email {
		if emailExists, err := s.userRepo.EmailExists(ctx, *req.Email); err != nil {
			s.logger.Error(err, "Failed to check email existence during update")
			return nil, ginext.NewInternalServerError("Failed to validate email")
		} else if emailExists {
			return nil, ginext.NewConflictError("email already exists")
		}
		user.Email = *req.Email
	}

	// Validate and update username if provided
	if req.Username != nil && *req.Username != user.Username {
		if usernameExists, err := s.userRepo.UsernameExists(ctx, *req.Username); err != nil {
			s.logger.Error(err, "Failed to check username existence during update")
			return nil, ginext.NewInternalServerError("Failed to validate username")
		} else if usernameExists {
			return nil, ginext.NewConflictError("username already exists")
		}
		user.Username = *req.Username
	}

	// Update other fields
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

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error(err, "Failed to update user in database")
		return nil, ginext.NewInternalServerError("Failed to update user")
	}

	return user.ToResponse(), nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Check if user exists
	if _, err := s.userRepo.GetByID(ctx, id); err != nil {
		return err // Repository already returns proper ginext errors
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.logger.Error(err, "Failed to delete user from database")
		return ginext.NewInternalServerError("Failed to delete user")
	}

	return nil
}

// ListUsers gets a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*model.UserResponse, int64, error) {
	users, total, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, ginext.NewInternalServerError("Failed to list users")
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

// ListUsersByRole gets a paginated list of users by role
func (s *UserService) ListUsersByRole(ctx context.Context, role string, limit, offset int) ([]*model.UserResponse, int64, error) {
	// Convert string role to UserRole
	var userRole model.UserRole
	if roleValue, err := strconv.Atoi(role); err == nil {
		userRole = model.UserRole(roleValue)
	} else {
		switch role {
		case "passenger":
			userRole = model.RolePassenger
		case "admin":
			userRole = model.RoleAdmin
		case "operator":
			userRole = model.RoleOperator
		case "support":
			userRole = model.RoleSupport
		default:
			return nil, 0, ginext.NewBadRequestError("invalid role")
		}
	}

	users, total, err := s.userRepo.ListByRole(ctx, userRole, limit, offset)
	if err != nil {
		return nil, 0, ginext.NewInternalServerError("Failed to list users by role")
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

// UpdateUserStatus updates user status
func (s *UserService) UpdateUserStatus(ctx context.Context, id uuid.UUID, status string) error {
	if _, err := s.userRepo.GetByID(ctx, id); err != nil {
		return err
	}

	if err := s.userRepo.UpdateStatus(ctx, id, status); err != nil {
		log.Err(err).Msg("Failed to update user status in database")
		return err
	}

	return nil
}
