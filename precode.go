package main

import (
	"net/http"
	"strconv"
	"strings"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}


package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const WRONH_CITY_ERROR = "wrong city value"

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=9&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	statusCode := responseRecorder.Code
	require.NotEqual(t, 0, statusCode)
	require.Equal(t, http.StatusOK, statusCode)

	responseBody := strings.Split(responseRecorder.Body.String(), ",")
	assert.Equal(t, totalCount, len(responseBody))
}

func TestMainHandlerWhenOK(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	statusCode := responseRecorder.Code

	require.Equal(t, http.StatusOK, statusCode)
	require.NotEmpty(t, responseRecorder.Body.String())
}

func TestMainHandlerWhenMissingCount(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=rostov", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	statusCode := responseRecorder.Code

	require.Equal(t, http.StatusBadRequest, statusCode)
	require.Equal(t, WRONH_CITY_ERROR, responseRecorder.Body.String())
}