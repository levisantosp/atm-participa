package dtos

import (
	"time"

	"github.com/levisantosp/atm-participa/api/ent/generated"
)

type User struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	IsAdmin     bool      `json:"isAdmin"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func UserFrom(user *generated.User) User {
	return User{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		IsAdmin:     user.IsAdmin,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
