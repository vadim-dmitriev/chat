package storage

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

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

	if _, err := s.Exec(`INSERT INTO cookies (user_Id) VALUES ($1)`, newUserPrimaryKey); err != nil {
		return err
	}

	return nil
}

func (s Sqlite) GetUser(username string) (model.User, error) {
	var Id, password string
	if err := s.DB.QueryRow(`SELECT * FROM USERS WHERE username = $1`, username).Scan(&Id, &username, &password); err == sql.ErrNoRows {
		return model.User{}, fmt.Errorf("user '%s' not found", username)
	}

	return model.User{Id: Id, Name: username, Password: password}, nil
}

func (s Sqlite) GetUserByToken(token string) (model.User, error) {
	var Id, username, password string
	if err := s.DB.QueryRow(`
		SELECT user_Id, username, password FROM USERS JOIN COOKIES USING(user_Id) WHERE value = $1`, token).Scan(&Id, &username, &password); err == sql.ErrNoRows {
		return model.User{}, fmt.Errorf("user not found")
	}

	return model.User{Id: Id, Name: username, Password: password}, nil
}

func (s Sqlite) GetUserByTokenChat(token string) (model.User, error) {
	return s.GetUserByToken(token)
}

func (s Sqlite) SetUserToken(user model.User, token string) error {
	_, err := s.DB.Exec(`
		UPDATE COOKIES SET value=$1 WHERE user_Id=$2
	`, token, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s Sqlite) SaveMessage(message model.Message) error {
	_, err := s.DB.Exec(`
		INSERT INTO MESSAGES(value, user_Id, conversation_Id) VALUES($1, $2, $3)
	`, message.Text, message.From.Id, message.To.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s Sqlite) GetConversations(user model.User) ([]model.Conversation, error) {
	rows, err := s.DB.Query(`
		SELECT conversation_Id, name, is_dialog, member_Id
		FROM CONVERSATIONS JOIN MEMBERS USING(conversation_Id)
		WHERE user_Id = $1
	`, user.Id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var conversationId, memberId string
	var isDialog int
	var isDialogBool bool
	var conversationName interface{}
	var conversations = make([]model.Conversation, 0)
	for rows.Next() {
		if err := rows.Scan(
			&conversationId,
			&conversationName,
			&isDialog,
			&memberId,
		); err != nil {
			return nil, err
		}

		if isDialog == 1 {
			isDialogBool = true
			if err := s.QueryRow(
				`SELECT username FROM USERS WHERE user_Id = $1

			`, memberId).Scan(&conversationName); err != nil {
				return nil, err
			}
		}

		conv := model.Conversation{
			Id:       conversationId,
			Name:     conversationName.(string),
			IsDialog: isDialogBool,
		}

		// convLastMessages, err := s.GetMessages(conv, 0, 10)
		// if err != nil {
		// 	return nil, err
		// }

		// conv.Messages = convLastMessages
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (s Sqlite) GetMessages(conv model.Conversation, offset, limit int) ([]model.Message, error) {
	rows, err := s.DB.Query(`
		SELECT conversation_Id, value, user_Id, username, time
		FROM MESSAGES JOIN USERS USING(user_Id)
		WHERE conversation_Id = $1
		ORDER BY time DESC
		LIMIT $2
		OFFSET $3
	`, conv.Id, limit, offset)
	if err != nil {
		return nil, err
	}

	var messages = make([]model.Message, 0)
	var messageId, value, senderId, senderName, datetime string
	for rows.Next() {
		if err := rows.Scan(
			&messageId,
			&value,
			&senderId,
			&senderName,
			&datetime,
		); err != nil {
			return nil, err
		}

		messages = append(messages, model.Message{
			Id:   messageId,
			Text: value,
			From: &model.User{
				Id:   senderId,
				Name: senderName,
			},
			To:       &conv,
			Datetime: time.Time{}.String(),
		})
	}

	return messages, nil
}
