package auth

import (
	"context"
	errs "errors"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	acc "github.com/levisantosp/atm-participa/apps/api/ent/generated/account"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/utils"
	"github.com/levisantosp/atm-participa/pkg/infra/errors"
	"github.com/levisantosp/atm-participa/pkg/infra/redis"
)

func ResetPassord(
	ctx context.Context,
	input *struct {
		Body struct {
			Token    string `json:"token" maxLength:"100"`
			Password string `json:"password" minLength:"8" maxLength:"72"`
		}
	},
) (*struct{}, error) {
	token, hash, err := GenerateToken(&input.Body.Token)
	if err != nil {
		return nil, err
	}
	if *token != input.Body.Token {
		return nil, huma.Error401Unauthorized(errors.Unauthorized)
	}

	raw, err := redis.Client.GetDel(ctx, "reset-password:"+*hash).Result()
	if err != nil {
		if errs.Is(err, redis.Nil) {
			return nil, huma.Error401Unauthorized("Token inválido ou expirado.")
		}
	}

	userID, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, err
	}

	passHash, err := utils.GeneratePasswordHash(input.Body.Password)
	if err != nil {
		return nil, err
	}

	_, err = db.WithTx(ctx, func(tx *generated.Tx) (*struct{}, error) {
		_, err := tx.Account.Update().
			Where(acc.UserIDEQ(userID), acc.ProviderEQ(acc.ProviderEmail)).
			SetPassword(*passHash).
			Save(ctx)
		if err != nil {
			return nil, err
		}

		keys := []string{"user_sessions:" + raw}
		sessionsIDs, err := redis.Client.ZRange(ctx, keys[0], 0, -1).Result()
		if err != nil {
			return nil, err
		}

		for _, id := range sessionsIDs {
			keys = append(keys, "session:"+id)
		}

		if err := redis.Client.Unlink(ctx, keys...).Err(); err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
