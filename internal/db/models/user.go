package models

import (
	"context"
	"database/sql"
	"irule-api/data"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uuid.UUID      `json:"id"`
	Name           sql.NullString `json:"name"`
	LastName       sql.NullString `json:"last_name"`
	Password       string         `json:"-"`
	Email          string         `json:"email"`
	CreatedAt      time.Time      `json:"created_at"`
	Role           string         `json:"role"`
	OrganizationId uuid.UUID      `json:"organization_id"`
}

func FindByEmail(dbPool *pgxpool.Pool, email string) (*User, error) {
	var user User
	err := dbPool.QueryRow(context.Background(), data.QueryUser, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.OrganizationId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) Create(dbPoll *pgxpool.Pool) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	_, err = dbPoll.Exec(context.Background(), data.QueryCreateUser, bytes, u.Email, u.Role)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CreateUser(dbPoll *pgxpool.Pool) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	_, err = dbPoll.Exec(context.Background(), data.QueryCreateUserOrg, bytes, u.Email, u.Role, u.OrganizationId)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
