package auth

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/ent/generated/user"
	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/tests"
	"github.com/levisantosp/atm-participa/infra/redis"
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

		var sessionID string
		for _, cookie := range res.Result().Cookies() {
			if cookie.Name == "session" {
				sessionID = cookie.Value
			}
		}
		if sessionID == "" {
			t.Fatal("expected a session cookie in the response")
		}

		u, err := db.Client.User.Query().
			Where(user.EmailEQ(email)).
			Only(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		if err := redis.Client.Get(
			t.Context(),
			"session:"+sessionID,
		).Err(); err != nil {
			t.Fatalf("expected session to be stored: %v", err)
		}

		if err := redis.Client.ZScore(
			t.Context(),
			"user_sessions:"+strconv.FormatInt(u.ID, 10),
			sessionID,
		).Err(); err != nil {
			t.Fatalf("expected session to be tracked: %v", err)
		}
	})
}
