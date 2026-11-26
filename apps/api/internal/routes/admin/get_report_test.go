package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/report"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"
)

func TestGetReport(t *testing.T) {
	tests.Setup(t)
	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	admin := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(admin, t)
	tests.CreateReport(t, admin, report.EntityUser, user.ID, "Spam")
	target := tests.CreateReport(
		t,
		user,
		report.EntityUser,
		admin.ID,
		"Assédio",
	)
	path := fmt.Sprintf("/admin/reports/%d", target.ID)

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

	t.Run(
		"should return the requested report with its author",
		func(t *testing.T) {
			res := api.Get(path, tests.GetCookie(adminSession.ID))
			if res.Code != http.StatusOK {
				t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
			}

			var body dtos.Report
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body.ID != fmt.Sprint(target.ID) ||
				body.Reason != target.Reason {
				t.Fatalf("unexpected report: %+v", body)
			}
			if body.Author == nil || body.Author.ID != fmt.Sprint(user.ID) {
				t.Fatalf("expected author %d got %v", user.ID, body.Author)
			}
		},
	)

	t.Run("should return not found for a missing report", func(t *testing.T) {
		res := api.Get(
			"/admin/reports/999999",
			tests.GetCookie(adminSession.ID),
		)
		if res.Code != http.StatusNotFound {
			t.Fatalf("expected status %d got %d", http.StatusNotFound, res.Code)
		}
	})
}
