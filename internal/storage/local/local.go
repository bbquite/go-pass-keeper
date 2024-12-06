package local

import (
	"encoding/json"
	"sync"

	"github.com/bbquite/go-pass-keeper/internal/models"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
)

type ClientStorage struct {
	UserID       *uint32             `json:"user_id"`
	Token        *jwttoken.JWT       `json:"token"`
	PairsList    []models.PairData   `json:"pairs_list"`
	TextsList    []models.TextData   `json:"texts_list"`
	BinariesList []models.BinaryData `json:"binary_list"`
	CardsList    []models.CardData   `json:"cards_list"`
	mx           sync.RWMutex        `json:"-"`
}

func NewClientStorage() *ClientStorage {
	return &ClientStorage{
		UserID: nil,
		Token:  nil,
	}
}

func (storage *ClientStorage) IsAuth() bool {
	return storage.Token != nil
}

func (storage *ClientStorage) SetUserID(userID *uint32) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.UserID = userID
	return nil
}

func (storage *ClientStorage) SetToken(token *jwttoken.JWT) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.Token = token
	return nil
}

func (storage *ClientStorage) GetToken() string {
	if storage.Token != nil {
		return storage.Token.Token
	}
	return ""
}

func (storage *ClientStorage) AddPairs(data models.PairData) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.PairsList = append(storage.PairsList, data)
	return nil
}

func (storage *ClientStorage) GetPairs() ([]models.PairData, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	return storage.PairsList, nil
}

func (storage *ClientStorage) AddTexts(data models.TextData) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.TextsList = append(storage.TextsList, data)
	return nil
}

func (storage *ClientStorage) GetTexts() ([]models.TextData, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	return storage.TextsList, nil
}

func (storage *ClientStorage) AddBinaries(data models.BinaryData) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.BinariesList = append(storage.BinariesList, data)
	return nil
}

func (storage *ClientStorage) GetBinary() ([]models.BinaryData, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	return storage.BinariesList, nil
}

func (storage *ClientStorage) AddCards(data models.CardData) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.CardsList = append(storage.CardsList, data)
	return nil
}

func (storage *ClientStorage) GetCards() ([]models.CardData, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	return storage.CardsList, nil
}

func (storage *ClientStorage) ClearStorage() {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.PairsList = nil
	storage.TextsList = nil
	storage.BinariesList = nil
	storage.CardsList = nil
}

func (storage *ClientStorage) Debug() ([]byte, error) {
	test, err := json.Marshal(storage)
	if err != nil {
		return nil, err
	}
	return test, nil
}
