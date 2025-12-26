package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"bus-booking/shared/constants"
	db_mocks "bus-booking/shared/db/mocks"
	"bus-booking/shared/utils"
	"bus-booking/user-service/config"
	client_mocks "bus-booking/user-service/internal/client/mocks"
	"bus-booking/user-service/internal/model"
	repo_mocks "bus-booking/user-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthService(t *testing.T) (
	*AuthServiceImpl,
	*gomock.Controller,
	*repo_mocks.MockUserRepository,
	*db_mocks.MockRedisManager,
	*client_mocks.MockNotificationClient,
) {
	ctrl := gomock.NewController(t)

	mockUserRepo := repo_mocks.NewMockUserRepository(ctrl)
	mockRedis := db_mocks.NewMockRedisManager(ctrl)
	mockNotification := client_mocks.NewMockNotificationClient(ctrl)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:        "test-secret",
			RefreshSecretKey: "test-refresh-secret",
			AccessTokenTTL:   15 * time.Minute,
			RefreshTokenTTL:  24 * time.Hour,
			Issuer:           "test-issuer",
			Audience:         "test-audience",
		},
	}

	jwtManager := NewJWTManager(&cfg.JWT)
	tokenManager := NewTokenManager(mockRedis, jwtManager)
	firebaseAuth := NewFirebaseAuth(nil) // nil client for testing

	service := NewAuthService(
		cfg,
		jwtManager,
		firebaseAuth,
		tokenManager,
		mockUserRepo,
		mockRedis,
		mockNotification,
	).(*AuthServiceImpl)

	return service, ctrl, mockUserRepo, mockRedis, mockNotification
}

func TestNewAuthService(t *testing.T) {
	service, ctrl, _, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	assert.NotNil(t, service)
	assert.IsType(t, &AuthServiceImpl{}, service)
}

func TestRegister_Success(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
	}

	// Mock: email doesn't exist
	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(nil, assert.AnError).
		Times(1)

	// Mock: create user
	mockUserRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, user *model.User) error {
			// Verify password was hashed
			assert.NotNil(t, user.PasswordHash)
			assert.NotEqual(t, req.Password, *user.PasswordHash)
			assert.Equal(t, req.Email, user.Email)
			assert.Equal(t, req.FullName, user.FullName)
			assert.Equal(t, constants.RolePassenger, user.Role)
			return nil
		}).
		Times(1)

	result, err := service.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, req.Email, result.User.Email)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	existingUser := &model.User{
		BaseModel: model.BaseModel{ID: uuid.New()},
		Email:     req.Email,
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(existingUser, nil).
		Times(1)

	result, err := service.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Email đã được đăng ký")
}

func TestRegister_CreateUserFails(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(nil, assert.AnError).
		Times(1)

	mockUserRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(assert.AnError).
		Times(1)

	result, err := service.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Không thể tạo tài khoản")
}

func TestLogin_Success(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	password := "testpassword123"
	passwordHash, _ := utils.HashPassword(password)

	user := &model.User{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		Email:        "test@example.com",
		FullName:     "Test User",
		PasswordHash: &passwordHash,
		Role:         constants.RolePassenger,
		Status:       constants.UserStatusActive,
	}

	req := &model.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(user, nil).
		Times(1)

	result, err := service.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, user.Email, result.User.Email)
}

func TestLogin_UserNotFound(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Email hoặc mật khẩu không đúng")
}

