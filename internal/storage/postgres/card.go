package postgres

import (
	"context"
	"github.com/bbquite/go-pass-keeper/internal/utils"
	"time"

	"github.com/bbquite/go-pass-keeper/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (storage *DBStorage) CreateCardData(ctx context.Context, data *models.CardData) (models.CardData, error) {
	sqlString := `
		INSERT INTO public.card_data (card_num, card_owner, card_exp, card_cvv, meta, account_id) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, card_num, card_owner, card_exp, card_cvv, meta, uploaded_at
	`

	var resultData models.CardData
	accountID := ctx.Value(utils.UserIDKey).(uint32)

	row := storage.DB.QueryRowContext(
		ctx, sqlString, data.CardNum, data.CardOwner, data.CardExp, data.CardCvv, data.Meta, accountID)

	err := row.Scan(
		&resultData.ID,
		&resultData.CardNum,
		&resultData.CardOwner,
		&resultData.CardExp,
		&resultData.CardCvv,
		&resultData.Meta,
		&resultData.UploadedAt,
	)
	if err != nil {
		return resultData, err
	}

	return resultData, nil
}

func (storage *DBStorage) GetCardDataList(ctx context.Context) ([]models.CardData, error) {
	sqlStringSelect := `
		SELECT id, card_num, card_owner, card_exp, card_cvv, meta, uploaded_at
		FROM public.card_data
		WHERE account_id = $1;
	`

	var result []models.CardData
	accountID := ctx.Value(utils.UserIDKey).(uint32)

	rows, err := storage.DB.QueryContext(ctx, sqlStringSelect, accountID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	for rows.Next() {
		var cardItem models.CardData

		err := rows.Scan(
			&cardItem.ID,
			&cardItem.CardNum,
			&cardItem.CardOwner,
			&cardItem.CardExp,
			&cardItem.CardCvv,
			&cardItem.Meta,
			&cardItem.UploadedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, cardItem)
	}

	return result, nil
}

func (storage *DBStorage) UpdateCardData(ctx context.Context, data *models.CardData) error {
	sqlString := `
		UPDATE public.card_data 
		SET card_num = $1, card_owner = $2, card_exp = $3, card_cvv = $4, meta = $5, uploaded_at = $6
		WHERE id = $7
	`

	_, err := storage.DB.ExecContext(
		ctx, sqlString, data.CardNum, data.CardOwner, data.CardExp, data.CardCvv, data.Meta, time.Now(), data.ID)
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) DeleteCardData(ctx context.Context, cardDataID uint32) error {
	sqlString := `
		DELETE FROM public.card_data
		WHERE id = $1 and account_id = $2
	`

	accountID := ctx.Value(utils.UserIDKey).(uint32)

	_, err := storage.DB.ExecContext(ctx, sqlString, cardDataID, accountID)
	if err != nil {
		return err
	}

	return nil
}
