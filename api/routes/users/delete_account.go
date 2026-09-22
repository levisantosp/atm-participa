package users

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/middlewares"
	"github.com/levisantosp/atm-participa/api/redis"
	"github.com/levisantosp/atm-participa/api/routes/auth"
	"github.com/levisantosp/atm-participa/api/utils"
)

func DeleteAccount(
	ctx context.Context,
	input *struct{},
) (*auth.SignInOutput, error) {
	user := middlewares.MustGetUserFromContext(ctx)
	session := middlewares.MustGetSessionFromContext(ctx)

	if err := db.Client.User.DeleteOneID(user.ID).Exec(ctx); err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	if err := redis.Client.Unlink(ctx, "session:"+session.ID).
		Err(); err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
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
