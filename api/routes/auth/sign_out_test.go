package auth

import (
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/tests"
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
	})
}
