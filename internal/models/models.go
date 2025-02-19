package models

import "time"

type DataTypeEnum string

const (
	DataTypeUNDEFINE DataTypeEnum = "UNDEFINE"
	DataTypePAIR     DataTypeEnum = "PAIR"
	DataTypeTEXT     DataTypeEnum = "TEXT"
	DataTypeBINARY   DataTypeEnum = "BINARY"
	DataTypeCARD     DataTypeEnum = "CARD"
)

type DataStoreFormat struct {
	ID         uint32       `json:"id"`
	DataType   DataTypeEnum `json:"data_type"`
	DataInfo   string       `json:"data_info"`
	Meta       string       `json:"meta"`
	UploadedAt time.Time    `json:"uploaded_at"`
}

type UserAccountData struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type Account struct {
	ID        uint32    `json:"id"`
	Username  string    `json:"login"`
	Password  string    `json:"password"`
	CreatedOn time.Time `json:"created_on"`
}

type PairData struct {
	ID   uint32 `json:"id,omitempty"`
	Key  string `json:"key"`
	Pwd  string `json:"pwd"`
	Meta string `json:"meta,omitempty"`
}

type TextData struct {
	ID   uint32 `json:"id,omitempty"`
	Text string `json:"text"`
	Meta string `json:"meta,omitempty"`
}

type BinaryData struct {
	ID       uint32 `json:"id,omitempty"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Binary   []byte `json:"binary"`
	Meta     string `json:"meta,omitempty"`
}

type CardData struct {
	ID        uint32 `json:"id,omitempty"`
	CardNum   string `json:"card_num"`
	CardCvv   string `json:"card_cvv"`
	CardOwner string `json:"card_owner"`
	CardExp   string `json:"card_exp"`
	Meta      string `json:"meta,omitempty"`
}
