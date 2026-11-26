package users

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/report"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/middlewares"
)

func Report(
	ctx context.Context,
	input *struct {
		UserID int64 `path:"userId" minimum:"1"`
	},
) (*struct{}, error) {
	user := middlewares.MustGetUserFromContext(ctx)
	if user.ID == input.UserID {
		return nil, huma.Error422UnprocessableEntity(
			"Você não pode denunciar a si mesmo",
		)
	}

	_, err := db.WithTx(ctx, func(tx *generated.Tx) (*generated.Report, error) {
		target, err := tx.User.Get(ctx, input.UserID)
		if err != nil {
			return nil, err
		}

		return tx.Report.Create().
			SetUserID(user.ID).
			SetEntity(report.EntityUser).
			SetEntityID(target.ID).
			SetTargetUser(target).
			Save(ctx)
	})
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Usuário não encontrado")
		}

		return nil, err
	}

	return nil, nil
}
