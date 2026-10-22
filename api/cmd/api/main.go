package main

import (
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/levisantosp/atm-participa/api/internal/db"
	"github.com/levisantosp/atm-participa/api/internal/routes/admin"
	"github.com/levisantosp/atm-participa/api/internal/routes/auth"
	"github.com/levisantosp/atm-participa/api/internal/routes/issues"
	"github.com/levisantosp/atm-participa/api/internal/routes/users"
	"github.com/levisantosp/atm-participa/api/internal/utils"
	"github.com/levisantosp/atm-participa/infra/errors"
	"github.com/levisantosp/atm-participa/infra/redis"

	_ "github.com/levisantosp/atm-participa/api/ent/generated/runtime"
)

func main() {
	utils.LoadEnv(".env")
	db.Connect()
	redis.Connect(utils.Env.RedisAddr, utils.Env.RedisPassword)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RedirectSlashes)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: utils.Env.TrustedOrigins,
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"PATCH",
			"OPTIONS",
		},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world!"))
	})

	huma.NewErrorWithContext = func(_ huma.Context, status int, msg string, errs ...error) huma.StatusError {
		if status >= http.StatusInternalServerError {
			for _, err := range errs {
				if err != nil {
					log.Println(err)
				}
			}

			return huma.NewError(status, errors.InternalServerError)
		}

		return huma.NewError(status, msg, errs...)
	}

	api := humachi.New(r, huma.DefaultConfig("api docs", "0.0.0"))

	auth.Routes(api)
	users.Routes(api)
	issues.Routes(api)
	admin.Routes(api)

	log.Println("HTTP server running at http://localhost:3333")
	if err := http.ListenAndServe(":3333", r); err != nil {
		log.Fatal(err)
	}
}
