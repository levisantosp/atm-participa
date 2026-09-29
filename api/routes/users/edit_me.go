package users

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/errors"
	"github.com/levisantosp/atm-participa/api/middlewares"
)

type EditMeOutput struct {
	Body dtos.User
}

func EditMe(
	ctx context.Context,
	input *struct {
		Body struct {
			Username    *string `json:"username" minLength:"3" maxLength:"32" pattern:"^[a-zA-Z0-9_]+$" required:"false"`
			DisplayName *string `json:"displayName" minLength:"1" maxLength:"100" required:"false"`
		}
	},
) (*EditMeOutput, error) {
	if input.Body.Username == nil && input.Body.DisplayName == nil {
		return nil, huma.Error400BadRequest(
			"Informe ao menos um campo para atualizar",
		)
	}

	userCtx := middlewares.MustGetUserFromContext(ctx)
	user, err := db.Client.User.UpdateOneID(userCtx.ID).
		SetNillableDisplayName(input.Body.DisplayName).
		SetNillableUsername(input.Body.Username).
		Save(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound(errors.NotFound)
		}

		return nil, err
	}

	return &EditMeOutput{
		Body: dtos.UserFrom(user),
	}, nil
}
