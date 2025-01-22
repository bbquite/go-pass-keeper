package postgres

import (
	"context"
	"fmt"

	"github.com/bbquite/go-pass-keeper/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (storage *DBStorage) CreateData(ctx context.Context, accountID uint32, data *models.DataStoreFormat) (models.DataStoreFormat, error) {
	sqlString := `
		INSERT INTO public.pass_keeper_data (data_type, data_info, meta, account_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, uploaded_at
	`

	resultData := data
	args := []any{data.DataType, data.DataInfo, data.Meta, accountID}

	row := storage.DB.QueryRowContext(ctx, sqlString, args...)
	err := row.Scan(
		&resultData.ID,
		&resultData.UploadedAt,
	)
	if err != nil {
		return *resultData, err
	}

	return *resultData, nil
}

func (storage *DBStorage) GetDataList(ctx context.Context, accountID uint32) ([]models.DataStoreFormat, error) {
	sqlStringSelect := `
		SELECT id, data_type, data_info, meta, uploaded_at
		FROM public.pass_keeper_data
		WHERE account_id = $1;
	`

	var result []models.DataStoreFormat

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

	defer rows.Close()

	return result, nil
}

func (storage *DBStorage) GetDataByIDForUser(ctx context.Context, accountID uint32, storedDataID uint32) (models.DataStoreFormat, error) {
	var result models.DataStoreFormat

	sqlString := `
		SELECT id, data_type, data_info, meta, uploaded_at 
		FROM public.pass_keeper_data 
		WHERE id = $1 AND account_id = $2
	`

	row := storage.DB.QueryRowContext(ctx, sqlString, storedDataID, accountID)
	err := row.Scan(&result.ID, &result.DataType, &result.DataInfo, &result.Meta, &result.UploadedAt)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (storage *DBStorage) UpdateData(ctx context.Context, accountID uint32, data *models.DataStoreFormat) error {
	sqlString := `
		UPDATE public.pass_keeper_data 
		SET data_info = $1, meta = $2, uploaded_at = NOW()
		WHERE id = $3 AND account_id = $4 AND data_type = $5
	`

	args := []any{data.DataInfo, data.Meta, data.ID, accountID, data.DataType}

	result, err := storage.DB.ExecContext(ctx, sqlString, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for data ID %d", data.ID)
	}

	return nil
}

func (storage *DBStorage) DeleteData(ctx context.Context, accountID uint32, storedDataID uint32) error {
	sqlString := `
		DELETE FROM public.pass_keeper_data
		WHERE id = $1 and account_id = $2
	`

	result, err := storage.DB.ExecContext(ctx, sqlString, storedDataID, accountID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for data ID %d", storedDataID)
	}

	return nil
}
