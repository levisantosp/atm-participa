package users

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/predicate"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/upvote"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/middlewares"
	"github.com/levisantosp/atm-participa/apps/api/internal/utils"
	"github.com/levisantosp/atm-participa/pkg/infra/errors"
)

func GetUpvotes(
	ctx context.Context,
	input *struct {
		UserID    int64   `path:"userId"`
		IssuesIDs []int64 `query:"issuesIds" maxItems:"100"`
		Cursor    int64   `query:"cursor" minimum:"1"`
		Limit     int     `query:"limit" minimum:"1" maximum:"100" default:"10"`
	},
) (*Output[utils.CursorPaginatedResponse[dtos.Upvote]], error) {
	user := middlewares.MustGetUserFromContext(ctx)
	if user.ID != input.UserID {
		return nil, huma.Error403Forbidden(errors.Forbidden)
	}

	conditions := []predicate.Upvote{upvote.UserIDEQ(input.UserID)}

	if len(input.IssuesIDs) > 0 {
		conditions = append(conditions, upvote.IssueIDIn(input.IssuesIDs...))
		input.Limit = len(input.IssuesIDs)
	}

	query := db.Client.Upvote.Query().
		Where(conditions...).
		Order(upvote.ByID(sql.OrderDesc())).
		Limit(input.Limit + 1)

	if input.Cursor > 0 {
		query = query.Where(upvote.IDLT(input.Cursor))
	}

	upvotes, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return &Output[utils.CursorPaginatedResponse[dtos.Upvote]]{
		Body: utils.CursorPaginatedResponseFrom(
			dtos.UpvotesFrom(upvotes),
			input.Limit,
			func(u dtos.Upvote) string {
				return u.ID
			},
		),
	}, nil
}
