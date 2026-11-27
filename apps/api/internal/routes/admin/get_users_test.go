package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"
	"github.com/levisantosp/atm-participa/apps/api/internal/utils"
)

func TestGetUsers(t *testing.T) {
	tests.Setup(t)

	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	admin := tests.CreateAdminUser(t)
	session := tests.CreateSession(user, t)
	adminSession := tests.CreateAdminSession(admin, t)
	path := "/admin/users"

	t.Run("should reject unauthenticated request", func(t *testing.T) {
		res := api.Get(path)
		if res.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusUnauthorized,
				res.Code,
			)
		}
	})

	t.Run("should reject non-admin request", func(t *testing.T) {
		res := api.Get(path, tests.GetCookie(session.ID))
		if res.Code != http.StatusForbidden {
			t.Fatalf(
				"expected status %d got %d",
				http.StatusForbidden,
				res.Code,
			)
		}
	})

	for _, tc := range []struct {
		name  string
		query string
		ids   []int64
	}{
		{
			"should list users newest first",
			"",
			[]int64{admin.ID, user.ID},
		},
		{
			"should accept a single-character search",
			"?search=a",
			[]int64{admin.ID, user.ID},
		},
		{
			"should search by user ID",
			fmt.Sprintf("?search=%d", user.ID),
			[]int64{user.ID},
		},
		{
			"should search by username ignoring case",
			"?search=ST-AD",
			[]int64{admin.ID},
		},
		{
			"should search by email ignoring case",
			"?search=STUSER@",
			[]int64{user.ID},
		},
		{
			"should return an empty page when no user matches",
			"?search=absent",
			nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := api.Get(path+tc.query, tests.GetCookie(adminSession.ID))
			if res.Code != http.StatusOK {
				t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
			}
			var body utils.CursorPaginatedResponse[dtos.User]
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body.NextCursor != nil || body.Items == nil ||
				len(body.Items) != len(tc.ids) {
				t.Fatalf("unexpected page: %+v", body)
			}
			for i, id := range tc.ids {
				if body.Items[i].ID != fmt.Sprint(id) {
					t.Fatalf("expected user %d got %s", id, body.Items[i].ID)
				}
			}
		})
	}

	t.Run("should paginate with a cursor", func(t *testing.T) {
		res := api.Get(path+"?limit=1", tests.GetCookie(adminSession.ID))
		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}
		var first utils.CursorPaginatedResponse[dtos.User]
		if err := json.NewDecoder(res.Body).Decode(&first); err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || first.Items[0].ID != fmt.Sprint(admin.ID) ||
			first.NextCursor == nil || *first.NextCursor != admin.ID {
			t.Fatalf("unexpected first page: %+v", first)
		}

		nextPath := fmt.Sprintf("%s?limit=1&cursor=%d", path, *first.NextCursor)
		res = api.Get(nextPath, tests.GetCookie(adminSession.ID))
		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}
		var next utils.CursorPaginatedResponse[dtos.User]
		if err := json.NewDecoder(res.Body).Decode(&next); err != nil {
			t.Fatal(err)
		}
		if len(next.Items) != 1 || next.Items[0].ID != fmt.Sprint(user.ID) ||
			next.NextCursor != nil {
			t.Fatalf("unexpected final page: %+v", next)
		}
	})

	for _, query := range []string{
		"limit=0", "limit=101", "cursor=0",
		"search=" + strings.Repeat("a", 101),
	} {
		t.Run("should reject "+query, func(t *testing.T) {
			res := api.Get(path+"?"+query, tests.GetCookie(adminSession.ID))
			if res.Code != http.StatusUnprocessableEntity {
				t.Fatalf(
					"expected status %d got %d",
					http.StatusUnprocessableEntity,
					res.Code,
				)
			}
		})
	}
}
