package admin

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/utils"
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
			return nil, huma.Error404NotFound("Demanda não encontrada")
		}

		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	return &UpdateIssueOutput{
		Body: dtos.IssueFrom(issue),
	}, nil
}
