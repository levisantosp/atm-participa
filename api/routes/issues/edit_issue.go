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
)

type EditIssueOutput struct {
	Body dtos.Issue
}

func EditIssue(ctx context.Context, input *struct {
	ID   int64 `path:"id"`
	Body struct {
		Title       string `json:"title" maxLength:"72" minLength:"3"`
		Description string `json:"description" maxLength:"65000" minLength:"10"`
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
			return nil, huma.Error404NotFound("Ocorrência não encontrada")
		}

		return nil, err
	}

	return &EditIssueOutput{
		Body: dtos.IssueFrom(issueEntity),
	}, nil
}
