package tests

import (
	"testing"

	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/report"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
)

func CreateReport(
	t *testing.T,
	user *generated.User,
	entity report.Entity,
	targetID int64,
	reason string,
) *generated.Report {
	t.Helper()

	create := db.Client.Report.Create().
		SetUserID(user.ID).
		SetEntity(entity).
		SetReason(reason)
	if entity == report.EntityIssue {
		create.SetTargetIssueID(targetID)
	}
	if entity == report.EntityUser {
		create.SetTargetUserID(targetID)
	}

	return create.SaveX(t.Context())
}
