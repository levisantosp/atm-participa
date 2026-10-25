package users

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/internal/middlewares"
	"github.com/levisantosp/atm-participa/api/internal/tests"

	_ "github.com/levisantosp/atm-participa/api/ent/generated/runtime"
)

func TestEditMe(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		res := api.Put("/users/me", map[string]any{
			"displayName": "Updated Name",
		})

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should reject request without fields to update", func(t *testing.T) {
		res := api.Put(
			"/users/me",
			tests.GetCookie(session.ID),
		)

		if res.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusBadRequest,
				res.Code,
			)
		}
	})

	t.Run("should update only the provided fields", func(t *testing.T) {
		res := api.Put("/users/me", tests.GetCookie(session.ID), map[string]any{
			"displayName": "Updated Name",
		})

		if res.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusOK,
				res.Code,
			)
		}

		var updatedSession middlewares.Session
		if err := json.NewDecoder(res.Body).
			Decode(&updatedSession); err != nil {
			t.Fatal(err)
		}

		if updatedSession.DisplayName != "Updated Name" {
			t.Errorf(
				"expected display name %q got %q",
				"Updated Name",
				updatedSession.DisplayName,
			)
		}

		if updatedSession.Username != user.Username {
			t.Errorf(
				"expected username %q to remain unchanged, got %q",
				user.Username,
				updatedSession.Username,
			)
		}

		if updatedSession.UserId != session.UserId {
			t.Errorf(
				"expected user id %q got %q",
				session.UserId,
				updatedSession.UserId,
			)
		}

		if updatedSession.ID != session.ID {
			t.Errorf(
				"expected session id %q got %q",
				session.ID,
				updatedSession.ID,
			)
		}

		if got := res.Header().Get("Set-Cookie"); got == "" {
			t.Error("expected Set-Cookie response header")
		}
	})
}
