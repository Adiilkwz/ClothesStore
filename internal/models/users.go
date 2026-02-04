package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
)

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"-"`
	Created        time.Time `json:"created_at"`
	Address        string    `json:"address"`
	Role           string    `json:"role"` // "user" or "admin"
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password, address string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, password_hash, role, address, created_at)
	VALUES ($1, $2, $3, 'customer', $4, NOW())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword), address)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPasswd []byte

	stmt := "SELECT id, password_hash FROM users WHERE email = $1"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPasswd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPasswd, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (m *UserModel) GetByID(id int) (*User, error) {
	stmt := `SELECT id, name, email, address, created_at, FROM users WHERE id = $1`
	u := &User{}
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Address, &u.Created)
	return u, err
}

func (m *UserModel) Update(id int, name, address string) error {
	stmt := `UPDATE users SET name = $1, address = $2 WHERE id = $3`
	_, err := m.DB.Exec(stmt, name, address, id)
	return err
}
