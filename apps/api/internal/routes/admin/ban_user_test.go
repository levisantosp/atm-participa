package admin

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"

	_ "github.com/levisantosp/atm-participa/apps/api/ent/generated/runtime"
)

func TestBanUser(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	admin := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(admin, t)
	path := fmt.Sprintf("/admin/users/%d/ban", user.ID)
	endsAt := time.Now().UTC().Add(time.Hour).Truncate(time.Second)

	t.Run("should reject unauthenticated requests", func(t *testing.T) {
		res := api.Post(path)
		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should reject non-admin requests", func(t *testing.T) {
		res := api.Post(
			path,
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

	t.Run("should forbid banning your own account", func(t *testing.T) {
		res := api.Post(
			fmt.Sprintf("/admin/users/%d/ban", admin.ID),
			tests.GetCookie(adminSession.ID),
			map[string]any{"banEndsAt": endsAt},
		)
		if res.Code != http.StatusForbidden {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusForbidden,
				res.Code,
			)
		}

		saved := db.Client.User.GetX(t.Context(), admin.ID)
		if saved.Banned != nil || saved.BannedAt != nil ||
			saved.BanEndsAt != nil {
			t.Fatal("expected ban fields to remain unset")
		}
	})

	t.Run("should return not found for a missing user", func(t *testing.T) {
		res := api.Post(
			"/admin/users/999999999/ban",
			tests.GetCookie(adminSession.ID),
		)
		if res.Code != http.StatusNotFound {
			t.Fatalf("expected status %d got %d", http.StatusNotFound, res.Code)
		}
	})

	t.Run("should reject an expired ban", func(t *testing.T) {
		res := api.Post(path, tests.GetCookie(adminSession.ID), map[string]any{
			"banEndsAt": time.Now().Add(-time.Hour),
		})
		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should ban a user until the given date", func(t *testing.T) {
		res := api.Post(
			path,
			tests.GetCookie(adminSession.ID),
			map[string]any{"banEndsAt": endsAt},
		)
		if res.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusOK,
				res.Code,
			)
		}
		if res := api.Get(
			"/admin/users",
			tests.GetCookie(session.ID),
		); res.Code != http.StatusUnauthorized {
			t.Fatalf("expected revoked session to return 401, got %d", res.Code)
		}
		saved := db.Client.User.GetX(t.Context(), user.ID)
		if saved.Banned == nil || !*saved.Banned || saved.BannedAt == nil ||
			saved.BanEndsAt == nil || !saved.BanEndsAt.Equal(endsAt) {
			t.Fatalf("expected temporary ban until %s, got %v", endsAt, saved)
		}
	})

	t.Run(
		"should keep the existing ban end date when omitted",
		func(t *testing.T) {
			db.Client.User.UpdateOneID(user.ID).
				SetBanEndsAt(endsAt).
				ExecX(t.Context())
			res := api.Post(
				path,
				tests.GetCookie(adminSession.ID),
			)
			if res.Code != http.StatusOK {
				t.Fatalf(
					"expected status %d got %d",
					http.StatusOK,
					res.Code,
				)
			}
			saved := db.Client.User.GetX(t.Context(), user.ID)
			if saved.Banned == nil || !*saved.Banned || saved.BannedAt == nil ||
				saved.BanEndsAt == nil || !saved.BanEndsAt.Equal(endsAt) {
				t.Fatalf("expected ban end date %s, got %v", endsAt, saved)
			}
		},
	)
}
