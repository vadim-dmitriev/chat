package storage

type Storager interface {
	RegisterUser(string, string) error
	AuthUser(string, string) bool
	GetUserConversations(string) ([]string, error)
	IsUserExists(string) bool
	UpdateUserSessionCookie(string, string) error
	GetUsernameByCookie(string) (string, error)
}
