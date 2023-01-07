// Code generated by MockGen. DO NOT EDIT.
// Source: snapshot.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	messagebroker "github.com/gobox-preegnees/connection_controller/internal/adapter/message_broker"
	gomock "github.com/golang/mock/gomock"
)

// MockISnapshotMessageBroker is a mock of ISnapshotMessageBroker interface.
type MockISnapshotMessageBroker struct {
	ctrl     *gomock.Controller
	recorder *MockISnapshotMessageBrokerMockRecorder
}

// MockISnapshotMessageBrokerMockRecorder is the mock recorder for MockISnapshotMessageBroker.
type MockISnapshotMessageBrokerMockRecorder struct {
	mock *MockISnapshotMessageBroker
}

// NewMockISnapshotMessageBroker creates a new mock instance.
func NewMockISnapshotMessageBroker(ctrl *gomock.Controller) *MockISnapshotMessageBroker {
	mock := &MockISnapshotMessageBroker{ctrl: ctrl}
	mock.recorder = &MockISnapshotMessageBrokerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISnapshotMessageBroker) EXPECT() *MockISnapshotMessageBrokerMockRecorder {
	return m.recorder
}

// CreateOneSnapshot mocks base method.
func (m *MockISnapshotMessageBroker) CreateOneSnapshot(arg0 messagebroker.PublishSnapshotReqDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOneSnapshot", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOneSnapshot indicates an expected call of CreateOneSnapshot.
func (mr *MockISnapshotMessageBrokerMockRecorder) CreateOneSnapshot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOneSnapshot", reflect.TypeOf((*MockISnapshotMessageBroker)(nil).CreateOneSnapshot), arg0)
}
