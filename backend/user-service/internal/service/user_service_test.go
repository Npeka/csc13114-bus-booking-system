package service

import (
	"context"
	"testing"
	"time"

	"bus-booking/shared/constants"
	"bus-booking/shared/db/mocks"
	storage_mocks "bus-booking/shared/storage/mocks"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/model"
	repo_mocks "bus-booking/user-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	assert.NotNil(t, service)
	assert.IsType(t, &UserServiceImpl{}, service)
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	firebaseUID := "firebase-uid-123"
	req := &model.UserCreateRequest{
		FirebaseUID: firebaseUID,
		Email:       "test@example.com",
		Phone:       "0123456789",
		FullName:    "Test User",
		Avatar:      "https://example.com/avatar.jpg",
		Role:        constants.RolePassenger,
	}

	mockRepo.EXPECT().EmailExists(ctx, req.Email).Return(false, nil).Times(1)
	mockRepo.EXPECT().GetByFirebaseUID(ctx, firebaseUID).Return(nil, assert.AnError).Times(1)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)

	result, err := service.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Email, result.Email)
	assert.Equal(t, req.FullName, result.FullName)
}

func TestCreateUser_EmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	req := &model.UserCreateRequest{
		FirebaseUID: "firebase-uid",
		Email:       "existing@example.com",
		FullName:    "Test User",
	}

	mockRepo.EXPECT().EmailExists(ctx, req.Email).Return(true, nil).Times(1)

	result, err := service.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "email đã tồn tại")
}

func TestCreateUser_FirebaseUIDExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	firebaseUID := "existing-firebase-uid"
	req := &model.UserCreateRequest{
		FirebaseUID: firebaseUID,
		Email:       "test@example.com",
		FullName:    "Test User",
	}

	existingUser := &model.User{BaseModel: model.BaseModel{ID: uuid.New()}}

	mockRepo.EXPECT().EmailExists(ctx, req.Email).Return(false, nil).Times(1)
	mockRepo.EXPECT().GetByFirebaseUID(ctx, firebaseUID).Return(existingUser, nil).Times(1)

	result, err := service.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Firebase UID")
}

func TestGetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()

	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Email:     "test@example.com",
		FullName:  "Test User",
		Role:      constants.RolePassenger,
		Status:    constants.UserStatusActive,
	}

	mockRepo.EXPECT().GetByID(ctx, userID).Return(user, nil).Times(1)

	result, err := service.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestGetUserByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()

	mockRepo.EXPECT().GetByID(ctx, userID).Return(nil, assert.AnError).Times(1)

	result, err := service.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	req := model.UserListQuery{
		PaginationRequest: model.PaginationRequest{Page: 1, PageSize: 10},
	}

	users := []*model.User{
		{BaseModel: model.BaseModel{ID: uuid.New()}, Email: "user1@example.com"},
		{BaseModel: model.BaseModel{ID: uuid.New()}, Email: "user2@example.com"},
	}

	mockRepo.EXPECT().List(ctx, req).Return(users, int64(2), nil).Times(1)

	result, total, err := service.ListUsers(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
}

func TestUpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()

	existingUser := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Email:     "old@example.com",
		FullName:  "Old Name",
	}

	newFullName := "New Name"
	req := &model.UserUpdateRequest{
		FullName: &newFullName,
	}

	mockRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil).Times(1)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Do(func(_ context.Context, u *model.User) {
		assert.Equal(t, "New Name", u.FullName)
	}).Return(nil).Times(1)

	result, err := service.UpdateUser(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.FullName)
}

func TestUpdateUser_EmailChange_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()

	existingUser := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Email:     "old@example.com",
	}

	newEmail := "taken@example.com"
	req := &model.UserUpdateRequest{
		Email: &newEmail,
	}

	mockRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil).Times(1)
	mockRepo.EXPECT().EmailExists(ctx, newEmail).Return(true, nil).Times(1)

	result, err := service.UpdateUser(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "email đã tồn tại")
}

func TestUpdateUser_AllFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()

	existingUser := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Email:     "old@example.com",
	}

	newEmail := "new@example.com"
	newFullName := "New Name"
	newPhone := "0987654321"
	newAvatar := "https://new-avatar.com"
	newRole := constants.RoleAdmin
	newStatus := constants.UserStatusInactive

	req := &model.UserUpdateRequest{
		Email:    &newEmail,
		FullName: &newFullName,
		Phone:    &newPhone,
		Avatar:   &newAvatar,
		Role:     &newRole,
		Status:   &newStatus,
	}

	mockRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil).Times(1)
	mockRepo.EXPECT().EmailExists(ctx, newEmail).Return(false, nil).Times(1)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Do(func(_ context.Context, u *model.User) {
		assert.Equal(t, newEmail, u.Email)
		assert.Equal(t, newFullName, u.FullName)
		assert.Equal(t, newPhone, u.Phone)
		assert.Equal(t, newAvatar, u.Avatar)
		assert.Equal(t, newRole, u.Role)
		assert.Equal(t, newStatus, u.Status)
	}).Return(nil).Times(1)

	result, err := service.UpdateUser(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()

	user := &model.User{BaseModel: model.BaseModel{ID: userID}}

	mockRepo.EXPECT().GetByID(ctx, userID).Return(user, nil).Times(1)
	mockRepo.EXPECT().Delete(ctx, userID).Return(nil).Times(1)

	err := service.DeleteUser(ctx, userID)

	assert.NoError(t, err)
}

func TestDeleteUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()

	mockRepo.EXPECT().GetByID(ctx, userID).Return(nil, assert.AnError).Times(1)

	err := service.DeleteUser(ctx, userID)

	assert.Error(t, err)
}

