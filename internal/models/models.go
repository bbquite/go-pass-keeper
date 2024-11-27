package models

import "time"

type UserLoginData struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type UserRegisterData struct {
	Username string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Account struct {
	ID        uint32    `json:"id"`
	Username  string    `json:"login"`
	Password  string    `json:"password"`
	CreatedOn time.Time `json:"created_on"`
}
