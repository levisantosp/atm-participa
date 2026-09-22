package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/tests"
)

func TestSignIn(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	t.Run("should reject invalid email", func(t *testing.T) {
		res := api.Post("/auth/sign-in/email", map[string]any{
			"email":    strings.Repeat("a", 10),
			"password": strings.Repeat("a", 10),
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	const email = "testuser@email.com"

	t.Run("should reject password less than 8", func(t *testing.T) {
		res := api.Post("/auth/sign-in/email", map[string]any{
			"email":    email,
			"password": strings.Repeat("a", 7),
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject password greater than 72", func(t *testing.T) {
		res := api.Post("/auth/sign-in/email", map[string]any{
			"email":    email,
			"password": strings.Repeat("a", 73),
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should return 404", func(t *testing.T) {
		res := api.Post("/auth/sign-in/email", map[string]any{
			"email":    email,
			"password": strings.Repeat("a", 10),
		})

		if res.Code != http.StatusNotFound {
			t.Fatalf("expected status %d got %d", http.StatusNotFound, res.Code)
		}
	})

	res := api.Post("/auth/sign-up/email", map[string]any{
		"email":       email,
		"password":    strings.Repeat("a", 10),
		"displayName": "Test User",
		"username":    "testuser",
	})

	if res.Code != http.StatusNoContent {
		t.Fatalf("expected status %d got %d", http.StatusNoContent, res.Code)
	}

	t.Run("should sign in", func(t *testing.T) {
		res := api.Post("/auth/sign-in/email", map[string]any{
			"email":    email,
			"password": strings.Repeat("a", 10),
		})

		if res.Code != http.StatusNoContent {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusNoContent,
				res.Code,
			)
		}
	})
}
