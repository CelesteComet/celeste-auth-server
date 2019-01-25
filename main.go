package main

import (
	"database/sql"
	"fmt"
	"github.com/CelesteComet/celeste-auth-service/app"
	"github.com/CelesteComet/celeste-auth-service/app/mhttp"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

// Declare the database
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

func SayHello() string {
	return "HELLO"
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
	us := mhttp.UserHandler{DB: db}

	// Routes
	s.Router.Handle("/user", us.CreateUser()).Methods("POST")

	log.Println("Service is now running on port 1337")
	http.ListenAndServe(s.Port, s.Router)
}
