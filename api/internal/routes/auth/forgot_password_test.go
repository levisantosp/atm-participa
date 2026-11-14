package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/internal/tests"
	"github.com/levisantosp/atm-participa/infra/redis"
)

func TestForgotPassword(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	t.Run("should reject invalid email", func(t *testing.T) {
		res := api.Post("/auth/forgot-password", map[string]any{
			"email": strings.Repeat("a", 10),
		})

		if res.Code != http.StatusUnprocessableEntity {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnprocessableEntity,
				res.Code,
			)
		}
	})

	t.Run("should not reveal whether an email exists", func(t *testing.T) {
		res := api.Post("/auth/forgot-password", map[string]any{
			"email": "unknown@email.com",
		})

		if res.Code != http.StatusNoContent {
			t.Fatalf("expected status %d got %d", http.StatusNoContent, res.Code)
		}
	})

	t.Run("should send an email and store a short-lived reset token", func(t *testing.T) {
		user := tests.CreateUser(t)
		keysBefore, err := redis.Client.Keys(
			t.Context(),
			"reset-password:*",
		).Result()
		if err != nil {
			t.Fatal(err)
		}
		keysBeforeSet := make(map[string]struct{}, len(keysBefore))
		for _, key := range keysBefore {
			keysBeforeSet[key] = struct{}{}
		}

		res := api.Post("/auth/forgot-password", map[string]any{
			"email": user.Email,
		})

		if res.Code != http.StatusNoContent {
			t.Fatalf("expected status %d got %d", http.StatusNoContent, res.Code)
		}

		keysAfter, err := redis.Client.Keys(
			t.Context(),
			"reset-password:*",
		).Result()
		if err != nil {
			t.Fatal(err)
		}

		for _, key := range keysAfter {
			if _, exists := keysBeforeSet[key]; exists {
				continue
			}

			userID, err := redis.Client.Get(t.Context(), key).Int64()
			if err != nil {
				t.Fatalf("expected reset token to contain the user ID: %v", err)
			}
			if userID != user.ID {
				t.Fatalf("expected user ID %d got %d", user.ID, userID)
			}

			ttl, err := redis.Client.TTL(t.Context(), key).Result()
			if err != nil {
				t.Fatal(err)
			}
			if ttl <= 0 || ttl > 10*time.Minute {
				t.Fatalf("expected token TTL up to 10 minutes, got %s", ttl)
			}

			return
		}

		t.Fatal("expected a new password reset token in Redis")
	})
}
