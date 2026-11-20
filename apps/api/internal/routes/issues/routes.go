package issues

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/internal/middlewares"
)

type Output[T any] struct {
	Body T
}

func Routes(api huma.API) {
	group := huma.NewGroup(api, "/issues")
	group.UseMiddleware(middlewares.RateLimit(api, 60))
	group.UseMiddleware(middlewares.Auth(api, false))
	group.UseMiddleware(middlewares.Cache(api, 10*time.Second))

	huma.Register(group, huma.Operation{
		OperationID:   "get-issues",
		Method:        http.MethodGet,
		Path:          "",
		DefaultStatus: http.StatusOK,
	}, GetIssues)
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
