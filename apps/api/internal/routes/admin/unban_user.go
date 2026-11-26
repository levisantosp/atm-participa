package admin

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/middlewares"
)

func UnbanUser(
	ctx context.Context,
	input *struct {
		ID int64 `path:"id" minimum:"1"`
	},
) (*struct{}, error) {
	user := middlewares.MustGetUserFromContext(ctx)
	if user.ID == input.ID {
		return nil, huma.Error403Forbidden(
			"Você não pode desbanir a própria conta",
		)
	}

	err := db.Client.User.UpdateOneID(input.ID).
		ClearBanned().
		ClearBannedAt().
		ClearBanEndsAt().
		Exec(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Usuário não encontrado")
		}

		return nil, err
	}

	return nil, nil
}
