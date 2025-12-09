package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bus-booking/shared/constants"
	"bus-booking/shared/db"
	"bus-booking/shared/ginext"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/client"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/repository"
	"bus-booking/user-service/internal/utils"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// Redis key prefixes for password reset flow
const (
	redisKeyResetOTP        = "reset:otp:"            // Stores OTP -> email mapping
	redisKeyResetEmailToOTP = "reset:email_to_otp:"   // Stores email -> OTP mapping
	redisKeyResetRateLimit  = "reset:otp_rate_limit:" // Rate limit for OTP requests
	redisKeyResetVerified   = "reset:verified:"       // Verified email after OTP check
)

type AuthService interface {
	VerifyToken(ctx context.Context, token string) (*model.TokenVerifyResponse, error)
	FirebaseAuth(ctx context.Context, req *model.FirebaseAuthRequest) (*model.AuthResponse, error)
	Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error)
	Logout(ctx context.Context, req model.LogoutRequest, userID uuid.UUID) error
	ForgotPassword(ctx context.Context, req *model.ForgotPasswordRequest) error
	VerifyOTP(ctx context.Context, otp string) error
	ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error
	RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.AuthResponse, error)
	CreateGuestAccount(ctx context.Context, req *model.CreateGuestAccountRequest) (*model.UserResponse, error)
}

type AuthServiceImpl struct {
	config             *config.Config
	userRepo           repository.UserRepository
	jwtManager         utils.JWTManager
	firebaseAuth       FirebaseAuthClient
	tokenManager       TokenManager
	redisClient        *db.RedisManager
	notificationClient client.NotificationClient
}

func NewAuthService(
	config *config.Config,
	jwtManager utils.JWTManager,
	firebaseAuth FirebaseAuthClient,
	tokenManager TokenManager,
	userRepo repository.UserRepository,
	redisClient *db.RedisManager,
	notificationClient client.NotificationClient,
) AuthService {
	return &AuthServiceImpl{
		config:             config,
		jwtManager:         jwtManager,
		firebaseAuth:       firebaseAuth,
		tokenManager:       tokenManager,
		userRepo:           userRepo,
		redisClient:        redisClient,
		notificationClient: notificationClient,
	}
}

func (s *AuthServiceImpl) VerifyToken(ctx context.Context, accessToken string) (*model.TokenVerifyResponse, error) {
	claims, err := s.jwtManager.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, ginext.NewUnauthorizedError("token không hợp lệ")
	}

	if s.tokenManager.IsBlacklisted(ctx, accessToken) {
		return nil, ginext.NewUnauthorizedError("token đã bị blacklisted")
	}

	if s.tokenManager.IsUserTokensBlacklisted(ctx, claims.UserID, claims.IssuedAt.Unix()) {
		return nil, ginext.NewUnauthorizedError("token của người dùng đã bị blacklisted")
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, ginext.NewUnauthorizedError("không tìm thấy người dùng")
	}

	if user.Status != constants.UserStatusActive && user.Status != constants.UserStatusVerified {
		return nil, ginext.NewUnauthorizedError("tài khoản không hoạt động")
	}

	return &model.TokenVerifyResponse{
		UserID: claims.UserID.String(),
		Email:  user.Email,
		Role:   user.Role,
		Name:   user.FullName,
	}, nil
}

