package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vadim-dmitriev/chat/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	SigningMethodHS256 = jwt.SigningMethodHS256
)

// JWT имплементирует интерфейс IAuth. Реализует аутентификацию
// пользователя при помощи JWT
type JWT struct {
	Repo IUserRepository

	Method       jwt.SigningMethod
	ExpiringTime time.Duration
	Secret       []byte
}

type authClaims struct {
	*jwt.StandardClaims
	User model.User `json:"user"`
}

// SignIn аутентифицирует пользователя
func (a JWT) SignIn(username, password string) (string, error) {
	user, err := a.Repo.GetUser(username)
	if err != nil {
		return "", err
	}

	if !isPasswordMatch(password, user.Password) {
		return "", fmt.Errorf("password does not match")
	}

	claims := authClaims{
		User: user,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().Add(a.ExpiringTime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	return token.SignedString(a.Secret)
}

// SignUp регистрирует нового пользователя
func (a JWT) SignUp(username, password string) error {
	encryptedPassword, err := encryptPassword(password)
	if err != nil {
		return err
	}

	user := model.User{
		Name:     username,
		Password: encryptedPassword,
	}

	return a.Repo.CreateUser(user)
}

// ParseToken валидирует токен пользователя
func (a JWT) ParseToken(tokenString string) (model.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != a.Method.Alg() {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return a.Secret, nil

	})
	if err != nil {
		return model.User{}, err
	}

	if err := token.Claims.Valid(); err != nil {
		return model.User{}, fmt.Errorf("claims is not valid")
	}
	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.User{}, fmt.Errorf("claims is not valid")
	}
	userMap := claimsMap["user"].(map[string]interface{})
	user := model.User{
		ID:       userMap["ID"].(string),
		Name:     userMap["Name"].(string),
		Password: userMap["Password"].(string),
	}

	return user, nil
}

func encryptPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPassword), err
}

func isPasswordMatch(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}
	return true
}
