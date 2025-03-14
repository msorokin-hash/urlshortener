package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ошибка проверки соединения с базой данных: %w", err)
	}

	return &PostgresStorage{pool: pool}, nil
}

func (ps *PostgresStorage) Add(shortURL, originalURL string) error {
	return nil
}

func (ps *PostgresStorage) Lookup(shortURL string) (string, error) {
	return "", nil
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
	if ps.pool != nil {
		ps.pool.Close()
		log.Println("пул соединений с базой данных закрыт")
	}
}
