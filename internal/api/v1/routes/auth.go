package routes

import (
	"encoding/json"
	"irule-api/internal/config"
	"irule-api/internal/constant"
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
		dUser, err := models.FindByEmail(dbPool, body.Email)
		if err != nil && err.Error() != "no rows in result set" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if dUser != nil {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}
		user := models.User{
			Email:    body.Email,
			Password: body.Password,
			Role:     "admin",
		}
		err = user.Create(dbPool)
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

func Me(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(constant.UserKey).(*svc.UserClaims)
		render.JSON(w, r, user)
	}
}

func CreateUser(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(constant.UserKey).(*svc.UserClaims)
		if user.Role != "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var body LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		dbUser := &models.User{
			Role:           "user",
			Email:          body.Email,
			Password:       body.Password,
			OrganizationId: user.OrganizationId,
		}
		err := dbUser.CreateUser(dbPool)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, map[string]string{"status": "ok"})
	}
}
