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
	entityID int64,
	reason string,
) *generated.Report {
	t.Helper()

	create := db.Client.Report.Create().
		SetUserID(user.ID).
		SetEntity(entity).
		SetEntityID(entityID).
		SetReason(reason)
	if entity == report.EntityIssue {
		create.SetTargetIssueID(entityID)
	}
	if entity == report.EntityUser {
		create.SetTargetUserID(entityID)
	}

	return create.SaveX(t.Context())
}
