package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/report"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/apps/api/internal/dtos"
	"github.com/levisantosp/atm-participa/apps/api/internal/tests"
	"github.com/levisantosp/atm-participa/apps/api/internal/utils"
)

func TestGetReports(t *testing.T) {
	tests.Setup(t)
	_, api := humatest.New(t)
	Routes(api)

	user := tests.CreateUser(t)
	admin := tests.CreateAdminUser(t)
	adminSession := tests.CreateAdminSession(admin, t)
	issue, err := tests.CreateIssue(t, user)
	if err != nil {
		t.Fatal(err)
	}
	first := tests.CreateReport(t, admin, report.EntityIssue, issue.ID, "Spam")
	middle := tests.CreateReport(
		t,
		admin,
		report.EntityUser,
		user.ID,
		"Assédio",
	)
	last := tests.CreateReport(
		t,
		user,
		report.EntityIssue,
		issue.ID,
		"Conteúdo ofensivo",
	)
	if err := db.Client.Issue.DeleteOne(issue).Exec(t.Context()); err != nil {
		t.Fatal(err)
	}

	t.Run("should list both entities newest first", func(t *testing.T) {
		res := api.Get("/admin/reports", tests.GetCookie(adminSession.ID))
		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
		}
		var body utils.CursorPaginatedResponse[dtos.Report]
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Items) != 3 || body.NextCursor != nil {
			t.Fatalf("unexpected page: %+v", body)
		}
		for i, id := range []int64{last.ID, middle.ID, first.ID} {
			if body.Items[i].ID != fmt.Sprint(id) {
				t.Fatalf("expected report %d got %s", id, body.Items[i].ID)
			}
		}
		if body.Items[1].Reason != middle.Reason {
			t.Fatalf(
				"expected reason %q got %q",
				middle.Reason,
				body.Items[1].Reason,
			)
		}
		author := body.Items[1].Author
		if author == nil || author.ID != fmt.Sprint(admin.ID) {
			t.Fatalf("expected report author %d, got %v", admin.ID, author)
		}
	})

	t.Run(
		"should paginate filtered reports after target deletion",
		func(t *testing.T) {
			path := "/admin/reports?entity=issue&limit=1"
			for i, id := range []int64{last.ID, first.ID} {
				res := api.Get(path, tests.GetCookie(adminSession.ID))
				if res.Code != http.StatusOK {
					t.Fatalf(
						"expected status %d got %d",
						http.StatusOK,
						res.Code,
					)
				}
				var body utils.CursorPaginatedResponse[dtos.Report]
				if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if len(body.Items) != 1 {
					t.Fatalf("expected one report got %d", len(body.Items))
				}
				item := body.Items[0]
				if item.ID != fmt.Sprint(id) ||
					item.Entity != report.EntityIssue ||
					item.EntityID != fmt.Sprint(issue.ID) {
					t.Fatalf("unexpected report: %+v", item)
				}
				if i == 0 {
					if body.NextCursor == nil || *body.NextCursor != id {
						t.Fatalf("expected cursor %d", id)
					}
					path += fmt.Sprintf("&cursor=%d", *body.NextCursor)
				} else if body.NextCursor != nil {
					t.Fatal("expected no next page")
				}
			}
		},
	)
}
