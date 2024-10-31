package models

import "github.com/google/uuid"

type Tag struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
}
