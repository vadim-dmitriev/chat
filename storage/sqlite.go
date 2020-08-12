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

	_, err = db.Exec(`CREATE table users (
		login text NOT NULL PRIMARY KEY,
		password text NOT NULL
	);`)
	// if err != nil {
	// 	panic(err)
	// }

	_, err = db.Exec(`CREATE TABLE conversations (
		id int not null primary key,
		name text NOT NULL
	);`)
	// if err != nil {
	// 	panic(err)
	// }
	_, err = db.Exec(`CREATE table messages (
		id int NOT NULL PRIMARY KEY,
		value text,
		sender text NOT NULL,
		receiver int NOT NULL,
		foreign key (sender) references users (login),
		foreign key (receiver) references conversations (id)
	);`)
	// if err != nil {
	// 	panic(err)
	// }
	_, err = db.Exec(`CREATE table members (
		user int not null,
		conversation int not null,
		foreign key (user) references users (id),
		foreign key (conversation) references conversations (id)
	);`)
	// if err != nil {
	// 	panic(err)
	// }

	return Sqlite{
		db,
	}
}

func (s Sqlite) RegisterUser(login, passowrd string) error {
	// TODO: Проверить, есть ли пользователь с таким же именем
	result, err := s.Exec("insert into users (login, password) values ($1, $2)", login, passowrd)
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}

func (s Sqlite) AuthUser(login, password string) bool {
	// select password from users WHERE login = "123";
	row := s.QueryRow(`SELECT password FROM users WHERE login = $1`, login)

	var passwordFromDB string
	if err := row.Scan(&passwordFromDB); err != nil {
		panic(err)
	}

	if password != passwordFromDB {
		return false
	}
	return true
}

func (s Sqlite) GetUserConversations(login string) []interface{} {

	return nil
}
