package issues

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/api/tests"

	_ "github.com/levisantosp/atm-participa/api/ent/generated/runtime"
)

func TestDeleteIssue(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)

	Routes(api)

	user := tests.CreateUser(t)
	session := tests.CreateSession(user, t)

	issue, err := tests.CreateIssue(t, user)
	if err != nil {
		t.Fatal(err)
	}

	issueWithImage, err := tests.CreateIssueWithImage(t, user)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		path := fmt.Sprintf("/issues/%d", issue.ID)
		res := api.Delete(path)

		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should delete issue", func(t *testing.T) {
		path := fmt.Sprintf("/issues/%d", issue.ID)
		res := api.Delete(path, tests.GetCookie(session.ID))

		if res.Code != http.StatusNoContent {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusNoContent,
				res.Code,
			)
		}
	})

	t.Run("should delete issue with image", func(t *testing.T) {
		path := fmt.Sprintf("/issues/%d", issueWithImage.ID)
		res := api.Delete(path, tests.GetCookie(session.ID))

		if res.Code != http.StatusNoContent {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusNoContent,
				res.Code,
			)
		}
	})
}
