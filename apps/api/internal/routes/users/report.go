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
	reporter := middlewares.MustGetUserFromContext(ctx)
	if reporter.ID == input.UserID {
		return nil, huma.Error422UnprocessableEntity(
			"Você não pode denunciar a si mesmo",
		)
	}

	_, err := db.WithTx(ctx, func(tx *generated.Tx) (*generated.Report, error) {
		user, err := tx.User.Get(ctx, input.UserID)
		if err != nil {
			return nil, err
		}

		return tx.Report.Create().
			SetReporterID(reporter.ID).
			SetEntity(report.EntityUser).
			SetEntityID(user.ID).
			SetTargetUser(user).
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
