package storage

type Storager interface {
	RegisterUser(string, string) error
	AuthUser(string, string) bool
	GetUserConversations(string) []interface{}
	IsUserExists(string) bool
	UpdateUserSessionCookie(string, string) error
	GetUsernameByCookie(string) (string, error)
}
