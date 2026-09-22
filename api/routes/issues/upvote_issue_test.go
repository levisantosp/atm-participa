package issues

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/tests"
	"github.com/levisantosp/atm-participa/api/utils"
)

func TestUpvoteIssue(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)

	issue, err := tests.CreateIssue(t, user)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		path := fmt.Sprintf("/issues/%d/upvote", issue.ID)
		res := api.Post(path)

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should upvote", func(t *testing.T) {
		path := fmt.Sprintf("/issues/%d/upvote", issue.ID)
		res := api.Post(path, tests.GetCookie(session.ID))

		if res.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusOK,
				res.Code,
			)
		}

		var body utils.CursorPaginatedResponse[dtos.Issue]
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should reject duplicate upvote", func(t *testing.T) {
		path := fmt.Sprintf("/issues/%d/upvote", issue.ID)
		res := api.Post(path, tests.GetCookie(session.ID))

		if res.Code != http.StatusConflict {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusConflict,
				res.Code,
			)
		}
	})

	t.Run("should reject invalid issue", func(t *testing.T) {
		res := api.Post("/issues/999999999/upvote", tests.GetCookie(session.ID))

		if res.Code != http.StatusNotFound {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusNotFound,
				res.Code,
			)
		}
	})
}
