// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-servicenow-virtual-agent/server/plugin (interfaces: Client)

// Package mock_plugin is a generated GoMock package.
package mock_plugin

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/mattermost/mattermost-server/v5/model"

	serializer "github.com/mattermost/mattermost-plugin-servicenow-virtual-agent/server/serializer"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetMe mocks base method.
func (m *MockClient) GetMe(arg0 string) (*serializer.ServiceNowUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMe", arg0)
	ret0, _ := ret[0].(*serializer.ServiceNowUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMe indicates an expected call of GetMe.
func (mr *MockClientMockRecorder) GetMe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMe", reflect.TypeOf((*MockClient)(nil).GetMe), arg0)
}

// OpenDialogRequest mocks base method.
func (m *MockClient) OpenDialogRequest(arg0 *model.OpenDialogRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenDialogRequest", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// OpenDialogRequest indicates an expected call of OpenDialogRequest.
func (mr *MockClientMockRecorder) OpenDialogRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenDialogRequest", reflect.TypeOf((*MockClient)(nil).OpenDialogRequest), arg0)
}

// SendMessageToVirtualAgentAPI mocks base method.
func (m *MockClient) SendMessageToVirtualAgentAPI(arg0, arg1 string, arg2 bool, arg3 *serializer.MessageAttachment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessageToVirtualAgentAPI", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessageToVirtualAgentAPI indicates an expected call of SendMessageToVirtualAgentAPI.
func (mr *MockClientMockRecorder) SendMessageToVirtualAgentAPI(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessageToVirtualAgentAPI", reflect.TypeOf((*MockClient)(nil).SendMessageToVirtualAgentAPI), arg0, arg1, arg2, arg3)
}

// StartConverstaionWithVirtualAgent mocks base method.
func (m *MockClient) StartConverstaionWithVirtualAgent(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartConverstaionWithVirtualAgent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartConverstaionWithVirtualAgent indicates an expected call of StartConverstaionWithVirtualAgent.
func (mr *MockClientMockRecorder) StartConverstaionWithVirtualAgent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartConverstaionWithVirtualAgent", reflect.TypeOf((*MockClient)(nil).StartConverstaionWithVirtualAgent), arg0)
}
