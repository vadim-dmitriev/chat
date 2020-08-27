package storage

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	initSQLScriptName = "init.sql"
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

// IsUserExists проверяет существует ли пользователь с именем username
func (s Sqlite) IsUserExists(username string) bool {
	if err := s.QueryRow(`SELECT username FROM users WHERE username = $1;`, username).Scan(); err == sql.ErrNoRows {
		return false
	}

	return true
}

// RegisterUser заносит новую запись в таблицу users и таблицу cookies
func (s Sqlite) RegisterUser(username, password string) error {
	if s.IsUserExists(username) {
		return fmt.Errorf("username %s already exists", username)
	}

	fingerprint, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	res, err := s.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, username, string(fingerprint))
	if err != nil {
		return err
	}
	newUserPrimaryKey, err := res.LastInsertId()

	if _, err := s.Exec(`INSERT INTO cookies (user_id) VALUES ($1)`, newUserPrimaryKey); err != nil {
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

/* GetUserConversations возвращает список бесед пользователя
Структура ответа:
	[
		{
			"name": "{conversation_name}",
			"is_dialog": boolean,
			"messages": [
				{
					"value": string({message_value}),
					"sender": string({username}),
					"time": string({datetime})
				},
			]
		},
		{
			...
		}
	]
*/
func (s Sqlite) GetUserConversations(username string) ([]map[string]interface{}, error) {
	rows, err := s.Query(`
		SELECT conversation_id, name, is_dialog FROM members JOIN conversations USING (conversation_id) WHERE user_id = (
			SELECT user_id FROM users where username = $1
		); 
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations = make([]map[string]interface{}, 0)
	var conversationName interface{}
	var conversationID int
	var isDialog int
	for rows.Next() {
		if err := rows.Scan(&conversationID, &conversationName, &isDialog); err != nil {
			return nil, err
		}
		conversation := make(map[string]interface{})
		conversation["is_dialog"] = isDialog

		if isDialog == 1 {
			s.DB.QueryRow(`
				SELECT username FROM members JOIN users USING (user_id) WHERE conversation_id = $1 AND user_id != 
					(SELECT user_id FROM users WHERE username = $2)
			`, conversationID, username).Scan(&conversationName)
		}
		conversation["name"] = conversationName.(string)
		messagesRows, err := s.Query(`
			SELECT value, user_id, time FROM  conversations JOIN messages USING (conversation_id)
			WHERE conversation_id = $1
			ORDER BY time DESC
			LIMIT 10
		`, conversationID)
		if err == sql.ErrNoRows {
			conversation["messages"] = []map[string]interface{}{
				{
					"value":  "Нет сообщений...",
					"sender": "",
					"time":   "",
				},
			}

			conversations = append(conversations, conversation)
			continue
		}

		messages := make([]map[string]interface{}, 0)
		var value, time string
		var userID int
		for messagesRows.Next() {
			if err := messagesRows.Scan(&value, &userID, &time); err != nil {
				return nil, err
			}
			messages = append(messages, map[string]interface{}{
				"value":  value,
				"sender": userID,
				"time":   time,
			})
		}
		conversation["messages"] = messages
		conversations = append(conversations, conversation)

	}

	return conversations, nil
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

// GetUsernameByCookie возвращает имя пользователя (username) исходя из куки
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

// UpdateUserSessionCookie устанавливает новое значение сессионной куки пользователя
func (s Sqlite) UpdateUserSessionCookie(newSessionCookieValue, username string) error {
	if _, err := s.Exec(`UPDATE cookies SET value = $1 WHERE user_id = (
		SELECT user_id FROM users WHERE username = $2
	);`, newSessionCookieValue, username); err != nil {
		return err
	}

	return nil
}

func (s Sqlite) SetMessage(value, senderName string, conversationID string) error {
	_, err := s.DB.Exec(`
		INSERT INTO
		messages(value, user_id, conversation_id)
		VALUES (
				$1,
			   (SELECT user_id FROM users WHERE username = $2),
			   (SELECT conversation_id FROM conversations JOIN members USING(conversation_id) WHERE user_id = (SELECT user_id FROM users WHERE username == $2))
			   );
	`, value, senderName, conversationID)

	return err
}

func (s Sqlite) SetDialog(memberOne, memberTwo string) error {
	result, err := s.Exec(`
		INSERT INTO conversations(name, is_dialog) VALUES($1, $2)
	`, nil, 1)
	if err != nil {
		return err
	}
	newDialogFK, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if _, err := s.Exec(`
		INSERT INTO members(user_id, conversation_id) VALUES(
			(SELECT user_id FROM users WHERE username = $1),
			$2
		);
	`, memberOne, newDialogFK); err != nil {
		return err
	}

	if _, err := s.Exec(`
		INSERT INTO members(user_id, conversation_id) VALUES(
			(SELECT user_id FROM users WHERE username = $1),
			$2
		);
	`, memberTwo, newDialogFK); err != nil {
		return err
	}

	return nil
}
