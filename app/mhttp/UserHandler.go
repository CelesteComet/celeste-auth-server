package mhttp

import (
	"database/sql"
	"encoding/json"
	"github.com/CelesteComet/celeste-auth-server/app"
	"github.com/dgrijalva/jwt-go"
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
		log.Println("Attempting to create a new user...")
		decoder := json.NewDecoder(r.Body)
		user := app.User{}
		err := decoder.Decode(&user)
		if err != nil {
			log.Println("Error decoding")
			log.Println(err)
		}

		if user.Email == "" || user.Password == "" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 1)
		if err != nil {
			log.Println(err)
			http.Error(w, "bcrypt error", 500)
		}

		// Check if user exists already
		// Otherwise create a user and save to database
		id := 0
		h.DB.QueryRow("insert into member (email, password) values ($1, $2) RETURNING id", user.Email, hash).Scan(&id)
		if id == 0 {
			http.Error(w, "bad post", http.StatusUnauthorized)
			return 
		}
		user.Id = id 
		tokenString := h.ProvideToken(&user)
		w.Header().Set("JWT", tokenString)

		json.NewEncoder(w).Encode(&user)
	})
}

func (h *UserHandler) FindUserByCredentials() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user := app.User{}
		err := decoder.Decode(&user)
		if err != nil {
			log.Println("Error decoding")
			log.Println(err)
		}		
		// Look for the user in the database
		dbUser := app.User{}
		h.DB.QueryRow("select * from member where email = $1", user.Email).Scan(&dbUser.Id, &dbUser.Email, &dbUser.Password)

		if (&dbUser.Email == nil) {
			http.Error(w, "user doesnt exist", http.StatusUnauthorized)
			return
		}

		// decode password
		log.Println(user.Password)
		log.Println(dbUser.Email)
		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if (err != nil) {
			// Bad password
			http.Error(w, "bad password", http.StatusUnauthorized)
			return
		}
		tokenString := h.ProvideToken(&dbUser)
		w.Header().Set("JWT", tokenString)
		json.NewEncoder(w).Encode(&dbUser)		
	})
}

func (h *UserHandler) ProvideToken(u *app.User) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": u.Id,
		"email": u.Email,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte("secret"))
	
	return tokenString
}