func TestLogin_WrongPassword(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	correctPassword := "correct123"
	passwordHash, _ := utils.HashPassword(correctPassword)

	user := &model.User{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		Email:        "test@example.com",
		PasswordHash: &passwordHash,
		Status:       constants.UserStatusActive,
	}

	req := &model.LoginRequest{
		Email:    user.Email,
		Password: "wrongpassword",
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(user, nil).
		Times(1)

	result, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Email hoặc mật khẩu không đúng")
}

func TestLogin_NoPasswordSet(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// User without password (Firebase-only)
	user := &model.User{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		Email:        "firebase@example.com",
		PasswordHash: nil, // No password
		Status:       constants.UserStatusActive,
	}

	req := &model.LoginRequest{
		Email:    user.Email,
		Password: "anypassword",
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(user, nil).
		Times(1)

	result, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Email hoặc mật khẩu không đúng")
}

func TestLogin_InactiveUser(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	password := "password123"
	passwordHash, _ := utils.HashPassword(password)

	user := &model.User{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		Email:        "inactive@example.com",
		PasswordHash: &passwordHash,
		Status:       constants.UserStatusInactive,
	}

	req := &model.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(user, nil).
		Times(1)

	result, err := service.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Tài khoản không hoạt động")
}

func TestVerifyToken_Success(t *testing.T) {
	service, ctrl, mockUserRepo, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()

	// Generate a valid token
	accessToken, err := service.jwtManager.GenerateAccessToken(
		userID,
		"test@example.com",
		fmt.Sprintf("%d", constants.RolePassenger),
	)
	require.NoError(t, err)

	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Email:     "test@example.com",
		FullName:  "Test User",
		Role:      constants.RolePassenger,
		Status:    constants.UserStatusActive,
	}

	// Mock token not blacklisted
	mockRedis.EXPECT().
		Exists(ctx, gomock.Any()).
		Return(int64(0), nil).
		Times(1)

	// Mock user blacklist check
	mockRedis.EXPECT().
		Get(ctx, gomock.Any()).
		Return("", assert.AnError). // No user blacklist
		Times(1)

	mockUserRepo.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil).
		Times(1)

	result, err := service.VerifyToken(ctx, accessToken)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID.String(), result.UserID)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.Role, result.Role)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	service, ctrl, _, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	result, err := service.VerifyToken(ctx, "invalid.token.here")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "token không hợp lệ")
}

func TestVerifyToken_BlacklistedToken(t *testing.T) {
	service, ctrl, _, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()

	accessToken, err := service.jwtManager.GenerateAccessToken(
		userID,
		"test@example.com",
		fmt.Sprintf("%d", constants.RolePassenger),
	)
	require.NoError(t, err)

	// Mock token IS blacklisted
	mockRedis.EXPECT().
		Exists(ctx, gomock.Any()).
		Return(int64(1), nil). // Token blacklisted
		Times(1)

	result, err := service.VerifyToken(ctx, accessToken)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "token đã bị blacklisted")
}

func TestVerifyToken_UserNotFound(t *testing.T) {
	service, ctrl, mockUserRepo, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()

	accessToken, err := service.jwtManager.GenerateAccessToken(
		userID,
		"test@example.com",
		fmt.Sprintf("%d", constants.RolePassenger),
	)
	require.NoError(t, err)

	// Mock token not blacklisted
	mockRedis.EXPECT().
		Exists(ctx, gomock.Any()).
		Return(int64(0), nil).
		Times(1)

	// Mock user blacklist check
	mockRedis.EXPECT().
		Get(ctx, gomock.Any()).
		Return("", assert.AnError).
		Times(1)

	// User not found
	mockUserRepo.EXPECT().
		GetByID(ctx, userID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.VerifyToken(ctx, accessToken)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "không tìm thấy người dùng")
}

func TestVerifyToken_InactiveUser(t *testing.T) {
	service, ctrl, mockUserRepo, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()

	accessToken, err := service.jwtManager.GenerateAccessToken(
		userID,
		"test@example.com",
		fmt.Sprintf("%d", constants.RolePassenger),
	)
	require.NoError(t, err)

	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Email:     "test@example.com",
		Status:    constants.UserStatusInactive, // Inactive!
	}

	mockRedis.EXPECT().
		Exists(ctx, gomock.Any()).
		Return(int64(0), nil).
		Times(1)

	mockRedis.EXPECT().
		Get(ctx, gomock.Any()).
		Return("", assert.AnError).
		Times(1)

	mockUserRepo.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil).
		Times(1)

	result, err := service.VerifyToken(ctx, accessToken)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "tài khoản không hoạt động")
}

