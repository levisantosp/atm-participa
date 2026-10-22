package admin

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/dtos"
	"github.com/levisantosp/atm-participa/api/internal/r2"
	"github.com/levisantosp/atm-participa/api/internal/routes"
)

type GetIssueOutput struct {
	Body dtos.Issue
}

func GetIssue(
	ctx context.Context,
	input *struct {
		ID int64 `path:"id"`
	},
) (*GetIssueOutput, error) {
	issue, err := db.Client.Issue.Query().
		Where(issue.IDEQ(input.ID)).
		WithIssueFiles().
		First(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Ocorrência não encontrada")
		}
		return nil, err
	}

	item, err := routes.IssueResponseFrom(ctx, issue, r2.Client)
	if err != nil {
		return nil, err
	}

	return &GetIssueOutput{
		Body: *item,
	}, nil
}
