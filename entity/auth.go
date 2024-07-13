package entity

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
	PasswordConfirmation string `json:"password_confirmation" validate:"required"`
}

type CreateUserRequest struct {
	Username string `json:"username"  validate:"required" `
	Email    string `json:"email"     validate:"required,email,unique=users;email"`
	Password string `json:"password"  validate:"required,min=8,max=100"`
}

type UpdateUserRequest struct {
	ID       int    `json:"username" validate:"required"`
	Username string `json:"username" validate:"required,max=200,min=2"`
	Email    string `json:"email"    validate:"required,email,unique=users;email;id"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UserParams struct {
	UserId int `params:"userId" validate:"required"`
}

type DeleteBatchUserRequest struct {
	ID []int `json:"id" validate:"required,notEmptyIntSlice"`
}

type UserQueryFilter struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Username  string `form:"username"`
	Email     string `form:"email"`
	Sort      string `form:"sort"`
}
