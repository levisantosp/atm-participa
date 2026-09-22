package issues

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/middlewares"
)

func Routes(api huma.API) {
	group := huma.NewGroup(api, "/issues")
	group.UseMiddleware(middlewares.Auth(api, false))

	huma.Get(group, "", GetIssues)
	huma.Register(group, huma.Operation{
		OperationID:   "create-issue",
		Method:        http.MethodPost,
		Path:          "",
		DefaultStatus: http.StatusCreated,
		MaxBodyBytes:  16_000_000,
	}, CreateIssue)
	huma.Post(group, "/{id}/upvote", UpvoteIssue)
	huma.Delete(group, "/{id}", DeleteIssue)
	huma.Put(group, "/{id}", EditIssue)
}
