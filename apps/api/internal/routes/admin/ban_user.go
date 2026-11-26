package admin

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
)

func BanUser(
	ctx context.Context,
	input *struct {
		ID   int64 `path:"id" minimum:"1"`
		Body *struct {
			BanEndsAt *time.Time `json:"banEndsAt" required:"false"`
		}
	},
) (*struct{}, error) {
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

	if err := update.Exec(ctx); err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Usuário não encontrado")
		}

		return nil, err
	}

	return nil, nil
}
