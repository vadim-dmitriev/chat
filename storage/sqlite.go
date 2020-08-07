package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	*sql.DB
}

func NewSqlite() Storager {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	db.Exec("create table users (id integer PRIMARY KEY, login text, password text);")

	return Sqlite{
		db,
	}
}

func (s Sqlite) RegisterUser(login, passowrd string) error {
	result, err := s.Exec("insert into users (login, password) values ($1, $2)", login, passowrd)
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}

func (s Sqlite) AuthUser(login, password string) bool {
	// select password from users WHERE login = "123";
	fmt.Println(login)
	row := s.QueryRow(`SELECT password FROM users WHERE login = $1`, login)

	var passwordFromDB string
	if err := row.Scan(&passwordFromDB); err != nil {
		panic(err)
	}
	fmt.Println(row, "!"+passwordFromDB+"!")

	if password != passwordFromDB {
		return false
	}
	return true
}
