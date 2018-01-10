// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lucas-clemente/quic-go/internal/handshake (interfaces: MintTLS)

// Package mockhandshake is a generated GoMock package.
package mockhandshake

import (
	io "io"
	reflect "reflect"

	mint "github.com/bifurcation/mint"
	gomock "github.com/golang/mock/gomock"
)

// MockMintTLS is a mock of MintTLS interface
type MockMintTLS struct {
	ctrl     *gomock.Controller
	recorder *MockMintTLSMockRecorder
}

// MockMintTLSMockRecorder is the mock recorder for MockMintTLS
type MockMintTLSMockRecorder struct {
	mock *MockMintTLS
}

// NewMockMintTLS creates a new mock instance
func NewMockMintTLS(ctrl *gomock.Controller) *MockMintTLS {
	mock := &MockMintTLS{ctrl: ctrl}
	mock.recorder = &MockMintTLSMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMintTLS) EXPECT() *MockMintTLSMockRecorder {
	return m.recorder
}

// ComputeExporter mocks base method
func (m *MockMintTLS) ComputeExporter(arg0 string, arg1 []byte, arg2 int) ([]byte, error) {
	ret := m.ctrl.Call(m, "ComputeExporter", arg0, arg1, arg2)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ComputeExporter indicates an expected call of ComputeExporter
func (mr *MockMintTLSMockRecorder) ComputeExporter(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ComputeExporter", reflect.TypeOf((*MockMintTLS)(nil).ComputeExporter), arg0, arg1, arg2)
}

// ConnectionState mocks base method
func (m *MockMintTLS) ConnectionState() mint.ConnectionState {
	ret := m.ctrl.Call(m, "ConnectionState")
	ret0, _ := ret[0].(mint.ConnectionState)
	return ret0
}

// ConnectionState indicates an expected call of ConnectionState
func (mr *MockMintTLSMockRecorder) ConnectionState() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectionState", reflect.TypeOf((*MockMintTLS)(nil).ConnectionState))
}

// GetCipherSuite mocks base method
func (m *MockMintTLS) GetCipherSuite() mint.CipherSuiteParams {
	ret := m.ctrl.Call(m, "GetCipherSuite")
	ret0, _ := ret[0].(mint.CipherSuiteParams)
	return ret0
}

// GetCipherSuite indicates an expected call of GetCipherSuite
func (mr *MockMintTLSMockRecorder) GetCipherSuite() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCipherSuite", reflect.TypeOf((*MockMintTLS)(nil).GetCipherSuite))
}

// Handshake mocks base method
func (m *MockMintTLS) Handshake() mint.Alert {
	ret := m.ctrl.Call(m, "Handshake")
	ret0, _ := ret[0].(mint.Alert)
	return ret0
}

// Handshake indicates an expected call of Handshake
func (mr *MockMintTLSMockRecorder) Handshake() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handshake", reflect.TypeOf((*MockMintTLS)(nil).Handshake))
}

// SetCryptoStream mocks base method
func (m *MockMintTLS) SetCryptoStream(arg0 io.ReadWriter) {
	m.ctrl.Call(m, "SetCryptoStream", arg0)
}

// SetCryptoStream indicates an expected call of SetCryptoStream
func (mr *MockMintTLSMockRecorder) SetCryptoStream(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCryptoStream", reflect.TypeOf((*MockMintTLS)(nil).SetCryptoStream), arg0)
}

// SetExtensionHandler mocks base method
func (m *MockMintTLS) SetExtensionHandler(arg0 mint.AppExtensionHandler) error {
	ret := m.ctrl.Call(m, "SetExtensionHandler", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetExtensionHandler indicates an expected call of SetExtensionHandler
func (mr *MockMintTLSMockRecorder) SetExtensionHandler(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetExtensionHandler", reflect.TypeOf((*MockMintTLS)(nil).SetExtensionHandler), arg0)
}

// State mocks base method
func (m *MockMintTLS) State() mint.State {
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(mint.State)
	return ret0
}

// State indicates an expected call of State
func (mr *MockMintTLSMockRecorder) State() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockMintTLS)(nil).State))
}
