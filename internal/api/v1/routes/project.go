package routes

// import (
// 	"context"
// 	"encoding/json"
// 	"irule-api/internal/config"
// 	"irule-api/internal/constant"
// 	"irule-api/internal/db/models"
// 	middlewares "irule-api/internal/middleware"
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func Project(dbpool *pgxpool.Pool, cfg *config.Config) chi.Router {
// 	r := chi.NewRouter()

// 	r.With(middlewares.AdminOnly).Post("/projects", CreateProject(dbpool))
// 	r.With(middlewares.AdminOnly).Delete("/projects/{id}", DeleteProject(dbpool))
// 	r.Get("/projects/{id}", GetProjectByID(dbpool))
// 	r.Get("/organizations/{orgID}/projects", GetProjectsByOrganization(dbpool))
// 	return r
// }

// type CreateProjectRequest struct {
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// }

// func CreateProject(dbpool *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req CreateProjectRequest
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		_, err := dbpool.Exec(context.Background(), "INSERT INTO projects (name, description, organization_id) VALUES ($1, $2, $3)", req.Name, req.Description, r.Context().Value(constant.UserKey).(*models.User).OrganizationId)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusCreated)
// 	}
// }

// func DeleteProject(dbpool *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		projectID := chi.URLParam(r, "id")

// 		_, err := dbpool.Exec(context.Background(), "DELETE FROM projects WHERE id = $1", projectID)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusNoContent)
// 	}
// }

// func GetProjectByID(dbpool *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		projectID := chi.URLParam(r, "id")
// 		var project CreateProjectRequest
// 		err := dbpool.QueryRow(context.Background(), "SELECT name, description, organization_id FROM projects WHERE id = $1", projectID).Scan(&project.Name, &project.Description, r.Context().Value(constant.UserKey).(*models.User).OrganizationId)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		json.NewEncoder(w).Encode(project)
// 	}
// }

// func GetProjectsByOrganization(dbpool *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		orgID := chi.URLParam(r, "orgID")
// 		rows, err := dbpool.Query(context.Background(), "SELECT name, description, organization_id FROM projects WHERE organization_id = $1", orgID)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		defer rows.Close()

// 		var projects []*models.Project
// 		for rows.Next() {
// 			var project *models.Project
// 			if err := rows.Scan(project.Name, project.Description, project.OrganizationID); err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			projects = append(projects, project)
// 		}
// 		json.NewEncoder(w).Encode(projects)
// 	}
// }
