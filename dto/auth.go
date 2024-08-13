package dto

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"         `
	Password string `json:"password" validate:"required,min=2,max=100" `
}

type ForgotPasswordRequest struct {
	Email string `json:"email"  validate:"required,email"`
}

type CheckOtpRequest struct {
	Otp int `json:"otp" validate:"required"`
}

type ResetPasswordRequest struct {
	Otp                  int    `json:"otp"                   validate:"required"`
	Password             string `json:"password"              validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,equal=Password"`
}
