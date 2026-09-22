package issues

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/middlewares"
	"github.com/levisantosp/atm-participa/api/utils"
)

func UpvoteIssue(
	ctx context.Context,
	input *struct {
		IssueID int64 `path:"id"`
	},
) (*CreateIssueOutput, error) {
	user := middlewares.MustGetUserFromContext(ctx)

	issue, err := db.WithTx(
		ctx,
		func(tx *generated.Tx) (*generated.Issue, error) {
			issue, err := tx.Issue.Get(ctx, input.IssueID)
			if err != nil {
				return nil, err
			}

			_, err = tx.Upvote.Create().
				SetUserID(user.ID).
				SetIssueID(issue.ID).
				Save(ctx)
			if err != nil {
				return nil, err
			}

			issue, err = tx.Issue.UpdateOneID(issue.ID).AddUpvotes(1).Save(ctx)
			if err != nil {
				return nil, err
			}

			return issue, nil
		},
	)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Demanda não encontrada")
		}

		if generated.IsConstraintError(err) {
			return nil, huma.Error409Conflict(
				"Você só pode apoiar uma demanda uma unica vez.",
			)
		}

		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	return &CreateIssueOutput{
		Body: dtos.IssueFrom(issue),
	}, nil
}
