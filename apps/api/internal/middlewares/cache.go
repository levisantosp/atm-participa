package middlewares

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	errs "github.com/levisantosp/atm-participa/pkg/infra/errors"
	"github.com/levisantosp/atm-participa/pkg/infra/redis"
)

type Response struct {
	Status      int    `json:"status"`
	ContentType string `json:"content_type"`
	Body        []byte `json:"body"`
}

type Recorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (r *Recorder) WriteHeader(status int) {
	if r.status != 0 {
		return
	}

	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *Recorder) Write(body []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}

	n, err := r.ResponseWriter.Write(body)
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
				_, writer := humachi.Unwrap(ctx)

				if res.ContentType != "" {
					writer.Header().Set("content-type", res.ContentType)
				}

				writer.Header().Set("x-cache", "hit")
				writer.WriteHeader(res.Status)

				if _, err := writer.Write(res.Body); err != nil {
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

		req, writer := humachi.Unwrap(ctx)
		rec := &Recorder{
			ResponseWriter: writer,
			status:         http.StatusOK,
		}
		rec.Header().Set("x-cache", "miss")
		wrappedCtx := humachi.NewContext(ctx.Operation(), req, rec)
		wrappedCtx = huma.WithContext(wrappedCtx, ctx.Context())

		next(wrappedCtx)

		if rec.status != http.StatusOK {
			return
		}

		data, err := json.Marshal(Response{
			Status:      rec.status,
			ContentType: rec.Header().Get("content-type"),
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
