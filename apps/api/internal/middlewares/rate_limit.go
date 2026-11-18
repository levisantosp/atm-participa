package middlewares

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/pkg/infra/errors"
	"github.com/levisantosp/atm-participa/pkg/infra/redis"
)

func RateLimit(
	api huma.API,
	limit int64,
) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		count, ttl, err := incr(
			ctx.Context(),
			fmt.Sprintf(
				"ratelimit:%s:%s:%s",
				getClientIP(ctx),
				ctx.Method(),
				ctx.Operation().Path,
			),
		)
		if err != nil {
			huma.WriteErr(
				api,
				ctx,
				http.StatusInternalServerError,
				errors.InternalServerError,
				err,
			)
			return
		}

		ctx.SetHeader("x-ratelimit-limit", strconv.FormatInt(limit, 10))
		ctx.SetHeader(
			"x-rate-limit-remaining",
			strconv.FormatInt(max(0, limit-*count), 10),
		)
		ctx.SetHeader(
			"x-ratelimit-reset",
			strconv.FormatInt(max(0, int64(*ttl)), 10),
		)

		if *count > int64(limit) {
			ctx.SetHeader(
				"retry-after",
				strconv.FormatInt(max(0, int64(*ttl)), 10),
			)

			huma.WriteErr(
				api,
				ctx,
				http.StatusTooManyRequests,
				"Muitas requisições. Tente novamente em alguns segundos.",
			)

			return
		}

		next(ctx)
	}
}

func incr(ctx context.Context, key string) (*int64, *time.Duration, error) {
	count, err := redis.Client.Incr(ctx, key).Result()
	if err != nil {
		return nil, nil, err
	}

	if count == 1 {
		if err := redis.Client.Expire(ctx, key, time.Minute).Err(); err != nil {
			return nil, nil, err
		}
	}

	ttl, err := redis.Client.TTL(ctx, key).Result()
	if err != nil {
		return nil, nil, err
	}

	return &count, &ttl, nil
}

func getClientIP(ctx huma.Context) string {
	if forwarded := ctx.Header("x-forwarded-for"); forwarded != "" {
		if ip := strings.TrimSpace(strings.Split(forwarded, ",")[0]); ip != "" {
			return ip
		}
	}

	if realIP := strings.TrimSpace(ctx.Header("x-real-ip")); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(ctx.RemoteAddr())
	if err == nil {
		return host
	}

	return ctx.RemoteAddr()
}
