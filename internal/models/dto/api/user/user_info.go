package user

import (
	"time"
)

type UserInfo struct {
	Id          string    `json:"id"`
	UpdateAt    time.Time `json:"updated_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	UserName    string    `json:"username"`
	CreateAt    time.Time `json:"created_at"`
}
