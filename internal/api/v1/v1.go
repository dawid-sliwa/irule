package v1

import (
	"irule-api/data"
	"irule-api/internal/api/v1/routes"
	"irule-api/internal/config"
	"irule-api/internal/constant"
	middlewares "irule-api/internal/middleware"
	"irule-api/internal/svc"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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
		r.Get("/user-stats", UserStats(dbPool))
		r.Mount("/documentation", routes.Documentations(dbPool, cfg))
		r.Mount("/tag", routes.Tags(dbPool))
		r.Get("/me", routes.Me(dbPool))

		r.Post("/create-user", routes.CreateUser(dbPool))
	})

	return r
}

type UserStatsResponse struct {
	ID               uuid.UUID `json:"id"`
	TagsCreatedCount int       `json:"tags_created_count"`
	Email            string    `json:"email"`
}

func UserStats(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := dbPool.Query(r.Context(), data.GetUserSttats, r.Context().Value(constant.UserKey).(*svc.UserClaims).OrganizationId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []UserStatsResponse

		for rows.Next() {
			var user UserStatsResponse
			if err := rows.Scan(&user.ID, &user.TagsCreatedCount, &user.Email); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}
		render.JSON(w, r, users)
	}
}
