package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/groob/ape/models"
)

// TODO set from flag
var secret = []byte("insecure-signing-key")

// verify using jwt token
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		accept := acceptHeader(r)
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			respondError(rw, http.StatusUnauthorized, accept,
				fmt.Errorf("Authentication error: %v", err))
			return
		}
		// successful, pass
		next.ServeHTTP(rw, r)
	})
}

func handleBasicAuth() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		accept := acceptHeader(r)
		username, password, ok := r.BasicAuth()
		if !ok {
			rw.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			respondError(rw, http.StatusUnauthorized, accept,
				errors.New("Authentication Error: Unable to get username and password from header"))
			return
		}
		// check credentials here
		fmt.Println(username, password)
		// credentials ok, issue a token
		token, err := newToken()
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept, err)
			return
		}
		respondOK(rw, &jwtToken{token}, accept)
	}
}

// takes a valid token and returns a new one
func handleTokenRefresh() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// validate request
		accept := acceptHeader(r)
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		// TODO some steps to validate token based on expiration or claims...
		if err != nil || !token.Valid {
			respondError(rw, http.StatusUnauthorized, accept,
				fmt.Errorf("Authentication error: %v", err))
			return
		}
		// issue new token
		updatedToken, err := newToken()
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept, err)
			return
		}
		respondOK(rw, &jwtToken{updatedToken}, accept)
	}
}

type jwtToken struct {
	Token string `json:"token"`
}

func newToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	return token.SignedString(secret)
}

func (t *jwtToken) View(accept string) (*models.Response, error) {
	data, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return nil, err
	}
	resp := &models.Response{
		Data: data,
	}
	return resp, nil
}
