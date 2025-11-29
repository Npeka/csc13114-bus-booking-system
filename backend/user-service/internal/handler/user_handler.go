package handler

import (
	"strconv"

	"bus-booking/shared/constants"
	"bus-booking/shared/context"
	"bus-booking/shared/ginext"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/service"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserHandler interface {
	GetProfile(r *ginext.Request) (*ginext.Response, error)
	CreateUser(r *ginext.Request) (*ginext.Response, error)
	GetUser(r *ginext.Request) (*ginext.Response, error)
	UpdateUser(r *ginext.Request) (*ginext.Response, error)
	DeleteUser(r *ginext.Request) (*ginext.Response, error)
	ListUsers(r *ginext.Request) (*ginext.Response, error)
	UpdateUserStatus(r *ginext.Request) (*ginext.Response, error)
}

type UserHandlerImpl struct {
	us service.UserService
}

func NewUserHandler(us service.UserService) UserHandler {
	return &UserHandlerImpl{
		us: us,
	}
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Retrieves the profile information of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ginext.Response{data=model.UserResponse} "Profile retrieved successfully"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /users/profile [get]
func (h *UserHandlerImpl) GetProfile(r *ginext.Request) (*ginext.Response, error) {
	userID := context.GetUserID(r.GinCtx)

	user, err := h.us.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user profile")
		return nil, err
	}

	return ginext.NewSuccessResponse(user), nil
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UserCreateRequest true "User creation request"
// @Success 201 {object} ginext.Response{data=model.UserResponse} "User created successfully"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 403 {object} ginext.Response "Forbidden"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /users [post]
func (h *UserHandlerImpl) CreateUser(r *ginext.Request) (*ginext.Response, error) {
	var createReq model.UserCreateRequest
	if err := r.GinCtx.ShouldBind(&createReq); err != nil {
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	user, err := h.us.CreateUser(r.Context(), &createReq)
	if err != nil {
		return nil, err
	}

	return ginext.NewCreatedResponse(user), nil
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieves a user's information by their ID (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} ginext.Response{data=model.UserResponse} "User retrieved successfully"
// @Failure 400 {object} ginext.Response "Invalid user ID"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 403 {object} ginext.Response "Forbidden"
// @Failure 404 {object} ginext.Response "User not found"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandlerImpl) GetUser(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	user, err := h.us.GetUserByID(r.Context(), id)
	if err != nil {
		return nil, err
	}

	return ginext.NewSuccessResponse(user), nil
}

// UpdateUser godoc
// @Summary Update user information
// @Description Updates a user's information (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Param request body model.UserUpdateRequest true "User update request"
// @Success 200 {object} ginext.Response{data=model.UserResponse} "User updated successfully"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 403 {object} ginext.Response "Forbidden"
// @Failure 404 {object} ginext.Response "User not found"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /users/{id} [put]
func (h *UserHandlerImpl) UpdateUser(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	var updateReq model.UserUpdateRequest
	if err := r.GinCtx.ShouldBind(&updateReq); err != nil {
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	user, err := h.us.UpdateUser(r.Context(), id, &updateReq)
	if err != nil {
		return nil, err
	}

	return ginext.NewSuccessResponse(user), nil
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes a user account (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} ginext.Response "User deleted successfully"
// @Failure 400 {object} ginext.Response "Invalid user ID"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 403 {object} ginext.Response "Forbidden"
// @Failure 404 {object} ginext.Response "User not found"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /users/{id} [delete]
func (h *UserHandlerImpl) DeleteUser(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	if err := h.us.DeleteUser(r.Context(), id); err != nil {
		return nil, err
	}

	return ginext.NewSuccessResponse("User deleted successfully"), nil
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieves a paginated list of users with optional filtering (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param search query string false "Search by name or email"
// @Param role query int false "Filter by role (1=Passenger, 2=Driver, 3=Admin)"
// @Param status query string false "Filter by status (active, suspended, etc.)"
// @Success 200 {object} ginext.Response{data=object} "Users retrieved successfully"
// @Failure 400 {object} ginext.Response "Invalid query parameters"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 403 {object} ginext.Response "Forbidden"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /users [get]
func (h *UserHandlerImpl) ListUsers(r *ginext.Request) (*ginext.Response, error) {
	var query model.UserListQuery
	if err := r.GinCtx.ShouldBindQuery(&query); err != nil {
		return nil, ginext.NewBadRequestError("Invalid query parameters")
	}

	// Set defaults
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}

	// Calculate offset from page
	offset := (query.Page - 1) * query.Limit

	users, total, err := h.us.ListUsers(r.Context(), query.Limit, offset)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"users":       users,
		"total":       total,
		"page":        query.Page,
		"limit":       query.Limit,
		"total_pages": (total + int64(query.Limit) - 1) / int64(query.Limit),
		"search":      query.Search,
		"role":        query.Role,
		"status":      query.Status,
	}

	return ginext.NewSuccessResponse(result), nil
}

func (h *UserHandlerImpl) ListUsersByRole(r *ginext.Request) (*ginext.Response, error) {
	roleParam := r.Param("role")
	if roleParam == "" {
		return nil, ginext.NewBadRequestError("role parameter is required")
	}

	roleInt, err := strconv.Atoi(roleParam)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid role parameter")
	}

	role := constants.UserRole(roleInt)

	limitStr := r.DefaultQuery("limit", "20")
	offsetStr := r.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return nil, ginext.NewBadRequestError("invalid limit parameter")
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return nil, ginext.NewBadRequestError("invalid offset parameter")
	}

	users, total, err := h.us.ListUsersByRole(r.Context(), role, limit, offset)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"users":  users,
		"total":  total,
		"role":   role,
		"limit":  limit,
		"offset": offset,
	}

	return ginext.NewSuccessResponse(result), nil
}

// UpdateUserStatus godoc
// @Summary Update user status
// @Description Updates a user's status (e.g., active, suspended) (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Param request body model.UserStatusUpdateRequest true "Status update request"
// @Success 200 {object} ginext.Response "User status updated successfully"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 403 {object} ginext.Response "Forbidden"
// @Failure 404 {object} ginext.Response "User not found"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /users/{id}/status [patch]
func (h *UserHandlerImpl) UpdateUserStatus(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	var statusReq model.UserStatusUpdateRequest
	if err := r.GinCtx.ShouldBind(&statusReq); err != nil {
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	if err := h.us.UpdateUserStatus(r.Context(), id, statusReq.Status); err != nil {
		return nil, err
	}

	return ginext.NewSuccessResponse("User status updated successfully"), nil
}
