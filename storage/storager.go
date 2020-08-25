package storage

type Storager interface {
	RegisterUser(string, string) error
	AuthUser(string, string) bool
	GetUserConversations(string) ([]map[string]interface{}, error)
	IsUserExists(string) bool
	UpdateUserSessionCookie(string, string) error
	GetUsernameByCookie(string) (string, error)
	SetMessage(string, string, string) error
	SetDialog(string, string) error
}
