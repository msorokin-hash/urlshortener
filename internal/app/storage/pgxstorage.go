package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	return &PostgresStorage{conn: conn}, nil
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

	err := ps.conn.Ping(ctx)
	if err != nil {
		return errors.New("не удалось подключиться к базе данных")
	}
	return nil
}

func (ps *PostgresStorage) Close() {
	if ps.conn != nil {
		ps.conn.Close(context.Background())
		log.Println("cоединение с базой данных закрыто")
	}
}
