package app

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Server struct {
	Port   string
	DB     *sql.DB
	Router *mux.Router
}

type User struct {
	Id       int
	Email    string
	Password string
}

type UserService interface {
	CreateUser(email string, password string) (u *User, err error)
	FindUserByCredentials(email string, password string) (u *User, err error)
}
