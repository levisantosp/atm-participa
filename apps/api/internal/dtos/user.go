package dtos

import (
	"strconv"
	"time"

	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
)

type User struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	DisplayName string     `json:"displayName"`
	IsAdmin     bool       `json:"isAdmin"`
	Banned      *bool      `json:"banned"`
	BannedAt    *time.Time `json:"bannedAt"`
	BanEndsAt   *time.Time `json:"banEndsAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func UserFrom(user *generated.User) User {
	return User{
		ID:          strconv.FormatInt(user.ID, 10),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		IsAdmin:     user.IsAdmin,
		Banned:      user.Banned,
		BannedAt:    user.BannedAt,
		BanEndsAt:   user.BanEndsAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func UsersFrom(users []*generated.User) []User {
	items := make([]User, 0, len(users))
	for _, user := range users {
		items = append(items, UserFrom(user))
	}
	return items
}
