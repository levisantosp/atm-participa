package admin

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/user"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/utils"
)

func GetUsers(
	ctx context.Context,
	input *struct {
		Limit  int    `query:"limit" minimum:"1" maximum:"100" default:"10"`
		Cursor int64  `query:"cursor" minimum:"1"`
		Search string `query:"search" minLength:"2" maxLength:"100"`
	},
) (*Output[utils.CursorPaginatedResponse[dtos.User]], error) {
	query := db.Client.User.Query().
		Order(user.ByID(sql.OrderDesc())).
		Limit(input.Limit + 1)

	if input.Cursor > 0 {
		query = query.Where(user.IDLT(input.Cursor))
	}

	if input.Search != "" {
		query = query.Where(
			user.Or(
				user.UsernameContainsFold(input.Search),
				user.EmailContainsFold(input.Search),
			),
		)
	}

	users, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return &Output[utils.CursorPaginatedResponse[dtos.User]]{
		Body: utils.CursorPaginatedResponseFrom(
			dtos.UsersFrom(users),
			input.Limit,
			func(item dtos.User) string {
				return item.ID
			},
		),
	}, nil
}
