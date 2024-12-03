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

type PairData struct {
	ID         uint32    `json:"id"`
	Key        string    `json:"key"`
	Pwd        string    `json:"pwd"`
	Meta       string    `json:"meta"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type TextData struct {
	ID         uint32    `json:"id"`
	Text       string    `json:"text"`
	Meta       string    `json:"meta"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type BinaryData struct {
	ID         uint32    `json:"id"`
	Binary     []byte    `json:"binary"`
	Meta       string    `json:"meta"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type CardData struct {
	ID         uint32    `json:"id"`
	CardNum    string    `json:"card_num"`
	CardCvv    string    `json:"card_cvv"`
	CardOwner  string    `json:"card_owner"`
	CardExp    string    `json:"card_exp"`
	Meta       string    `json:"meta"`
	UploadedAt time.Time `json:"uploaded_at"`
}
