package main

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/levisantosp/atm-participa/cdn/internal/env"
	"github.com/levisantosp/atm-participa/cdn/internal/r2"
	"github.com/levisantosp/atm-participa/cdn/internal/routes"
	"github.com/levisantosp/atm-participa/infra/errors"
)

func main() {
	env.LoadEnv(".env")

	if err := r2.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RedirectSlashes)

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

	routes.CreateRoutes(api)

	log.Println("HTTP server running at http://localhost:3336")
	if err := http.ListenAndServe(":3336", r); err != nil {
		log.Fatal(err)
	}
}
