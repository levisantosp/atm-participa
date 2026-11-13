package auth

import (
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/internal/tests"
	"github.com/levisantosp/atm-participa/infra/redis"
)

func TestSignOut(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		res := api.Post("/auth/sign-out")

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should sign out", func(t *testing.T) {
		res := api.Post("/auth/sign-out", tests.GetCookie(session.ID))

		if res.Code != http.StatusNoContent {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusNoContent,
				res.Code,
			)
		}

		if err := redis.Client.Get(
			t.Context(),
			"session:"+session.ID,
		).Err(); err != redis.Nil {
			t.Fatalf("expected session to be removed, got: %v", err)
		}

		if err := redis.Client.ZScore(
			t.Context(),
			"user_sessions:"+session.UserId,
			session.ID,
		).Err(); err != redis.Nil {
			t.Fatalf("expected session to be untracked, got: %v", err)
		}
	})
}