func (s *AuthServiceImpl) FirebaseAuth(ctx context.Context, req *model.FirebaseAuthRequest) (*model.AuthResponse, error) {
	token, err := s.firebaseAuth.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		return nil, ginext.NewUnauthorizedError("firebase token không hợp lệ")
	}

	user, err := s.userRepo.GetByFirebaseUID(ctx, token.UID)
	if err != nil {
		return nil, ginext.NewInternalServerError(err.Error())
	}
	if user != nil {
		if user.Status != constants.UserStatusActive && user.Status != constants.UserStatusVerified {
			return nil, ginext.NewForbiddenError("tài khoản không hoạt động")
		}
		return s.generateAuthResponse(user)
	}

	// Extract claims from Firebase token
	email := ""
	phone := ""
	fullName := ""
	avatar := ""

	if emailClaim, exists := token.Claims["email"]; exists && emailClaim != nil {
		email, _ = emailClaim.(string)
	}
	if phoneClaim, exists := token.Claims["phone_number"]; exists && phoneClaim != nil {
		phone, _ = phoneClaim.(string)
	}
	if nameClaim, exists := token.Claims["name"]; exists && nameClaim != nil {
		fullName, _ = nameClaim.(string)
	}
	if pictureClaim, exists := token.Claims["picture"]; exists && pictureClaim != nil {
		avatar, _ = pictureClaim.(string)
	}

	// Generate full name from email if not provided
	if fullName == "" && email != "" {
		fullName = strings.Split(email, "@")[0]
	}
	// Fallback to phone number for full name
	if fullName == "" && phone != "" {
		fullName = phone
	}
	// If still empty, use firebase UID as fallback
	if fullName == "" {
		fullName = token.UID[:12]
	}

	// Check email verification status
	emailVerified := false
	if emailVerifyClaim, exists := token.Claims["email_verified"]; exists && emailVerifyClaim != nil {
		emailVerified, _ = emailVerifyClaim.(bool)
	}

	// Check phone verification status
	phoneVerified := phone != ""

	// Create new user
	user = &model.User{
		Email:         email,
		Phone:         phone,
		FullName:      fullName,
		Avatar:        avatar,
		Role:          constants.RolePassenger,
		Status:        constants.UserStatusVerified,
		FirebaseUID:   &token.UID,
		EmailVerified: emailVerified,
		PhoneVerified: phoneVerified,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to create Firebase user")
		return nil, ginext.NewInternalServerError("Không thể tạo người dùng")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		log.Warn().Str("email", req.Email).Msg("Email already registered")
		return nil, ginext.NewBadRequestError("Email đã được đăng ký")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, ginext.NewInternalServerError("Không thể tạo tài khoản")
	}

	// Create new user
	user := &model.User{
		Email:         req.Email,
		FullName:      req.FullName,
		PasswordHash:  &passwordHash,
		Role:          constants.RolePassenger,
		Status:        constants.UserStatusActive,
		FirebaseUID:   nil, // Empty for email/password users
		EmailVerified: false,
		PhoneVerified: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, ginext.NewInternalServerError("Không thể tạo tài khoản")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("User not found")
		return nil, ginext.NewUnauthorizedError("Email hoặc mật khẩu không đúng")
	}

	// Check if user has password set (not a Firebase-only user)
	if user.PasswordHash == nil {
		log.Warn().Str("email", req.Email).Msg("User does not have password set")
		return nil, ginext.NewUnauthorizedError("Email hoặc mật khẩu không đúng")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, *user.PasswordHash) {
		log.Error().Str("email", req.Email).Msg("Password verification failed")
		return nil, ginext.NewUnauthorizedError("Email hoặc mật khẩu không đúng")
	}

	// Check user status
	if user.Status != constants.UserStatusActive && user.Status != constants.UserStatusVerified {
		return nil, ginext.NewForbiddenError("Tài khoản không hoạt động")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthServiceImpl) ForgotPassword(ctx context.Context, req *model.ForgotPasswordRequest) error {
	// Check if user exists with email/password auth
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return ginext.NewInternalServerError("Không thể xử lý yêu cầu đặt lại mật khẩu")
	}
	if user == nil {
		return ginext.NewBadRequestError("Tài khoản không tồn tại")
	}

	// Check if user has password (not Firebase-only user)
	if user.PasswordHash == nil {
		log.Warn().Str("email", req.Email).Msg("Password reset requested for Firebase-only user")
		return ginext.NewBadRequestError("Tài khoảng chưa đặt mật khẩu")
	}

	// Check rate limit FIRST before doing anything
	rateLimitKey := redisKeyResetRateLimit + req.Email
	_, err = s.redisClient.Get(ctx, rateLimitKey)
	if err == nil {
		// Key exists, rate limit active - user must wait
		ttl, err := s.redisClient.TTL(ctx, rateLimitKey)
		if err != nil {
			ttl = 30 * time.Second // Default fallback
		}
		return ginext.NewBadRequestError(fmt.Sprintf("Vui lòng đợi %d giây trước khi gửi lại", int(ttl.Seconds())))
	}

	// Rate limit passed, now blacklist old OTP if exists
	emailOTPKey := redisKeyResetEmailToOTP + req.Email
	oldOTP, err := s.redisClient.Get(ctx, emailOTPKey)
	if err == nil && oldOTP != "" {
		// Delete the old OTP key to blacklist it
		oldOTPKey := redisKeyResetOTP + oldOTP
		if err := s.redisClient.Del(ctx, oldOTPKey); err != nil {
			log.Warn().Err(err).Str("otp", oldOTP).Msg("Failed to blacklist old OTP")
		}
	}

	// Generate 6-digit OTP
	otp, err := utils.GenerateOTP(6)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate OTP")
		return ginext.NewInternalServerError("Không thể xử lý yêu cầu đặt lại mật khẩu")
	}

	// Store OTP in Redis (OTP itself is the "token")
	key := redisKeyResetOTP + otp
	if err := s.redisClient.Set(ctx, key, req.Email, 15*time.Minute); err != nil {
		log.Error().Err(err).Msg("Failed to store OTP")
		return ginext.NewInternalServerError("Không thể xử lý yêu cầu đặt lại mật khẩu")
	}

	// Store email-to-OTP mapping for blacklisting old OTPs
	if err := s.redisClient.Set(ctx, emailOTPKey, otp, 15*time.Minute); err != nil {
		log.Warn().Err(err).Msg("Failed to store email-OTP mapping")
	}

	// Store rate limit key for this email (30 seconds)
	if err := s.redisClient.Set(ctx, rateLimitKey, "1", 30*time.Second); err != nil {
		log.Warn().Err(err).Msg("Failed to set rate limit")
	}

	// Send OTP email via notification service
	go func() {
		// Use background context to avoid cancellation when request completes
		bgCtx := context.Background()
		if err := s.notificationClient.Send(bgCtx, req.Email, user.FullName, otp); err != nil {
			log.Error().Err(err).Msg("Failed to send OTP email")
			// Continue even if email fails - don't reveal to user
		}
	}()

	return nil
}

