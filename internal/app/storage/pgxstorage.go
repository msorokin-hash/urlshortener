package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pgInstance *PostgresStorage
	pgOnce     sync.Once
	mu         sync.Mutex
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	var err error

	pgOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pool, connErr := pgxpool.New(ctx, dsn)
		if connErr != nil {
			err = fmt.Errorf("ошибка подключения к базе данных: %w", connErr)
			return
		}

		if pingErr := pool.Ping(ctx); pingErr != nil {
			pool.Close()
			err = fmt.Errorf("ошибка проверки соединения с базой данных: %w", pingErr)
			return
		}

		if tableErr := createTables(ctx, pool); tableErr != nil {
			pool.Close()
			err = fmt.Errorf("ошибка при создании таблиц: %w", tableErr)
			return
		}

		pgInstance = &PostgresStorage{pool: pool}
	})

	if err != nil {
		return nil, err
	}

	return pgInstance, nil
}

func (ps *PostgresStorage) Add(shortURL, originalURL string) error {
	return nil
}

func (ps *PostgresStorage) Lookup(shortURL string) (string, error) {
	return "", nil
}

func (ps *PostgresStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := ps.pool.Ping(ctx)
	if err != nil {
		return errors.New("не удалось подключиться к базе данных")
	}

	return nil
}

func (ps *PostgresStorage) Close() {
	mu.Lock()
	defer mu.Unlock()

	if pgInstance != nil && pgInstance.pool != nil {
		pgInstance.pool.Close()
		pgInstance = nil
		log.Println("пул соединений с базой данных закрыт")
	}
}

func createTables(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS shortened_urls (
		id SERIAL PRIMARY KEY,
		uuid TEXT UNIQUE NOT NULL,
		short_url TEXT UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы shortened_urls: %w", err)
	}

	return nil
}
