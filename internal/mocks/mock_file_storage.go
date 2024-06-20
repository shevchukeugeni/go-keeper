// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/eugene/Work/Projects/learning/Yandex/keeper-project/internal/store/file/storage.go

// Package mock_file is a generated GoMock package.
package mocks

import (
	context "context"
	types "keeper-project/types"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// CreateFile mocks base method.
func (m *MockStorage) CreateFile(ctx context.Context, bucketName string, file *types.File) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFile", ctx, bucketName, file)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFile indicates an expected call of CreateFile.
func (mr *MockStorageMockRecorder) CreateFile(ctx, bucketName, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFile", reflect.TypeOf((*MockStorage)(nil).CreateFile), ctx, bucketName, file)
}

// DeleteFile mocks base method.
func (m *MockStorage) DeleteFile(ctx context.Context, bucketName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", ctx, bucketName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockStorageMockRecorder) DeleteFile(ctx, bucketName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockStorage)(nil).DeleteFile), ctx, bucketName, fileName)
}

// GetFile mocks base method.
func (m *MockStorage) GetFile(ctx context.Context, bucketName, fileName string) (*types.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFile", ctx, bucketName, fileName)
	ret0, _ := ret[0].(*types.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFile indicates an expected call of GetFile.
func (mr *MockStorageMockRecorder) GetFile(ctx, bucketName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockStorage)(nil).GetFile), ctx, bucketName, fileName)
}

// GetFilesList mocks base method.
func (m *MockStorage) GetFilesList(ctx context.Context, bucketName string) ([]*types.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilesList", ctx, bucketName)
	ret0, _ := ret[0].([]*types.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilesList indicates an expected call of GetFilesList.
func (mr *MockStorageMockRecorder) GetFilesList(ctx, bucketName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilesList", reflect.TypeOf((*MockStorage)(nil).GetFilesList), ctx, bucketName)
}
