package users

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/report"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"

	_ "github.com/levisantosp/atm-participa/apps/api/ent/generated/runtime"
)

func TestReportUser(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	target := tests.CreateUser(t)
	user := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	path := fmt.Sprintf("/users/%d/report", target.ID)

	const reason = "Conteúdo ofensivo"

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

	t.Run("should reject an invalid user id", func(t *testing.T) {
		res := api.Post(
			"/users/0/report",
			tests.GetCookie(session.ID),
			map[string]any{"reason": reason},
		)
		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should return not found for a missing user", func(t *testing.T) {
		res := api.Post(
			"/users/999999999/report",
			tests.GetCookie(session.ID),
			map[string]any{"reason": reason},
		)
		if res.Code != http.StatusNotFound {
			t.Fatalf("expected status %d got %d", http.StatusNotFound, res.Code)
		}
	})

	t.Run("should reject reporting yourself", func(t *testing.T) {
		res := api.Post(
			fmt.Sprintf("/users/%d/report", user.ID),
			tests.GetCookie(session.ID),
			map[string]any{"reason": reason},
		)
		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject missing or blank reason", func(t *testing.T) {
		for _, body := range []map[string]any{{}, {"reason": "   "}} {
			res := api.Post(path, tests.GetCookie(session.ID), body)
			if res.Code != http.StatusUnprocessableEntity {
				t.Fatalf(
					"expected status %d got %d",
					http.StatusUnprocessableEntity,
					res.Code,
				)
			}
		}
	})

	t.Run("should report another user", func(t *testing.T) {
		res := api.Post(
			path,
			tests.GetCookie(session.ID),
			map[string]any{"reason": reason},
		)
		if res.Code != http.StatusCreated {
			t.Fatalf("expected status %d got %d", http.StatusCreated, res.Code)
		}

		saved, err := user.QueryReports().Only(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		if saved.Reason != reason {
			t.Errorf("expected reason %q got %q", reason, saved.Reason)
		}
		if saved.UserID == nil || *saved.UserID != user.ID {
			t.Errorf("expected user id %d, got %v", user.ID, saved)
		}
		if saved.Entity != report.EntityUser || saved.EntityID != target.ID {
			t.Errorf(
				"expected user %d, got %s %d",
				target.ID,
				saved.Entity,
				saved.EntityID,
			)
		}
		if saved.TargetUserID == nil || *saved.TargetUserID != target.ID {
			t.Errorf("expected target user id %d, got %v", target.ID, saved)
		}
		if saved.TargetIssueID != nil {
			t.Errorf("unexpected target issue id %d", *saved.TargetIssueID)
		}
	})
}
