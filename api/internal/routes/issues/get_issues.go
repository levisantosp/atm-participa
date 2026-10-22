package issues

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/dtos"
	"github.com/levisantosp/atm-participa/api/internal/utils"
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
	query := db.Client.Issue.Query().
		WithIssueFiles().
		Order(issue.ByID(sql.OrderDesc())).
		Limit(input.Limit + 1)

	if input.Cursor > 0 {
		query = query.Where(issue.IDLT(input.Cursor))
	}

	issues, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return &GetIssuesOutput{
		Body: utils.CursorPaginatedResponseFrom(
			dtos.IssuesFrom(issues),
			input.Limit,
			func(item dtos.Issue) string {
				return item.ID
			},
		),
	}, nil
}
