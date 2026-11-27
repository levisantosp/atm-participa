package users

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/user"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
)

func GetUser(
	ctx context.Context,
	input *struct {
		Username string `path:"username"`
	},
) (*Output[dtos.User], error) {
	user, err := db.Client.User.Query().
		Where(
			user.UsernameEQ(input.Username),
			user.Or(user.BannedIsNil(), user.BannedNEQ(true)),
		).
		First(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("Usuário não encontrado")
		}

		return nil, err
	}

	return &Output[dtos.User]{
		Body: dtos.UserFrom(user),
	}, nil
}
