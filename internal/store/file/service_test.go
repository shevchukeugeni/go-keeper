package file

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"keeper-project/internal/mocks"
	"keeper-project/types"
)

func TestService_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileStorage := mocks.NewMockStorage(mockCtrl)

	fs, err := NewService(mockFileStorage, zap.L())
	require.NoError(t, err)

	filepath := t.TempDir() + "/test.txt"
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	require.NoError(t, err)

	mockFileStorage.EXPECT().CreateFile(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", gomock.Any()).Return(nil).Times(1)
	mockFileStorage.EXPECT().CreateFile(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", gomock.Any()).Return(testErr).Times(1)

	tests := []struct {
		name    string
		bucket  string
		fileDTO types.CreateFileDTO
		wantErr bool
		err     error
	}{
		{
			name:   "Positive test Create",
			bucket: "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
			fileDTO: types.CreateFileDTO{
				Name:     "123321",
				Size:     1000,
				Metadata: "test_meta",
				Reader:   file,
			},
			wantErr: false,
		},
		{
			name:   "Failed test #1 Failed to create new file",
			bucket: "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
			fileDTO: types.CreateFileDTO{
				Name:     "123321",
				Size:     1000,
				Metadata: "test_meta",
				Reader:   errReader(0),
			},
			wantErr: true,
			err:     testErr,
		},
		{
			name:   "Failed test #2 Client err",
			bucket: "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
			fileDTO: types.CreateFileDTO{
				Name:     "123321",
				Size:     1000,
				Metadata: "test_meta",
				Reader:   file,
			},
			wantErr: true,
			err:     testErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = fs.Create(context.Background(), tt.bucket, tt.fileDTO)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileStorage := mocks.NewMockStorage(mockCtrl)

	fs, err := NewService(mockFileStorage, zap.L())
	require.NoError(t, err)

	fileTest := &types.File{
		Name:     "123321",
		Size:     1000,
		Bytes:    []byte{1},
		Metadata: "test_meta",
	}

	mockFileStorage.EXPECT().GetFile(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(fileTest, nil).Times(1)
	mockFileStorage.EXPECT().GetFile(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(nil, notFound).Times(1)

	tests := []struct {
		name     string
		bucket   string
		fileName string
		wantErr  bool
		err      error
	}{
		{
			name:     "Positive test Get",
			bucket:   "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
			fileName: "test",
			wantErr:  false,
		},
		{
			name:     "Failed test #1 Client err",
			bucket:   "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
			fileName: "test",
			wantErr:  true,
			err:      notFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fs.GetFile(context.Background(), tt.bucket, tt.fileName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetFilesList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileStorage := mocks.NewMockStorage(mockCtrl)

	fs, err := NewService(mockFileStorage, zap.L())
	require.NoError(t, err)

	filesList := []*types.Key{{
		Id:  "test",
		Key: "test_key"},
	}

	mockFileStorage.EXPECT().GetFilesList(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83").Return(filesList, nil).Times(1)
	mockFileStorage.EXPECT().GetFilesList(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83").Return(nil, notFound).Times(1)

	tests := []struct {
		name    string
		bucket  string
		wantErr bool
		err     error
	}{
		{
			name:    "Positive test Get",
			bucket:  "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
			wantErr: false,
		},
		{
			name:    "Failed test #1 Client err",
			bucket:  "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
			wantErr: true,
			err:     notFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fs.GetFilesList(context.Background(), tt.bucket)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileStorage := mocks.NewMockStorage(mockCtrl)

	fs, err := NewService(mockFileStorage, zap.L())
	require.NoError(t, err)

	mockFileStorage.EXPECT().DeleteFile(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(nil).Times(1)
	mockFileStorage.EXPECT().DeleteFile(gomock.Any(), "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83", "test").Return(testErr).Times(1)

	tests := []struct {
		name     string
		bucket   string
		fileName string
		wantErr  bool
		err      error
	}{
		{
			name:     "Positive test Get",
			bucket:   "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
			fileName: "test",
			wantErr:  false,
		},
		{
			name:     "Failed test #1 Client err",
			bucket:   "40d3289b-cc0c-4e2d-81b1-51ec81aa2e83",
			fileName: "test",
			wantErr:  true,
			err:      testErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fs.Delete(context.Background(), tt.bucket, tt.fileName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type errReader int

var testErr = errors.New("test error")
var notFound = errors.New("not found")

func (errReader) Read(p []byte) (n int, err error) {
	return 0, testErr
}
