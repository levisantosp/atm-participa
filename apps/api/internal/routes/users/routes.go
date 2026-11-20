package users

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/internal/middlewares"
)

type Output[T any] struct {
	Body T
}

func Routes(api huma.API) {
	group := huma.NewGroup(api, "/users")
	group.UseMiddleware(middlewares.RateLimit(api, 60))
	group.UseMiddleware(middlewares.Auth(api, false))
	group.UseMiddleware(middlewares.Cache(api, 10*time.Second))

	huma.Get(group, "/{userId}/issues", GetIssues)
	huma.Delete(group, "", DeleteAccount)
	huma.Put(group, "/me", EditMe)
}
