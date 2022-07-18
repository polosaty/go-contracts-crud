package handlers

import (
	"github.com/pashagolub/pgxmock"
	"github.com/polosaty/go-contracts-crud/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t require.TestingT, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestCreateCompany(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	pgRepository := storage.NewStorageFromPool(mock)

	type want struct {
		body        string
		contentType string
		statusCode  int
	}
	type request struct {
		body   string
		path   string
		method string
	}

	tests := []struct {
		name string
		want
		request
		mock func()
	}{
		{
			name: "Test #1 Create company",
			request: request{
				body:   `{"name": "some company 1", "code": "some_company_1_code"}`,
				method: "POST",
				path:   "/api/company",
			},
			want: want{
				body:        `{"id": 123, "name": "some company 1", "code": "some_company_1_code"}`,
				statusCode:  http.StatusOK,
				contentType: "application/json; charset=utf-8",
			},
			mock: func() {
				name := "some company 1"
				code := "some_company_1_code"
				mock.ExpectQuery(`INSERT INTO "company" \(name, code\) VALUES\(\$1, \$2\) `+
					` RETURNING id`).
					WithArgs(name, code).
					WillReturnRows(
						mock.NewRows([]string{"id"}).
							AddRow(int64(123)))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			h := NewMainHandler(pgRepository)
			ts := httptest.NewServer(h)
			defer ts.Close()
			body := strings.NewReader(tt.request.body)
			resp, respBody := testRequest(t, ts, tt.method, tt.path, body)
			resp.Body.Close()
			// проверяем код ответа
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			// проверяем код заголовок
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			// получаем и проверяем тело запроса
			switch resp.Header.Get("Content-Type") {
			case "application/json; charset=utf-8":
				assert.JSONEq(t, tt.want.body, respBody,
					"Expected body [%s], got [%s]", tt.want.body, respBody)
			default:
				assert.Equal(t, tt.want.body, respBody,
					"Expected body [%s], got [%s]", tt.want.body, respBody)
			}
		})
	}
}

func BenchmarkCreateCompany(b *testing.B) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	pgRepository := storage.NewStorageFromPool(mock)

	type request struct {
		body   string
		path   string
		method string
	}

	tests := []struct {
		name string
		request
		mock func()
	}{
		{
			name: "Benchmark #1 Create company",
			request: request{
				body:   `{"name": "some company 1", "code": "some_company_1_code"}`,
				method: "POST",
				path:   "/api/company",
			},
			mock: func() {
				name := "some company 1"
				code := "some_company_1_code"
				mock.ExpectQuery(`INSERT INTO "company" \(name, code\) VALUES\(\$1, \$2\) `+
					` RETURNING id`).
					WithArgs(name, code).
					WillReturnRows(
						mock.NewRows([]string{"id"}).
							AddRow(int64(123)))
			},
		},
	}
	h := NewMainHandler(pgRepository)
	ts := httptest.NewServer(h)
	defer ts.Close()

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tt.mock()
				body := strings.NewReader(tt.request.body)
				resp, _ := testRequest(b, ts, tt.method, tt.path, body)
				resp.Body.Close()
			}

		})
	}
}
