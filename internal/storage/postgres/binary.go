package postgres

import (
	"context"
	"github.com/bbquite/go-pass-keeper/internal/utils"
	"time"

	"github.com/bbquite/go-pass-keeper/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (storage *DBStorage) CreateBinaryData(ctx context.Context, data *models.BinaryData) (models.BinaryData, error) {
	sqlString := `
		INSERT INTO public.binary_data (binary_data, meta, account_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, binary_data, meta, uploaded_at
	`

	var resultData models.BinaryData
	accountID := ctx.Value(utils.UserIDKey).(uint32)

	row := storage.DB.QueryRowContext(ctx, sqlString, data.Binary, data.Meta, accountID)
	err := row.Scan(
		&resultData.ID,
		&resultData.Binary,
		&resultData.Meta,
		&resultData.UploadedAt,
	)
	if err != nil {
		return resultData, err
	}

	return resultData, nil
}

func (storage *DBStorage) GetBinaryDataList(ctx context.Context) ([]models.BinaryData, error) {
	sqlStringSelect := `
		SELECT id, binary_data, meta, uploaded_at
		FROM public.binary_data
		WHERE account_id = $1;
	`

	var result []models.BinaryData
	accountID := ctx.Value(utils.UserIDKey).(uint32)

	rows, err := storage.DB.QueryContext(ctx, sqlStringSelect, accountID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	for rows.Next() {
		var binaryItem models.BinaryData

		err := rows.Scan(&binaryItem.ID, &binaryItem.Binary, &binaryItem.Meta, &binaryItem.UploadedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, binaryItem)
	}

	return result, nil
}

func (storage *DBStorage) UpdateBinaryData(ctx context.Context, data *models.BinaryData) error {
	sqlString := `
		UPDATE public.binary_data 
		SET binary_data = $1, meta = $2, uploaded_at = $3
		WHERE id = $4
	`

	_, err := storage.DB.ExecContext(ctx, sqlString, data.Binary, data.Meta, time.Now(), data.ID)
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) DeleteBinaryData(ctx context.Context, binaryDataID uint32) error {
	sqlString := `
		DELETE FROM public.binary_data
		WHERE id = $1 and account_id = $2
	`

	_, err := storage.DB.ExecContext(ctx, sqlString, binaryDataID)
	if err != nil {
		return err
	}

	return nil
}
