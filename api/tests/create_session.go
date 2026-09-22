package tests

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/middlewares"
	"github.com/levisantosp/atm-participa/api/redis"
)

func CreateSession(user *generated.User, t *testing.T) middlewares.Session {
	session := middlewares.Session{
		ID:          "test-session",
		UserId:      strconv.FormatInt(user.ID, 10),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		IsAdmin:     user.IsAdmin,
	}

	sessionJson, err := json.Marshal(session)
	if err != nil {
		t.Fatal(err)
	}

	if err := redis.Client.Set(
		t.Context(),
		"session:"+session.ID,
		sessionJson,
		time.Minute,
	).Err(); err != nil {
		t.Fatal(err)
	}

	return session
}

func CreateAdminSession(
	user *generated.User,
	t *testing.T,
) middlewares.Session {
	session := middlewares.Session{
		ID:          "admin-test-session",
		UserId:      strconv.FormatInt(user.ID, 10),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		IsAdmin:     user.IsAdmin,
	}

	sessionJson, err := json.Marshal(session)
	if err != nil {
		t.Fatal(err)
	}

	if err := redis.Client.Set(
		t.Context(),
		"session:"+session.ID,
		sessionJson,
		time.Minute,
	).Err(); err != nil {
		t.Fatal(err)
	}

	return session
}