func (s *AuthServiceImpl) VerifyOTP(ctx context.Context, otp string) error {
	// Validate OTP exists in Redis
	otpKey := redisKeyResetOTP + otp
	email, err := s.redisClient.Get(ctx, otpKey)
	if err != nil {
		log.Error().Err(err).Msg("Invalid or expired OTP")
		return ginext.NewBadRequestError("Mã OTP không hợp lệ hoặc đã hết hạn")
	}

	// Blacklist this OTP after verification to prevent reuse
	// Delete the OTP key so it can't be used again
	if err := s.redisClient.Del(ctx, otpKey); err != nil {
		log.Warn().Err(err).Str("otp", otp).Msg("Failed to blacklist verified OTP")
	}

	// Store the OTP->email mapping with verified key
	// This gives user 5 minutes to complete password reset after OTP verification
	// Key is OTP (not email) so ResetPassword can look it up by OTP
	verifiedKey := redisKeyResetVerified + otp
	if err := s.redisClient.Set(ctx, verifiedKey, email, 5*time.Minute); err != nil {
		log.Warn().Err(err).Msg("Failed to store verified OTP")
	}

	return nil
}

func (s *AuthServiceImpl) ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error {
	// Note: OTP has already been verified and deleted in VerifyOTP step
	// We now validate using the verified key which was created after OTP verification

	// Try to get email from verified key (created in VerifyOTP)
	// The req.Token here is still the OTP that user entered
	verifiedKey := redisKeyResetVerified + req.Token
	email, err := s.redisClient.Get(ctx, verifiedKey)

	// If verified key doesn't exist, fallback to checking OTP key
	// (for backward compatibility or if OTP wasn't verified yet)
	if err != nil {
		otpKey := redisKeyResetOTP + req.Token
		email, err = s.redisClient.Get(ctx, otpKey)
		if err != nil {
			log.Error().Err(err).Msg("Invalid or expired OTP")
			return ginext.NewBadRequestError("Token đặt lại mật khẩu không hợp lệ hoặc đã hết hạn")
		}
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("User not found for reset")
		return ginext.NewBadRequestError("Token đặt lại mật khẩu không hợp lệ")
	}

	// Hash new password
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash new password")
		return ginext.NewInternalServerError("Không thể đặt lại mật khẩu")
	}

	// Update user password
	user.PasswordHash = &newPasswordHash
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update password")
		return ginext.NewInternalServerError("Không thể đặt lại mật khẩu")
	}

	// Cleanup tasks asynchronously (password already updated successfully)
	go func() {
		// Use background context to avoid cancellation when request completes
		bgCtx := context.Background()

		// Delete verified key
		if err := s.redisClient.Del(bgCtx, verifiedKey); err != nil {
			log.Warn().Err(err).Msg("Failed to delete verified key")
		}

		// Delete OTP key if it still exists
		otpKey := redisKeyResetOTP + req.Token
		if err := s.redisClient.Del(bgCtx, otpKey); err != nil {
			log.Warn().Err(err).Msg("Failed to delete OTP key")
		}

		// Delete email-to-OTP mapping
		emailOTPKey := redisKeyResetEmailToOTP + email
		if err := s.redisClient.Del(bgCtx, emailOTPKey); err != nil {
			log.Warn().Err(err).Msg("Failed to delete email-OTP mapping")
		}

		// Invalidate all user sessions (blacklist all tokens issued before now)
		if !s.tokenManager.BlacklistUserTokens(bgCtx, user.ID) {
			log.Error().Msg("Failed to blacklist user tokens")
		}
	}()

	log.Info().Str("email", email).Msg("Password reset successful")
	return nil
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.AuthResponse, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ginext.NewUnauthorizedError("refresh token không hợp lệ")
	}

	// Check blacklist - đơn giản
	if s.tokenManager.IsBlacklisted(ctx, req.RefreshToken) {
		return nil, ginext.NewUnauthorizedError("refresh token đã bị thu hồi")
	}

	if s.tokenManager.IsUserTokensBlacklisted(ctx, claims.UserID, claims.IssuedAt.Unix()) {
		return nil, ginext.NewUnauthorizedError("tất cả token người dùng đã bị thu hồi")
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, ginext.NewInternalServerError("không tìm thấy người dùng")
	}

	if user.Status != constants.UserStatusActive && user.Status != constants.UserStatusVerified {
		return nil, ginext.NewForbiddenError("tài khoản không hoạt động")
	}

	// Blacklist old refresh token
	s.tokenManager.Blacklist(ctx, req.RefreshToken)

	return s.generateAuthResponse(user)
}

