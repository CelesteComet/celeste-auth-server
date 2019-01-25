package mhttp

import (
	"database/sql"
	"encoding/json"
	"github.com/CelesteComet/celeste-auth-service/app"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type UserHandler struct {
	DB *sql.DB
}

var _ app.UserHandler = &UserHandler{}

func (h *UserHandler) CreateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Creating a new user")
		decoder := json.NewDecoder(r.Body)
		user := app.User{}
		err := decoder.Decode(&user)
		if err != nil {
			log.Println(err)
		}

		// Hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 1)
		if err != nil {
			log.Println(err)
		}

		log.Println(user)
		// Check if user exists already
		// Otherwise create a user and save to database
		rows, err := h.DB.Query("insert into member (email, password) values ($1, $2)", user.Email, hash)
		if err != nil {
			log.Println(err)
		}
		json.NewEncoder(w).Encode(rows)
	})
}

func (h *UserHandler) FindUserByCredentials() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
