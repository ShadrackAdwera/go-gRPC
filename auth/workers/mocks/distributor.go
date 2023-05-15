// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ShadrackAdwera/go-gRPC/workers (interfaces: Distributor)

// Package mockworkers is a generated GoMock package.
package mockworkers

import (
	context "context"
	reflect "reflect"

	workers "github.com/ShadrackAdwera/go-gRPC/workers"
	gomock "github.com/golang/mock/gomock"
	asynq "github.com/hibiken/asynq"
)

// MockDistributor is a mock of Distributor interface.
type MockDistributor struct {
	ctrl     *gomock.Controller
	recorder *MockDistributorMockRecorder
}

// MockDistributorMockRecorder is the mock recorder for MockDistributor.
type MockDistributorMockRecorder struct {
	mock *MockDistributor
}

// NewMockDistributor creates a new mock instance.
func NewMockDistributor(ctrl *gomock.Controller) *MockDistributor {
	mock := &MockDistributor{ctrl: ctrl}
	mock.recorder = &MockDistributorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDistributor) EXPECT() *MockDistributorMockRecorder {
	return m.recorder
}

// DistributeUser mocks base method.
func (m *MockDistributor) DistributeUser(arg0 context.Context, arg1 workers.UserPayload, arg2 ...asynq.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DistributeUser", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DistributeUser indicates an expected call of DistributeUser.
func (mr *MockDistributorMockRecorder) DistributeUser(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DistributeUser", reflect.TypeOf((*MockDistributor)(nil).DistributeUser), varargs...)
}
