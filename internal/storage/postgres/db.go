package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

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
