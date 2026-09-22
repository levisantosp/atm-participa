package issues

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/ent/generated/user"
	"github.com/levisantosp/atm-participa/api/middlewares"
	"github.com/levisantosp/atm-participa/api/utils"
)

type EditIssueOutput struct {
	Body dtos.Issue
}

func EditIssue(ctx context.Context, input *struct {
	ID   int64 `path:"id"`
	Body struct {
		Title       string `json:"title" maxLength:"72" minLength:"3" required:"true"`
		Description string `json:"description" maxLength:"65000" minLength:"10" required:"true"`
	}
},
) (*EditIssueOutput, error) {
	userCtx := middlewares.MustGetUserFromContext(ctx)

	issueEntity, err := db.Client.Issue.UpdateOneID(input.ID).
		Where(issue.HasUserWith(user.IDEQ(userCtx.ID))).
		SetTitle(input.Body.Title).
		SetDescription(input.Body.Description).
		Save(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, utils.LogErr(
				huma.Error404NotFound("Demanda não encontrada"),
				err,
			)
		}

		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	return &EditIssueOutput{
		Body: dtos.IssueFrom(issueEntity),
	}, nil
}
