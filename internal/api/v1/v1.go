package v1

import (
	"irule-api/internal/api/v1/routes"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(dbPool *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Mount("/auth", routes.Auth(dbPool))

	return r
}
