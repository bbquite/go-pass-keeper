package server

import (
	"context"
	"fmt"
	encryptor "github.com/bbquite/go-pass-keeper/internal/encryption"
	"github.com/bbquite/go-pass-keeper/internal/models"
	"github.com/bbquite/go-pass-keeper/internal/utils"
	"go.uber.org/zap"
)

type dataStorageRepo interface {
	CreateData(ctx context.Context, accountID uint32, data *models.DataStoreFormat) (models.DataStoreFormat, error)
	GetDataList(ctx context.Context, accountID uint32) ([]models.DataStoreFormat, error)
	UpdateData(ctx context.Context, accountID uint32, data *models.DataStoreFormat) error
	DeleteData(ctx context.Context, accountID uint32, dataID uint32) error

	GetDataByIDForUser(ctx context.Context, accountID uint32, storedDataID uint32) (models.DataStoreFormat, error)
}

type DataService struct {
	store     dataStorageRepo
	encryptor *encryptor.Encryptor
	logger    *zap.SugaredLogger
}

func NewDataService(store dataStorageRepo, encryptorManager *encryptor.Encryptor, logger *zap.SugaredLogger) *DataService {
	return &DataService{
		store:     store,
		encryptor: encryptorManager,
		logger:    logger.Named("DATA"),
	}
}

func (s *DataService) EncryptData(data *models.DataStoreFormat) (*models.DataStoreFormat, error) {
	encryptedInfo, err := s.encryptor.Encrypt(data.DataInfo)
	if err != nil {
		return nil, err
	}

	encryptedMeta, err := s.encryptor.Encrypt(data.Meta)
	if err != nil {
		return nil, err
	}

	data.DataInfo = encryptedInfo
	data.Meta = encryptedMeta

	return data, nil
}

func (s *DataService) DecryptData(data *models.DataStoreFormat) (*models.DataStoreFormat, error) {
	decryptedInfo, err := s.encryptor.Decrypt(data.DataInfo)
	if err != nil {
		return nil, err
	}

	decryptedMeta, err := s.encryptor.Decrypt(data.Meta)
	if err != nil {
		return nil, err
	}

	data.DataInfo = decryptedInfo
	data.Meta = decryptedMeta

	return data, nil
}

func (s *DataService) CreateData(ctx context.Context, data *models.DataStoreFormat) (models.DataStoreFormat, error) {
	var resultData models.DataStoreFormat
	accountID := ctx.Value(utils.UserIDKey).(uint32)

	encryptedData, err := s.EncryptData(data)
	if err != nil {
		return resultData, fmt.Errorf("encryption error: %v", err)
	}

	resultData, err = s.store.CreateData(ctx, accountID, encryptedData)
	if err != nil {
		return resultData, err
	}

	return resultData, nil
}

func (s *DataService) GetDataList(ctx context.Context) ([]models.DataStoreFormat, error) {
	var resultDataList []models.DataStoreFormat
	accountID := ctx.Value(utils.UserIDKey).(uint32)

	resultDataList, err := s.store.GetDataList(ctx, accountID)
	if err != nil {
		return resultDataList, err
	}

	var decryptedDataList []models.DataStoreFormat

	for _, item := range resultDataList {
		dd, err := s.DecryptData(&item)
		if err != nil {
			return resultDataList, fmt.Errorf("decryption error: %v", err)
		}

		decryptedDataList = append(decryptedDataList, *dd)
	}

	return decryptedDataList, nil
}

func (s *DataService) UpdateData(ctx context.Context, data *models.DataStoreFormat) error {

	accountID := ctx.Value(utils.UserIDKey).(uint32)

	_, err := s.store.GetDataByIDForUser(ctx, accountID, data.ID)
	if err != nil {
		return err
	}

	encryptedData, err := s.EncryptData(data)
	if err != nil {
		return fmt.Errorf("encryption error: %v", err)
	}

	err = s.store.UpdateData(ctx, accountID, encryptedData)
	if err != nil {
		return err
	}

	return nil
}

func (s *DataService) DeleteData(ctx context.Context, dataID uint32) error {

	accountID := ctx.Value(utils.UserIDKey).(uint32)

	_, err := s.store.GetDataByIDForUser(ctx, accountID, dataID)
	if err != nil {
		return err
	}

	err = s.store.DeleteData(ctx, accountID, dataID)
	if err != nil {
		return err
	}

	return nil
}
