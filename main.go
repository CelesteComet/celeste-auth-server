package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CelesteComet/celeste-auth-server/pkg/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/matryer/respond.v1"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	DB *sql.DB
}

type ErrorMessages struct {
	ServerError string
	BadCreds    string
}

var errorMessages ErrorMessages = ErrorMessages{
	ServerError: "Server error, please try again later",
	BadCreds:    "Username or Password is Incorrect",
}

func (h *UserHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (h *UserHandler) FindByCredentials() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromRequest(r)
		if err != nil {
			errors := []string{errorMessages.BadCreds}
			respond.With(w, r, http.StatusUnauthorized, errors)
			return
		}

		// Look for the user in the database
		dbUser := User{}
		h.DB.QueryRow("select * from member where email = $1", user.Email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password)

		if &dbUser.Email == nil {
			errors := []string{errorMessages.BadCreds}
			respond.With(w, r, http.StatusUnauthorized, errors)
			return
		}

		// decode password
		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if err != nil {
			errors := []string{"Email or Password is Incorrect"}
			respond.With(w, r, http.StatusUnauthorized, errors)
			return
		}

		tokenString := h.ProvideToken(&dbUser, w)
		// Set HttpOnly To Prevent Future Tampering
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			HttpOnly: true,
			Path:     "/",
		})
		respond.With(w, r, http.StatusOK, &dbUser)
	})
}

func (h *UserHandler) FindOrCreate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromRequest(r)
		log.Println(r.Body)
		log.Println(user)
		if err != nil {
			errors := []string{errorMessages.BadCreds}
			respond.With(w, r, http.StatusUnauthorized, errors)
			return
		}

		// Look for the user in the database
		dbUser := User{}
		h.DB.QueryRow("select * from member where email = $1", user.Email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password)

		if &dbUser.Email == nil {
			errors := []string{errorMessages.BadCreds}
			respond.With(w, r, http.StatusUnauthorized, errors)
			return
		}

		// decode password
		// err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		// if (err != nil) {
		//   errors := []string{"Email or Password is Incorrect"}
		//   respond.With(w, r, http.StatusUnauthorized, errors)
		//   return
		// }

		h.ProvideToken(&dbUser, w)
		respond.With(w, r, http.StatusOK, &dbUser)
	})
}

func (h *UserHandler) ProvideToken(u *User, w http.ResponseWriter) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    u.ID,
		"email": u.Email,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte("secret"))
	w.Header().Set("JWT", tokenString)
	return tokenString
}

func (h *UserHandler) LoginWithGoogle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := gothic.CompleteUserAuth(w, r); err == nil {
			// t, _ := template.New("foo").Parse(userTemplate)
			// t.Execute(res, gothUser)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	})
}

func (h *UserHandler) handleOAuthCallback() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("HANDLE THAT GOOGLE CALLBACK")
		if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {

			// Look for the user in the database, if the user exists, return the user
			dbUser := User{}
			h.DB.QueryRow("select * from member where email = $1", gothUser.Email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password)

			if &dbUser.Email == nil {
				errors := []string{errorMessages.BadCreds}
				respond.With(w, r, http.StatusUnauthorized, errors)
				return
			}
			h.ProvideToken(&dbUser, w)
			respond.With(w, r, http.StatusOK, &dbUser)
		}
	})
}

// HELPER METHODS
func getUserFromRequest(r *http.Request) (User, error) {
	decoder := json.NewDecoder(r.Body)
	user := User{}
	err := decoder.Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Database Setup
var (
	host     = "celestecomet.c7bjz8zer8ha.us-east-1.rds.amazonaws.com"
	port     = 5432
	user     = os.Getenv("AWS_DB_USERNAME")
	password = os.Getenv("AWS_DB_PASSWORD")
	dbname   = "CelesteComet"
)

var (
	connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
)

func main() {

	// Initialize Database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server connection successful")
	defer db.Close()

	// Model Handlers
	userHandler := UserHandler{DB: db}

	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.Cors)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	// OAuth
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:1337/auth/google/callback"),
	)

	// Routes
	router.Handle("/users", userHandler.Create()).Methods("POST")
	router.Handle("/login", userHandler.FindByCredentials()).Methods("POST", "OPTIONS")
	router.Handle("/oauth", userHandler.FindOrCreate()).Methods("POST")

	router.Handle("/auth/{provider}", userHandler.LoginWithGoogle()).Methods("GET")
	router.Handle("/auth/{provider}/callback", userHandler.handleOAuthCallback()).Methods("GET")

	server := &http.Server{
		Handler: loggedRouter,
		Addr:    ":1337",
	}

	log.Println("Server is now running on port 1337")
	server.ListenAndServe()
}
