package models

import "time"

type dataTypeEnum string

const (
	DataTypePAIR   dataTypeEnum = "PAIR"
	DataTypeTEXT   dataTypeEnum = "TEXT"
	DataTypeBINARY dataTypeEnum = "BINARY"
	DataTypeCARD   dataTypeEnum = "CARD"
)

type DataStoreFormat struct {
	ID         uint32       `json:"id"`
	DataType   dataTypeEnum `json:"data_type"`
	DataInfo   string       `json:"data_info"`
	Meta       string       `json:"meta"`
	UploadedAt time.Time    `json:"uploaded_at"`
}

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
	Key  string `json:"key"`
	Pwd  string `json:"pwd"`
	Meta string `json:"meta"`
}

type TextData struct {
	Text string `json:"text"`
	Meta string `json:"meta"`
}

type BinaryData struct {
	FileName string `json:"file_name"`
	Binary   []byte `json:"binary"`
	Meta     string `json:"meta"`
}

type CardData struct {
	CardNum   string `json:"card_num"`
	CardCvv   string `json:"card_cvv"`
	CardOwner string `json:"card_owner"`
	CardExp   string `json:"card_exp"`
	Meta      string `json:"meta"`
}
