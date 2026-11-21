package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"
	"github.com/levisantosp/atm-participa/apps/api/internal/utils"
)

func TestGetUpvotes(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)
	var upvotes []*generated.Upvote

	for range 3 {
		issue, err := tests.CreateIssue(t, user)
		if err != nil {
			t.Fatal(err)
		}

		upvote, err := db.Client.Upvote.Create().
			SetUserID(user.ID).
			SetIssueID(issue.ID).
			Save(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		upvotes = append(upvotes, upvote)
	}

	t.Run(
		"should return the user's upvotes with default pagination",
		func(t *testing.T) {
			path := fmt.Sprintf("/users/%d/upvotes?limit=20", user.ID)

			res := api.Get(path, tests.GetCookie(session.ID))
			if res.Code != http.StatusOK {
				t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
			}

			var body utils.CursorPaginatedResponse[dtos.Upvote]
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if len(body.Items) != len(upvotes) {
				t.Fatalf(
					"expected %d upvotes got %d",
					len(upvotes),
					len(body.Items),
				)
			}
			if body.NextCursor != nil {
				t.Fatal("expected no next page")
			}
		},
	)

	t.Run("should filter by issue IDs", func(t *testing.T) {
		path := fmt.Sprintf(
			"/users/%d/upvotes?issuesIds=%d",
			user.ID,
			upvotes[1].IssueID,
		)

		res := api.Get(path, tests.GetCookie(session.ID))
		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}

		var body utils.CursorPaginatedResponse[dtos.Upvote]
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Items) != 1 {
			t.Fatalf("expected 1 upvote got %d", len(body.Items))
		}
		if body.Items[0].IssueID != fmt.Sprint(upvotes[1].IssueID) {
			t.Fatalf(
				"expected issue ID %d got %s",
				upvotes[1].IssueID,
				body.Items[0].IssueID,
			)
		}
	})

	t.Run("should paginate upvotes with a cursor", func(t *testing.T) {
		path := fmt.Sprintf("/users/%d/upvotes?limit=1", user.ID)

		res := api.Get(path, tests.GetCookie(session.ID))
		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}

		var firstPage utils.CursorPaginatedResponse[dtos.Upvote]
		if err := json.NewDecoder(res.Body).Decode(&firstPage); err != nil {
			t.Fatal(err)
		}
		if len(firstPage.Items) != 1 {
			t.Fatalf("expected 1 upvote got %d", len(firstPage.Items))
		}
		if firstPage.NextCursor == nil {
			t.Fatal("expected a next cursor")
		}

		path = fmt.Sprintf(
			"/users/%d/upvotes?limit=1&cursor=%d",
			user.ID,
			*firstPage.NextCursor,
		)

		res = api.Get(path, tests.GetCookie(session.ID))
		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}

		var secondPage utils.CursorPaginatedResponse[dtos.Upvote]
		if err := json.NewDecoder(res.Body).Decode(&secondPage); err != nil {
			t.Fatal(err)
		}
		if len(secondPage.Items) != 1 {
			t.Fatalf("expected 1 upvote got %d", len(secondPage.Items))
		}
		firstID, err := strconv.ParseInt(firstPage.Items[0].ID, 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		secondID, err := strconv.ParseInt(secondPage.Items[0].ID, 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		if secondID >= firstID {
			t.Fatalf(
				"expected next upvote ID below %s, got %s",
				firstPage.Items[0].ID,
				secondPage.Items[0].ID,
			)
		}
	})

	t.Run("should forbid access to another user's upvotes", func(t *testing.T) {
		user := tests.CreateAdminUser(t)
		path := fmt.Sprintf("/users/%d/upvotes?limit=20", user.ID)

		res := api.Get(path, tests.GetCookie(session.ID))
		if res.Code != http.StatusForbidden {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusForbidden,
				res.Code,
			)
		}
	})

	t.Run("should reject unauthenticated requests", func(t *testing.T) {
		path := fmt.Sprintf("/users/%d/upvotes?limit=20", user.ID)

		res := api.Get(path)
		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})
}
