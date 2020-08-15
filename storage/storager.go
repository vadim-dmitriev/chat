package storage

type Storager interface {
	RegisterUser(string, string) error
	AuthUser(string, string) bool
	GetUserConversations(string) []interface{}
	IsUserExists(string) bool
	GetUserSessionCookie(string) string
	SetUserSessionCookie(string, string) error
	GetUsernameByCookie(string) (string, error)
}
