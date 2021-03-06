// Code generated by MockGen. DO NOT EDIT.
// Source: data.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	internal "usermanagement/app/internal"

	gomock "github.com/golang/mock/gomock"
)

// MockUserData is a mock of UserData interface.
type MockUserData struct {
	ctrl     *gomock.Controller
	recorder *MockUserDataMockRecorder
}

// MockUserDataMockRecorder is the mock recorder for MockUserData.
type MockUserDataMockRecorder struct {
	mock *MockUserData
}

// NewMockUserData creates a new mock instance.
func NewMockUserData(ctrl *gomock.Controller) *MockUserData {
	mock := &MockUserData{ctrl: ctrl}
	mock.recorder = &MockUserDataMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserData) EXPECT() *MockUserDataMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockUserData) ChangePassword(userID uint, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", userID, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockUserDataMockRecorder) ChangePassword(userID, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockUserData)(nil).ChangePassword), userID, password)
}

// CreateUser mocks base method.
func (m *MockUserData) CreateUser(request internal.UserRequest) (internal.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", request)
	ret0, _ := ret[0].(internal.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserDataMockRecorder) CreateUser(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserData)(nil).CreateUser), request)
}

// DeleteUser mocks base method.
func (m *MockUserData) DeleteUser(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserDataMockRecorder) DeleteUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserData)(nil).DeleteUser), id)
}

// GetUsers mocks base method.
func (m *MockUserData) GetUsers(offset, limit uint) (internal.UsersResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", offset, limit)
	ret0, _ := ret[0].(internal.UsersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUserDataMockRecorder) GetUsers(offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUserData)(nil).GetUsers), offset, limit)
}

// UpdateUser mocks base method.
func (m *MockUserData) UpdateUser(request internal.UpdateUserRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserDataMockRecorder) UpdateUser(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserData)(nil).UpdateUser), request)
}

// MockGroupData is a mock of GroupData interface.
type MockGroupData struct {
	ctrl     *gomock.Controller
	recorder *MockGroupDataMockRecorder
}

// MockGroupDataMockRecorder is the mock recorder for MockGroupData.
type MockGroupDataMockRecorder struct {
	mock *MockGroupData
}

// NewMockGroupData creates a new mock instance.
func NewMockGroupData(ctrl *gomock.Controller) *MockGroupData {
	mock := &MockGroupData{ctrl: ctrl}
	mock.recorder = &MockGroupDataMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupData) EXPECT() *MockGroupDataMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockGroupData) AddUser(userID, groupID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", userID, groupID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockGroupDataMockRecorder) AddUser(userID, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockGroupData)(nil).AddUser), userID, groupID)
}

// CreateGroup mocks base method.
func (m *MockGroupData) CreateGroup(request internal.GroupRequest) (internal.GroupResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", request)
	ret0, _ := ret[0].(internal.GroupResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockGroupDataMockRecorder) CreateGroup(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockGroupData)(nil).CreateGroup), request)
}

// DeleteGroup mocks base method.
func (m *MockGroupData) DeleteGroup(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroup", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup.
func (mr *MockGroupDataMockRecorder) DeleteGroup(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockGroupData)(nil).DeleteGroup), id)
}

// GetGroups mocks base method.
func (m *MockGroupData) GetGroups(offset, limit uint) (internal.GroupsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroups", offset, limit)
	ret0, _ := ret[0].(internal.GroupsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroups indicates an expected call of GetGroups.
func (mr *MockGroupDataMockRecorder) GetGroups(offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroups", reflect.TypeOf((*MockGroupData)(nil).GetGroups), offset, limit)
}

// GetUsersByGroupID mocks base method.
func (m *MockGroupData) GetUsersByGroupID(groupID, offset, limit uint) (internal.UsersResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersByGroupID", groupID, offset, limit)
	ret0, _ := ret[0].(internal.UsersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersByGroupID indicates an expected call of GetUsersByGroupID.
func (mr *MockGroupDataMockRecorder) GetUsersByGroupID(groupID, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersByGroupID", reflect.TypeOf((*MockGroupData)(nil).GetUsersByGroupID), groupID, offset, limit)
}

// RemoveUser mocks base method.
func (m *MockGroupData) RemoveUser(groupID, userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUser", groupID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUser indicates an expected call of RemoveUser.
func (mr *MockGroupDataMockRecorder) RemoveUser(groupID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUser", reflect.TypeOf((*MockGroupData)(nil).RemoveUser), groupID, userID)
}

// UpdateGroup mocks base method.
func (m *MockGroupData) UpdateGroup(request internal.UpdateGroupRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroup", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGroup indicates an expected call of UpdateGroup.
func (mr *MockGroupDataMockRecorder) UpdateGroup(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroup", reflect.TypeOf((*MockGroupData)(nil).UpdateGroup), request)
}
