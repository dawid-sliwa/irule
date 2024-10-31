package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"irule-api/data"
	"irule-api/internal/config"
	"irule-api/internal/constant"
	"irule-api/internal/db/models"
	"irule-api/internal/svc"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Documentations(dbpool *pgxpool.Pool, cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	r.Get("/", GetDocumentations(dbpool))
	r.Get("/{id}", GetDocumentation(dbpool))
	r.Post("/", CreateDocumentation(dbpool))
	r.Put("/{id}", UpdateDocumentation(dbpool))
	r.Delete("/{id}", DeleteDocumentation(dbpool))

	return r
}

type CreateDocumentationRequest struct {
	Title   string `json:"title" validate:"required,min=5,max=100"`
	Content string `json:"content" validate:"required,min=20"`
}

type DocumentationResponse struct {
	models.Documentation
	TagCount int           `json:"tag_count"`
	Tags     []*models.Tag `json:"tags"`
}

func GetDocumentations(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := dbpool.Query(context.Background(), data.QueryDocumentationsByOrg, r.Context().Value(constant.UserKey).(*svc.UserClaims).OrganizationId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var documentations []*DocumentationResponse
		for rows.Next() {
			var documentation DocumentationResponse
			if err := rows.Scan(&documentation.ID, &documentation.Name, &documentation.Content, &documentation.TagCount); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tagRows, err := dbpool.Query(context.Background(), data.QueryTagsByDoc, documentation.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer tagRows.Close()

			for tagRows.Next() {
				var tag models.Tag
				if err := tagRows.Scan(&tag.ID, &tag.Name); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				documentation.Tags = append(documentation.Tags, &tag)
			}

			documentations = append(documentations, &documentation)
		}
		render.JSON(w, r, documentations)
	}
}

func GetDocumentation(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		documentationID := chi.URLParam(r, "id")
		var documentation DocumentationResponse
		err := dbpool.QueryRow(context.Background(), data.QueryDocumentationByID, documentationID, r.Context().Value(constant.UserKey).(*svc.UserClaims).OrganizationId).Scan(&documentation.ID, &documentation.Name, &documentation.Content, &documentation.TagCount)
		if err != nil && err.Error() == "no rows in result set" {
			http.Error(w, err.Error(), http.StatusNotFound)
			render.JSON(w, r, map[string]string{"status": "not found"})
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tagRows, err := dbpool.Query(context.Background(), data.QueryTagsByDoc, documentation.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tagRows.Close()

		for tagRows.Next() {
			var tag models.Tag
			if err := tagRows.Scan(&tag.ID, &tag.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			documentation.Tags = append(documentation.Tags, &tag)
		}

		render.JSON(w, r, documentation)
	}
}

func CreateDocumentation(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateDocumentationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
			return
		}

		var doc models.Documentation
		err := dbpool.QueryRow(
			context.Background(),
			data.InsertDocumentation,
			req.Title,
			req.Content,
			r.Context().Value(constant.UserKey).(*svc.UserClaims).OrganizationId,
		).Scan(&doc.ID, &doc.Name, &doc.Content)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, doc)
	}
}

func UpdateDocumentation(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		documentationID := chi.URLParam(r, "id")
		var req CreateDocumentationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
			return
		}
		var doc models.Documentation

		err := dbpool.QueryRow(
			context.Background(),
			data.UpdateDocumentation,
			req.Title,
			req.Content,
			documentationID,
		).Scan(&doc.ID, &doc.Name, &doc.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, doc)
	}
}

func DeleteDocumentation(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		documentationID := chi.URLParam(r, "id")

		_, err := dbpool.Exec(context.Background(), "DELETE FROM documentations WHERE id = $1", documentationID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, map[string]string{"status": "ok"})
	}
}
