package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/ent/generated"
	"github.com/levisantosp/atm-participa/api/middlewares"
	"github.com/levisantosp/atm-participa/api/redis"
	"github.com/levisantosp/atm-participa/api/utils"
)

func SignUpWithEmail(
	ctx context.Context,
	input *struct {
		Body struct {
			Email       string `json:"email" format:"email" required:"true"`
			Password    string `json:"password" minLength:"8" maxLength:"72" required:"true"`
			DisplayName string `json:"displayName" minLength:"1" maxLength:"100" required:"true"`
			Username    string `json:"username" minLength:"3" maxLength:"32" pattern:"^[a-zA-Z0-9_]+$" required:"true"`
		}
	},
) (*SignInOutput, error) {
	hash, err := utils.GeneratePasswordHash(input.Body.Password)
	if err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	tx, err := db.Client.Tx(ctx)
	if err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	user, err := tx.User.Create().
		SetEmail(input.Body.Email).
		SetUsername(input.Body.Username).
		SetDisplayName(input.Body.DisplayName).
		Save(ctx)
	if err != nil {
		if generated.IsConstraintError(err) {
			return nil, huma.Error409Conflict(
				"O nome de usuário ou e-mail informado já está cadastrado no sistema.",
			)
		}

		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	_, err = tx.Account.Create().
		SetPassword(*hash).
		SetProvider("email").
		SetUser(user).
		Save(ctx)
	if err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	sessionHash := make([]byte, 32)
	_, err = rand.Read(sessionHash)
	if err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	sessionId := hex.EncodeToString(sessionHash)

	session, err := json.Marshal(middlewares.Session{
		ID:          sessionId,
		UserId:      strconv.FormatInt(user.ID, 10),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		IsAdmin:     user.IsAdmin,
	})
	if err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	ttl := 7 * 24 * time.Hour
	err = redis.Client.Set(ctx, "session:"+sessionId, session, ttl).Err()
	if err != nil {
		return nil, utils.LogErr(
			huma.Error500InternalServerError("Internal Server Error"),
			err,
		)
	}

	sessionCookie := http.Cookie{
		Name:     "session",
		Value:    sessionId,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
		Path:     "/",
	}

	return &SignInOutput{
		SetCookie: []http.Cookie{sessionCookie},
	}, nil
}
