package postgres

import (
	"context"
	"errors"
	"github.com/bbquite/go-pass-keeper/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (storage *DBStorage) GetAccountByUsername(ctx context.Context, username string) (models.Account, error) {
	var account models.Account

	sqlString := `
		SELECT id 
		FROM public.account 
		WHERE username = $1 
	`

	row := storage.DB.QueryRowContext(ctx, sqlString, username)
	err := row.Scan(&account.ID)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (storage *DBStorage) GetAccountByLoginData(ctx context.Context, username string, password string) (models.Account, error) {
	var account models.Account

	sqlString := `
		SELECT id 
		FROM public.account 
		WHERE username = $1 AND password = $2
	`

	row := storage.DB.QueryRowContext(ctx, sqlString, username, password)
	err := row.Scan(&account.ID)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (storage *DBStorage) CreateAccount(ctx context.Context, username string, password string, email string) (uint32, error) {
	sqlString := `
		INSERT INTO public.account (username, password, email) 
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var userID uint32
	row := storage.DB.QueryRowContext(ctx, sqlString, username, password, email)
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}

	if userID == 0 {
		return 0, errors.New("unspecified error while creating record")
	}

	return userID, nil
}
