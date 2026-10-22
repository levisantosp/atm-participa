package tests

import (
	"context"
	"testing"

	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/env"
	"github.com/levisantosp/atm-participa/api/internal/r2"
	"github.com/levisantosp/atm-participa/infra/redis"
)

func Setup(t *testing.T) {
	env.LoadEnv("../../../.env.test")
	if err := r2.Connect(t.Context()); err != nil {
		t.Fatal(err)
	}
	redis.Connect(env.Env.RedisAddr, env.Env.RedisPassword)
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
