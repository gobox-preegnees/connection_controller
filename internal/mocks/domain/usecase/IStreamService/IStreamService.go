// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockIStreamService is a mock of IStreamService interface.
type MockIStreamService struct {
	ctrl     *gomock.Controller
	recorder *MockIStreamServiceMockRecorder
}

// MockIStreamServiceMockRecorder is the mock recorder for MockIStreamService.
type MockIStreamServiceMockRecorder struct {
	mock *MockIStreamService
}

// NewMockIStreamService creates a new mock instance.
func NewMockIStreamService(ctrl *gomock.Controller) *MockIStreamService {
	mock := &MockIStreamService{ctrl: ctrl}
	mock.recorder = &MockIStreamServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIStreamService) EXPECT() *MockIStreamServiceMockRecorder {
	return m.recorder
}

// DeleteStream mocks base method.
func (m *MockIStreamService) DeleteStream(ctx context.Context, stream entity.Stream) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStream", ctx, stream)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteStream indicates an expected call of DeleteStream.
func (mr *MockIStreamServiceMockRecorder) DeleteStream(ctx, stream interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStream", reflect.TypeOf((*MockIStreamService)(nil).DeleteStream), ctx, stream)
}

// SaveStream mocks base method.
func (m *MockIStreamService) SaveStream(ctx context.Context, stream entity.Stream) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveStream", ctx, stream)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveStream indicates an expected call of SaveStream.
func (mr *MockIStreamServiceMockRecorder) SaveStream(ctx, stream interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveStream", reflect.TypeOf((*MockIStreamService)(nil).SaveStream), ctx, stream)
}

// MockIConsistencyService is a mock of IConsistencyService interface.
type MockIConsistencyService struct {
	ctrl     *gomock.Controller
	recorder *MockIConsistencyServiceMockRecorder
}

// MockIConsistencyServiceMockRecorder is the mock recorder for MockIConsistencyService.
type MockIConsistencyServiceMockRecorder struct {
	mock *MockIConsistencyService
}

// NewMockIConsistencyService creates a new mock instance.
func NewMockIConsistencyService(ctrl *gomock.Controller) *MockIConsistencyService {
	mock := &MockIConsistencyService{ctrl: ctrl}
	mock.recorder = &MockIConsistencyServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIConsistencyService) EXPECT() *MockIConsistencyServiceMockRecorder {
	return m.recorder
}

// GetConsistency mocks base method.
func (m *MockIConsistencyService) GetConsistency(ctx context.Context) (entity.Consistency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConsistency", ctx)
	ret0, _ := ret[0].(entity.Consistency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConsistency indicates an expected call of GetConsistency.
func (mr *MockIConsistencyServiceMockRecorder) GetConsistency(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConsistency", reflect.TypeOf((*MockIConsistencyService)(nil).GetConsistency), ctx)
}

// SaveConsistency mocks base method.
func (m *MockIConsistencyService) SaveConsistency(ctx context.Context, consistency entity.Consistency) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveConsistency", ctx, consistency)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveConsistency indicates an expected call of SaveConsistency.
func (mr *MockIConsistencyServiceMockRecorder) SaveConsistency(ctx, consistency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveConsistency", reflect.TypeOf((*MockIConsistencyService)(nil).SaveConsistency), ctx, consistency)
}

// MockISnapshotService is a mock of ISnapshotService interface.
type MockISnapshotService struct {
	ctrl     *gomock.Controller
	recorder *MockISnapshotServiceMockRecorder
}

// MockISnapshotServiceMockRecorder is the mock recorder for MockISnapshotService.
type MockISnapshotServiceMockRecorder struct {
	mock *MockISnapshotService
}

// NewMockISnapshotService creates a new mock instance.
func NewMockISnapshotService(ctrl *gomock.Controller) *MockISnapshotService {
	mock := &MockISnapshotService{ctrl: ctrl}
	mock.recorder = &MockISnapshotServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISnapshotService) EXPECT() *MockISnapshotServiceMockRecorder {
	return m.recorder
}

// SaveSnapshot mocks base method.
func (m *MockISnapshotService) SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSnapshot", ctx, snapshot)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSnapshot indicates an expected call of SaveSnapshot.
func (mr *MockISnapshotServiceMockRecorder) SaveSnapshot(ctx, snapshot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSnapshot", reflect.TypeOf((*MockISnapshotService)(nil).SaveSnapshot), ctx, snapshot)
}
