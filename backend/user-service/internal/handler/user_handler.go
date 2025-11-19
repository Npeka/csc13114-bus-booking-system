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

func (h *UserHandlerImpl) GetProfile(r *ginext.Request) (*ginext.Response, error) {
	userID := context.GetUserID(r.GinCtx)

	user, err := h.us.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user profile")
		return nil, err
	}

	return ginext.NewSuccessResponse(user, "Profile retrieved successfully"), nil
}

func (h *UserHandlerImpl) CreateUser(r *ginext.Request) (*ginext.Response, error) {
	var createReq model.UserCreateRequest
	if err := r.GinCtx.ShouldBind(&createReq); err != nil {
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	user, err := h.us.CreateUser(r.Context(), &createReq)
	if err != nil {
		return nil, err
	}

	return ginext.NewCreatedResponse(user, "User created successfully"), nil
}

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

	return ginext.NewSuccessResponse(user, "User retrieved successfully"), nil
}

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

	return ginext.NewSuccessResponse(user, "User updated successfully"), nil
}

func (h *UserHandlerImpl) DeleteUser(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	if err := h.us.DeleteUser(r.Context(), id); err != nil {
		return nil, err
	}

	return ginext.NewSuccessResponse(nil, "User deleted successfully"), nil
}

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

	return ginext.NewSuccessResponse(result, "Users retrieved successfully"), nil
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

	return ginext.NewSuccessResponse(result, "Users retrieved successfully"), nil
}

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

	return ginext.NewSuccessResponse(nil, "User status updated successfully"), nil
}
