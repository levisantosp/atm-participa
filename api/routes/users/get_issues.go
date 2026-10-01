package users

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/ent/generated/user"
	"github.com/levisantosp/atm-participa/api/r2"
	"github.com/levisantosp/atm-participa/api/routes"
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
		WithIssueFiles().
		Order(issue.ByID(sql.OrderDesc())).
		Limit(input.Limit + 1)

	if input.Cursor > 0 {
		query = query.Where(issue.IDLT(input.Cursor))
	}

	if input.Status != "" {
		query = query.Where(issue.StatusEQ(issue.Status(input.Status)))
	}

	issues, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	client, err := r2.New(ctx)
	if err != nil {
		return nil, err
	}

	items, err := routes.IssueResponsesFrom(ctx, issues, client)
	if err != nil {
		return nil, err
	}

	return &GetUserIssuesOutput{
		Body: utils.CursorPaginatedResponseFrom(
			items,
			input.Limit,
			func(item dtos.Issue) int64 { return item.ID },
		),
	}, nil
}
