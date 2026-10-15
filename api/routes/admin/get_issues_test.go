package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/tests"
	"github.com/levisantosp/atm-participa/api/utils"
)

func TestGetIssues(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	adminUser := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(adminUser, t)

	issues := make([]int64, 0, 3)
	for range 3 {
		createdIssue, err := tests.CreateIssue(t, user)
		if err != nil {
			t.Fatal(err)
		}
		issues = append(issues, createdIssue.ID)
	}

	for id, fields := range map[int64][2]string{
		issues[0]: {"Broken streetlight", "Lamp near the park is out"},
		issues[1]: {"Flooded crossing", "Water leaks near school"},
	} {
		if _, err := db.Client.Issue.UpdateOneID(id).
			SetTitle(fields[0]).
			SetDescription(fields[1]).
			Save(t.Context()); err != nil {
			t.Fatal(err)
		}
	}

	for _, status := range []issue.Status{
		issue.StatusClosed,
		issue.StatusInReview,
	} {
		if _, err := tests.CreateIssueWithStatus(t, user, status); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		res := api.Get("/admin/issues")

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should reject non-admin request", func(t *testing.T) {
		res := api.Get("/admin/issues", tests.GetCookie(session.ID))

		if res.Code != http.StatusForbidden {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusForbidden,
				res.Code,
			)
		}
	})

	t.Run("should return open issues by default", func(t *testing.T) {
		res := api.Get("/admin/issues", tests.GetCookie(adminSession.ID))

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

		if body.NextCursor != nil {
			t.Fatal("expected no next page")
		}

		if len(body.Items) != 3 {
			t.Fatalf("expected %d items got %d", 3, len(body.Items))
		}
	})

	t.Run("should filter by status", func(t *testing.T) {
		res := api.Get(
			"/admin/issues?status=in_review",
			tests.GetCookie(adminSession.ID),
		)

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

		if len(body.Items) != 1 {
			t.Fatalf("expected %d item got %d", 1, len(body.Items))
		}

		if body.Items[0].Status != issue.StatusInReview {
			t.Fatalf(
				"expected status %s got %s",
				issue.StatusInReview,
				body.Items[0].Status,
			)
		}
	})

	t.Run("should reject invalid status", func(t *testing.T) {
		res := api.Get(
			"/admin/issues?status=invalid",
			tests.GetCookie(adminSession.ID),
		)

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run(
		"should search issue title and description case-insensitively",
		func(t *testing.T) {
			for _, testCase := range []struct {
				search string
				id     int64
			}{
				{search: "STREETLIGHT", id: issues[0]},
				{search: "LEAKS", id: issues[1]},
			} {
				res := api.Get(
					"/admin/issues?search="+testCase.search,
					tests.GetCookie(adminSession.ID),
				)

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

				if len(body.Items) != 1 {
					t.Fatalf(
						"expected %d item got %d for search %q",
						1,
						len(body.Items),
						testCase.search,
					)
				}

				if body.Items[0].ID != testCase.id {
					t.Fatalf(
						"expected issue %d got %d for search %q",
						testCase.id,
						body.Items[0].ID,
						testCase.search,
					)
				}
			}
		},
	)

	t.Run("should paginate with limit", func(t *testing.T) {
		res := api.Get(
			"/admin/issues?limit=2",
			tests.GetCookie(adminSession.ID),
		)

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

		if body.NextCursor == nil {
			t.Fatal("expected next page")
		}

		if len(body.Items) != 2 {
			t.Fatalf("expected %d items got %d", 2, len(body.Items))
		}

		if *body.NextCursor != body.Items[1].ID {
			t.Fatal("expected cursor to match the last returned item")
		}

		path := fmt.Sprintf("/admin/issues?limit=2&cursor=%d", *body.NextCursor)
		res = api.Get(path, tests.GetCookie(adminSession.ID))
		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}

		var nextBody utils.CursorPaginatedResponse[dtos.Issue]
		if err := json.NewDecoder(res.Body).Decode(&nextBody); err != nil {
			t.Fatal(err)
		}

		if nextBody.NextCursor != nil {
			t.Fatal("expected no next page")
		}

		if len(nextBody.Items) != 1 {
			t.Fatalf("expected 1 item got %d", len(nextBody.Items))
		}

		if nextBody.Items[0].ID >= *body.NextCursor {
			t.Fatal("expected the next page to contain older issues")
		}
	})

	t.Run("should reject limit greater than 100", func(t *testing.T) {
		res := api.Get(
			"/admin/issues?limit=101",
			tests.GetCookie(adminSession.ID),
		)

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject cursor less than 1", func(t *testing.T) {
		res := api.Get(
			"/admin/issues?cursor=-1",
			tests.GetCookie(adminSession.ID),
		)

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject search shorter than 2 characters", func(t *testing.T) {
		res := api.Get(
			"/admin/issues?search=a",
			tests.GetCookie(adminSession.ID),
		)

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run(
		"should reject search longer than 100 characters",
		func(t *testing.T) {
			res := api.Get(
				"/admin/issues?search="+strings.Repeat("a", 101),
				tests.GetCookie(adminSession.ID),
			)

			if res.Code != http.StatusUnprocessableEntity {
				t.Fatalf(
					"expected status %d got %d",
					http.StatusUnprocessableEntity,
					res.Code,
				)
			}
		},
	)
}
