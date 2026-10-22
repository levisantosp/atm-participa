package admin

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/r2"
	"github.com/levisantosp/atm-participa/api/routes"
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

	client, err := r2.New(ctx)
	if err != nil {
		return nil, err
	}

	item, err := routes.IssueResponseFrom(ctx, issue, client)
	if err != nil {
		return nil, err
	}

	return &GetIssueOutput{
		Body: *item,
	}, nil
}
