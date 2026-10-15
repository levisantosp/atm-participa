package admin

import (
	"context"

	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
)

func GetInReviewIssuesCount(
	ctx context.Context,
	_ *struct{},
) (*GetOpenIssuesCountOutput, error) {
	count, err := db.Client.Issue.Query().
		Where(issue.StatusEQ(issue.StatusInReview)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	return &GetOpenIssuesCountOutput{
		Body: count,
	}, nil
}
