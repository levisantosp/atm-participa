package admin

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/dtos"
	"github.com/levisantosp/atm-participa/api/internal/utils"
)

func GetIssues(
	ctx context.Context,
	input *struct {
		Limit  int          `query:"limit" minimum:"1" maximum:"100" default:"10"`
		Cursor int64        `query:"cursor" minimum:"1"`
		Status issue.Status `query:"status" default:"open"`
		Search string       `query:"search" minLength:"2" maxLength:"100"`
	},
) (*Output[utils.CursorPaginatedResponse[dtos.Issue]], error) {
	query := db.Client.Issue.Query().
		WithIssueFiles().
		Where(issue.StatusEQ(input.Status)).
		Order(issue.ByID(sql.OrderDesc())).
		Limit(input.Limit + 1)

	if input.Cursor > 0 {
		query = query.Where(issue.IDLT(input.Cursor))
	}

	if input.Search != "" {
		query = query.Where(
			issue.Or(
				issue.TitleContainsFold(input.Search),
				issue.DescriptionContainsFold(input.Search),
			),
		)
	}

	issues, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	items := dtos.IssuesFrom(issues)

	return &Output[utils.CursorPaginatedResponse[dtos.Issue]]{
		Body: utils.CursorPaginatedResponseFrom(
			items,
			input.Limit,
			func(item dtos.Issue) string {
				return item.ID
			},
		),
	}, nil
}
