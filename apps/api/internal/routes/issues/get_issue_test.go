package issues

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"
)

func TestGetIssue(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	huma.Register(api, huma.Operation{
		OperationID: "get-issue",
		Method:      http.MethodGet,
		Path:        "/issues/{id}",
	}, GetIssue)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)
	issue, err := tests.CreateIssue(t, user)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should return the issue", func(t *testing.T) {
		path := fmt.Sprintf("/issues/%d", issue.ID)
		res := api.Get(path, tests.GetCookie(session.ID))
		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}

		var body dtos.Issue
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.ID != fmt.Sprint(issue.ID) {
			t.Fatalf("expected issue ID %d got %s", issue.ID, body.ID)
		}
		if body.Title != issue.Title {
			t.Fatalf("expected title %q got %q", issue.Title, body.Title)
		}
		if body.Description != issue.Description {
			t.Fatalf(
				"expected description %q got %q",
				issue.Description,
				body.Description,
			)
		}
		if body.Author == nil || body.Author.ID != fmt.Sprint(user.ID) {
			t.Fatalf("expected author ID %d", user.ID)
		}
	})

	t.Run("should return 404 for non-existent issue", func(t *testing.T) {
		res := api.Get("/issues/999999", tests.GetCookie(session.ID))
		if res.Code != http.StatusNotFound {
			t.Fatalf("expected status %d got %d", http.StatusNotFound, res.Code)
		}
	})
}
