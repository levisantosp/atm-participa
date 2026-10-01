package admin

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
	adminUser := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(adminUser, t)

	for range 3 {
		if _, err := tests.CreateIssue(t, user); err != nil {
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

	t.Run("should return all issues", func(t *testing.T) {
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
}
