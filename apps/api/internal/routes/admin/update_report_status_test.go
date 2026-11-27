package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/report"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"

	_ "github.com/levisantosp/atm-participa/apps/api/ent/generated/runtime"
)

func TestUpdateReportStatus(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	adminUser := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(adminUser, t)

	item := tests.CreateReport(t, user, report.EntityUser, adminUser.ID, "Spam")
	path := fmt.Sprintf("/admin/reports/%d", item.ID)

	t.Run("should reject unauthenticated request", func(t *testing.T) {
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

	t.Run("should close and reopen a report", func(t *testing.T) {
		statuses := []report.Status{report.StatusClosed, report.StatusOpen}
		for _, status := range statuses {
			res := api.Patch(
				path,
				tests.GetCookie(adminSession.ID),
				map[string]any{
					"status": status,
				},
			)

			if res.Code != http.StatusOK {
				t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
			}

			var body dtos.Report
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}

			if body.ID != fmt.Sprint(item.ID) || body.Status != status {
				t.Fatalf("unexpected report: %+v", body)
			}

			stored := db.Client.Report.GetX(t.Context(), item.ID)
			if stored.Status != status {
				t.Fatalf(
					"expected stored status %s got %s",
					status,
					stored.Status,
				)
			}
		}
	})

	t.Run("should reject invalid status", func(t *testing.T) {
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

	t.Run("should return 404 for non-existent report", func(t *testing.T) {
		res := api.Patch(
			"/admin/reports/999999",
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
