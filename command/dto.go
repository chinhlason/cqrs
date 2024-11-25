package command

import (
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Order struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	Product   string    `json:"product"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserId;references:Id"`
}
