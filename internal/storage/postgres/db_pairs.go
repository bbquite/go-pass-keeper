package postgres

import (
	"context"
	"github.com/bbquite/go-pass-keeper/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (storage *DBStorage) CreatePairData(ctx context.Context, data *models.PairData) (models.PairData, error) {
	sqlString := `
		INSERT INTO public.pairs_data (key, pwd, meta) 
		VALUES ($1, $2, $3)
		RETURNING id, key, pwd, meta, uploaded_at
	`

	var resultPairData models.PairData
	row := storage.DB.QueryRowContext(ctx, sqlString, data.Key, data.Pwd, data.Meta)
	err := row.Scan(
		&resultPairData.ID,
		&resultPairData.Key,
		&resultPairData.Pwd,
		&resultPairData.Meta,
		resultPairData.UploadedAt,
	)
	if err != nil {
		return resultPairData, err
	}

	return resultPairData, nil
}

func (storage *DBStorage) GetPairsDataList(ctx context.Context) ([]models.PairData, error) {
	var result []models.PairData

	sqlStringSelect := `
		SELECT id, key, pwd, meta
		FROM public.pairs_data;
	`

	rows, err := storage.DB.QueryContext(ctx, sqlStringSelect)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	for rows.Next() {
		var pairItem models.PairData

		err := rows.Scan(&pairItem.ID, &pairItem.Key, &pairItem.Pwd, &pairItem.Meta, &pairItem.UploadedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, pairItem)
	}

	return result, nil
}

func (storage *DBStorage) UpdatePairData(ctx context.Context, data *models.PairData) error {
	sqlString := `
		UPDATE public.pairs_data 
		SET key=$1, pwd = $2, meta = $3
		WHERE id = $4
	`

	_, err := storage.DB.ExecContext(ctx, sqlString, data.Key, data.Pwd, data.Meta)
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) DeletePairData(ctx context.Context, pairID uint32) error {
	sqlString := `
		DELETE FROM public.pairs_data
		WHERE id = $1
	`

	_, err := storage.DB.ExecContext(ctx, sqlString, pairID)
	if err != nil {
		return err
	}

	return nil
}
