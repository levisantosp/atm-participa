package users

import (
	"context"
	"net/http"
	"strconv"

	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/middlewares"
	"github.com/levisantosp/atm-participa/api/internal/routes/auth"
	"github.com/levisantosp/atm-participa/infra/redis"
)

func DeleteAccount(
	ctx context.Context,
	input *struct{},
) (*auth.OutputAuth, error) {
	user := middlewares.MustGetUserFromContext(ctx)

	if err := db.Client.User.DeleteOneID(user.ID).Exec(ctx); err != nil {
		return nil, err
	}

	key := "user_sessions:" + strconv.FormatInt(user.ID, 10)
	ids, err := redis.Client.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	if len(ids) > 0 {
		keys := make([]string, 0, len(ids))
		for _, id := range ids {
			keys = append(keys, "session:"+id)
		}
		if err := redis.Client.Unlink(ctx, keys...).Err(); err != nil {
			return nil, err
		}
	}
	if err := redis.Client.Unlink(ctx, key).Err(); err != nil {
		return nil, err
	}

	return &auth.OutputAuth{
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
