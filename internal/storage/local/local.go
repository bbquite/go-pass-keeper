package local

import (
	"encoding/json"
	"sync"

	"github.com/bbquite/go-pass-keeper/internal/models"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
)

type ClientStorage struct {
	UserID    *uint32           `json:"user_id"`
	Token     *jwttoken.JWT     `json:"token"`
	PairsList []models.PairData `json:"pairs_list"`
	mx        sync.RWMutex      `json:"-"`
}

func NewClientStorage() *ClientStorage {
	return &ClientStorage{
		UserID:    nil,
		Token:     nil,
		PairsList: nil,
	}
}

func (storage *ClientStorage) IsAuth() bool {
	if storage.Token != nil {
		return true
	}
	return false
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

func (storage *ClientStorage) AddPairs(data models.PairData) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.PairsList = append(storage.PairsList, data)
	return nil
}

func (storage *ClientStorage) Debug() ([]byte, error) {
	test, err := json.Marshal(storage)
	if err != nil {
		return nil, err
	}
	return test, nil
}
