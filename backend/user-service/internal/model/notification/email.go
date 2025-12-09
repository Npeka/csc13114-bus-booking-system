package notification

type SendOTPEmailRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
	Name  string `json:"name"`
}

type SendOTPEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GenericNotificationRequest struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}
