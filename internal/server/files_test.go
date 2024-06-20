package server

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"keeper-project/internal/mocks"
	"keeper-project/types"
)

func Test_router_files_create(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileService := mocks.NewMockFileService(mockCtrl)

	filepath := t.TempDir() + "/test.txt"
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	require.NoError(t, err)
	defer file.Close()

	mockFileService.EXPECT().Create(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", gomock.Any()).Return(nil).Times(1)
	mockFileService.EXPECT().Create(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", gomock.Any()).Return(minio.ToErrorResponse(errors.New("failed to store"))).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, nil, nil, nil, nil, mockFileService))
	defer ts.Close()

	tests := []struct {
		name   string
		method string
		target string
		token  string
		body   map[string]io.Reader
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodPost,
			target: "/api/secret/file",
			token:  validToken,
			body: map[string]io.Reader{
				"file":     file,
				"Metadata": strings.NewReader("some_test_meta"),
			},
			want: want{
				code:          201,
				emptyResponse: false,
				contentType:   "form/json",
			},
		},
		{
			name:   "failed test #1 invalid token",
			method: http.MethodPost,
			target: "/api/secret/file",
			token:  invalidToken,
			body: map[string]io.Reader{
				"file":     file,
				"Metadata": strings.NewReader("some_test_meta"),
			},
			want: want{
				code:          401,
				emptyResponse: false,
				response:      "Unauthorized: invalid token\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #2 invalid request",
			method: http.MethodPost,
			target: "/api/secret/file",
			token:  validToken,
			body:   map[string]io.Reader{},
			want: want{
				code:          400,
				emptyResponse: false,
				response:      "file required\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #4 s3 error",
			method: http.MethodPost,
			target: "/api/secret/file",
			token:  validToken,
			body: map[string]io.Reader{
				"file":     file,
				"Metadata": strings.NewReader("some_test_meta"),
			},
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "unable to store file: Error response code .\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testAuthorizedRequestMultipartForm(t, ts, tt.method, tt.target, tt.token, tt.body)
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.emptyResponse {
				require.Empty(t, body)
			} else {
				assert.Equal(t, tt.want.response, body)
			}

			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_router_file_get(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileService := mocks.NewMockFileService(mockCtrl)

	fileTest := &types.File{
		Name:     "123321",
		Size:     1000,
		Bytes:    []byte{1},
		Metadata: "test_meta",
	}

	mockFileService.EXPECT().GetFile(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(fileTest, nil).Times(1)
	mockFileService.EXPECT().GetFile(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(nil, errors.New("not found")).Times(1)
	mockFileService.EXPECT().GetFile(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(nil, minio.ToErrorResponse(errors.New("failed request"))).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, nil, nil, nil, nil, mockFileService))
	defer ts.Close()

	tests := []struct {
		name   string
		method string
		target string
		token  string
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodGet,
			target: "/api/secret/file/test",
			token:  validToken,
			want: want{
				code:          200,
				emptyResponse: false,
				response:      "\x01",
				contentType:   "",
			},
		},
		{
			name:   "failed test #1 invalid token",
			method: http.MethodGet,
			target: "/api/secret/file/test",
			token:  invalidToken,
			want: want{
				code:          401,
				emptyResponse: false,
				response:      "Unauthorized: invalid token\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #2 invalid id",
			method: http.MethodGet,
			target: "/api/secret/file/",
			token:  validToken,
			want: want{
				code:          404,
				emptyResponse: false,
				response:      "404 page not found\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #3 not found",
			method: http.MethodGet,
			target: "/api/secret/file/test",
			token:  validToken,
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "not found\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #4 sql error",
			method: http.MethodGet,
			target: "/api/secret/file/test",
			token:  validToken,
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "Error response code .\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testAuthorizedRequest(t, ts, tt.method, tt.target, tt.token, nil)
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.emptyResponse {
				require.Empty(t, body)
			} else {
				assert.Equal(t, tt.want.response, body)
			}

			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_router_file_list(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileService := mocks.NewMockFileService(mockCtrl)

	filesList := []*types.Key{{
		Id:  "test",
		Key: "test_key"},
	}

	mockFileService.EXPECT().GetFilesList(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83").Return(filesList, nil).Times(1)
	mockFileService.EXPECT().GetFilesList(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83").Return(nil, nil).Times(1)
	mockFileService.EXPECT().GetFilesList(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83").Return(nil, minio.ToErrorResponse(errors.New("failed request"))).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, nil, nil, nil, nil, mockFileService))
	defer ts.Close()

	tests := []struct {
		name   string
		method string
		target string
		token  string
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodGet,
			target: "/api/secret/files",
			token:  validToken,
			want: want{
				code:          200,
				emptyResponse: false,
				response:      "[{\"id\":\"test\",\"key\":\"test_key\"}]\n",
				contentType:   "application/json",
			},
		},
		{
			name:   "failed test #1 invalid token",
			method: http.MethodGet,
			target: "/api/secret/files",
			token:  invalidToken,
			want: want{
				code:          401,
				emptyResponse: false,
				response:      "Unauthorized: invalid token\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #2 not found",
			method: http.MethodGet,
			target: "/api/secret/files",
			token:  validToken,
			want: want{
				code:          404,
				emptyResponse: false,
				response:      "no files\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #3 sql error",
			method: http.MethodGet,
			target: "/api/secret/files",
			token:  validToken,
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "Error response code .\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testAuthorizedRequest(t, ts, tt.method, tt.target, tt.token, nil)
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.emptyResponse {
				require.Empty(t, body)
			} else {
				assert.Equal(t, tt.want.response, body)
			}

			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_router_file_delete(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileService := mocks.NewMockFileService(mockCtrl)

	mockFileService.EXPECT().Delete(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(nil).Times(1)
	mockFileService.EXPECT().Delete(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(errors.New("deletion failed")).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, nil, nil, nil, nil, mockFileService))
	defer ts.Close()

	tests := []struct {
		name   string
		method string
		target string
		token  string
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodDelete,
			target: "/api/secret/file/test",
			token:  validToken,
			want: want{
				code:          204,
				emptyResponse: true,
			},
		},
		{
			name:   "failed test #1 invalid token",
			method: http.MethodDelete,
			target: "/api/secret/file/test",
			token:  invalidToken,
			want: want{
				code:          401,
				emptyResponse: false,
				response:      "Unauthorized: invalid token\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #2 invalid id",
			method: http.MethodDelete,
			target: "/api/secret/file/",
			token:  validToken,
			want: want{
				code:          404,
				emptyResponse: false,
				response:      "404 page not found\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #3 sql error",
			method: http.MethodDelete,
			target: "/api/secret/file/test",
			token:  validToken,
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "unable to delete: deletion failed\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testAuthorizedRequest(t, ts, tt.method, tt.target, tt.token, nil)
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.emptyResponse {
				require.Empty(t, body)
			} else {
				assert.Equal(t, tt.want.response, body)
			}

			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func testAuthorizedRequestMultipartForm(t *testing.T, ts *httptest.Server,
	method, path, token string, values map[string]io.Reader) (*http.Response, string) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := r.(*os.File); ok {
			if fw, _ = w.CreateFormFile(key, x.Name()); err != nil {
			}
		} else {
			if fw, _ = w.CreateFormField(key); err != nil {
			}
		}
		if _, _ = io.Copy(fw, r); err != nil {
		}
	}
	w.Close()

	req, err := http.NewRequest(method, ts.URL+path, &b)
	require.NoError(t, err)

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
