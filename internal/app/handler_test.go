package app

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempDB() (*sql.DB, error) {
	dbConn, err := InitDB()
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

		url := "https://ya.ru"
		hash := HashStringData(url)
		_ = CreateURL(dbConn, hash, url)

		req := httptest.NewRequest("GET", "/"+hash, nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetURLHandler)

		handler(rr, req)

		result := rr.Result()

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

		req := httptest.NewRequest("POST", "/", bytes.NewBufferString("https://yandex.ru"))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(AddURLHandler)

		handler(rr, req)

		result := rr.Result()

		assert.Equal(t, http.StatusCreated, result.StatusCode, "Ожидается статус 201")
		assert.NotNil(t, result.Body, "Ответ не может быть пустым")
	})
}
