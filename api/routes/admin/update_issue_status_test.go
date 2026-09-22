package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/dtos"
	"github.com/levisantosp/atm-participa/api/tests"

	_ "github.com/levisantosp/atm-participa/api/ent/generated/runtime"
)

func TestUpdateIssueStatus(t *testing.T) {
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

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		path := fmt.Sprintf("/admin/issues/%d", issue.ID)
		res := api.Patch(path)

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should reject non-admin request", func(t *testing.T) {
		path := fmt.Sprintf("/admin/issues/%d", issue.ID)
		res := api.Patch(path, tests.GetCookie(session.ID), map[string]any{
			"status": "closed",
		})

		if res.Code != http.StatusForbidden {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusForbidden,
				res.Code,
			)
		}
	})

	t.Run("should update issue status", func(t *testing.T) {
		path := fmt.Sprintf("/admin/issues/%d", issue.ID)
		res := api.Patch(path, tests.GetCookie(adminSession.ID), map[string]any{
			"status": "closed",
		})

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

		if body.Status != "closed" {
			t.Fatalf("expected status closed got %s", body.Status)
		}
	})

	t.Run("should reject invalid status", func(t *testing.T) {
		path := fmt.Sprintf("/admin/issues/%d", issue.ID)
		res := api.Patch(path, tests.GetCookie(adminSession.ID), map[string]any{
			"status": "banana",
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should return 404 for non-existent issue", func(t *testing.T) {
		res := api.Patch(
			"/admin/issues/999999",
			tests.GetCookie(adminSession.ID),
			map[string]any{
				"status": "closed",
			},
		)

		if res.Code != http.StatusNotFound {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusNotFound,
				res.Code,
			)
		}
	})
}
