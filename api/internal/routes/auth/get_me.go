package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/internal/middlewares"
	"github.com/levisantosp/atm-participa/infra/errors"
	"github.com/levisantosp/atm-participa/infra/redis"
)

func GetMe(ctx context.Context, input *struct {
	Session http.Cookie `cookie:"session"`
},
) (*Output[middlewares.Session], error) {
	if input.Session.Name == "" {
		return nil, huma.Error404NotFound(errors.NotFound)
	}

	raw, err := redis.Client.Get(ctx, "session:"+input.Session.Value).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, huma.Error404NotFound(errors.NotFound)
		}
		return nil, err
	}

	var session middlewares.Session
	if err := json.Unmarshal([]byte(raw), &session); err != nil {
		return nil, err
	}

	return &Output[middlewares.Session]{
		Body: session,
	}, nil
}
