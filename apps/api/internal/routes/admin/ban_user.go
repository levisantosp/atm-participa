package admin

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/middlewares"
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
	update := db.Client.User.UpdateOneID(input.ID).
		SetBanned(true).
		SetBannedAt(now)

	if input.Body != nil && input.Body.BanEndsAt != nil {
		if !input.Body.BanEndsAt.After(now) {
			return nil, huma.Error422UnprocessableEntity(
				"A data de término do banimento deve ser futura",
			)
		}
		update.SetBanEndsAt(*input.Body.BanEndsAt)
	}

	u, err := update.Save(ctx)
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
