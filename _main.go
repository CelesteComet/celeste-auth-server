package main

import (
	"database/sql"
	"fmt"
	"github.com/CelesteComet/celeste-auth-server/app"
	"github.com/CelesteComet/celeste-auth-server/app/mhttp"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
  "github.com/markbates/goth"
  "github.com/markbates/goth/gothic"
  "github.com/markbates/goth/providers/google"	
	"log"
	"net/http"
	"os"
	"github.com/gorilla/handlers"
)

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

type corsHandler struct {
	next http.Handler
}

func (handler *corsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "JWT")
	w.Header().Set("Access-Control-Expose-Headers", "Jwt")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")

  if (r.Method == "OPTIONS") {
  	w.WriteHeader(http.StatusOK)
  	return
  } else {
  	handler.next.ServeHTTP(w, r)
  }
}

func withCors(h http.Handler) http.Handler {
	return &corsHandler{next: h}
}

// Oauth Middleware
type oauthHandler struct {
	next http.Handler 
}

var us mhttp.UserHandler

func (handler *oauthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Println("DOING ITS THING")
  user, err := gothic.CompleteUserAuth(res, req)
  if err != nil {
    fmt.Fprintln(res, err)
    return
  }		
  log.Println(user)	
  log.Println(res)
  us.CreateOrFindUser(user)
}


func main() {
	log.Println("Starting Authentication Service")
	log.Println("Connecting to AWS RDS Postgresql server")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server connection successful")
	defer db.Close()


	router := &mux.Router{}
	s := &app.Server{Port: ":1337", DB: db, Router: router}
	us = mhttp.UserHandler{DB: db}

	// Routes
	s.Router.Handle("/users", us.CreateUser()).Methods("POST")
	s.Router.Handle("/login", us.FindUserByCredentials()).Methods("POST")

	// OAuth
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:1337/auth/google/callback"),	
	)


	s.Router.HandleFunc("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
    if _ , err := gothic.CompleteUserAuth(res, req); err == nil {
      // t, _ := template.New("foo").Parse(userTemplate)
      // t.Execute(res, gothUser)
    } else {
      gothic.BeginAuthHandler(res, req)
    }		
	})	


	// Middleware
	withCorsRouter := withCors(s.Router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, withCorsRouter)

	http.ListenAndServe(s.Port, loggedRouter)
	log.Println("Service is now running on port 1337")
}
