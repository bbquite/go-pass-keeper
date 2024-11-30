package local

import (
	"github.com/bbquite/go-pass-keeper/internal/models"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
	"sync"
)

type ClientStorage struct {
	UserID    *uint32
	Token     *jwttoken.JWT
	PairsList []models.PairsData
	mx        sync.RWMutex
}

func NewClientStorage() *ClientStorage {
	return &ClientStorage{
		UserID:    nil,
		Token:     nil,
		PairsList: nil,
	}
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

func (storage *ClientStorage) AddPairs(data models.PairsData) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.PairsList = append(storage.PairsList, data)
	return nil
}
