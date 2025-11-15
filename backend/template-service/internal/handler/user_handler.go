package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"bus-booking/shared/constants"
	sharedcontext "bus-booking/shared/context"
	"bus-booking/shared/response"
	"bus-booking/shared/validator"
	"bus-booking/template-service/internal/model"
	"bus-booking/template-service/internal/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService service.UserServiceInterface
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser handles user creation
func (h *UserHandler) CreateUser(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)

	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgFailedToBindRequest)
		response.BadRequestResponse(c, constants.ErrInvalidRequestData)
		return
	}

	// Validate request
	if validationErrors := validator.ValidatorInstance.ValidateStructDetailed(&req); len(validationErrors) > 0 {
		log.Error().Str("request_id", requestID).Interface("errors", validationErrors).Msg(MsgValidationFailed)
		response.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Create user
	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgFailedToCreateUser)

		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			response.ConflictResponse(c, err.Error())
			return
		}

		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	response.CreatedResponse(c, constants.MsgCreatedSuccess, user)
}

// GetUser handles getting user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Str("request_id", requestID).Str("id", idStr).Str("error", err.Error()).Msg(MsgInvalidUserID)
		response.BadRequestResponse(c, MsgInvalidUserID)
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		log.Error().Str("request_id", requestID).Uint64("id", id).Str("error", err.Error()).Msg(MsgFailedToGetUser)

		if err.Error() == constants.ErrNotFound {
			response.NotFoundResponse(c, MsgUserNotFound)
			return
		}

		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	response.SuccessResponse(c, "User retrieved successfully", user)
}

// UpdateUser handles user update
func (h *UserHandler) UpdateUser(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Str("request_id", requestID).Str("id", idStr).Str("error", err.Error()).Msg(MsgInvalidUserID)
		response.BadRequestResponse(c, MsgInvalidUserID)
		return
	}

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgFailedToBindRequest)
		response.BadRequestResponse(c, constants.ErrInvalidRequestData)
		return
	}

	// Validate request
	if validationErrors := validator.ValidatorInstance.ValidateStructDetailed(&req); len(validationErrors) > 0 {
		log.Error().Str("request_id", requestID).Interface("errors", validationErrors).Msg(MsgValidationFailed)
		response.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Update user
	user, err := h.userService.UpdateUser(c.Request.Context(), uint(id), &req)
	if err != nil {
		log.Error().Str("request_id", requestID).Uint64("id", id).Str("error", err.Error()).Msg(MsgFailedToUpdateUser)

		if err.Error() == constants.ErrNotFound {
			response.NotFoundResponse(c, MsgUserNotFound)
			return
		}

		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			response.ConflictResponse(c, err.Error())
			return
		}

		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	response.SuccessResponse(c, constants.MsgUpdatedSuccess, user)
}

// DeleteUser handles user deletion
func (h *UserHandler) DeleteUser(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Str("request_id", requestID).Str("id", idStr).Str("error", err.Error()).Msg(MsgInvalidUserID)
		response.BadRequestResponse(c, MsgInvalidUserID)
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), uint(id))
	if err != nil {
		log.Error().Str("request_id", requestID).Uint64("id", id).Str("error", err.Error()).Msg(MsgFailedToDeleteUser)

		if err.Error() == constants.ErrNotFound {
			response.NotFoundResponse(c, MsgUserNotFound)
			return
		}

		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	response.SuccessResponse(c, constants.MsgDeletedSuccess, nil)
}

// ListUsers handles getting paginated list of users
func (h *UserHandler) ListUsers(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)

	// Parse pagination parameters
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get users
	users, total, err := h.userService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgFailedToListUsers)
		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	pagination := &response.Pagination{
		Page:       offset/limit + 1,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}
	response.PaginatedResponse(c, "Users retrieved successfully", users, pagination)
}

// UpdateUserStatus handles updating user status
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Str("request_id", requestID).Str("id", idStr).Str("error", err.Error()).Msg(MsgInvalidUserID)
		response.BadRequestResponse(c, MsgInvalidUserID)
		return
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=active inactive suspended"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgFailedToBindRequest)
		response.BadRequestResponse(c, constants.ErrInvalidRequestData)
		return
	}

	// Validate request
	if validationErrors := validator.ValidatorInstance.ValidateStructDetailed(&req); len(validationErrors) > 0 {
		log.Error().Str("request_id", requestID).Interface("errors", validationErrors).Msg(MsgValidationFailed)
		response.ValidationErrorResponse(c, validationErrors)
		return
	}

	err = h.userService.UpdateUserStatus(c.Request.Context(), uint(id), req.Status)
	if err != nil {
		log.Error().Str("request_id", requestID).Uint64("id", id).Str("error", err.Error()).Msg(MsgFailedToUpdateStatus)

		if err.Error() == constants.ErrNotFound {
			response.NotFoundResponse(c, MsgUserNotFound)
			return
		}

		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	response.SuccessResponse(c, constants.MsgUpdatedSuccess, nil)
}
