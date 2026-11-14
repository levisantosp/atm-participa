package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	acc "github.com/levisantosp/atm-participa/api/ent/generated/account"
	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/tests"
	"github.com/levisantosp/atm-participa/api/internal/utils"
	"github.com/levisantosp/atm-participa/infra/redis"
)

func TestResetPassword(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	oldPasswordHash, err := utils.GeneratePasswordHash("old-password")
	if err != nil {
		t.Fatal(err)
	}
	db.Client.Account.Create().
		SetProvider(acc.ProviderEmail).
		SetPassword(*oldPasswordHash).
		SetUser(user).
		SaveX(t.Context())

	t.Run("should reject a password shorter than 8 characters", func(t *testing.T) {
		res := api.Post("/auth/reset-password", map[string]any{
			"token":    "test-token",
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

	t.Run("should reject a missing or expired token", func(t *testing.T) {
		res := api.Post("/auth/reset-password", map[string]any{
			"token":    "missing-token",
			"password": "new-password",
		})

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should update the password and consume the token", func(t *testing.T) {
		token := "valid-reset-token"
		const newPassword = "new-password"
		_, tokenHash, err := GenerateToken(&token)
		if err != nil {
			t.Fatal(err)
		}
		key := "reset-password:" + *tokenHash
		if err := redis.Client.Set(
			t.Context(),
			key,
			user.ID,
			10*time.Minute,
		).Err(); err != nil {
			t.Fatal(err)
		}

		res := api.Post("/auth/reset-password", map[string]any{
			"token":    token,
			"password": newPassword,
		})

		if res.Code != http.StatusNoContent {
			t.Fatalf("expected status %d got %d", http.StatusNoContent, res.Code)
		}

		updatedAccount, err := db.Client.Account.Query().
			Where(acc.UserIDEQ(user.ID), acc.ProviderEQ(acc.ProviderEmail)).
			Only(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		if !utils.VerifyPassword(newPassword, updatedAccount.Password) {
			t.Fatal("expected account password to be updated")
		}
		if utils.VerifyPassword("old-password", updatedAccount.Password) {
			t.Fatal("expected previous password to stop working")
		}

		if err := redis.Client.Get(t.Context(), key).Err(); err != redis.Nil {
			t.Fatalf("expected reset token to be consumed, got: %v", err)
		}
	})
}
