package users

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/middlewares"
	"github.com/levisantosp/atm-participa/infra/errors"
	"github.com/levisantosp/atm-participa/infra/redis"
	goredis "github.com/redis/go-redis/v9"
)

type EditMeOutput struct {
	Body      middlewares.Session
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

func EditMe(
	ctx context.Context,
	input *struct {
		Body struct {
			Username    *string `json:"username" minLength:"3" maxLength:"32" pattern:"^[a-zA-Z0-9_]+$" required:"false"`
			DisplayName *string `json:"displayName" minLength:"1" maxLength:"100" required:"false"`
		}
	},
) (*EditMeOutput, error) {
	if input.Body.Username == nil && input.Body.DisplayName == nil {
		return nil, huma.Error400BadRequest(
			"Informe ao menos um campo para atualizar",
		)
	}

	userCtx := middlewares.MustGetUserFromContext(ctx)
	session := middlewares.MustGetSessionFromContext(ctx)
	user, err := db.Client.User.UpdateOneID(userCtx.ID).
		SetNillableDisplayName(input.Body.DisplayName).
		SetNillableUsername(input.Body.Username).
		Save(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound(errors.NotFound)
		}

		return nil, err
	}

	session = middlewares.Session{
		ID:          session.ID,
		UserId:      strconv.FormatInt(user.ID, 10),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		IsAdmin:     user.IsAdmin,
	}

	sessionJson, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	ttl := 7 * 24 * time.Hour
	err = redis.Client.SetArgs(ctx, "session:"+session.ID, sessionJson, goredis.SetArgs{
		Mode: "XX",
		TTL:  ttl,
	}).
		Err()
	if err != nil {
		if err == redis.Nil {
			return nil, huma.Error401Unauthorized(errors.Unauthorized)
		}
		return nil, err
	}

	sessionCookie := http.Cookie{
		Name:     "session",
		Value:    session.ID,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
		Path:     "/",
	}

	return &EditMeOutput{
		Body:      session,
		SetCookie: []http.Cookie{sessionCookie},
	}, nil
}
