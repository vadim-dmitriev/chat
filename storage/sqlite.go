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
			'name' text,
			'is_dialog' INTEGER NOT NULL
		);
	`)

	db.Exec(`
		CREATE table 'messages' (
			'message_id' INTEGER PRIMARY KEY,
			'value' text NOT NULL,
			'user_id' INTEGER NOT NULL,
			'conversation_id' INTEGER NOT NULL,
			'time' datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY ('user_id') REFERENCES users ('user_id'),
			FOREIGN KEY ('conversation_id') REFERENCES conversations ('conversation_id')
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
	{
		"{conversation_name}": {
			"is_dialog": boolean,
			"last_message": {
				"value": string({message_value}),
				"sender": string({username}),
				"time": string({datetime})
			}
		},
		"{conversation_name}": {
			...
		}
	}
*/
func (s Sqlite) GetUserConversations(username string) (map[string]interface{}, error) {
	rows, err := s.Query(`
		SELECT conversation_id, name, is_dialog FROM members JOIN conversations USING (conversation_id) WHERE user_id = (
			SELECT user_id FROM users where username = $1
		); 
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations = make(map[string]interface{})
	var conversationName interface{}
	var conversationID int
	var isDialog int
	for rows.Next() {
		if err := rows.Scan(&conversationID, &conversationName, &isDialog); err != nil {
			fmt.Println(err)
		}
		conversation := make(map[string]interface{})
		conversation["is_dialog"] = isDialog

		var value, time string
		var userID int
		if err := s.QueryRow(`
			SELECT value, user_id, time FROM  conversations JOIN messages USING (conversation_id)
			WHERE conversation_id = $1
			ORDER BY time DESC
			LIMIT 1
		`, conversationID).Scan(&value, &userID, &time); err == sql.ErrNoRows {
			conversation["last_message"] = nil
			conversationName = "conversationName"
			conversations["conversationName"] = conversation
			continue
		}
		lastMessage := make(map[string]interface{})
		lastMessage["value"] = value
		lastMessage["sender"] = userID
		lastMessage["time"] = time

		conversation["last_message"] = lastMessage

		if isDialog == 1 {
			s.DB.QueryRow(`
				SELECT username FROM members JOIN users USING (user_id) WHERE conversation_id = $1 AND user_id != 
					(SELECT user_id FROM users WHERE username = $2)
			`, conversationID, username).Scan(&conversationName)
		}
		conversations[conversationName.(string)] = conversation

	}
	fmt.Println(conversations)
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
	fmt.Println(err)

	return err
}
