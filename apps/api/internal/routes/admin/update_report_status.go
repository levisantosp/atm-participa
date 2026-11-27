package admin

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/report"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
)

func UpdateReportStatus(
	ctx context.Context,
	input *struct {
		ID   int64 `path:"id"`
		Body struct {
			Status report.Status `json:"status" enum:"open,closed"`
		}
	},
) (*Output[dtos.Report], error) {
	report, err := db.Client.Report.UpdateOneID(input.ID).
		SetStatus(input.Body.Status).
		Save(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Denúncia não encontrada")
		}

		return nil, err
	}

	return &Output[dtos.Report]{
		Body: dtos.ReportFrom(report),
	}, nil
}
