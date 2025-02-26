package handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/msorokin-hash/urlshortener/internal/app/config"
	"github.com/msorokin-hash/urlshortener/internal/app/storage"
	"github.com/msorokin-hash/urlshortener/internal/app/util"
	"github.com/stretchr/testify/assert"
)

func createTempDB() (*sql.DB, error) {
	dbConn, err := storage.InitDB()
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func deleteTempDB(dbFile string) {
	err := os.Remove(dbFile)
	if err != nil {
		fmt.Println("Ошибка удаления файла базы данных:", err)
	}
}

func TestGetURLHandler(t *testing.T) {
	t.Run("test get handler", func(t *testing.T) {
		dbConn, err := createTempDB()
		assert.NoError(t, err, "Ошибка при создании временной базы данных")
		defer dbConn.Close()
		defer deleteTempDB("./urls.db")

		app := &App{
			Config: &config.Config{
				Address:      "localhost:8080",
				BaseShortURL: "http://localhost:8080",
			},
			DB: dbConn,
		}
		url := "https://ya.ru"
		hash := util.HashStringData(url)
		_ = storage.CreateURL(dbConn, hash, url)

		req := httptest.NewRequest("GET", "/"+hash, nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.GetURLHandler)

		handler(rr, req)

		result := rr.Result()
		defer result.Body.Close()

		assert.NotEmpty(t, result.Header.Get("Location"), "Заголовок Location должен присутствовать")
		assert.Equal(t, url, result.Header.Get("Location"), "Значение заголовка Location=%v", url)
		assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode, "Ожидается статус 307")
	})
}

func TestAddURLHandler(t *testing.T) {
	t.Run("test add handler", func(t *testing.T) {
		dbConn, err := createTempDB()
		assert.NoError(t, err, "Ошибка при создании временной базы данных")
		defer dbConn.Close()
		defer deleteTempDB("./urls.db")

		app := &App{
			Config: &config.Config{
				Address:      "localhost:8080",
				BaseShortURL: "http://localhost:8080",
			},
			DB: dbConn,
		}

		req := httptest.NewRequest("POST", "/", bytes.NewBufferString("https://yandex.ru"))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.AddURLHandler)

		handler(rr, req)

		result := rr.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusCreated, result.StatusCode, "Ожидается статус 201")
		assert.NotNil(t, result.Body, "Ответ не может быть пустым")
	})
}
