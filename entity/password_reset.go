package entity

import "time"

type PasswordReset struct {
	ID        int       `json:"id" gorm:"type:int;primary_key"`
	Email     string    `json:"email"`
	Otp       int       `json:"otp"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:Email;references:Email"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}
