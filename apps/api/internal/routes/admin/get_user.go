package admin

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
)

func GetUser(
	ctx context.Context,
	input *struct {
		ID int64 `path:"id" minimum:"1"`
	},
) (*Output[dtos.User], error) {
	user, err := db.Client.User.Get(ctx, input.ID)
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
