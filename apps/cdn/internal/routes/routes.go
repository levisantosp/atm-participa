package routes

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/cdn/internal/middlewares"
)

func CreateRoutes(api huma.API) {
	group := huma.NewGroup(api, "/issues")
	group.UseMiddleware(middlewares.Cache(api, time.Hour))

	huma.Get(group, "/{id}/{key}", GetIssueImage)
}
