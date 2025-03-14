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
