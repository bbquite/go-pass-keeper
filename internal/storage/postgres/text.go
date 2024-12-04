package postgres

import (
	"context"
	"github.com/bbquite/go-pass-keeper/internal/utils"
	"time"

	"github.com/bbquite/go-pass-keeper/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (storage *DBStorage) CreateTextData(ctx context.Context, data *models.TextData) (models.TextData, error) {
	sqlString := `
		INSERT INTO public.simple_data (text_data, meta, account_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, text_data, meta, uploaded_at
	`

	var resultData models.TextData
	accountID := ctx.Value(utils.UserIDKey).(uint32)

	row := storage.DB.QueryRowContext(ctx, sqlString, data.Text, data.Meta, accountID)
	err := row.Scan(
		&resultData.ID,
		&resultData.Text,
		&resultData.Meta,
		&resultData.UploadedAt,
	)
	if err != nil {
		return resultData, err
	}

	return resultData, nil
}

func (storage *DBStorage) GetTextDataList(ctx context.Context) ([]models.TextData, error) {
	sqlStringSelect := `
		SELECT id, text_data, meta, uploaded_at
		FROM public.simple_data
		WHERE account_id = $1;
	`

	var result []models.TextData
	accountID := ctx.Value(utils.UserIDKey).(uint32)

	rows, err := storage.DB.QueryContext(ctx, sqlStringSelect, accountID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	for rows.Next() {
		var textItem models.TextData

		err := rows.Scan(&textItem.ID, &textItem.Text, &textItem.Meta, &textItem.UploadedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, textItem)
	}

	return result, nil
}

func (storage *DBStorage) UpdateTextData(ctx context.Context, data *models.TextData) error {
	sqlString := `
		UPDATE public.simple_data 
		SET text_data = $1, meta = $2, uploaded_at = $3
		WHERE id = $4
	`

	_, err := storage.DB.ExecContext(ctx, sqlString, data.Text, data.Meta, time.Now(), data.ID)
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) DeleteTextData(ctx context.Context, textDataID uint32) error {
	sqlString := `
		DELETE FROM public.simple_data
		WHERE id = $1 and account_id = $2
	`

	accountID := ctx.Value(utils.UserIDKey).(uint32)

	_, err := storage.DB.ExecContext(ctx, sqlString, textDataID, accountID)
	if err != nil {
		return err
	}

	return nil
}
