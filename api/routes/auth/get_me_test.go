package auth

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/middlewares"
	"github.com/levisantosp/atm-participa/api/tests"

	_ "github.com/levisantosp/atm-participa/api/ent/generated/runtime"
)

func TestGetMe(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)

	t.Run("should return 404", func(t *testing.T) {
		res := api.Get("/auth/me")

		if res.Code != http.StatusNotFound {
			t.Fatalf("expected status %d got %d", http.StatusNotFound, res.Code)
		}
	})

	t.Run("should return the current session", func(t *testing.T) {
		res := api.Get("/auth/me", tests.GetCookie(session.ID))

		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}

		var sessionJson middlewares.Session
		if err := json.NewDecoder(res.Body).Decode(&sessionJson); err != nil {
			t.Fatal(err)
		}
	})
}
