package constants

// Common message constants for handlers
const (
	// Request/Response messages
	MSG_INVALID_REQUEST_BODY = "Invalid request body"
	MSG_VALIDATION_FAILED    = "Validation failed"
	MSG_INVALID_USER_ID      = "Invalid user ID"
	MSG_USER_NOT_FOUND       = "User not found"
	MSG_USER_ID_REQUIRED     = "User ID required"

	// Auth messages
	MSG_SIGNUP_SUCCESS            = "User registered successfully"
	MSG_SIGNUP_FAILED             = "Registration failed"
	MSG_SIGNIN_SUCCESS            = "User signed in successfully"
	MSG_SIGNIN_FAILED             = "Sign in failed"
	MSG_SIGNOUT_SUCCESS           = "User signed out successfully"
	MSG_SIGNOUT_FAILED            = "Sign out failed"
	MSG_INVALID_CREDENTIALS       = "Invalid email or password"
	MSG_OAUTH2_SIGNIN_SUCCESS     = "OAuth2 sign in successful"
	MSG_OAUTH2_SIGNIN_FAILED      = "OAuth2 sign in failed"
	MSG_INVALID_TOKEN             = "Invalid token"
	MSG_TOKEN_VALID               = "Token is valid"
	MSG_TOKEN_INVALID             = "Token is invalid"
	MSG_TOKEN_VERIFICATION_FAILED = "Token verification failed"
	MSG_TOKEN_REFRESH_SUCCESS     = "Token refreshed successfully"
	MSG_TOKEN_REFRESH_FAILED      = "Token refresh failed"
	MSG_INVALID_REFRESH_TOKEN     = "Invalid refresh token"
	MSG_TOKEN_VERIFY_SUCCESS      = "Token verified successfully"
	MSG_TOKEN_VERIFY_FAILED       = "Token verification failed"

	// User management messages
	MSG_USER_CREATED        = "User created successfully"
	MSG_CREATE_USER_FAILED  = "Failed to create user"
	MSG_GET_USER_SUCCESS    = "User retrieved successfully"
	MSG_GET_USER_FAILED     = "Failed to get user"
	MSG_UPDATE_USER_SUCCESS = "User updated successfully"
	MSG_UPDATE_USER_FAILED  = "Failed to update user"
	MSG_DELETE_USER_SUCCESS = "User deleted successfully"
	MSG_DELETE_USER_FAILED  = "Failed to delete user"
	MSG_LIST_USERS_SUCCESS  = "Users listed successfully"
	MSG_LIST_USERS_FAILED   = "Failed to list users"

	// Legacy message constants for user_handler.go
	MsgFailedToBindRequest  = "Failed to bind request"
	MsgValidationFailed     = "Validation failed"
	MsgFailedToCreateUser   = "Failed to create user"
	MsgInvalidUserID        = "Invalid user ID"
	MsgFailedToGetUser      = "Failed to get user"
	MsgUserNotFound         = "User not found"
	MsgFailedToUpdateUser   = "Failed to update user"
	MsgFailedToDeleteUser   = "Failed to delete user"
	MsgFailedToListUsers    = "Failed to list users"
	MsgFailedToUpdateStatus = "Failed to update status"
)
