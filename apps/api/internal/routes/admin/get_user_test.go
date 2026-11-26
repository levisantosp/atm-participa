package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"
)

func TestGetUser(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	admin := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(admin, t)
	path := fmt.Sprintf("/admin/users/%d", user.ID)

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

	t.Run("should reject invalid user id", func(t *testing.T) {
		res := api.Get("/admin/users/0", tests.GetCookie(adminSession.ID))
		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should return not found for a missing user", func(t *testing.T) {
		res := api.Get(
			"/admin/users/999999999",
			tests.GetCookie(adminSession.ID),
		)
		if res.Code != http.StatusNotFound {
			t.Fatalf("expected status %d got %d", http.StatusNotFound, res.Code)
		}
	})

	t.Run(
		"should return the requested user and ban details",
		func(t *testing.T) {
			bannedAt := time.Now().UTC().Truncate(time.Second)
			endsAt := bannedAt.Add(time.Hour)
			if err := user.Update().
				SetBanned(true).
				SetBannedAt(bannedAt).
				SetBanEndsAt(endsAt).
				Exec(t.Context()); err != nil {
				t.Fatal(err)
			}

			res := api.Get(path, tests.GetCookie(adminSession.ID))
			if res.Code != http.StatusOK {
				t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
			}

			var body dtos.User
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body.ID != fmt.Sprint(user.ID) ||
				body.Username != user.Username ||
				body.DisplayName != user.DisplayName ||
				body.IsAdmin != user.IsAdmin {
				t.Fatalf("unexpected user: %+v", body)
			}
			if body.Banned == nil || !*body.Banned ||
				body.BannedAt == nil || !body.BannedAt.Equal(bannedAt) ||
				body.BanEndsAt == nil || !body.BanEndsAt.Equal(endsAt) {
				t.Fatalf("unexpected ban details: %+v", body)
			}
		},
	)
}
