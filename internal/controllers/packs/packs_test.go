package packs

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/assert"

	"github.com/pkg/errors"

	"github.com/JackFazackerley/complete-packs/internal/database/sqlite"

	"github.com/JackFazackerley/complete-packs/pkg/cache"
	"github.com/gin-gonic/gin"
)

const (
	tableDDL = `
	CREATE TABLE "packs" (
   		size INTEGER NOT NULL UNIQUE
	);`

	insert = `INSERT INTO packs ("size") VALUES (250),(500),(1000),(2000),(5000);`
)

func dbInit(create bool, t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("creating db: %s", err)
		return nil
	}

	if create {
		if _, err := db.Exec(tableDDL); err != nil {
			t.Errorf("creating table: %s", err)
			return nil
		}

		if _, err := db.Exec(insert); err != nil {
			t.Errorf("inserting into table: %s", err)
			return nil
		}
	}

	return db
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

func TestController_Read(t *testing.T) {
	tests := []struct {
		name             string
		sizes            []float64
		expectedResponse []byte
		expectedCode     int
		dbCreate         bool
	}{
		{
			name:             "returns as expected",
			sizes:            []float64{250, 500, 1000, 2000, 5000},
			expectedResponse: []byte(`[250,500,1000,2000,5000]`),
			expectedCode:     http.StatusOK,
			dbCreate:         true,
		},
		{
			name:             "returns error",
			sizes:            []float64{250, 500, 1000, 2000, 5000},
			expectedResponse: []byte(`{"error":"reading packs"}`),
			expectedCode:     http.StatusInternalServerError,
			dbCreate:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := dbInit(tt.dbCreate, t)
			dbInterface := sqlite.New(db)
			r := &Controller{
				db:         sqlite.New(db),
				sizesCache: cache.New(dbInterface),
			}

			_ = r.sizesCache.Load() //ignore this error so we can see errors from read

			recorder, err := httpTest(t, nil, http.MethodGet, "/read", r.Read)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedResponse, recorder.Body.Bytes())
		})
	}
}

func TestController_Write(t *testing.T) {
	tests := []struct {
		name         string
		sizes        []float64
		body         []byte
		dbCreate     bool
		expectedCode int
		expectedBody []byte
	}{
		{
			name:         "returns as expected",
			sizes:        []float64{250, 500, 1000, 2000, 5000},
			body:         []byte(`{"size": 10}`),
			dbCreate:     true,
			expectedCode: http.StatusOK,
			expectedBody: nil,
		},
		{
			name:         "incorrect request body",
			sizes:        []float64{250, 500, 1000, 2000, 5000},
			body:         []byte(`[{}]`),
			dbCreate:     true,
			expectedCode: http.StatusBadRequest,
			expectedBody: []byte(`{"error":"parsing request body"}`),
		},
		{
			name:         "size exists",
			sizes:        []float64{250, 500, 1000, 2000, 5000},
			body:         []byte(`{"size": 250}`),
			dbCreate:     true,
			expectedCode: http.StatusInternalServerError,
			expectedBody: []byte(`{"error":"duplicate key"}`),
		},
		{
			name:         "unable to write to table",
			sizes:        []float64{250, 500, 1000, 2000, 5000},
			body:         []byte(`{"size": 10}`),
			dbCreate:     false,
			expectedCode: http.StatusInternalServerError,
			expectedBody: []byte(`{"error":"writing pack"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := dbInit(tt.dbCreate, t)
			dbInterface := sqlite.New(db)
			r := &Controller{
				db:         sqlite.New(db),
				sizesCache: cache.New(dbInterface),
			}

			_ = r.sizesCache.Load() //ignore this error so we can see errors from read

			recorder, err := httpTest(t, tt.body, http.MethodPost, "/write", r.Write)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, recorder.Body.Bytes())
		})
	}
}

func TestController_Delete(t *testing.T) {
	tests := []struct {
		name         string
		sizes        []float64
		body         []byte
		dbCreate     bool
		expectedCode int
	}{
		{
			name:         "returns as expected",
			sizes:        []float64{250, 500, 1000, 2000, 5000},
			body:         []byte(`{"size": 250}`),
			dbCreate:     true,
			expectedCode: http.StatusOK,
		},
		{
			name:         "incorrect request body",
			sizes:        []float64{250, 500, 1000, 2000, 5000},
			body:         []byte(`[{}]`),
			dbCreate:     true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "unable to delete from table",
			sizes:        []float64{250, 500, 1000, 2000, 5000},
			body:         []byte(`{"size": 10}`),
			dbCreate:     false,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := dbInit(tt.dbCreate, t)
			dbInterface := sqlite.New(db)
			r := &Controller{
				db:         sqlite.New(db),
				sizesCache: cache.New(dbInterface),
			}

			_ = r.sizesCache.Load() //ignore this error so we can see errors from read

			recorder, err := httpTest(t, tt.body, http.MethodDelete, "/delete", r.Delete)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, recorder.Code)
		})
	}
}
