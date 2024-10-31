package routes

import (
	"context"
	"encoding/json"
	"irule-api/data"
	"irule-api/internal/constant"
	"irule-api/internal/db/models"
	"irule-api/internal/svc"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Tags(dbpool *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Get("/", GetTags(dbpool))
	r.Get("/{id}", GetTag(dbpool))
	r.Post("/", CreateTag(dbpool))
	r.Put("/{id}", UpdateTag(dbpool))
	r.Delete("/{id}", DeleteTag(dbpool))

	return r
}

type CreateTagRequest struct {
	Name            string    `json:"name"`
	DocumentationId uuid.UUID `json:"documentation_id"`
}

func GetTags(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID := chi.URLParam(r, "docID")
		rows, err := dbpool.Query(context.Background(), data.QueryTagsByDoc, docID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tags []*models.Tag
		for rows.Next() {
			var tag models.Tag
			if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tags = append(tags, &tag)
		}
		render.JSON(w, r, tags)
	}
}

func GetTag(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID := chi.URLParam(r, "id")
		docID := chi.URLParam(r, "docID")
		var tag models.Tag
		err := dbpool.QueryRow(context.Background(), data.QueryTagByID, tagID, docID).Scan(&tag.ID, &tag.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		render.JSON(w, r, tag)
	}
}

func CreateTag(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tag CreateTagRequest
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var dbTag models.Tag
		err := dbpool.QueryRow(context.Background(), data.InsertTag, tag.Name, tag.DocumentationId, r.Context().Value(constant.UserKey).(*svc.UserClaims).UserID).Scan(&dbTag.ID, &dbTag.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, dbTag)
	}
}

func UpdateTag(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID := chi.URLParam(r, "id")
		var tag models.Tag
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := dbpool.QueryRow(context.Background(), data.UpdateTag, tag.Name, tagID).Scan(&tag.ID, &tag.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, tag)
	}
}

func DeleteTag(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID := chi.URLParam(r, "id")
		_, err := dbpool.Exec(context.Background(), data.DeleteTag, tagID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, map[string]string{"status": "ok"})
	}
}
