package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Auth(dbPool *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Post("/login", Login(dbPool))
	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r
}

func Login(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		row := dbPool.QueryRow(context.Background(), "SELECT 1")

		var result int
		err := row.Scan(&result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		fmt.Println(result)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
