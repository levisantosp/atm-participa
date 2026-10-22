package users

import (
	"context"
	"net/http"

	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/middlewares"
	"github.com/levisantosp/atm-participa/api/internal/routes/auth"
	"github.com/levisantosp/atm-participa/infra/redis"
)

func DeleteAccount(
	ctx context.Context,
	input *struct{},
) (*auth.SignInOutput, error) {
	user := middlewares.MustGetUserFromContext(ctx)
	session := middlewares.MustGetSessionFromContext(ctx)

	if err := db.Client.User.DeleteOneID(user.ID).Exec(ctx); err != nil {
		return nil, err
	}

	if err := redis.Client.Unlink(ctx, "session:"+session.ID).
		Err(); err != nil {
		return nil, err
	}

	return &auth.SignInOutput{
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
