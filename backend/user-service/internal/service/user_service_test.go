package service

import (
	"context"
	"errors"
	"testing"

	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/service/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	req := &model.UserCreateRequest{
		FirebaseUID: "firebase123",
		Email:       "test@example.com",
		Phone:       "1234567890",
		FullName:    "Test User",
		Avatar:      "https://example.com/avatar.jpg",
		Role:        constants.RolePassenger,
	}

	mockRepo.On("EmailExists", ctx, req.Email).Return(false, nil)
	mockRepo.On("GetByFirebaseUID", ctx, req.FirebaseUID).Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Email, result.Email)
	assert.Equal(t, req.FullName, result.FullName)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_EmailExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	req := &model.UserCreateRequest{
		FirebaseUID: "firebase123",
		Email:       "existing@example.com",
		FullName:    "Test User",
		Role:        constants.RolePassenger,
	}

	mockRepo.On("EmailExists", ctx, req.Email).Return(true, nil)

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "email already exists")
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_FirebaseUIDExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	req := &model.UserCreateRequest{
		FirebaseUID: "firebase123",
		Email:       "test@example.com",
		FullName:    "Test User",
		Role:        constants.RolePassenger,
	}

	existingUser := &model.User{
		ID:          uuid.New(),
		FirebaseUID: req.FirebaseUID,
	}

	mockRepo.On("EmailExists", ctx, req.Email).Return(false, nil)
	mockRepo.On("GetByFirebaseUID", ctx, req.FirebaseUID).Return(existingUser, nil)

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Firebase UID already exists")
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	req := &model.UserCreateRequest{
		FirebaseUID: "firebase123",
		Email:       "test@example.com",
		FullName:    "Test User",
		Role:        constants.RolePassenger,
	}

	mockRepo.On("EmailExists", ctx, req.Email).Return(false, nil)
	mockRepo.On("GetByFirebaseUID", ctx, req.FirebaseUID).Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(errors.New("database error"))

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	expectedUser := &model.User{
		ID:       userID,
		Email:    "test@example.com",
		FullName: "Test User",
		Role:     constants.RolePassenger,
		Status:   "active",
	}

	mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

	// Act
	result, err := service.GetUserByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	mockRepo.On("GetByID", ctx, userID).Return(nil, errors.New("user not found"))

	// Act
	result, err := service.GetUserByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	newEmail := "newemail@example.com"
	newFullName := "Updated Name"

	existingUser := &model.User{
		ID:       userID,
		Email:    "old@example.com",
		FullName: "Old Name",
		Role:     constants.RolePassenger,
		Status:   "active",
	}

	req := &model.UserUpdateRequest{
		Email:    &newEmail,
		FullName: &newFullName,
	}

	mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("EmailExists", ctx, newEmail).Return(false, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	// Act
	result, err := service.UpdateUser(ctx, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newEmail, result.Email)
	assert.Equal(t, newFullName, result.FullName)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_EmailConflict(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	newEmail := "existing@example.com"

	existingUser := &model.User{
		ID:    userID,
		Email: "old@example.com",
		Role:  constants.RolePassenger,
	}

	req := &model.UserUpdateRequest{
		Email: &newEmail,
	}

	mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("EmailExists", ctx, newEmail).Return(true, nil)

	// Act
	result, err := service.UpdateUser(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "email already exists")
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	newEmail := "new@example.com"

	req := &model.UserUpdateRequest{
		Email: &newEmail,
	}

	mockRepo.On("GetByID", ctx, userID).Return(nil, ginext.NewNotFoundError("user not found"))

	// Act
	result, err := service.UpdateUser(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	existingUser := &model.User{
		ID:    userID,
		Email: "test@example.com",
	}

	mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("Delete", ctx, userID).Return(nil)

	// Act
	err := service.DeleteUser(ctx, userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	mockRepo.On("GetByID", ctx, userID).Return(nil, ginext.NewNotFoundError("user not found"))

	// Act
	err := service.DeleteUser(ctx, userID)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsers_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	users := []*model.User{
		{ID: uuid.New(), Email: "user1@example.com", FullName: "User 1"},
		{ID: uuid.New(), Email: "user2@example.com", FullName: "User 2"},
	}

	mockRepo.On("List", ctx, 10, 0).Return(users, int64(2), nil)

	// Act
	result, total, err := service.ListUsers(ctx, 10, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsersByRole_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	users := []*model.User{
		{ID: uuid.New(), Email: "driver1@example.com", Role: constants.RolePassenger},
		{ID: uuid.New(), Email: "driver2@example.com", Role: constants.RolePassenger},
	}

	mockRepo.On("ListByRole", ctx, constants.RolePassenger, 10, 0).Return(users, int64(2), nil)

	// Act
	result, total, err := service.ListUsersByRole(ctx, constants.RolePassenger, 10, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserStatus_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	existingUser := &model.User{
		ID:     userID,
		Status: "active",
	}

	mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("UpdateStatus", ctx, userID, "suspended").Return(nil)

	// Act
	err := service.UpdateUserStatus(ctx, userID, "suspended")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserStatus_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	mockRepo.On("GetByID", ctx, userID).Return(nil, ginext.NewNotFoundError("user not found"))

	// Act
	err := service.UpdateUserStatus(ctx, userID, "suspended")

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
