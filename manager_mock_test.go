// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ory/ladon (interfaces: Manager)

// Package ladon_test is a generated GoMock package.
package ladon_test

import (
	context "context"
	"reflect"

	gomock "github.com/golang/mock/gomock"

	ladon "github.com/paullesiak/ladon"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockManager) Create(arg0 context.Context, arg1 ladon.Policy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockManagerMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockManager)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockManager) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockManagerMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockManager)(nil).Delete), arg0, arg1)
}

// FindPoliciesForResource mocks base method.
func (m *MockManager) FindPoliciesForResource(arg0 context.Context, arg1 string) (ladon.Policies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPoliciesForResource", arg0, arg1)
	ret0, _ := ret[0].(ladon.Policies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPoliciesForResource indicates an expected call of FindPoliciesForResource.
func (mr *MockManagerMockRecorder) FindPoliciesForResource(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPoliciesForResource", reflect.TypeOf((*MockManager)(nil).FindPoliciesForResource), arg0, arg1)
}

// FindPoliciesForSubject mocks base method.
func (m *MockManager) FindPoliciesForSubject(arg0 context.Context, arg1 string) (ladon.Policies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPoliciesForSubject", arg0, arg1)
	ret0, _ := ret[0].(ladon.Policies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPoliciesForSubject indicates an expected call of FindPoliciesForSubject.
func (mr *MockManagerMockRecorder) FindPoliciesForSubject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPoliciesForSubject", reflect.TypeOf((*MockManager)(nil).FindPoliciesForSubject), arg0, arg1)
}

// FindRequestCandidates mocks base method.
func (m *MockManager) FindRequestCandidates(arg0 context.Context, arg1 *ladon.Request) (ladon.Policies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRequestCandidates", arg0, arg1)
	ret0, _ := ret[0].(ladon.Policies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindRequestCandidates indicates an expected call of FindRequestCandidates.
func (mr *MockManagerMockRecorder) FindRequestCandidates(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRequestCandidates", reflect.TypeOf((*MockManager)(nil).FindRequestCandidates), arg0, arg1)
}

// Get mocks base method.
func (m *MockManager) Get(arg0 context.Context, arg1 string) (ladon.Policy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(ladon.Policy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockManagerMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockManager)(nil).Get), arg0, arg1)
}

// GetAll mocks base method.
func (m *MockManager) GetAll(arg0 context.Context, arg1, arg2 int64) (ladon.Policies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0, arg1, arg2)
	ret0, _ := ret[0].(ladon.Policies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockManagerMockRecorder) GetAll(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockManager)(nil).GetAll), arg0, arg1, arg2)
}

// Update mocks base method.
func (m *MockManager) Update(arg0 context.Context, arg1 ladon.Policy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockManagerMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockManager)(nil).Update), arg0, arg1)
}
