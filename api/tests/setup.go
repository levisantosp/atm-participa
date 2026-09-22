package tests

import (
	"context"
	"testing"

	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/redis"
	"github.com/levisantosp/atm-participa/api/utils"
)

func Setup(t *testing.T) {
	utils.LoadEnv("../../../.env.test")
	redis.Connect()
	sqlDB := db.Connect()

	if _, err := sqlDB.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`); err != nil {
		t.Fatal(err)
	}

	if err := db.Client.Schema.Create(context.Background()); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := db.Client.Close(); err != nil {
			t.Error(err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Error(err)
		}
	})
}