func (s *AuthServiceImpl) generateAuthResponse(user *model.User) (*model.AuthResponse, error) {
	var (
		accessToken  string
		refreshToken string
	)

	// Generate both tokens in parallel using errgroup
	g := new(errgroup.Group)

	// Generate access token
	g.Go(func() error {
		token, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, fmt.Sprintf("%d", user.Role))
		if err != nil {
			return ginext.NewInternalServerError("Không thể tạo token truy cập")
		}
		accessToken = token
		return nil
	})

	// Generate refresh token
	g.Go(func() error {
		token, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, fmt.Sprintf("%d", user.Role))
		if err != nil {
			return ginext.NewInternalServerError("Không thể tạo refresh token")
		}
		refreshToken = token
		return nil
	})

	// Wait for both to complete
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessTokenTTL.Seconds()),
	}, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, req model.LogoutRequest, userID uuid.UUID) error {
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return ginext.NewUnauthorizedError("refresh token không hợp lệ")
	}

	if claims.UserID != userID {
		return ginext.NewUnauthorizedError("refresh token không khớp với người dùng")
	}

	// Blacklist tokens asynchronously (user already logged out from client)
	go func() {
		// Use background context to avoid cancellation when request completes
		bgCtx := context.Background()
		s.tokenManager.Blacklist(bgCtx, req.AccessToken)
		s.tokenManager.Blacklist(bgCtx, req.RefreshToken)
	}()

	return nil
}

// CreateGuestAccount creates a guest user account for non-authenticated bookings
func (s *AuthServiceImpl) CreateGuestAccount(ctx context.Context, req *model.CreateGuestAccountRequest) (*model.UserResponse, error) {
	// Validate: at least one contact method (email or phone) must be provided
	if req.Email == "" && req.Phone == "" {
		return nil, ginext.NewBadRequestError("Phải cung cấp email hoặc số điện thoại")
	}

	// Check if guest already exists by email or phone
	if req.Email != "" {
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			return nil, ginext.NewInternalServerError("Không thể kiểm tra người dùng")
		}
		if existingUser != nil {
			return nil, ginext.NewBadRequestError("Email đã được đăng ký")
		}
	}

	if req.Phone != "" {
		existingUser, err := s.userRepo.GetByPhone(ctx, req.Phone)
		if err != nil {
			return nil, ginext.NewInternalServerError("Không thể kiểm tra người dùng")
		}
		if existingUser != nil {
			return nil, ginext.NewBadRequestError("Số điện thoại đã được đăng ký")
		}
	}

	// Create new guest user
	guestUser := &model.User{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     constants.RoleGuest,
		Status:   constants.UserStatusActive,
	}

	if err := s.userRepo.Create(ctx, guestUser); err != nil {
		log.Error().Err(err).Msg("Failed to create guest account")
		return nil, ginext.NewInternalServerError("Không thể tạo tài khoản khách")
	}

	return guestUser.ToResponse(), nil
}
