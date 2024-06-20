package server

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"keeper-project/types"

	"keeper-project/internal/mocks"
)

var (
	logger = zap.L()
)

func Test_router_register(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUsers := mocks.NewMockUser(mockCtrl)

	mockUsers.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockUsers.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(types.ErrUserAlreadyExists).Times(1)
	mockUsers.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(sql.ErrConnDone).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, mockUsers, nil, nil, nil, nil))
	defer ts.Close()

	tests := []struct {
		name   string
		method string
		target string
		body   []byte
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodPost,
			target: "/api/user/register",
			body:   []byte(`{"login":"test","password":"test"}`),
			want: want{
				code:          200,
				emptyResponse: true,
			},
		},
		{
			name:   "failed test #1 empty login",
			method: http.MethodPost,
			target: "/api/user/register",
			body:   []byte(`{"password":"test"}`),
			want: want{
				code:          400,
				emptyResponse: false,
				response:      "Missing login or password.\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #2 invalid json",
			method: http.MethodPost,
			target: "/api/user/register",
			body:   []byte(`{"password":"test"`),
			want: want{
				code:          400,
				emptyResponse: false,
				response:      "Unable to decode json: unexpected EOF\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #3 already exists",
			method: http.MethodPost,
			target: "/api/user/register",
			body:   []byte(`{"login":"test","password":"test"}`),
			want: want{
				code:          409,
				emptyResponse: false,
				response:      "Unable to create user: user already exists\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #4 sql error",
			method: http.MethodPost,
			target: "/api/user/register",
			body:   []byte(`{"login":"test","password":"test"}`),
			want: want{
				code:          500,
				emptyResponse: false,
				response:      "Unable to create user: sql: connection is already closed\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testRequest(t, ts, tt.method, tt.target, tt.body)
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

func Test_router_auth(t *testing.T) {
	type want struct {
		code          int
		emptyResponse bool
		response      string
		contentType   string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUsers := mocks.NewMockUser(mockCtrl)

	user := &types.User{Login: "test", Password: "ns4IbpusSR+sXB0QRsoR1ze5KisuvZPwBde3EBEMCmeCiBZuf755aIOk8umzyp9IT1IdDORkNFzBrslneRScFA=="}

	mockUsers.EXPECT().GetByLogin(gomock.Any(), "test").Return(user, nil).Times(2)
	mockUsers.EXPECT().GetByLogin(gomock.Any(), "test").Return(nil, sql.ErrNoRows).Times(1)

	ts := httptest.NewServer(SetupRouter(logger, mockUsers, nil, nil, nil, nil))
	defer ts.Close()

	tests := []struct {
		name   string
		method string
		target string
		body   []byte
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodPost,
			target: "/api/user/login",
			body:   []byte(`{"login":"test","password":"test"}`),
			want: want{
				code:          200,
				emptyResponse: true,
			},
		},
		{
			name:   "failed test #1 empty login",
			method: http.MethodPost,
			target: "/api/user/login",
			body:   []byte(`{"password":"test"}`),
			want: want{
				code:          400,
				emptyResponse: false,
				response:      "Missing login or password.\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #2 invalid json",
			method: http.MethodPost,
			target: "/api/user/login",
			body:   []byte(`{"password":"test"`),
			want: want{
				code:          400,
				emptyResponse: false,
				response:      "Unable to decode json: unexpected EOF\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #3 invalid password",
			method: http.MethodPost,
			target: "/api/user/login",
			body:   []byte(`{"login":"test","password":"invalid"}`),
			want: want{
				code:          401,
				emptyResponse: false,
				response:      "Unauthorized\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test #4 not found",
			method: http.MethodPost,
			target: "/api/user/login",
			body:   []byte(`{"login":"test","password":"test"}`),
			want: want{
				code:          400,
				emptyResponse: false,
				response:      "Unable to find user: sql: no rows in result set\n",
				contentType:   "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testRequest(t, ts, tt.method, tt.target, tt.body)
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

func testRequest(t *testing.T, ts *httptest.Server,
	method, path string, body []byte) (*http.Response, string) {
	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest(method, ts.URL+path, bodyReader)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
