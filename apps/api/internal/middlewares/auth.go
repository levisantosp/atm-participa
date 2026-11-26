package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/pkg/infra/errors"
	"github.com/levisantosp/atm-participa/pkg/infra/redis"
)

type Session struct {
	ID          string     `json:"id"`
	UserId      string     `json:"userId"`
	Username    string     `json:"username"`
	DisplayName string     `json:"displayName"`
	IsAdmin     bool       `json:"isAdmin"`
	Banned      *bool      `json:"banned"`
	BannedAt    *time.Time `json:"bannedAt"`
	BanEndsAt   *time.Time `json:"banEndsAt"`
}

type ctxKey int

const (
	userCtxKey ctxKey = iota
	sessionCtxKey
)

func Auth(
	api huma.API,
	adminOnly bool,
) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		cookie, err := huma.ReadCookie(ctx, "session")
		if err != nil {
			huma.WriteErr(
				api,
				ctx,
				http.StatusUnauthorized,
				errors.Unauthorized,
			)
			return
		}

		raw, err := redis.Client.Get(ctx.Context(), "session:"+cookie.Value).
			Result()
		if err != nil {
			if err == redis.Nil {
				huma.WriteErr(
					api,
					ctx,
					http.StatusUnauthorized,
					errors.Unauthorized,
				)
				return
			}
			huma.WriteErr(
				api,
				ctx,
				http.StatusInternalServerError,
				errors.InternalServerError,
				err,
			)
			return
		}

		var session Session
		if err := json.Unmarshal([]byte(raw), &session); err != nil {
			huma.WriteErr(
				api,
				ctx,
				http.StatusInternalServerError,
				errors.InternalServerError,
				err,
			)
			return
		}

		userId, err := strconv.ParseInt(session.UserId, 10, 64)
		if err != nil {
			huma.WriteErr(
				api,
				ctx,
				http.StatusInternalServerError,
				errors.InternalServerError,
				err,
			)
			return
		}

		user, err := db.Client.User.Get(ctx.Context(), userId)
		if err != nil {
			if generated.IsNotFound(err) {
				huma.WriteErr(
					api,
					ctx,
					http.StatusUnauthorized,
					errors.Unauthorized,
				)
				return
			}

			huma.WriteErr(
				api,
				ctx,
				http.StatusInternalServerError,
				errors.InternalServerError,
				err,
			)
			return
		}

		if (user.Banned != nil && *user.Banned) ||
			(adminOnly && !user.IsAdmin) {
			huma.WriteErr(
				api,
				ctx,
				http.StatusForbidden,
				errors.Forbidden,
			)
			return
		}

		session.Username = user.Username
		session.DisplayName = user.DisplayName
		session.IsAdmin = user.IsAdmin
		session.Banned = user.Banned
		session.BannedAt = user.BannedAt
		session.BanEndsAt = user.BanEndsAt

		ctx = huma.WithValue(ctx, userCtxKey, user)
		ctx = huma.WithValue(ctx, sessionCtxKey, session)

		next(ctx)
	}
}

func MustGetUserFromContext(ctx context.Context) *generated.User {
	return ctx.Value(userCtxKey).(*generated.User)
}

func MustGetSessionFromContext(ctx context.Context) Session {
	return ctx.Value(sessionCtxKey).(Session)
}
