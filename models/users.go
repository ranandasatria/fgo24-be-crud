package models

import (
	"backend/utils"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID       int	`json:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func FindAllUser(search string) ([]User, error) {
	result := utils.RedisClient.Exists(context.Background(), "all-users")
	if result.Val() == 0 {
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

	encoded, err := json.Marshal(users)

	if err != nil {
		return []User{}, err
	}
	utils.RedisClient.Set(context.Background(), "all-users", string(encoded), 0)
	return users, nil
	} else {
		data := utils.RedisClient.Get(context.Background(), "all-users")
		str := data.Val()
		users := []User{}
		json.Unmarshal([]byte(str), &users)
		return users, nil
	}
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
