package storage

import (
	"context"
	"errors"
	"github.com/bbquite/go-pass-keeper/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (storage *DBStorage) CreatePairsData(ctx context.Context, data models.PairsData) (uint32, error) {
	sqlString := `
		INSERT INTO public.pairs_data (key, pwd, meta) 
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var pairsID uint32
	row := storage.DB.QueryRowContext(ctx, sqlString, data.Key, data.Pwd, data.Meta)
	err := row.Scan(&pairsID)
	if err != nil {
		return 0, err
	}

	if pairsID == 0 {
		return 0, errors.New("unspecified error while creating record")
	}

	return pairsID, nil
}

func (storage *DBStorage) GetPairsData(ctx context.Context) ([]models.PairsData, error) {
	var result []models.PairsData

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
		var pairItem models.PairsData

		err := rows.Scan(&pairItem.ID, &pairItem.Key, &pairItem.Pwd, &pairItem.Meta)
		if err != nil {
			return nil, err
		}

		result = append(result, pairItem)
	}

	return result, nil
}

func (storage *DBStorage) UpdatePairsData(ctx context.Context, data models.PairsData) error {
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

func (storage *DBStorage) DeletePairsData(ctx context.Context, id uint32) error {
	sqlString := `
		DELETE FROM public.pairs_data
		WHERE id = $1
	`

	_, err := storage.DB.ExecContext(ctx, sqlString, id)
	if err != nil {
		return err
	}

	return nil
}
