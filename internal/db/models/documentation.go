package models

import "github.com/google/uuid"

type Documentation struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
}
