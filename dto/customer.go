package dto

type CustomerResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
}

type CreateCustomerBatchRequest struct {
	Customers []CreateCustomerRequest `json:"customers" validate:"required,dive"`
}

type CreateCustomerRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,unique=customers;email"`
	Phone    string `json:"phone" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

type UpdateCustomerRequest struct {
	ID       int    `json:"id" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,unique=customers;email;id"`
	Phone    string `json:"phone" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

type DeleteBatchCustomerRequest struct {
	ID []int `json:"id" validate:"required,sliceInt"`
}

type CustomerParams struct {
	CustomerId int `params:"customerId" validate:"required"`
}

type CustomerQueryFilter struct {
	All       bool   `form:"all" example:"true"`
	Limit     int    `form:"limit"`
	Page      int    `form:"page"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Username  string `form:"username"`
	Email     string `form:"email"`
	Sort      string `form:"sort"`
}
