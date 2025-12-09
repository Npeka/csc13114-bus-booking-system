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
	UpdateProfile(r *ginext.Request) (*ginext.Response, error)

	GetUser(r *ginext.Request) (*ginext.Response, error)
	ListUsers(r *ginext.Request) (*ginext.Response, error)
	CreateUser(r *ginext.Request) (*ginext.Response, error)
	UpdateUser(r *ginext.Request) (*ginext.Response, error)
	DeleteUser(r *ginext.Request) (*ginext.Response, error)
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

// UpdateProfile godoc
// @Summary Update current user profile
// @Description Updates the profile information of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UserUpdateRequest true "Profile update request"
// @Success 200 {object} ginext.Response{data=model.UserResponse} "Profile updated successfully"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /users/profile [put]
func (h *UserHandlerImpl) UpdateProfile(r *ginext.Request) (*ginext.Response, error) {
	userID := context.GetUserID(r.GinCtx)

	var req model.UserUpdateRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	user, err := h.us.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to update profile")
		return nil, err
	}

	return ginext.NewSuccessResponse(user), nil
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
		log.Error().Err(err).Msg("Invalid user ID")
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	user, err := h.us.GetUserByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("user_id", idStr).Msg("Failed to get user")
		return nil, err
	}

	return ginext.NewSuccessResponse(user), nil
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
// @Success 200 {object} ginext.Response "Paginated users list"
// @Failure 400 {object} ginext.Response "Invalid query parameters"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 403 {object} ginext.Response "Forbidden"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /users [get]
func (h *UserHandlerImpl) ListUsers(r *ginext.Request) (*ginext.Response, error) {
	var req model.UserListQuery
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("Query binding failed")
		return nil, ginext.NewBadRequestError("Dữ liệu truy vấn không hợp lệ")
	}

	req.Normalize()

	users, total, err := h.us.ListUsers(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list users")
		return nil, err
	}

	return ginext.NewPaginatedResponse(users, req.Page, req.PageSize, total), nil
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
	var req model.UserCreateRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	user, err := h.us.CreateUser(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	return ginext.NewCreatedResponse(user), nil
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
		log.Error().Err(err).Msg("Invalid user ID")
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	var req model.UserUpdateRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	user, err := h.us.UpdateUser(r.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("user_id", idStr).Msg("Failed to update user")
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
		log.Error().Err(err).Msg("Invalid user ID")
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	if err = h.us.DeleteUser(r.Context(), id); err != nil {
		log.Error().Err(err).Str("user_id", idStr).Msg("Failed to delete user")
		return nil, err
	}

	return ginext.NewSuccessResponse("Xóa người dùng thành công"), nil
}

func (h *UserHandlerImpl) ListUsersByRole(r *ginext.Request) (*ginext.Response, error) {
	roleParam := r.Param("role")
	if roleParam == "" {
		return nil, ginext.NewBadRequestError("tham số vai trò là bắt buộc")
	}

	roleInt, err := strconv.Atoi(roleParam)
	if err != nil {
		return nil, ginext.NewBadRequestError("tham số vai trò không hợp lệ")
	}

	role := constants.UserRole(roleInt)

	limitStr := r.DefaultQuery("limit", "20")
	offsetStr := r.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return nil, ginext.NewBadRequestError("tham số limit không hợp lệ")
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return nil, ginext.NewBadRequestError("tham số offset không hợp lệ")
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
