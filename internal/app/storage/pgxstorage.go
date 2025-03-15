package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (ps *PostgresStorage) Insert(shortURL, originalURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO shortened_urls (uuid, short_url, original_url, created_at) 
	          VALUES (@uuid, @shortURL, @originalURL, NOW())`
	args := pgx.NamedArgs{
		"uuid":        uuid.New().String(),
		"shortURL":    shortURL,
		"originalURL": originalURL,
	}
	_, err := ps.pool.Exec(ctx, query, args)
	if err != nil {
		return errors.New("ошибка при добавлении записи в базу данных")
	}

	log.Printf("URL добавлен: %s -> %s", shortURL, originalURL)
	return nil
}

func (ps *PostgresStorage) Get(shortURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT original_url FROM shortened_urls WHERE short_url = @shortURL`
	args := pgx.NamedArgs{"shortURL": shortURL}
	row := ps.pool.QueryRow(ctx, query, args)

	var originalURL string
	err := row.Scan(&originalURL)
	if err != nil {
		return "", err
	}

	log.Printf("URL получен из БД: %s", originalURL)
	return originalURL, nil
}

func (ps *PostgresStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
