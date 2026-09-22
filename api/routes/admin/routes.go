package admin

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/middlewares"
)

func Routes(api huma.API) {
	group := huma.NewGroup(api, "/admin")
	group.UseMiddleware(middlewares.Auth(api, true))

	huma.Get(group, "/issues", GetIssues)
	huma.Delete(group, "/issues/{id}", DeleteIssue)
	huma.Patch(group, "/issues/{id}", UpdateIssueStatus)
}
