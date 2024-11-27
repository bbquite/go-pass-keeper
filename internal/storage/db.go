package storage

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/bbquite/go-pass-keeper/internal/models"
	"github.com/bbquite/go-pass-keeper/pkg/xretry"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	DB      *sql.DB
	retrier *xretry.Retrier
}

func NewDBStorage(databaseDSN string) (*DBStorage, error) {
	DB, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}

	retryPolicy := xretry.NewRetryPolicy(
		xretry.WithRetriesWithBackoff(3, 1*time.Second, 1.5),
	)
	retrier := xretry.NewRetrier(retryPolicy)

	store := &DBStorage{
		DB:      DB,
		retrier: retrier,
	}

	err = store.ApplyMigrations(context.Background())
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (storage *DBStorage) ApplyMigrations(ctx context.Context) error {
	err := storage.Ping(ctx)
	if err != nil {
		return err
	}

	path := filepath.Join("migrations", "01_init.sql")
	out, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sqlCreateString := string(out)

	_, err = storage.DB.ExecContext(ctx, sqlCreateString)
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) Ping(ctx context.Context) error {

	retryFunction := func() error { return storage.DB.PingContext(ctx) }

	err := storage.retrier.Retry(retryFunction)
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) GetAccountByUsername(ctx context.Context, username string) (models.Account, error) {
	var account models.Account

	sqlString := `
		SELECT id 
		FROM account 
		WHERE username = $1 
		LIMIT 1
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
		FROM account 
		WHERE username = $1 AND password = $2
		LIMIT 1
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
		INSERT INTO account (username, password, email) 
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
