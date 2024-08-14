package dto

type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
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
	ID []int `json:"id" validate:"required,sliceInt"`
}

type UserQueryFilter struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Username  string `form:"username"`
	Email     string `form:"email"`
	Sort      string `form:"sort"`
}
