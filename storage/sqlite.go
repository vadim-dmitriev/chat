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

	db.Exec(`CREATE table users (
		id int NOT NULL PRIMARY KEY,
		login text NOT NULL,
		password text NOT NULL
	);`)

	db.Exec(`CREATE TABLE conversations (
		id int not null primary key,
		name text,
		member int,
		FOREIGN KEY (member) references members (conversation)
	);`)

	db.Exec(`CREATE table messages (
		id int NOT NULL PRIMARY KEY,
		value text,
		sender int NOT NULL,
		receiver int NOT NULL,
		foreign key (sender) references users (id),
		foreign key (receiver) references users (id)
	);`)

	db.Exec(`CREATE table members (
		user int not null,
		conversation int not null,
		foreign key (user) references users (id),
		foreign key (conversation) references conversations (id)
	);`)

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
