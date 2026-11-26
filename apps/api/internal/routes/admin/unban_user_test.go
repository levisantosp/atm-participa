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

func TestUnbanUser(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	admin := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(admin, t)
	path := fmt.Sprintf("/admin/users/%d/unban", user.ID)

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
			fmt.Sprintf("/admin/users/%d/unban", admin.ID),
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

	t.Run("should reject invalid user id", func(t *testing.T) {
		res := api.Post(
			"/admin/users/0/unban",
			tests.GetCookie(adminSession.ID),
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
			"/admin/users/999999999/unban",
			tests.GetCookie(adminSession.ID),
		)
		if res.Code != http.StatusNotFound {
			t.Fatalf("expected status %d got %d", http.StatusNotFound, res.Code)
		}
	})

	t.Run(
		"should succeed when the user is already unbanned",
		func(t *testing.T) {
			res := api.Post(path, tests.GetCookie(adminSession.ID))
			if res.Code != http.StatusOK {
				t.Fatalf(
					"expected status %d got %d",
					http.StatusOK,
					res.Code,
				)
			}
		},
	)

	bannedAt := time.Now()
	endsAt := bannedAt.Add(time.Hour)
	for _, tc := range []struct {
		name   string
		endsAt *time.Time
	}{
		{"should lift a permanent ban", nil},
		{"should lift a temporary ban", &endsAt},
	} {
		t.Run(tc.name, func(t *testing.T) {
			user.Update().
				SetBanned(true).
				SetBannedAt(bannedAt).
				SetNillableBanEndsAt(tc.endsAt).
				ExecX(t.Context())

			session := tests.CreateSession(user, t)
			res := api.Post(path, tests.GetCookie(adminSession.ID))
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
				t.Fatalf(
					"expected revoked session to return 401, got %d",
					res.Code,
				)
			}
			saved := db.Client.User.GetX(t.Context(), user.ID)
			if saved.Banned != nil || saved.BannedAt != nil ||
				saved.BanEndsAt != nil {
				t.Fatal("expected all ban fields to be cleared")
			}
		})
	}

	t.Run("should forbid unbanning your own account", func(t *testing.T) {
		admin.Update().SetBanned(true).SetBannedAt(bannedAt).ExecX(t.Context())

		res := api.Post(
			fmt.Sprintf("/admin/users/%d/unban", admin.ID),
			tests.GetCookie(adminSession.ID),
		)
		if res.Code != http.StatusForbidden {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusForbidden,
				res.Code,
			)
		}
		saved := db.Client.User.GetX(t.Context(), admin.ID)
		if saved.Banned == nil || !*saved.Banned || saved.BannedAt == nil {
			t.Fatal("expected the account to remain banned")
		}
	})
}
