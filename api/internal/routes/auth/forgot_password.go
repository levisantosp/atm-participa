package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/levisantosp/atm-participa/api/ent/generated"
	_user "github.com/levisantosp/atm-participa/api/ent/generated/user"
	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/env"
	"github.com/levisantosp/atm-participa/infra/redis"
	"github.com/resend/resend-go/v4"
)

func ForgotPassword(
	ctx context.Context,
	input *struct {
		Body struct {
			Email string `json:"email" format:"email" maxLength:"254"`
		}
	},
) (*struct{}, error) {
	user, err := db.Client.User.Query().
		Where(_user.EmailEQ(input.Body.Email)).
		First(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(token))
	hash := hex.EncodeToString(h[:])

	if err = redis.Client.Set(
		ctx,
		"reset-password:"+hash,
		user.ID,
		10*time.Minute,
	).Err(); err != nil {
		return nil, err
	}

	client := resend.NewClient(env.Data.ResendAPIKey)

	_, err = client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		To:      []string{input.Body.Email},
		From:    fmt.Sprintf("Altamira Participa <%s>", env.Data.ResendEmail),
		Text:    "Recuperação de Senha",
		Subject: "Código para redefinição de senha",
		Html: fmt.Sprintf(`
			<p>Olá,</p>
			<p>Recebemos uma solicitação para redefinir a senha da sua conta. Entre no link abaixo para prosseguir com a redefinição:</p>
			<h1><strong>%s</strong></h1>
			<p>Este link expira em <strong>10 minutos</strong>.</p>
			<p>Se você não solicitou a redefinição de senha, ignore este e-mail. Nenhuma alteração foi realizada na sua conta.</p>
		`, env.Data.WebURL+"/redefinir-senha?token="+token),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
