package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/PeterM45/go-postgres-api/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int     `json:"id"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
}

func (db *DB) CreateUser(username, email, password string) (*User, error) {
	// Validate based on config
	if db.Config.User.RequireUsername && username == "" {
		return nil, errors.ErrInvalidInput // instead of errors.New("username required")
	}
	if db.Config.User.RequireEmail && email == "" {
		return nil, errors.ErrInvalidInput // instead of errors.New("email required")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Build dynamic query
	var fields, values []string
	var args []interface{}
	argCount := 1

	if db.Config.User.RequireUsername {
		fields = append(fields, "username")
		values = append(values, fmt.Sprintf("$%d", argCount))
		args = append(args, username)
		argCount++
	}
	if db.Config.User.RequireEmail {
		fields = append(fields, "email")
		values = append(values, fmt.Sprintf("$%d", argCount))
		args = append(args, email)
		argCount++
	}

	// Add password_hash
	fields = append(fields, "password_hash")
	values = append(values, fmt.Sprintf("$%d", argCount))
	args = append(args, hashedPassword)

	query := fmt.Sprintf(
		"INSERT INTO users (%s) VALUES (%s) RETURNING id, username, email",
		strings.Join(fields, ", "),
		strings.Join(values, ", "),
	)

	var user User
	row := db.Pool.QueryRow(context.Background(), query, args...) // Changed from Conn to Pool

	// Scan based on config
	scanArgs := []interface{}{&user.ID}
	if db.Config.User.RequireUsername {
		scanArgs = append(scanArgs, &user.Username)
	}
	if db.Config.User.RequireEmail {
		scanArgs = append(scanArgs, &user.Email)
	}

	if err := row.Scan(scanArgs...); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetUserByID(id int) (*User, error) {
	var user User
	scanArgs := []interface{}{&user.ID}

	fields := []string{"id"}
	if db.Config.User.RequireUsername {
		fields = append(fields, "username")
		scanArgs = append(scanArgs, &user.Username)
	}
	if db.Config.User.RequireEmail {
		fields = append(fields, "email")
		scanArgs = append(scanArgs, &user.Email)
	}

	query := fmt.Sprintf("SELECT %s FROM users WHERE id = $1", strings.Join(fields, ", "))

	if err := db.Pool.QueryRow(context.Background(), query, id).Scan(scanArgs...); err != nil { // Changed from Conn to Pool
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	var fields []string
	if db.Config.User.RequireUsername {
		fields = append(fields, "username")
	}
	if db.Config.User.RequireEmail {
		fields = append(fields, "email")
	}

	query := fmt.Sprintf("SELECT id, %s FROM users", strings.Join(fields, ", "))

	rows, err := db.Pool.Query(context.Background(), query) // Changed from Conn to Pool
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		scanArgs := []interface{}{&user.ID}

		if db.Config.User.RequireUsername {
			scanArgs = append(scanArgs, &user.Username)
		}
		if db.Config.User.RequireEmail {
			scanArgs = append(scanArgs, &user.Email)
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) VerifyUser(email, password string) (*User, error) {
	var user User
	var hashedPassword []byte

	query := "SELECT id, password_hash"
	if db.Config.User.RequireUsername {
		query += ", username"
	}
	query += " FROM users WHERE email = $1"

	scanArgs := []interface{}{&user.ID, &hashedPassword}
	if db.Config.User.RequireUsername {
		scanArgs = append(scanArgs, &user.Username)
	}

	if err := db.Pool.QueryRow(context.Background(), query, email).Scan(scanArgs...); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	return &user, nil
}

func (db *DB) UpdateUser(id int, username, email string) (*User, error) {
	var updates []string
	var args []interface{}
	argCount := 1

	if username != "" {
		updates = append(updates, fmt.Sprintf("username = $%d", argCount))
		args = append(args, username)
		argCount++
	}
	if email != "" {
		updates = append(updates, fmt.Sprintf("email = $%d", argCount))
		args = append(args, email)
		argCount++
	}

	if len(updates) == 0 {
		return nil, errors.ErrInternalServer
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d RETURNING id, username, email",
		strings.Join(updates, ", "),
		argCount,
	)

	var user User
	scanArgs := []interface{}{&user.ID}
	if db.Config.User.RequireUsername {
		scanArgs = append(scanArgs, &user.Username)
	}
	if db.Config.User.RequireEmail {
		scanArgs = append(scanArgs, &user.Email)
	}

	if err := db.Pool.QueryRow(context.Background(), query, args...).Scan(scanArgs...); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) DeleteUser(id int) error {
	result, err := db.Pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}
