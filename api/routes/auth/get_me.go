package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/errors"
	"github.com/levisantosp/atm-participa/api/middlewares"
	"github.com/levisantosp/atm-participa/api/redis"
)

type GetMeOutput struct {
	Body middlewares.Session
}

func GetMe(ctx context.Context, input *struct {
	Session http.Cookie `cookie:"session"`
},
) (*GetMeOutput, error) {
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

	return &GetMeOutput{
		Body: session,
	}, nil
}
