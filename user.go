package main

import (
	"log"
)

// type User struct {
//   Id       int    `json:"id"`
//   Email    string `json:"email"`
//   Password string `json:"password"`
// }

func (u *User) CreateUser() User {
	log.Println("creating a user")
	user := User{}
	return user
}

func (u *User) FindUserByCred() User {
	log.Println("finding user by credentials")
	user := User{}
	return user
}

func (u *User) FindById(id int) User {
	user := User{}
	return user
}
