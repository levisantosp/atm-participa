package users

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/danielgtaylor/huma/v2"

	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/ent/generated/user"
	"github.com/levisantosp/atm-participa/api/utils"
)

type GetUserIssuesOutputBody = utils.CursorPaginatedResponse[dtos.Issue]

type GetUserIssuesOutput struct {
	Body GetUserIssuesOutputBody
}

func GetIssues(
	ctx context.Context,
	input *struct {
		UserID int64  `path:"userId"`
		Status string `query:"status" enum:"open,closed,in_review"`
		Limit  int    `query:"limit" minimum:"1" maximum:"100" default:"10"`
		Cursor int64  `query:"cursor" minimum:"1"`
	},
) (*GetUserIssuesOutput, error) {
	query := db.Client.Issue.Query().
		Where(issue.HasUserWith(user.IDEQ(input.UserID))).
		Where(issue.IDGT(input.Cursor)).
		Order(issue.ByID(sql.OrderDesc())).
		Limit(input.Limit + 1)

	if input.Status != "" {
		query = query.Where(issue.StatusEQ(issue.Status(input.Status)))
	}

	issues, err := query.All(ctx)
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

	return &GetUserIssuesOutput{
		Body: utils.CursorPaginatedResponseFrom(items, input.Limit),
	}, nil
}
