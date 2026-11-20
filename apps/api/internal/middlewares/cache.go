package middlewares

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	errs "github.com/levisantosp/atm-participa/pkg/infra/errors"
	"github.com/levisantosp/atm-participa/pkg/infra/redis"
)

type Response struct {
	Status      int    `json:"status"`
	ContentType string `json:"content_type"`
	Body        []byte `json:"body"`
}

type humaContext huma.Context

type Recorder struct {
	humaContext
	status      int
	body        bytes.Buffer
	writer      io.Writer
	contentType string
}

func (r *Recorder) BodyWriter() io.Writer {
	return r
}

func (r *Recorder) SetStatus(status int) {
	r.status = status
	r.humaContext.SetStatus(status)
}

func (r *Recorder) SetHeader(name, value string) {
	r.humaContext.SetHeader(name, value)
	if strings.EqualFold(name, "content-type") {
		r.contentType = value
	}
}

func (r *Recorder) Write(body []byte) (int, error) {
	n, err := r.writer.Write(body)
	if n > 0 {
		_, _ = r.body.Write(body[:n])
	}
	return n, err
}

func Cache(
	api huma.API,
	ttl time.Duration,
) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if ctx.Method() != http.MethodGet {
			next(ctx)
			return
		}

		var userID *int64

		if user, ok := ctx.Context().Value(userCtxKey).(*generated.User); ok &&
			ctx.Operation().OperationID == "get-issues" {
			userID = &user.ID
		}

		url := ctx.URL()
		key := hash(url.RequestURI(), userID)
		raw, err := redis.Client.Get(ctx.Context(), key).Bytes()
		if err == nil {
			var res Response

			if json.Unmarshal(raw, &res) == nil {
				println("bateu no cache")
				ctx.SetHeader("content-type", res.ContentType)
				ctx.SetHeader("x-cache", "hit")
				ctx.SetStatus(res.Status)

				if _, err := ctx.BodyWriter().Write(res.Body); err != nil {
					huma.WriteErr(
						api,
						ctx,
						http.StatusInternalServerError,
						errs.InternalServerError,
						err,
					)
				}

				return
			}
		} else if !errors.Is(err, redis.Nil) {
			huma.WriteErr(
				api,
				ctx,
				http.StatusInternalServerError,
				errs.InternalServerError,
				err,
			)
			return
		}

		println("sem cache")

		ctx.SetHeader("x-cache", "miss")
		rec := &Recorder{
			humaContext: ctx,
			status:      http.StatusOK,
			writer:      ctx.BodyWriter(),
		}

		next(rec)

		if rec.status != http.StatusOK {
			return
		}

		data, err := json.Marshal(Response{
			Status:      rec.status,
			ContentType: rec.contentType,
			Body:        rec.body.Bytes(),
		})
		if err != nil {
			log.Println(err)
			return
		}

		err = redis.Client.Set(ctx.Context(), key, data, ttl).Err()
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func hash(url string, userID *int64) string {
	sum := sha256.Sum256([]byte(url))
	hash := hex.EncodeToString(sum[:])

	if userID == nil {
		return hash
	}

	return fmt.Sprintf("user:%d:%s", *userID, hash)
}
