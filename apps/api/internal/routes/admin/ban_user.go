package admin

import (
	"context"
	"strconv"
	"time"

	"github.com/levisantosp/atm-participa/pkg/infra/redis"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/middlewares"
	"github.com/levisantosp/atm-participa/apps/api/internal/queue"
)

func BanUser(
	ctx context.Context,
	input *struct {
		ID   int64 `path:"id" minimum:"1"`
		Body *struct {
			BanEndsAt *time.Time `json:"banEndsAt" required:"false"`
		}
	},
) (*Output[dtos.User], error) {
	user := middlewares.MustGetUserFromContext(ctx)
	if user.ID == input.ID {
		return nil, huma.Error403Forbidden(
			"Você não pode banir a própria conta",
		)
	}

	now := time.Now()
	if input.Body != nil && input.Body.BanEndsAt != nil &&
		!input.Body.BanEndsAt.After(now) {
		return nil, huma.Error422UnprocessableEntity(
			"A data de término do banimento deve ser futura",
		)
	}

	u, err := db.WithTx(ctx, func(tx *generated.Tx) (*generated.User, error) {
		update := tx.User.UpdateOneID(input.ID).SetBanned(true).SetBannedAt(now)
		if input.Body != nil {
			update.SetNillableBanEndsAt(input.Body.BanEndsAt)
		}
		u, err := update.Save(ctx)
		if err != nil {
			return nil, err
		}
		if u.BanEndsAt != nil {
			if err := queue.ScheduleUnban(ctx, u.ID, *u.BanEndsAt); err != nil {
				return nil, err
			}
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
