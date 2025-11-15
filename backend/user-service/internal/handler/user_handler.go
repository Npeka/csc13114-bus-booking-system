package handler

import (
	"strconv"

	"bus-booking/shared/ginext"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/service"

	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(r *ginext.Request) (*ginext.Response, error) {
	var createReq model.UserCreateRequest
	r.MustBind(&createReq)

	user, err := h.userService.CreateUser(r.Context(), &createReq)
	if err != nil {
		return nil, err
	}

	return ginext.NewCreatedResponse(user, "User created successfully"), nil
}

func (h *UserHandler) GetUser(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		return nil, err
	}

	return ginext.NewSuccessResponse(user, "User retrieved successfully"), nil
}

func (h *UserHandler) UpdateUser(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	var updateReq model.UserUpdateRequest
	r.MustBind(&updateReq)

	user, err := h.userService.UpdateUser(r.Context(), id, &updateReq)
	if err != nil {
		return nil, err
	}

	return ginext.NewSuccessResponse(user, "User updated successfully"), nil
}

func (h *UserHandler) DeleteUser(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	if err := h.userService.DeleteUser(r.Context(), id); err != nil {
		return nil, err
	}

	return ginext.NewSuccessResponse(nil, "User deleted successfully"), nil
}

func (h *UserHandler) ListUsers(r *ginext.Request) (*ginext.Response, error) {
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

	users, total, err := h.userService.ListUsers(r.Context(), limit, offset)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	return ginext.NewSuccessResponse(result, "Users retrieved successfully"), nil
}

func (h *UserHandler) ListUsersByRole(r *ginext.Request) (*ginext.Response, error) {
	role := r.Param("role")
	if role == "" {
		return nil, ginext.NewBadRequestError("role parameter is required")
	}

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

	users, total, err := h.userService.ListUsersByRole(r.Context(), role, limit, offset)
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

func (h *UserHandler) UpdateUserStatus(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	var body map[string]string
	r.MustBind(&body)

	status, ok := body["status"]
	if !ok || status == "" {
		return nil, ginext.NewBadRequestError("status is required")
	}

	if err := h.userService.UpdateUserStatus(r.Context(), id, status); err != nil {
		return nil, err
	}

	return ginext.NewSuccessResponse(nil, "User status updated successfully"), nil
}