func TestRefreshToken_Success(t *testing.T) {
	service, ctrl, mockUserRepo, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()

	refreshToken, err := service.jwtManager.GenerateRefreshToken(
		userID,
		"test@example.com",
		fmt.Sprintf("%d", constants.RolePassenger),
	)
	require.NoError(t, err)

	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Email:     "test@example.com",
		FullName:  "Test User",
		Role:      constants.RolePassenger,
		Status:    constants.UserStatusActive,
	}

	req := &model.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	// Mock token not blacklisted
	mockRedis.EXPECT().
		Exists(ctx, gomock.Any()).
		Return(int64(0), nil).
		Times(1)

	// Mock user blacklist check
	mockRedis.EXPECT().
		Get(ctx, gomock.Any()).
		Return("", assert.AnError).
		Times(1)

	mockUserRepo.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil).
		Times(1)

	// RefreshToken calls Blacklist on old refresh token
	mockRedis.EXPECT().
		Set(ctx, gomock.Any(), "1", gomock.Any()).
		Return(nil).
		Times(1)

	result, err := service.RefreshToken(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	service, ctrl, _, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.RefreshTokenRequest{
		RefreshToken: "invalid.refresh.token",
	}

	result, err := service.RefreshToken(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	service, ctrl, mockUserRepo, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()

	refreshToken, err := service.jwtManager.GenerateRefreshToken(
		userID,
		"test@example.com",
		fmt.Sprintf("%d", constants.RolePassenger),
	)
	require.NoError(t, err)

	req := &model.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	mockRedis.EXPECT().
		Exists(ctx, gomock.Any()).
		Return(int64(0), nil).
		Times(1)

	mockRedis.EXPECT().
		Get(ctx, gomock.Any()).
		Return("", assert.AnError).
		Times(1)

	// User not found
	mockUserRepo.EXPECT().
		GetByID(ctx, userID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.RefreshToken(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestLogout_Success(t *testing.T) {
	service, ctrl, _, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()

	accessToken, _ := service.jwtManager.GenerateAccessToken(
		userID,
		"test@example.com",
		fmt.Sprintf("%d", constants.RolePassenger),
	)

	// Need refresh token for Logout
	refreshToken, _ := service.jwtManager.GenerateRefreshToken(
		userID,
		"test@example.com",
		fmt.Sprintf("%d", constants.RolePassenger),
	)

	req := model.LogoutRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Logout blacklists tokens in goroutine - set up expectations
	// Use AnyTimes() to handle async goroutine
	mockRedis.EXPECT().
		Set(gomock.Any(), gomock.Any(), "1", gomock.Any()).
		Return(nil).
		AnyTimes()

	err := service.Logout(ctx, req, userID)

	assert.NoError(t, err)

	// Give goroutine time to complete
	time.Sleep(50 * time.Millisecond)
}

func TestCreateGuestAccount_Success(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.CreateGuestAccountRequest{
		Phone:    "0123456789",
		FullName: "Guest User",
	}

	// Implementation calls GetByPhone, not EmailExists
	mockUserRepo.EXPECT().
		GetByPhone(ctx, req.Phone).
		Return(nil, nil). // User doesn't exist - return nil, nil
		Times(1)

	mockUserRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, user *model.User) error {
			assert.Equal(t, req.Phone, user.Phone)
			assert.Equal(t, req.FullName, user.FullName)
			assert.Equal(t, constants.RoleGuest, user.Role)
			assert.Equal(t, constants.UserStatusActive, user.Status)
			return nil
		}).
		Times(1)

	result, err := service.CreateGuestAccount(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Phone, result.Phone)
	assert.Equal(t, req.FullName, result.FullName)
}

func TestCreateGuestAccount_WithEmail(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.CreateGuestAccountRequest{
		Email:    "guest@example.com",
		Phone:    "0987654321",
		FullName: "Guest User",
	}

	// Implementation calls GetByEmail, not EmailExists
	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(nil, nil). // User doesn't exist
		Times(1)

	// Then checks GetByPhone
	mockUserRepo.EXPECT().
		GetByPhone(ctx, req.Phone).
		Return(nil, nil). // User doesn't exist
		Times(1)

	mockUserRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil).
		Times(1)

	result, err := service.CreateGuestAccount(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreateGuestAccount_EmailExists(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.CreateGuestAccountRequest{
		Email:    "existing@example.com",
		Phone:    "0987654321",
		FullName: "Guest User",
	}

	// Implementation calls GetByEmail and returns existing user
	existingUser := &model.User{
		BaseModel: model.BaseModel{ID: uuid.New()},
		Email:     req.Email,
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(existingUser, nil).
		Times(1)

	result, err := service.CreateGuestAccount(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Email đã được đăng ký")
}

func TestCreateGuestAccount_PhoneExists(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.CreateGuestAccountRequest{
		Email:    "new@example.com",
		Phone:    "existing-phone",
		FullName: "Guest User",
	}

	existingUser := &model.User{
		BaseModel: model.BaseModel{ID: uuid.New()},
		Phone:     req.Phone,
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(nil, nil).
		Times(1)

	mockUserRepo.EXPECT().
		GetByPhone(ctx, req.Phone).
		Return(existingUser, nil). // Phone exists
		Times(1)

	result, err := service.CreateGuestAccount(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Số điện thoại đã được đăng ký")
}

func TestCreateGuestAccount_NoContactMethod(t *testing.T) {
	service, ctrl, _, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.CreateGuestAccountRequest{
		Email:    "", // No email
		Phone:    "", // No phone
		FullName: "Guest",
	}

	result, err := service.CreateGuestAccount(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Phải cung cấp email hoặc số điện thoại")
}

func TestForgotPassword_UserNotFound(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.ForgotPasswordRequest{
		Email: "nonexistent@example.com",
	}

	// User doesn't exist
	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(nil, assert.AnError).
		Times(1)

	err := service.ForgotPassword(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Không thể xử lý yêu cầu đặt lại mật khẩu")
}

func TestForgotPassword_UserNil(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.ForgotPasswordRequest{
		Email: "nonexistent@example.com",
	}

	// User is nil
	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(nil, nil).
		Times(1)

	err := service.ForgotPassword(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Tài khoản không tồn tại")
}

func TestForgotPassword_NoPassword(t *testing.T) {
	service, ctrl, mockUserRepo, _, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.ForgotPasswordRequest{
		Email: "firebase@example.com",
	}

	user := &model.User{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		Email:        req.Email,
		PasswordHash: nil, // Firebase-only user
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(user, nil).
		Times(1)

	err := service.ForgotPassword(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Tài khoảng chưa đặt mật khẩu")
}

// VerifyOTP - minimal tests
func TestVerifyOTP_InvalidOrExpired(t *testing.T) {
	service, ctrl, _, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	otp := "123456"

	// OTP doesn't exist in Redis
	mockRedis.EXPECT().
		Get(ctx, "reset:otp:"+otp).
		Return("", assert.AnError).
		Times(1)

	err := service.VerifyOTP(ctx, otp)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Mã OTP không hợp lệ hoặc đã hết hạn")
}

func TestVerifyOTP_Success(t *testing.T) {
	service, ctrl, _, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	otp := "123456"
	email := "test@example.com"

	// OTP exists
	mockRedis.EXPECT().
		Get(ctx, "reset:otp:"+otp).
		Return(email, nil).
		Times(1)

	// Delete OTP (blacklist)
	mockRedis.EXPECT().
		Del(ctx, "reset:otp:"+otp).
		Return(nil).
		Times(1)

	// Store verified key
	mockRedis.EXPECT().
		Set(ctx, "reset:verified:"+otp, email, 5*time.Minute).
		Return(nil).
		Times(1)

	err := service.VerifyOTP(ctx, otp)

	assert.NoError(t, err)
}

// ResetPassword - minimal tests
func TestResetPassword_InvalidToken(t *testing.T) {
	service, ctrl, _, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.ResetPasswordRequest{
		Token:       "invalid-otp",
		NewPassword: "newpassword123",
	}

	// Verified key doesn't exist
	mockRedis.EXPECT().
		Get(ctx, "reset:verified:"+req.Token).
		Return("", assert.AnError).
		Times(1)

	// OTP key doesn't exist either
	mockRedis.EXPECT().
		Get(ctx, "reset:otp:"+req.Token).
		Return("", assert.AnError).
		Times(1)

	err := service.ResetPassword(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Token đặt lại mật khẩu không hợp lệ hoặc đã hết hạn")
}

func TestResetPassword_UserNotFound(t *testing.T) {
	service, ctrl, mockUserRepo, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.ResetPasswordRequest{
		Token:       "valid-otp",
		NewPassword: "newpassword123",
	}
	email := "test@example.com"

	// Verified key exists
	mockRedis.EXPECT().
		Get(ctx, "reset:verified:"+req.Token).
		Return(email, nil).
		Times(1)

	// User not found
	mockUserRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, assert.AnError).
		Times(1)

	err := service.ResetPassword(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Token đặt lại mật khẩu không hợp lệ")
}

func TestForgotPassword_Success(t *testing.T) {
	service, ctrl, mockUserRepo, mockRedis, mockNotification := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.ForgotPasswordRequest{
		Email: "test@example.com",
	}

	passwordHash := "hashed-password"
	user := &model.User{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		Email:        req.Email,
		FullName:     "Test User",
		PasswordHash: &passwordHash,
	}

	// User exists
	mockUserRepo.EXPECT().
		GetByEmail(ctx, req.Email).
		Return(user, nil).
		Times(1)

	// Check rate limit - no limit
	mockRedis.EXPECT().
		Get(ctx, "reset:otp_rate_limit:"+req.Email).
		Return("", assert.AnError). // Key doesn't exist
		Times(1)

	// Get old OTP - none exists
	mockRedis.EXPECT().
		Get(ctx, "reset:email_to_otp:"+req.Email).
		Return("", assert.AnError).
		Times(1)

	// Store new OTP
	mockRedis.EXPECT().
		Set(ctx, gomock.Any(), req.Email, 15*time.Minute).
		Return(nil).
		Times(1)

	// Store email-to-OTP mapping
	mockRedis.EXPECT().
		Set(ctx, "reset:email_to_otp:"+req.Email, gomock.Any(), 15*time.Minute).
		Return(nil).
		Times(1)

	// Store rate limit
	mockRedis.EXPECT().
		Set(ctx, "reset:otp_rate_limit:"+req.Email, "1", 30*time.Second).
		Return(nil).
		Times(1)

	// Async email sending - use AnyTimes for goroutine
	mockNotification.EXPECT().
		Send(gomock.Any(), req.Email, user.FullName, gomock.Any()).
		Return(nil).
		AnyTimes()

	err := service.ForgotPassword(ctx, req)

	assert.NoError(t, err)

	// Wait for goroutine to complete
	time.Sleep(50 * time.Millisecond)
}

func TestResetPassword_Success(t *testing.T) {
	service, ctrl, mockUserRepo, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	otp := "123456"
	email := "test@example.com"

	req := &model.ResetPasswordRequest{
		Token:       otp,
		NewPassword: "newpassword123",
	}

	oldPasswordHash := "old-hash"
	user := &model.User{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		Email:        email,
		PasswordHash: &oldPasswordHash,
	}

	// Verified key exists
	mockRedis.EXPECT().
		Get(ctx, "reset:verified:"+otp).
		Return(email, nil).
		Times(1)

	// User found
	mockUserRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(user, nil).
		Times(1)

	// Update user password
	mockUserRepo.EXPECT().
		Update(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, u *model.User) error {
			// Verify password was changed and hashed
			assert.NotNil(t, u.PasswordHash)
			assert.NotEqual(t, oldPasswordHash, *u.PasswordHash)
			return nil
		}).
		Times(1)

	// Async cleanup - use AnyTimes for goroutine
	mockRedis.EXPECT().
		Del(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	// Async BlacklistUserTokens
	mockRedis.EXPECT().
		Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	err := service.ResetPassword(ctx, req)

	assert.NoError(t, err)

	// Wait for async cleanup
	time.Sleep(50 * time.Millisecond)
}

func TestResetPassword_WithOTPKey(t *testing.T) {
	service, ctrl, mockUserRepo, mockRedis, _ := setupAuthService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	otp := "123456"
	email := "test@example.com"

	req := &model.ResetPasswordRequest{
		Token:       otp,
		NewPassword: "newpassword123",
	}

	passwordHash := "old-hash"
	user := &model.User{
		BaseModel:    model.BaseModel{ID: uuid.New()},
		Email:        email,
		PasswordHash: &passwordHash,
	}

	// Verified key doesn't exist (OTP not verified yet)
	mockRedis.EXPECT().
		Get(ctx, "reset:verified:"+otp).
		Return("", assert.AnError).
		Times(1)

	// Fall back to OTP key - it exists
	mockRedis.EXPECT().
		Get(ctx, "reset:otp:"+otp).
		Return(email, nil).
		Times(1)

	// User found
	mockUserRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(user, nil).
		Times(1)

	// Update user
	mockUserRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil).
		Times(1)

	// Async cleanup
	mockRedis.EXPECT().
		Del(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	mockRedis.EXPECT().
		Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	err := service.ResetPassword(ctx, req)

	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
}
