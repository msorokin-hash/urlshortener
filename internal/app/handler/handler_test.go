package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msorokin-hash/urlshortener/internal/app/config"
	"github.com/msorokin-hash/urlshortener/internal/app/entity"
	"github.com/msorokin-hash/urlshortener/internal/app/storage"
	"github.com/msorokin-hash/urlshortener/internal/app/util"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *Handler {
	store := storage.NewStorage()
	config := config.Config{Address: "localhost:8080", BaseShortURL: "http://localhost:8080"}
	handler := NewHandler(&config, store)
	return handler
}

func TestApp_GetURLHandler(t *testing.T) {
	t.Run("test get url handler", func(t *testing.T) {
		app := setupTestApp()

		url := "https://ya.ru"
		hash := util.HashStringData(url)
		_ = app.Storage.Add(hash, url)

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

func TestApp_AddURLHandler(t *testing.T) {
	t.Run("test add url handler", func(t *testing.T) {
		app := setupTestApp()

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

func TestHandler_AddURLShortenHandler(t *testing.T) {
	t.Run("test add url handler", func(t *testing.T) {
		app := setupTestApp()

		url := entity.Request{URL: "https://yandex.ru"}
		jsonURL, _ := json.Marshal(url)

		req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonURL))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.AddURLShortenHandler)

		handler(rr, req)

		result := rr.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusCreated, result.StatusCode, "Ожидается статус 201")

		var response map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &response)

		assert.NoError(t, err, "Ответ должен быть валидным JSON")
		assert.Contains(t, response, "result", "Должен присутствовать ключ 'result'")
		assert.NotEmpty(t, response["result"], "Result не должен быть пустым")
	})
}
