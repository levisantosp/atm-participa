package issues

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/internal/dtos"
	"github.com/levisantosp/atm-participa/api/internal/tests"
	"github.com/levisantosp/atm-participa/api/internal/utils"
)

func TestGetIssues(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)

	t.Run("should use default values", func(t *testing.T) {
		openIssue, err := tests.CreateIssueWithStatus(t, user, issue.StatusOpen)
		if err != nil {
			t.Fatal(err)
		}
		unvotedIssue, err := tests.CreateIssueWithStatus(
			t,
			user,
			issue.StatusOpen,
		)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := tests.CreateIssueWithStatus(
			t,
			user,
			issue.StatusClosed,
		); err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("/issues/%d/upvote", openIssue.ID)
		upvoteRes := api.Post(path, tests.GetCookie(session.ID))
		if upvoteRes.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, upvoteRes.Code)
		}

		res := api.Get("/issues", tests.GetCookie(session.ID))
		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}

		var body utils.CursorPaginatedResponse[dtos.Issue]
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}

		if body.NextCursor != nil {
			t.Fatal("expected no next page")
		}
		if len(body.Items) != 2 {
			t.Fatalf("expected 2 open issues got %d", len(body.Items))
		}

		hasUpvotedByID := make(map[string]bool, len(body.Items))
		for _, item := range body.Items {
			hasUpvotedByID[item.ID] = item.HasUpvoted
		}
		if !hasUpvotedByID[fmt.Sprint(openIssue.ID)] {
			t.Fatal("expected hasUpvoted to be true for the upvoted issue")
		}
		if hasUpvotedByID[fmt.Sprint(unvotedIssue.ID)] {
			t.Fatal(
				"expected hasUpvoted to be false for the issue without an upvote",
			)
		}
	})

	t.Run("should paginate with cursor", func(t *testing.T) {
		var issueIDs []int64
		for range 3 {
			item, err := tests.CreateIssueWithStatus(
				t,
				user,
				issue.StatusInReview,
			)
			if err != nil {
				t.Fatal(err)
			}
			issueIDs = append(issueIDs, item.ID)
		}

		path := "/issues?limit=1&status=in_review"
		for i := len(issueIDs) - 1; i >= 0; i-- {
			res := api.Get(path, tests.GetCookie(session.ID))
			if res.Code != http.StatusOK {
				t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
			}

			var body utils.CursorPaginatedResponse[dtos.Issue]
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}

			if len(body.Items) != 1 {
				t.Fatalf("expected 1 item got %d", len(body.Items))
			}

			if body.Items[0].ID != fmt.Sprint(issueIDs[i]) {
				t.Fatalf(
					"expected issue %d got %s",
					issueIDs[i],
					body.Items[0].ID,
				)
			}

			if i == 0 {
				if body.NextCursor != nil {
					t.Fatal("expected no next page")
				}
				continue
			}

			if body.NextCursor == nil || *body.NextCursor != issueIDs[i] {
				t.Fatal("expected cursor to match the last returned item")
			}

			path = fmt.Sprintf(
				"/issues?limit=1&status=in_review&cursor=%d",
				*body.NextCursor,
			)
		}
	})

	t.Run("should reject invalid status", func(t *testing.T) {
		res := api.Get("/issues?status=invalid", tests.GetCookie(session.ID))

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject limit greater than 100", func(t *testing.T) {
		res := api.Get("/issues?limit=101", tests.GetCookie(session.ID))

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject cursor less than 1", func(t *testing.T) {
		res := api.Get("/issues?cursor=-1", tests.GetCookie(session.ID))

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		res := api.Get("/issues?cursor=-1")

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})
}
