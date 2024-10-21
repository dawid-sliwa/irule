package routes

import (
	"encoding/json"
	"irule-api/internal/config"
	"irule-api/internal/db/models"
	"irule-api/internal/svc"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Auth(dbPool *pgxpool.Pool, cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	r.Post("/login", Login(dbPool, cfg))
	r.Post("/register", Register(dbPool))

	return r
}

func Login(dbPool *pgxpool.Pool, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}
		user, err := models.FindByEmail(dbPool, body.Email)
		if err != nil {
			http.Error(w, "Email or password is invalid", http.StatusUnauthorized)
			return
		}
		valid := user.ComparePassword(body.Password)
		if !valid {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		tokenString, err := svc.CreateToken(user, cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, map[string]string{"token": tokenString})
	}
}

func Register(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := models.User{
			Email:    body.Email,
			Password: body.Password,
			Role:     "admin",
		}
		err := user.Create(dbPool)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, map[string]string{"status": "ok"})
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
