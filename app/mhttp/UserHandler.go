package mhttp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/CelesteComet/celeste-auth-service/app"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	DB *sql.DB
}

var _ app.UserHandler = &UserHandler{}

func (h *UserHandler) CreateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Attempting to create a new user...")
		decoder := json.NewDecoder(r.Body)
		user := app.User{}
		err := decoder.Decode(&user)
		if err != nil {
			log.Println("Error decoding")
			log.Println(err)
		}
		log.Println(user)

		if user.Email == "" || user.Password == "" {
			log.Println("invalid json")
			http.Error(w, "invalid json", 400)
			return
		}

		// Hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 1)
		if err != nil {
			log.Println(err)
			http.Error(w, "bcrypt error", 500)
		}

		log.Println(user)
		// Check if user exists already
		// Otherwise create a user and save to database
		rows, err := h.DB.Query("insert into member (email, password) values ($1, $2)", user.Email, hash)
		if err != nil {
			http.Error(w, "server error", 400)
			return
		}

		// Create a cookie and set it
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "Bruce", Value: "isCool", Expires: expiration}
		http.SetCookie(w, &cookie)

		h.ProvideToken(&user, &w)

		json.NewEncoder(w).Encode(&rows)
	})
}

func (h *UserHandler) FindUserByCredentials() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (h *UserHandler) ProvideToken(u *app.User, w *http.ResponseWriter) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": "wd",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("lipy"))

	fmt.Println(tokenString, err)

	// validate token
	vtoken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, nil
		}
		return []byte("liapy"), nil
	})

	if vtoken.Valid {
		log.Println("YOU GOOD TO GO!")
	} else {
		log.Println("Invalid token")
	}

	return tokenString
}
