package routes

import "github.com/danielgtaylor/huma/v2"

func CreateRoutes(api huma.API) {
	group := huma.NewGroup(api, "/issues")

	huma.Get(group, "/{id}/{key}", GetIssueImage)
}
