package handler

// Common message constants for handlers
const (
	MsgFailedToBindRequest  = "Failed to bind request"
	MsgValidationFailed     = "Validation failed"
	MsgInvalidUserID        = "Invalid user ID"
	MsgUserNotFound         = "User not found"
	MsgFailedToCreateUser   = "Failed to create user"
	MsgFailedToGetUser      = "Failed to get user"
	MsgFailedToUpdateUser   = "Failed to update user"
	MsgFailedToDeleteUser   = "Failed to delete user"
	MsgFailedToListUsers    = "Failed to list users"
	MsgFailedToUpdateStatus = "Failed to update user status"

	// Auth specific messages
	MsgLoginFailed            = "Login failed"
	MsgLoginSuccess           = "User logged in successfully"
	MsgTokenRefreshFailed     = "Token refresh failed"
	MsgTokenRefreshSuccess    = "Token refreshed successfully"
	MsgPasswordChangeFailed   = "Password change failed"
	MsgPasswordChangeSuccess  = "Password changed successfully"
	MsgPasswordResetFailed    = "Password reset request failed"
	MsgPasswordResetRequested = "Password reset requested"
	MsgConfirmResetFailed     = "Password reset confirmation failed"
)
