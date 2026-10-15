package admin

import (
	"context"

	"github.com/levisantosp/atm-participa/api/db"
)

type GetOpenIssuesCountOutput struct {
	Body int
}

func GetOpenIssuesCount(
	ctx context.Context,
	_ *struct{},
) (*GetOpenIssuesCountOutput, error) {
	count, err := db.Client.Issue.Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	return &GetOpenIssuesCountOutput{
		Body: count,
	}, nil
}
