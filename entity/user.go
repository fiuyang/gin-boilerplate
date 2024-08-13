package entity

import "time"

type User struct {
	ID        int       `json:"id"         gorm:"type:int;primary_key"`
	Username  string    `json:"username"   gorm:"type:varchar(255);not null"`
	Email     string    `json:"email"      gorm:"uniqueIndex;not null"`
	Password  string    `json:"password"   gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
