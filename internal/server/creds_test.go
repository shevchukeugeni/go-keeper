package server

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"keeper-project/internal/mocks"
	"keeper-project/types"
)

func Test_router_creds_create(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mocksSecret := mocks.NewMockCredsSecret(mockCtrl)

	creds := &types.Credentials{
		Site:     "123321",
		Login:    "test",
		Password: "123",
		Metadata: "test_meta",
	}

	mocksSecret.EXPECT().Create(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", gomock.Any(), creds).Return(nil).Times(1)
	mocksSecret.EXPECT().Create(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", gomock.Any(), creds).Return(sql.ErrConnDone).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, nil, nil, mocksSecret, nil, nil))
	defer ts.Close()

	tests := []struct {
		name   string
		method string
		target string
		token  string
		body   []byte
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodPost,
			target: "/api/secret/cred",
			token:  validToken,
			body:   []byte(`{"site":"123321","login":"test","password":"123","metadata":"test_meta"}`),
			want: want{
				code:          202,
				emptyResponse: true,
			},
		},
		{
			name:   "failed test #1 invalid token",
			method: http.MethodPost,
			target: "/api/secret/cred",
			token:  invalidToken,
			body:   []byte(`{"site":"123321","login":"test","password":"123","metadata":"test_meta"}`),
			want: want{
				code:          401,
				emptyResponse: false,
				response:      "Unauthorized: invalid token\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #2 invalid json",
			method: http.MethodPost,
			target: "/api/secret/cred",
			token:  validToken,
			body:   []byte(`{"`),
			want: want{
				code:          400,
				emptyResponse: false,
				response:      "Unable to decode json: unexpected EOF\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #3 invalid data",
			method: http.MethodPost,
			target: "/api/secret/cred",
			token:  validToken,
			body:   []byte(`{}`),
			want: want{
				code:          400,
				emptyResponse: false,
				response:      "incorrect data: incorrect key\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #4 sql error",
			method: http.MethodPost,
			target: "/api/secret/cred",
			token:  validToken,
			body:   []byte(`{"site":"123321","login":"test","password":"123","metadata":"test_meta"}`),
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "failed to create: sql: connection is already closed\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testAuthorizedRequest(t, ts, tt.method, tt.target, tt.token, tt.body)
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

func Test_router_creds_get(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mocksSecret := mocks.NewMockCredsSecret(mockCtrl)

	creds := &types.Credentials{
		Site:     "123321",
		Login:    "test",
		Password: "123",
		Metadata: "test_meta",
	}

	mocksSecret.EXPECT().Get(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(creds, nil).Times(1)
	mocksSecret.EXPECT().Get(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(nil, sql.ErrNoRows).Times(1)
	mocksSecret.EXPECT().Get(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(nil, sql.ErrConnDone).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, nil, nil, mocksSecret, nil, nil))
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
			target: "/api/secret/cred/test",
			token:  validToken,
			want: want{
				code:          200,
				emptyResponse: false,
				response:      "{\"site\":\"123321\",\"login\":\"test\",\"password\":\"123\",\"metadata\":\"test_meta\"}\n",
				contentType:   "application/json",
			},
		},
		{
			name:   "failed test #1 invalid token",
			method: http.MethodGet,
			target: "/api/secret/cred/test",
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
			target: "/api/secret/cred/",
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
			target: "/api/secret/cred/test",
			token:  validToken,
			want: want{
				code:          404,
				emptyResponse: false,
				response:      "no such record\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #4 sql error",
			method: http.MethodGet,
			target: "/api/secret/cred/test",
			token:  validToken,
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "failed to get from db: sql: connection is already closed\n",
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

func Test_router_creds_list(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mocksSecret := mocks.NewMockCredsSecret(mockCtrl)

	credsList := []types.Key{{
		Id:  "test",
		Key: "test_key"},
	}

	mocksSecret.EXPECT().GetKeysList(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83").Return(credsList, nil).Times(1)
	mocksSecret.EXPECT().GetKeysList(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83").Return(nil, sql.ErrNoRows).Times(1)
	mocksSecret.EXPECT().GetKeysList(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83").Return(nil, sql.ErrConnDone).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, nil, nil, mocksSecret, nil, nil))
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
			target: "/api/secret/creds",
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
			target: "/api/secret/creds",
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
			target: "/api/secret/creds",
			token:  validToken,
			want: want{
				code:          404,
				emptyResponse: false,
				response:      "no credentials info\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #3 sql error",
			method: http.MethodGet,
			target: "/api/secret/creds",
			token:  validToken,
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "failed to get from db: sql: connection is already closed\n",
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

func Test_router_cred_update(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mocksSecret := mocks.NewMockCredsSecret(mockCtrl)

	creds := &types.Credentials{
		Site:     "123321",
		Login:    "test",
		Password: "123",
		Metadata: "test_meta",
	}

	mocksSecret.EXPECT().Update(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", gomock.Any(), creds).Return(nil).Times(1)
	mocksSecret.EXPECT().Update(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", gomock.Any(), creds).Return(sql.ErrNoRows).Times(1)
	mocksSecret.EXPECT().Update(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", gomock.Any(), creds).Return(sql.ErrConnDone).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, nil, nil, mocksSecret, nil, nil))
	defer ts.Close()

	tests := []struct {
		name   string
		method string
		target string
		token  string
		body   []byte
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodPut,
			target: "/api/secret/cred",
			token:  validToken,
			body:   []byte(`{"id":"test","site":"123321","login":"test","password":"123","metadata":"test_meta"}`),
			want: want{
				code:          200,
				emptyResponse: true,
			},
		},
		{
			name:   "failed test #1 invalid token",
			method: http.MethodPut,
			target: "/api/secret/cred",
			token:  invalidToken,
			body:   []byte(`{"id":"test","site":"123321","login":"test","password":"123","metadata":"test_meta"}`),
			want: want{
				code:          401,
				emptyResponse: false,
				response:      "Unauthorized: invalid token\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #2 invalid json",
			method: http.MethodPut,
			target: "/api/secret/cred",
			token:  validToken,
			body:   []byte(`{"`),
			want: want{
				code:          400,
				emptyResponse: false,
				response:      "Unable to decode json: unexpected EOF\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #3 invalid data",
			method: http.MethodPut,
			target: "/api/secret/cred",
			token:  validToken,
			body:   []byte(`{}`),
			want: want{
				code:          400,
				emptyResponse: false,
				response:      "incorrect data: incorrect request\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #4 not found",
			method: http.MethodPut,
			target: "/api/secret/cred",
			token:  validToken,
			body:   []byte(`{"id":"test","site":"123321","login":"test","password":"123","metadata":"test_meta"}`),
			want: want{
				code:          404,
				emptyResponse: false,
				response:      "nothing to update\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #5 sql error",
			method: http.MethodPut,
			target: "/api/secret/cred",
			token:  validToken,
			body:   []byte(`{"id":"test","site":"123321","login":"test","password":"123","metadata":"test_meta"}`),
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "failed to update: sql: connection is already closed\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testAuthorizedRequest(t, ts, tt.method, tt.target, tt.token, tt.body)
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

func Test_router_cred_delete(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mocksSecret := mocks.NewMockCredsSecret(mockCtrl)

	mocksSecret.EXPECT().Delete(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(nil).Times(1)
	mocksSecret.EXPECT().Delete(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(sql.ErrNoRows).Times(1)
	mocksSecret.EXPECT().Delete(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(sql.ErrConnDone).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, nil, nil, mocksSecret, nil, nil))
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
			target: "/api/secret/cred/test",
			token:  validToken,
			want: want{
				code:          204,
				emptyResponse: true,
			},
		},
		{
			name:   "failed test #1 invalid token",
			method: http.MethodDelete,
			target: "/api/secret/cred/test",
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
			target: "/api/secret/cred/",
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
			method: http.MethodDelete,
			target: "/api/secret/cred/test",
			token:  validToken,
			want: want{
				code:          404,
				emptyResponse: false,
				response:      "nothing to delete\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #4 sql error",
			method: http.MethodDelete,
			target: "/api/secret/cred/test",
			token:  validToken,
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "failed to delete from db: sql: connection is already closed\n",
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
