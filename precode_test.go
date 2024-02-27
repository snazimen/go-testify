package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getCountAndCity(req *http.Request) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	return responseRecorder
}

// m
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	count := 5
	city := "moscow"
	responseRecorder := getCountAndCity(httptest.NewRequest(http.MethodGet, fmt.Sprintf("/cafe?count=%d&city=%s", count, city), nil))
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, []byte(strings.Join(cafeList["moscow"][:totalCount], ",")), responseRecorder.Body.Bytes())

}

func TestMainHandlerWhenOk(t *testing.T) {
	count := 2
	city := "moscow"
	responseRecorder := getCountAndCity(httptest.NewRequest(http.MethodGet, fmt.Sprintf("/cafe?count=%d&city=%s", count, city), nil))
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, []byte(strings.Join(cafeList["moscow"][:count], ",")), responseRecorder.Body.Bytes())
}

func TestWhenWrongCity(t *testing.T) {
	count := 2
	city := "london"
	responseRecorder := getCountAndCity(httptest.NewRequest(http.MethodGet, fmt.Sprintf("/cafe?count=%d&city=%s", count, city), nil))
	assert.Equal(t, responseRecorder.Code, http.StatusBadRequest)
	assert.Equal(t, []byte("wrong city value"), responseRecorder.Body.Bytes())
}
