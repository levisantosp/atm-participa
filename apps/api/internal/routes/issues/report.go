package issues

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
		ID int64 `path:"id" minimum:"1"`
	},
) (*struct{}, error) {
	user := middlewares.MustGetUserFromContext(ctx)

	_, err := db.WithTx(ctx, func(tx *generated.Tx) (*generated.Report, error) {
		issue, err := tx.Issue.Get(ctx, input.ID)
		if err != nil {
			return nil, err
		}

		return tx.Report.Create().
			SetUserID(user.ID).
			SetEntity(report.EntityIssue).
			SetEntityID(issue.ID).
			SetTargetIssue(issue).
			Save(ctx)
	})
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Ocorrência não encontrada")
		}

		return nil, err
	}

	return nil, nil
}
