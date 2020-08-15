package storage

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

// Sqlite имплементирует интерфейс Storager
type Sqlite struct {
	*sql.DB
}

// NewSqlite создает таблицы в БД и возвращает
// пул соединений к ней
func NewSqlite() Storager {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		panic(err)
	}

	db.Exec(`
		CREATE table 'users' (
			'user_id' INTEGER PRIMARY KEY,
			'username' text NOT NULL,
			'password' text NOT NULL
		);
	`)

	db.Exec(`
		CREATE TABLE 'conversations' (
			'conversation_id' INTEGER PRIMARY KEY,
			'name' text NOT NULL
		);
	`)

	db.Exec(`
		CREATE table 'messages' (
			'message_id' INTEGER PRIMARY KEY,
			'value' text NOT NULL,
			'sender' INTEGER NOT NULL,
			'receiver' INTEGER NOT NULL,
			'time' datetime NOT NULL DEFAULT GETDATE,
			FOREIGN KEY ('sender') REFERENCES users ('user_id'),
			FOREIGN KEY ('receiver') REFERENCES conversations ('conversation_id')
		);
	`)

	db.Exec(`
		CREATE table 'members' (
			'member_id' INTEGER PRIMARY KEY,
			'user_id' INTEGER NOT NULL,
			'conversation_id' INTEGER NOT NULL,
			FOREIGN KEY ('user_id') REFERENCES users ('user_id'),
			FOREIGN KEY ('conversation_id') REFERENCES conversations ('conversation_id')
		);
	`)

	db.Exec(`
		CREATE TABLE 'cookies' (
			'cookie_id' INTEGER PRIMARY KEY,
			'user_id' INTEGER NOT NULL,
			'value' text,
			FOREIGN KEY ('user_id') REFERENCES users ('user_id')
		);
	`)

	return Sqlite{
		db,
	}
}

// IsUserExists проверяет существует ли пользователь с именем username
func (s Sqlite) IsUserExists(username string) bool {
	if err := s.QueryRow(`SELECT username FROM users WHERE username = $1;`, username).Scan(); err == sql.ErrNoRows {
		return false
	}

	return true
}

// RegisterUser заносит новую запись в таблицу users
func (s Sqlite) RegisterUser(username, password string) error {
	if s.IsUserExists(username) {
		return fmt.Errorf("username %s already exists", username)
	}

	fingerprint, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if _, err := s.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, username, string(fingerprint)); err != nil {
		return err
	}

	return nil
}

// AuthUser проверяет совпадает ли username и password из БД с тем,
// что ввел пользователь на сайте
func (s Sqlite) AuthUser(username, password string) bool {
	var passwordFromDB string
	if err := s.QueryRow(`SELECT password FROM users WHERE username = $1`, username).Scan(&passwordFromDB); err == sql.ErrNoRows {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordFromDB), []byte(password)); err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}

	return true
}

// GetUserConversations возвращает список бесед пользователя
func (s Sqlite) GetUserConversations(username string) []interface{} {

	return nil
}

// GetUserSessionCookie возвращает значение сессионной куки пользователя
func (s Sqlite) GetUserSessionCookie(username string) string {
	var sessionCookie string
	if err := s.QueryRow(`
		SELECT value FROM cookies WHERE user_id = (
			SELECT user_id FROM users WHERE username = $1
		);
	`, username).Scan(&sessionCookie); err == sql.ErrNoRows {
		return ""
	}

	return sessionCookie
}

func (s Sqlite) GetUsernameByCookie(sessionCookie string) (string, error) {
	var username string

	if err := s.QueryRow(`
		SELECT username FROM users WHERE user_id  = (
			SELECT user_id FROM cookies WHERE value = $1
		);
	`, sessionCookie).Scan(&username); err == sql.ErrNoRows {
		return "", fmt.Errorf("username not found")
	}

	return username, nil
}

// SetUserSessionCookie устанавливает новое значение сессионной куки пользователя
func (s Sqlite) SetUserSessionCookie(newSessionCookieValue, username string) error {
	if _, err := s.Exec(`UPDATE cookies SET value = $1 WHERE user_id = (
		SELECT user_id FROM users WHERE username = $2
	);`, newSessionCookieValue, username); err != nil {
		return err
	}

	return nil
}
