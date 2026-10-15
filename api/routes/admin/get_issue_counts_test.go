package admin

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/ent/generated/issue"
	"github.com/levisantosp/atm-participa/api/tests"
)

func TestGetOpenIssuesCount(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	adminUser := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(adminUser, t)

	for _, status := range []issue.Status{
		issue.StatusOpen,
		issue.StatusOpen,
		issue.StatusClosed,
		issue.StatusInReview,
	} {
		if _, err := tests.CreateIssueWithStatus(t, user, status); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		res := api.Get("/admin/issues/count/open")

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should reject non-admin request", func(t *testing.T) {
		res := api.Get(
			"/admin/issues/count/open",
			tests.GetCookie(session.ID),
		)

		if res.Code != http.StatusForbidden {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusForbidden,
				res.Code,
			)
		}
	})

	t.Run("should count only open issues", func(t *testing.T) {
		res := api.Get(
			"/admin/issues/count/open",
			tests.GetCookie(adminSession.ID),
		)

		if res.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusOK,
				res.Code,
			)
		}

		var count int
		if err := json.NewDecoder(res.Body).Decode(&count); err != nil {
			t.Fatal(err)
		}

		if count != 2 {
			t.Fatalf("expected count %d got %d", 2, count)
		}
	})
}

func TestGetInReviewIssuesCount(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	adminUser := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(adminUser, t)

	for _, status := range []issue.Status{
		issue.StatusInReview,
		issue.StatusInReview,
		issue.StatusOpen,
		issue.StatusClosed,
	} {
		if _, err := tests.CreateIssueWithStatus(t, user, status); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		res := api.Get("/admin/issues/count/in-review")

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should reject non-admin request", func(t *testing.T) {
		res := api.Get(
			"/admin/issues/count/in-review",
			tests.GetCookie(session.ID),
		)

		if res.Code != http.StatusForbidden {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusForbidden,
				res.Code,
			)
		}
	})

	t.Run("should count only in-review issues", func(t *testing.T) {
		res := api.Get(
			"/admin/issues/count/in-review",
			tests.GetCookie(adminSession.ID),
		)

		if res.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusOK,
				res.Code,
			)
		}

		var count int
		if err := json.NewDecoder(res.Body).Decode(&count); err != nil {
			t.Fatal(err)
		}

		if count != 2 {
			t.Fatalf("expected count %d got %d", 2, count)
		}
	})
}
