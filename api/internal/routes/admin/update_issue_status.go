package admin

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/dtos"
)

func UpdateIssueStatus(
	ctx context.Context,
	input *struct {
		ID   int64 `path:"id"`
		Body struct {
			Status issue.Status `json:"status" enum:"open,closed,in_review"`
		}
	},
) (*Output[dtos.Issue], error) {
	issue, err := db.Client.Issue.UpdateOneID(input.ID).
		SetStatus(input.Body.Status).
		Save(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Ocorrência não encontrada")
		}

		return nil, err
	}

	return &Output[dtos.Issue]{
		Body: dtos.IssueFrom(issue),
	}, nil
}
