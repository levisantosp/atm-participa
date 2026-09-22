package admin

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/utils"
)

type GetIssuesOutputBody = utils.CursorPaginatedResponse[dtos.Issue]

type GetIssuesOutput struct {
	Body GetIssuesOutputBody
}

func GetIssues(
	ctx context.Context,
	input *struct {
		Limit  int   `query:"limit" minimum:"1" maximum:"100" default:"10"`
		Cursor int64 `query:"cursor" minimum:"1"`
	},
) (*GetIssuesOutput, error) {
	issues, err := db.Client.Issue.Query().
		Where(issue.IDGT(input.Cursor)).
		Order(issue.ByID(sql.OrderDesc())).
		Limit(input.Limit + 1).
		All(ctx)
	if err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	items := make([]dtos.Issue, 0, len(issues))
	for _, item := range issues {
		items = append(items, dtos.IssueFrom(item))
	}

	return &GetIssuesOutput{
		Body: utils.CursorPaginatedResponseFrom(items, input.Limit),
	}, nil
}
