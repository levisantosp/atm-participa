package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/levisantosp/atm-participa/api/internal/middlewares"
	"github.com/levisantosp/atm-participa/infra/redis"
)

func SignOut(ctx context.Context, input *struct{}) (*OutputAuth, error) {
	session := middlewares.MustGetSessionFromContext(ctx)
	user := middlewares.MustGetUserFromContext(ctx)

	if err := redis.Client.Unlink(
		ctx,
		"session:"+session.ID,
	).
		Err(); err != nil {
		return nil, err
	}

	if err := redis.Client.ZRem(
		ctx,
		"user_sessions:"+strconv.FormatInt(user.ID, 10),
		session.ID,
	).Err(); err != nil {
		return nil, err
	}

	return &OutputAuth{
		SetCookie: []http.Cookie{
			{
				Name:     "session",
				Value:    "",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   -1,
				Path:     "/",
			},
		},
	}, nil
}
