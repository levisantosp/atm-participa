package admin

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated"
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
	issue, err := db.Client.Issue.Get(ctx, input.ID)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Ocorrência não encontrada")
		}
		return nil, err
	}

	return &GetIssueOutput{
		Body: dtos.IssueFrom(issue),
	}, nil
}
