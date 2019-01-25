package app

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

type Server struct {
	Port   string
	DB     *sql.DB
	Router *mux.Router
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService interface {
	CreateUser() (u *User, err error)
	FindUserByCredentials() (u *User, err error)
}

type UserHandler interface {
	CreateUser() http.Handler
	FindUserByCredentials() http.Handler
}
