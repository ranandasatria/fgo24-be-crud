package models

import (
	"backend/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID       int
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func FindAllUser(search string) ([]User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return []User{}, err
	}

	defer conn.Release()

	rows, err := conn.Query(
		context.Background(),
		`
		SELECT id, name, email, password FROM users
		WHERE name ILIKE $1
		`,
		fmt.Sprintf("%%%s%%", search),
	)

	if err != nil {
		return []User{}, err
	}

	users, err := pgx.CollectRows[User](rows, pgx.RowToStructByName)
	if err != nil {
		return []User{}, err
	}
	return users, err
}

func FindOneUserByEmail(email string) (User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return User{}, err
	}

	defer conn.Release()

	rows, err := conn.Query(
		context.Background(),
		`
		SELECT id, name, email, password FROM users
		WHERE email = $1
		`,
		email,
	)

	if err != nil {
		return User{}, err
	}

	user, err := pgx.CollectOneRow[User](rows, pgx.RowToStructByName)
	if err != nil {
		return User{}, err
	}
	return user, err
}


func CreateUser(user User) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		context.Background(),
		`
		INSERT INTO users (name, email, password)
		VALUES
		($1, $2, $3)
		`, 
		user.Name, user.Email, user.Password,
	)
	return err
}
