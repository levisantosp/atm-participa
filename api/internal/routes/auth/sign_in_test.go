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
