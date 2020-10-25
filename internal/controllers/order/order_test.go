package order

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/JackFazackerley/complete-packs/internal/database/sqlite"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/JackFazackerley/complete-packs/pkg/cache"
)

type TestDatbase struct {
	shouldError bool
	sizes       []float64
}

func (t *TestDatbase) ReadPacks() (packs []float64, err error) {
	if t.shouldError {
		return nil, sqlite.ReadError
	}
	return t.sizes, nil
}

func (t *TestDatbase) WritePack(size int) error {
	return nil
}

func (t *TestDatbase) DeletePack(size int) error {
	return nil
}

func (t *TestDatbase) Close() error {
	return nil
}

func httpTest(t *testing.T, body []byte, method, path string, handleFunc gin.HandlerFunc) (*httptest.ResponseRecorder, error) {
	t.Helper()

	recorder := httptest.NewRecorder()

	router := gin.Default()
	router.Handle(method, path, handleFunc)

	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest(method, path, buf)
	if err != nil {
		return nil, errors.Wrap(err, "error making request")
	}

	router.ServeHTTP(recorder, req)

	return recorder, nil
}

func TestController_Best(t *testing.T) {
	tests := []struct {
		name             string
		sizeCache        cache.Cache
		body             []byte
		expectedResponse []byte
		expectedCode     int
	}{
		{
			name: "returns as expected",
			sizeCache: cache.New(&TestDatbase{
				shouldError: false,
				sizes:       []float64{250, 500, 1000, 2000, 5000},
			}),
			body:             []byte(`{"target": 10}`),
			expectedResponse: []byte(`[{"count":1,"size":250}]`),
			expectedCode:     http.StatusOK,
		},
		{
			name: "returns error",
			sizeCache: cache.New(&TestDatbase{
				shouldError: false,
				sizes:       []float64{250, 500, 1000, 2000, 5000},
			}),
			body:             []byte(`[{}]`),
			expectedResponse: []byte(`{"error":"parsing request body"}`),
			expectedCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Controller{
				sizeCache: tt.sizeCache,
			}

			if err := tt.sizeCache.Load(); err != nil {
				t.Fatalf("loading cache: %s", err)
			}

			recorder, err := httpTest(t, tt.body, http.MethodGet, "/best", o.Best)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedResponse, recorder.Body.Bytes())
		})
	}
}

func TestController_Fast(t *testing.T) {
	tests := []struct {
		name             string
		sizeCache        cache.Cache
		body             []byte
		expectedResponse []byte
		expectedCode     int
	}{
		{
			name: "returns as expected",
			sizeCache: cache.New(&TestDatbase{
				shouldError: false,
				sizes:       []float64{250, 500, 1000, 2000, 5000},
			}),
			body:             []byte(`{"target": 10}`),
			expectedResponse: []byte(`[{"count":1,"size":250}]`),
			expectedCode:     http.StatusOK,
		},
		{
			name: "returns error",
			sizeCache: cache.New(&TestDatbase{
				shouldError: false,
				sizes:       []float64{250, 500, 1000, 2000, 5000},
			}),
			body:             []byte(`[{}]`),
			expectedResponse: []byte(`{"error":"parsing request body"}`),
			expectedCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Controller{
				sizeCache: tt.sizeCache,
			}

			if err := tt.sizeCache.Load(); err != nil {
				t.Fatalf("loading cache: %s", err)
			}

			recorder, err := httpTest(t, tt.body, http.MethodGet, "/fast", o.Fast)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedResponse, recorder.Body.Bytes())
		})
	}
}
