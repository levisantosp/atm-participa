package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/tests"
)

func TestSignUp(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	const (
		email       string = "testuser@email.com"
		displayName string = "Test User"
	)

	t.Run("should reject invalid email", func(t *testing.T) {
		res := api.Post("/auth/sign-up/email", map[string]any{
			"email":       strings.Repeat("a", 10),
			"password":    strings.Repeat("a", 10),
			"displayName": displayName,
			"username":    "testuser",
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject password less than 8", func(t *testing.T) {
		res := api.Post("/auth/sign-up/email", map[string]any{
			"email":       email,
			"password":    strings.Repeat("a", 7),
			"displayName": displayName,
			"username":    "testuser",
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject display name less than 1", func(t *testing.T) {
		res := api.Post("/auth/sign-up/email", map[string]any{
			"email":       email,
			"password":    strings.Repeat("a", 10),
			"displayName": "",
			"username":    "testuser",
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject display name greater than 100", func(t *testing.T) {
		res := api.Post("/auth/sign-up/email", map[string]any{
			"email":       email,
			"password":    strings.Repeat("a", 10),
			"displayName": strings.Repeat("a", 101),
			"username":    "testuser",
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject username less than 3", func(t *testing.T) {
		res := api.Post("/auth/sign-up/email", map[string]any{
			"email":       email,
			"password":    strings.Repeat("a", 10),
			"displayName": displayName,
			"username":    "aa",
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should reject username greater than 32", func(t *testing.T) {
		res := api.Post("/auth/sign-up/email", map[string]any{
			"email":       email,
			"password":    strings.Repeat("a", 10),
			"displayName": displayName,
			"username":    strings.Repeat("a", 33),
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should sign up", func(t *testing.T) {
		res := api.Post("/auth/sign-up/email", map[string]any{
			"email":       email,
			"password":    strings.Repeat("a", 10),
			"displayName": displayName,
			"username":    strings.Repeat("a", 10),
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
