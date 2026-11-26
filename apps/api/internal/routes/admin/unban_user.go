package admin

import (
	"context"
	"strconv"

	"github.com/levisantosp/atm-participa/pkg/infra/redis"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/user"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/middlewares"
)

func UnbanUser(
	ctx context.Context,
	input *struct {
		ID int64 `path:"id" minimum:"1"`
	},
) (*Output[dtos.User], error) {
	admin := middlewares.MustGetUserFromContext(ctx)
	if admin.ID == input.ID {
		return nil, huma.Error403Forbidden(
			"Você não pode desbanir a própria conta",
		)
	}

	u, err := db.WithTx(ctx, func(tx *generated.Tx) (*generated.User, error) {
		u, err := tx.User.UpdateOneID(input.ID).
			Where(user.BannedEQ(true)).
			ClearBanned().
			ClearBannedAt().
			ClearBanEndsAt().
			Save(ctx)
		if generated.IsNotFound(err) {
			return tx.User.Get(ctx, input.ID)
		}
		if err != nil {
			return nil, err
		}

		raw := strconv.FormatInt(u.ID, 10)
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

		return u, nil
	})
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Usuário não encontrado")
		}

		return nil, err
	}

	return &Output[dtos.User]{
		Body: dtos.UserFrom(u),
	}, nil
}