func TestListUsersByRole_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	role := constants.RoleAdmin
	limit := 10
	offset := 0

	users := []*model.User{
		{BaseModel: model.BaseModel{ID: uuid.New()}, Role: role},
		{BaseModel: model.BaseModel{ID: uuid.New()}, Role: role},
	}

	mockRepo.EXPECT().ListByRole(ctx, role, limit, offset).Return(users, int64(2), nil).Times(1)

	result, total, err := service.ListUsersByRole(ctx, role, limit, offset)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
}

func TestListUsersByRole_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	role := constants.RolePassenger

	mockRepo.EXPECT().ListByRole(ctx, role, 10, 0).Return(nil, int64(0), assert.AnError).Times(1)

	result, total, err := service.ListUsersByRole(ctx, role, 10, 0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
}

func TestCreateUser_EmailCheckError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	req := &model.UserCreateRequest{
		FirebaseUID: "firebase-uid",
		Email:       "test@example.com",
		FullName:    "Test User",
	}

	// EmailExists returns error (DB error)
	mockRepo.EXPECT().
		EmailExists(ctx, req.Email).
		Return(false, assert.AnError).
		Times(1)

	result, err := service.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Không thể xác thực email")
}

func TestUpdateUser_EmailCheckError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()

	existingUser := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Email:     "old@example.com",
	}

	newEmail := "new@example.com"
	req := &model.UserUpdateRequest{
		Email: &newEmail,
	}

	mockRepo.EXPECT().
		GetByID(ctx, userID).
		Return(existingUser, nil).
		Times(1)

	// EmailExists check fails
	mockRepo.EXPECT().
		EmailExists(ctx, newEmail).
		Return(false, assert.AnError).
		Times(1)

	result, err := service.UpdateUser(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Không thể xác thực email")
}

func TestUpdateUser_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()

	existingUser := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Email:     "test@example.com",
	}

	newFullName := "Updated Name"
	req := &model.UserUpdateRequest{
		FullName: &newFullName,
	}

	mockRepo.EXPECT().
		GetByID(ctx, userID).
		Return(existingUser, nil).
		Times(1)

	// Update fails
	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(assert.AnError).
		Times(1)

	result, err := service.UpdateUser(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Không thể cập nhật người dùng")
}

func TestDeleteUser_DeleteFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	userID := uuid.New()

	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
	}

	mockRepo.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil).
		Times(1)

	// Delete operation fails
	mockRepo.EXPECT().
		Delete(ctx, userID).
		Return(assert.AnError).
		Times(1)

	err := service.DeleteUser(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Không thể xóa người dùng")
}

func TestListUsers_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockStorage := storage_mocks.NewMockStorageService(ctrl)
	service := NewUserService(mockRepo, mockStorage)

	ctx := context.Background()
	req := model.UserListQuery{
		PaginationRequest: model.PaginationRequest{Page: 1, PageSize: 10},
	}

	mockRepo.EXPECT().
		List(ctx, req).
		Return(nil, int64(0), assert.AnError).
		Times(1)

	result, total, err := service.ListUsers(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
}

// ============================================
// ADD TO token_manager_test.go
// ============================================

func TestCalculateTokenTTL_ValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)

	cfg := &config.JWTConfig{
		SecretKey:        "test-secret",
		RefreshSecretKey: "test-refresh",
		AccessTokenTTL:   15 * time.Minute,
		RefreshTokenTTL:  24 * time.Hour,
	}
	jwtManager := NewJWTManager(cfg)
	tokenManager := NewTokenManager(mockRedis, jwtManager).(*TokenBlacklistManagerImpl)

	// Generate a real token
	userID := uuid.New()
	validToken, err := jwtManager.GenerateAccessToken(userID, "test@example.com", "2")
	require.NoError(t, err)

	// Calculate TTL
	ttl := tokenManager.calculateTokenTTL(validToken)

	// TTL should be > 0 and roughly 15 minutes + 5 min buffer
	assert.Greater(t, ttl, 10*time.Minute)
	assert.Less(t, ttl, 25*time.Minute)
}

func TestCalculateTokenTTL_InvalidClaims(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	tokenManager := NewTokenManager(mockRedis, jwtManager).(*TokenBlacklistManagerImpl)

	// Token with no standard claims structure - will fail parsing
	malformedToken := "malformed.token.string"

	ttl := tokenManager.calculateTokenTTL(malformedToken)

	// Should return fallback TTL
	assert.Equal(t, 24*time.Hour, ttl)
}

func TestCalculateTokenTTL_ExpiredToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)

	// Create config with very short TTL
	cfg := &config.JWTConfig{
		SecretKey:        "test-secret",
		RefreshSecretKey: "test-refresh",
		AccessTokenTTL:   1 * time.Nanosecond, // Expires immediately
		RefreshTokenTTL:  1 * time.Nanosecond,
	}
	jwtManager := NewJWTManager(cfg)
	tokenManager := NewTokenManager(mockRedis, jwtManager).(*TokenBlacklistManagerImpl)

	userID := uuid.New()
	expiredToken, err := jwtManager.GenerateAccessToken(userID, "test@example.com", "2")
	require.NoError(t, err)

	// Wait to ensure expiry
	time.Sleep(10 * time.Millisecond)

	ttl := tokenManager.calculateTokenTTL(expiredToken)

	// calculateTokenTTL adds 5-minute buffer even for expired tokens
	// So TTL should be around 5 minutes, not 0
	assert.Greater(t, ttl, 4*time.Minute)
	assert.Less(t, ttl, 6*time.Minute)
}
