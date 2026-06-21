package database

import (
	"errors"

	"argument/internal/auth"
)

// User contains everything required for a user
type User struct {
	ID       int64
	Username string
	Password string // Should be encrypted at all times
}

// CreateUser attempts to insert a new user to the database
func CreateUser(user User) (int64, error) {
	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	result, err := db.Exec(query, user.Username, user.Password)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// VerifyLogin checks if login information is correct and returns a user ID
func VerifyLogin(username string, password string) (User, error) {
	user := User{Username: username}
	query := "SELECT id, password FROM users WHERE username = ?"
	err := db.QueryRow(query, &user.Username).Scan(&user.ID, &user.Password)

	if err != nil {
		return user, errors.New("Username does not exist")
	}

	if auth.CheckPassword(password, user.Password) {
		return user, nil
	}

	// Password does not match
	return user, errors.New("Password is invalid")
}

// GetUser finds a user by ID (if any) and returns it
func GetUser(id int64) (User, error) {
	var user User

	query := "SELECT id, username, password FROM users WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return user, err
	}

	return user, nil
}
