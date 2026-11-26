package admin

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/report"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/utils"
)

func GetReports(
	ctx context.Context,
	input *struct {
		Limit  int           `query:"limit" minimum:"1" maximum:"100" default:"10"`
		Cursor int64         `query:"cursor" minimum:"1"`
		Entity report.Entity `query:"entity" enum:"user,issue"`
	},
) (*Output[utils.CursorPaginatedResponse[dtos.Report]], error) {
	query := db.Client.Report.Query().
		WithUser().
		Order(report.ByID(sql.OrderDesc())).
		Limit(input.Limit + 1)

	if input.Cursor > 0 {
		query = query.Where(report.IDLT(input.Cursor))
	}

	if input.Entity != "" {
		query = query.Where(report.EntityEQ(input.Entity))
	}

	reports, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return &Output[utils.CursorPaginatedResponse[dtos.Report]]{
		Body: utils.CursorPaginatedResponseFrom(
			dtos.ReportsFrom(reports),
			input.Limit,
			func(item dtos.Report) string {
				return item.ID
			},
		),
	}, nil
}
