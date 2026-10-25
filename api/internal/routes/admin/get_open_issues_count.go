package admin

import (
	"context"

	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/internal/db"
)

func GetOpenIssuesCount(
	ctx context.Context,
	_ *struct{},
) (*Output[int], error) {
	count, err := db.Client.Issue.Query().
		Where(issue.StatusEQ(issue.StatusOpen)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	return &Output[int]{
		Body: count,
	}, nil
}
