// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gardener/gardener/pkg/mock/gardener/utils (interfaces: NewTimer)

// Package utils is a generated GoMock package.
package utils

import (
	utils "github.com/gardener/gardener/pkg/utils"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockNewTimer is a mock of NewTimer interface
type MockNewTimer struct {
	ctrl     *gomock.Controller
	recorder *MockNewTimerMockRecorder
}

// MockNewTimerMockRecorder is the mock recorder for MockNewTimer
type MockNewTimerMockRecorder struct {
	mock *MockNewTimer
}

// NewMockNewTimer creates a new mock instance
func NewMockNewTimer(ctrl *gomock.Controller) *MockNewTimer {
	mock := &MockNewTimer{ctrl: ctrl}
	mock.recorder = &MockNewTimerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNewTimer) EXPECT() *MockNewTimerMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockNewTimer) Do(arg0 time.Duration) utils.Timer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(utils.Timer)
	return ret0
}

// Do indicates an expected call of Do
func (mr *MockNewTimerMockRecorder) Do(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockNewTimer)(nil).Do), arg0)
}