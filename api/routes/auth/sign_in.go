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
	"github.com/levisantosp/atm-participa/api/ent/generated/account"
	"github.com/levisantosp/atm-participa/api/ent/generated/user"
	"github.com/levisantosp/atm-participa/api/middlewares"
	"github.com/levisantosp/atm-participa/api/redis"
	"github.com/levisantosp/atm-participa/api/utils"
)

type SignInOutput struct {
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

func SignInWithEmail(
	ctx context.Context,
	input *struct {
		Body struct {
			Email    string `json:"email" format:"email" required:"true"`
			Password string `json:"password" minLength:"8" maxLength:"72" required:"true"`
		}
	},
) (*SignInOutput, error) {
	account, err := db.Client.Account.Query().
		Where(account.HasUserWith(user.EmailEQ(input.Body.Email))).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(
			"Usuário não encontrado. Verifique o e-mail.",
		)
	}

	ok := utils.VerifyPassword(input.Body.Password, account.Password)
	if !ok {
		return nil, huma.Error401Unauthorized(
			"Credenciais inválidas. Verifique o e-mail e/ou senha.",
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
		UserId:      strconv.FormatInt(account.Edges.User.ID, 10),
		Username:    account.Edges.User.Username,
		DisplayName: account.Edges.User.DisplayName,
		IsAdmin:     account.Edges.User.IsAdmin,
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
