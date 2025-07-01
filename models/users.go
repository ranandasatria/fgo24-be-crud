package models

import (
	"backend/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func FindAllUser(search string, page int) ([]User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return []User{}, err
	}

	defer conn.Release()

	limit := 5
	offset := (page - 1) * limit

	rows, err := conn.Query(
		context.Background(),
		`
		SELECT id, name, email, password FROM users
		WHERE name ILIKE $1
		ORDER BY id
		LIMIT $2 OFFSET $3
		`,
		fmt.Sprintf("%%%s%%", search), limit, offset,
	)

	if err != nil {
		return []User{}, err
	}

	users, err := pgx.CollectRows[User](rows, pgx.RowToStructByName)

	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func FindUserByID(id int) (User, error) {
  conn, err := utils.ConnectDB()
  if err != nil {
    return User{}, err
  }
  defer conn.Release()

  row := conn.QueryRow(
    context.Background(),
    `SELECT id, name, email, password FROM users WHERE id = $1`,
    id,
  )

  var user User
  err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
  if err != nil {
    if err == pgx.ErrNoRows {
      return User{}, fmt.Errorf("user not found")
    }
    return User{}, err
  }

  return user, nil
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

func DeleteUser(id int) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	result, err := conn.Exec(
		context.Background(),
		`
		DELETE FROM users
		WHERE id = $1
		`,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func UpdateUser(id int, user User) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	result, err := conn.Exec(
		context.Background(),
		`UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4`,
		user.Name, user.Email, user.Password, id,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
