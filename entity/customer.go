package entity

import "time"

type Customer struct {
	ID        int       `json:"id" gorm:"type:int;primary_key"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Customer) TableName() string {
	return "customers"
}
