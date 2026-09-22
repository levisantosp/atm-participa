package auth

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/middlewares"
	"github.com/levisantosp/atm-participa/api/redis"
	"github.com/levisantosp/atm-participa/api/utils"
)

func SignOut(ctx context.Context, input *struct{}) (*SignInOutput, error) {
	session := middlewares.MustGetSessionFromContext(ctx)

	if err := redis.Client.Unlink(ctx, "session:"+session.ID).
		Err(); err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	return &SignInOutput{
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
