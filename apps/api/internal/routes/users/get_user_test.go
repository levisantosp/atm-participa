package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"
)

func TestGetUser(t *testing.T) {
	tests.Setup(t)
	_, api := humatest.New(t)
	huma.Get(api, "/users/{username}", GetUser)

	user := tests.CreateUser(t)
	banned := tests.CreateAdminUser(t)
	banned.Update().SetBanned(true).SaveX(t.Context())

	t.Run("should return the user by username", func(t *testing.T) {
		res := api.Get("/users/" + user.Username)
		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}

		var body dtos.User
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.ID != fmt.Sprint(user.ID) || body.Username != user.Username {
			t.Fatalf("unexpected user: %+v", body)
		}
	})

	for _, tc := range []struct {
		name     string
		username string
	}{
		{
			"unknown user",
			"missing-user",
		},
		{
			"banned user",
			banned.Username,
		},
	} {
		t.Run("should return 404 for "+tc.name, func(t *testing.T) {
			res := api.Get("/users/" + tc.username)
			if res.Code != http.StatusNotFound {
				t.Fatalf(
					"expected status %d got %d",
					http.StatusNotFound,
					res.Code,
				)
			}
		})
	}
}
