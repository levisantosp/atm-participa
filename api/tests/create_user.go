package tests

import (
	"testing"

	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/ent/generated"
)

func CreateUser(t *testing.T) *generated.User {
	return db.Client.User.
		Create().
		SetUsername("test-user").
		SetDisplayName("Test User").
		SetEmail("testuser@email.com").
		SaveX(t.Context())
}
