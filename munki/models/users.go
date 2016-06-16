package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

// User represents an API user
type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	APIKeys  []APIKey `json:"api_keys"`
}

// NewUser creates a new user
func NewUser(username, password string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("Username or Password empty")
	}
	pass := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user := &User{
		Username: username,
		Password: string(hashedPassword),
	}
	return user, nil
}

// AddKey adds an api key to user
func (u *User) AddKey(key APIKey) {
	u.APIKeys = append(u.APIKeys, key)
}

// RevokeKey invalidates an API Key
func (u *User) RevokeKey(key string) {
	// if there are no keys to revoke, return
	if len(u.APIKeys) < 1 {
		return
	}
	for _, userKey := range u.APIKeys {
		if userKey.Key == key {
			userKey.Valid = false
			return
		}
	}
}

// APIKey is an API Key
type APIKey struct {
	Key   string `json:"key"`
	Valid bool   `json:"valid"`
}

type keys map[string]bool

// TODO set from flag
var secret = []byte("insecure-signing-key")

// NewAPIKey creates a new API Key
func NewAPIKey() (*APIKey, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["iat"] = time.Now().Unix()
	tokString, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}
	key := &APIKey{
		Key:   tokString,
		Valid: true,
	}
	return key, nil
}
