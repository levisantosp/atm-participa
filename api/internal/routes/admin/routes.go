package admin

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/internal/middlewares"
)

type Output[T any] struct {
	Body T
}

func Routes(api huma.API) {
	group := huma.NewGroup(api, "/admin")
	group.UseMiddleware(middlewares.Auth(api, true))

	huma.Get(group, "/issues", GetIssues)
	huma.Get(group, "/issues/{id}", GetIssue)
	huma.Delete(group, "/issues/{id}", DeleteIssue)
	huma.Patch(group, "/issues/{id}", UpdateIssueStatus)
	huma.Get(group, "/issues/count/open", GetOpenIssuesCount)
	huma.Get(group, "/issues/count/in-review", GetInReviewIssuesCount)
}
