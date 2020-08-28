package storage

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vadim-dmitriev/chat/model"
)

const (
	initSQLScriptName = "storage/init.sql"
	dbName            = "app.db"
)

// Sqlite имплементирует интерфейс Storager
type Sqlite struct {
	*sql.DB
}

// NewSqlite создает таблицы в БД и возвращает
// пул соединений к ней
func NewSqlite() Storager {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	initSQLScript, err := ioutil.ReadFile(initSQLScriptName)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(string(initSQLScript))
	if err != nil {
		panic(err)
	}

	return Sqlite{db}
}

// CreateUser заносит новую запись в таблицу users и таблицу cookies
func (s Sqlite) CreateUser(user model.User) error {
	if err := s.DB.QueryRow(`SELECT * FROM USERS WHERE username = $1;`, user.Name).Scan(); err != sql.ErrNoRows {
		return fmt.Errorf("username %s already exists", user.Name)
	}

	fingerprint, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	res, err := s.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, user.Name, string(fingerprint))
	if err != nil {
		return err
	}
	newUserPrimaryKey, err := res.LastInsertId()

	if _, err := s.Exec(`INSERT INTO cookies (user_id) VALUES ($1)`, newUserPrimaryKey); err != nil {
		return err
	}

	return nil
}

func (s Sqlite) GetUser(username string) (model.User, error) {
	var id, password string
	if err := s.DB.QueryRow(`SELECT * FROM USERS WHERE username = $1`, username).Scan(&id, &username, &password); err == sql.ErrNoRows {
		return model.User{}, fmt.Errorf("user '%s' not found", username)
	}

	return model.User{ID: id, Name: username, Password: password}, nil
}

func (s Sqlite) SaveMessage(message model.Message, from model.User, to model.Conversation) error {
	return nil
}
