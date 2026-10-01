package admin

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
)

type UpdateIssueOutput struct {
	Body dtos.Issue
}

func UpdateIssueStatus(
	ctx context.Context,
	input *struct {
		ID   int64 `path:"id"`
		Body struct {
			Status issue.Status `json:"status" enum:"open,closed,in_review"`
		}
	},
) (*UpdateIssueOutput, error) {
	issue, err := db.Client.Issue.UpdateOneID(input.ID).
		SetStatus(input.Body.Status).
		Save(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Ocorrência não encontrada")
		}

		return nil, err
	}

	return &UpdateIssueOutput{
		Body: dtos.IssueFrom(issue),
	}, nil
}
