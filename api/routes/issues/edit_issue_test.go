package issues

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/tests"

	_ "github.com/levisantosp/atm-participa/api/ent/generated/runtime"
)

func TestEditIssue(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)

	issue, err := tests.CreateIssue(t, user)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		path := fmt.Sprintf("/issues/%d", issue.ID)
		res := api.Put(path)

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should edit issue", func(t *testing.T) {
		path := fmt.Sprintf("/issues/%d", issue.ID)
		res := api.Put(path, tests.GetCookie(session.ID), map[string]any{
			"title":       strings.Repeat("a", 3),
			"description": strings.Repeat("a", 10),
		})

		if res.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusOK,
				res.Code,
			)
		}
	})
}
