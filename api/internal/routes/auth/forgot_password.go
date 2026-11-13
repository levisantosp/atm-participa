package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
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

	var code string
	b := make([]byte, 4)

	for i := range 4 {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return nil, err
		}

		b[i] = byte('0' + n.Int64())
	}
	code = string(b)

	if err = redis.Client.Set(
		ctx,
		"reset-password:"+strconv.FormatInt(user.ID, 10),
		code,
		10*time.Minute,
	).Err(); err != nil {
		return nil, err
	}

	client := resend.NewClient(env.Env.ResendAPIKey)

	_, err = client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		To:      []string{input.Body.Email},
		From:    "Altamira Participa <atm-participa@levisantos.dev>",
		Text:    "Recuperação de Senha",
		Subject: "Código para redefinição de senha",
		Html: fmt.Sprintf(`
			<p>Olá,</p>
			<p>Recebemos uma solicitação para redefinir a senha da sua conta. Use o código abaixo para prosseguir com a alteração:</p>
			<h1><strong><code>%s</code></strong></h1>
			<p>Este código é de uso único e expira em <strong>10 minutos</strong>.</p>
			<p>Se você não solicitou a redefinição de senha, ignore este e-mail. Nenhuma alteração foi realizada na sua conta.</p>
		`, code),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
