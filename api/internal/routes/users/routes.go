package users

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/internal/middlewares"
)

type Output[T any] struct {
	Body T
}

func Routes(api huma.API) {
	group := huma.NewGroup(api, "/users")
	group.UseMiddleware(middlewares.Auth(api, false))

	huma.Get(group, "/{userId}/issues", GetIssues)
	huma.Delete(group, "", DeleteAccount)
	huma.Put(group, "/me", EditMe)
}
