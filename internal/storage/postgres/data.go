package postgres

import (
	"context"

	"github.com/bbquite/go-pass-keeper/internal/models"
	"github.com/bbquite/go-pass-keeper/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (storage *DBStorage) CreateData(ctx context.Context, data *models.DataStoreFormat) (models.DataStoreFormat, error) {
	sqlString := `
		INSERT INTO public.pass_keeper_data (data_type, data_info, meta, account_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, data_type, data_info, meta, uploaded_at
	`

	var resultData models.DataStoreFormat
	accountID := ctx.Value(utils.UserIDKey).(uint32)
	args := []any{data.DataType, data.DataInfo, data.Meta, accountID}

	row := storage.DB.QueryRowContext(ctx, sqlString, args...)
	err := row.Scan(
		&resultData.ID,
		&resultData.DataType,
		&resultData.DataInfo,
		&resultData.Meta,
		&resultData.UploadedAt,
	)
	if err != nil {
		return resultData, err
	}

	return resultData, nil
}

func (storage *DBStorage) GetDataList(ctx context.Context) ([]models.DataStoreFormat, error) {
	sqlStringSelect := `
		SELECT id, data_type, data_info, meta, uploaded_at
		FROM public.pass_keeper_data
		WHERE account_id = $1;
	`

	var result []models.DataStoreFormat
	accountID := ctx.Value(utils.UserIDKey).(uint32)

	rows, err := storage.DB.QueryContext(ctx, sqlStringSelect, accountID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var storedItem models.DataStoreFormat

		err := rows.Scan(&storedItem.ID, &storedItem.DataType, &storedItem.DataInfo, &storedItem.Meta, &storedItem.UploadedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, storedItem)
	}

	return result, nil
}

func (storage *DBStorage) UpdateData(ctx context.Context, data *models.DataStoreFormat) error {
	sqlString := `
		UPDATE public.pass_keeper_data 
		SET data_info = $1, meta = $2, uploaded_at = NOW()
		WHERE id = $3 AND account_id = $4 AND data_type = $5
	`

	accountID := ctx.Value(utils.UserIDKey).(uint32)
	args := []any{data.DataInfo, data.Meta, data.ID, accountID, data.DataType}

	_, err := storage.DB.ExecContext(ctx, sqlString, args...)
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) DeleteData(ctx context.Context, storedDataID uint32) error {
	sqlString := `
		DELETE FROM public.pass_keeper_data
		WHERE id = $1 and account_id = $2
	`

	accountID := ctx.Value(utils.UserIDKey).(uint32)

	_, err := storage.DB.ExecContext(ctx, sqlString, storedDataID, accountID)
	if err != nil {
		return err
	}

	return nil
}
