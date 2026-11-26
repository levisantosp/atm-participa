package issues

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/report"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"

	_ "github.com/levisantosp/atm-participa/apps/api/ent/generated/runtime"
)

func TestReportIssue(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	author := tests.CreateUser(t)
	user := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)

	issue, err := tests.CreateIssue(t, author)
	if err != nil {
		t.Fatal(err)
	}
	path := fmt.Sprintf("/issues/%d/report", issue.ID)

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		res := api.Post(path)
		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should reject an invalid issue id", func(t *testing.T) {
		res := api.Post("/issues/0/report", tests.GetCookie(session.ID))
		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should return not found for a missing issue", func(t *testing.T) {
		res := api.Post("/issues/999999999/report", tests.GetCookie(session.ID))
		if res.Code != http.StatusNotFound {
			t.Fatalf("expected status %d got %d", http.StatusNotFound, res.Code)
		}
	})

	t.Run("should report an issue", func(t *testing.T) {
		res := api.Post(path, tests.GetCookie(session.ID))
		if res.Code != http.StatusCreated {
			t.Fatalf("expected status %d got %d", http.StatusCreated, res.Code)
		}

		saved, err := db.Client.Report.Query().Only(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		if saved.UserID == nil || *saved.UserID != user.ID {
			t.Errorf("expected user id %d, got %v", user.ID, saved)
		}
		if saved.Entity != report.EntityIssue || saved.EntityID != issue.ID {
			t.Errorf(
				"expected issue %d, got %s %d",
				issue.ID,
				saved.Entity,
				saved.EntityID,
			)
		}
		if saved.TargetIssueID == nil || *saved.TargetIssueID != issue.ID {
			t.Errorf("expected target issue id %d, got %v", issue.ID, saved)
		}
		if saved.TargetUserID != nil {
			t.Errorf("unexpected target user id %d", *saved.TargetUserID)
		}
	})
}
