package models

import (
	"context"
	"database/sql"
	"fmt"
	"irule-api/data"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID
	Name      sql.NullString
	LastName  sql.NullString
	Password  string
	Email     string
	CreatedAt time.Time
	Role      string
}

func FindByEmail(dbPool *pgxpool.Pool, email string) (*User, error) {
	var user User
	err := dbPool.QueryRow(context.Background(), data.QueryUser, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		fmt.Println(err)
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

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
