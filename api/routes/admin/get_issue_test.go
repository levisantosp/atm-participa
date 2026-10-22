package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/tests"
)

func TestGetIssue(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	adminUser := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(adminUser, t)

	issue, err := tests.CreateIssue(t, user)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("/admin/issues/%d", issue.ID)

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		res := api.Get(path)

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should reject non-admin request", func(t *testing.T) {
		res := api.Get(path, tests.GetCookie(session.ID))

		if res.Code != http.StatusForbidden {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusForbidden,
				res.Code,
			)
		}
	})

	t.Run("should return the issue", func(t *testing.T) {
		res := api.Get(path, tests.GetCookie(adminSession.ID))

		if res.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusOK,
				res.Code,
			)
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

		if body.Status != issue.Status {
			t.Fatalf("expected status %s got %s", issue.Status, body.Status)
		}
	})

	t.Run("should return 404 for non-existent issue", func(t *testing.T) {
		res := api.Get("/admin/issues/999999", tests.GetCookie(adminSession.ID))

		if res.Code != http.StatusNotFound {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusNotFound,
				res.Code,
			)
		}
	})
}
