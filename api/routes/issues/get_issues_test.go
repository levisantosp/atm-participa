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

func TestGetIssues(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)

	t.Run("should use default values", func(t *testing.T) {
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
	})

	t.Run("should paginate with cursor", func(t *testing.T) {
		var issueIDs []int64
		for range 3 {
			item, err := tests.CreateIssue(t, user)
			if err != nil {
				t.Fatal(err)
			}
			issueIDs = append(issueIDs, item.ID)
		}

		path := "/issues?limit=1"
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

			if body.Items[0].ID != issueIDs[i] {
				t.Fatalf(
					"expected issue %d got %d",
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

			path = fmt.Sprintf("/issues?limit=1&cursor=%d", *body.NextCursor)
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
