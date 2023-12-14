package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mainHandle(t *testing.T) {
	testCases := []struct {
		name              string
		path              string
		method            string
		expectedCafeCount int
		expectedCode      int
		expectedBody      string
	}{
		{
			name:              "Successful request: count=3, city=moscow",
			path:              "/cafe?count=3&city=moscow",
			method:            http.MethodGet,
			expectedCafeCount: 3,
			expectedCode:      http.StatusOK,
			expectedBody:      "Мир кофе,Сладкоежка,Кофе и завтраки",
		},
		{
			name:              "Successful request: count=1, city=moscow",
			path:              "/cafe?count=1&city=moscow",
			method:            http.MethodGet,
			expectedCafeCount: 1,
			expectedCode:      http.StatusOK,
			expectedBody:      "Мир кофе",
		},
		{
			name:              "Successful request: cafe count more then exists, returns all",
			path:              "/cafe?count=99&city=moscow",
			method:            http.MethodGet,
			expectedCafeCount: 4,
			expectedCode:      http.StatusOK,
			expectedBody:      "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент",
		},
		{
			name:         "Bad request: count negative",
			path:         "/cafe?count=-1&city=moscow",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			expectedBody: "wrong count value",
		},
		{
			name:         "Bad request: count = 0",
			path:         "/cafe?count=0&city=moscow",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			expectedBody: "wrong count value",
		},
		{
			name:         "Bad request: count missing",
			path:         "/cafe?count=&city=moscow",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			expectedBody: "count missing",
		},
		{
			name:         "Bad request: invalid city (wrong city value)",
			path:         "/cafe?count=4&city=ansbmdanbs",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			expectedBody: "wrong city value",
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(mainHandle)
			handler.ServeHTTP(w, request)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			if test.expectedCode == http.StatusOK {
				assert.Equal(t, test.expectedCafeCount, len(strings.Split(w.Body.String(), ",")))
			}
		})
	}
}
