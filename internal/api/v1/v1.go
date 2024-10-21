package v1

import (
	"irule-api/internal/api/v1/routes"
	"irule-api/internal/config"
	middlewares "irule-api/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(dbPool *pgxpool.Pool, cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	r.Mount("/auth", routes.Auth(dbPool, cfg))

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(cfg))
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{"message": "protected"})
		})
	})

	return r
}
