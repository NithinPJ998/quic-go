// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lucas-clemente/quic-go/internal/ackhandler (interfaces: SentPacketHandler)

// Package mockackhandler is a generated GoMock package.
package mockackhandler

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	ackhandler "github.com/lucas-clemente/quic-go/internal/ackhandler"
	protocol "github.com/lucas-clemente/quic-go/internal/protocol"
	wire "github.com/lucas-clemente/quic-go/internal/wire"
)

// MockSentPacketHandler is a mock of SentPacketHandler interface
type MockSentPacketHandler struct {
	ctrl     *gomock.Controller
	recorder *MockSentPacketHandlerMockRecorder
}

// MockSentPacketHandlerMockRecorder is the mock recorder for MockSentPacketHandler
type MockSentPacketHandlerMockRecorder struct {
	mock *MockSentPacketHandler
}

// NewMockSentPacketHandler creates a new mock instance
func NewMockSentPacketHandler(ctrl *gomock.Controller) *MockSentPacketHandler {
	mock := &MockSentPacketHandler{ctrl: ctrl}
	mock.recorder = &MockSentPacketHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSentPacketHandler) EXPECT() *MockSentPacketHandlerMockRecorder {
	return m.recorder
}

// DequeuePacketForRetransmission mocks base method
func (m *MockSentPacketHandler) DequeuePacketForRetransmission() *ackhandler.Packet {
	ret := m.ctrl.Call(m, "DequeuePacketForRetransmission")
	ret0, _ := ret[0].(*ackhandler.Packet)
	return ret0
}

// DequeuePacketForRetransmission indicates an expected call of DequeuePacketForRetransmission
func (mr *MockSentPacketHandlerMockRecorder) DequeuePacketForRetransmission() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DequeuePacketForRetransmission", reflect.TypeOf((*MockSentPacketHandler)(nil).DequeuePacketForRetransmission))
}

// DequeueProbePacket mocks base method
func (m *MockSentPacketHandler) DequeueProbePacket() (*ackhandler.Packet, error) {
	ret := m.ctrl.Call(m, "DequeueProbePacket")
	ret0, _ := ret[0].(*ackhandler.Packet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DequeueProbePacket indicates an expected call of DequeueProbePacket
func (mr *MockSentPacketHandlerMockRecorder) DequeueProbePacket() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DequeueProbePacket", reflect.TypeOf((*MockSentPacketHandler)(nil).DequeueProbePacket))
}

// GetAlarmTimeout mocks base method
func (m *MockSentPacketHandler) GetAlarmTimeout() time.Time {
	ret := m.ctrl.Call(m, "GetAlarmTimeout")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetAlarmTimeout indicates an expected call of GetAlarmTimeout
func (mr *MockSentPacketHandlerMockRecorder) GetAlarmTimeout() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmTimeout", reflect.TypeOf((*MockSentPacketHandler)(nil).GetAlarmTimeout))
}

// GetLowestPacketNotConfirmedAcked mocks base method
func (m *MockSentPacketHandler) GetLowestPacketNotConfirmedAcked() protocol.PacketNumber {
	ret := m.ctrl.Call(m, "GetLowestPacketNotConfirmedAcked")
	ret0, _ := ret[0].(protocol.PacketNumber)
	return ret0
}

// GetLowestPacketNotConfirmedAcked indicates an expected call of GetLowestPacketNotConfirmedAcked
func (mr *MockSentPacketHandlerMockRecorder) GetLowestPacketNotConfirmedAcked() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLowestPacketNotConfirmedAcked", reflect.TypeOf((*MockSentPacketHandler)(nil).GetLowestPacketNotConfirmedAcked))
}

// GetPacketNumberLen mocks base method
func (m *MockSentPacketHandler) GetPacketNumberLen(arg0 protocol.PacketNumber) protocol.PacketNumberLen {
	ret := m.ctrl.Call(m, "GetPacketNumberLen", arg0)
	ret0, _ := ret[0].(protocol.PacketNumberLen)
	return ret0
}

// GetPacketNumberLen indicates an expected call of GetPacketNumberLen
func (mr *MockSentPacketHandlerMockRecorder) GetPacketNumberLen(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPacketNumberLen", reflect.TypeOf((*MockSentPacketHandler)(nil).GetPacketNumberLen), arg0)
}

// OnAlarm mocks base method
func (m *MockSentPacketHandler) OnAlarm() error {
	ret := m.ctrl.Call(m, "OnAlarm")
	ret0, _ := ret[0].(error)
	return ret0
}

// OnAlarm indicates an expected call of OnAlarm
func (mr *MockSentPacketHandlerMockRecorder) OnAlarm() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnAlarm", reflect.TypeOf((*MockSentPacketHandler)(nil).OnAlarm))
}

// ReceivedAck mocks base method
func (m *MockSentPacketHandler) ReceivedAck(arg0 *wire.AckFrame, arg1 protocol.PacketNumber, arg2 protocol.EncryptionLevel, arg3 time.Time) error {
	ret := m.ctrl.Call(m, "ReceivedAck", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReceivedAck indicates an expected call of ReceivedAck
func (mr *MockSentPacketHandlerMockRecorder) ReceivedAck(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceivedAck", reflect.TypeOf((*MockSentPacketHandler)(nil).ReceivedAck), arg0, arg1, arg2, arg3)
}

// SendMode mocks base method
func (m *MockSentPacketHandler) SendMode() ackhandler.SendMode {
	ret := m.ctrl.Call(m, "SendMode")
	ret0, _ := ret[0].(ackhandler.SendMode)
	return ret0
}

// SendMode indicates an expected call of SendMode
func (mr *MockSentPacketHandlerMockRecorder) SendMode() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMode", reflect.TypeOf((*MockSentPacketHandler)(nil).SendMode))
}

// SentPacket mocks base method
func (m *MockSentPacketHandler) SentPacket(arg0 *ackhandler.Packet) {
	m.ctrl.Call(m, "SentPacket", arg0)
}

// SentPacket indicates an expected call of SentPacket
func (mr *MockSentPacketHandlerMockRecorder) SentPacket(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SentPacket", reflect.TypeOf((*MockSentPacketHandler)(nil).SentPacket), arg0)
}

// SentPacketsAsRetransmission mocks base method
func (m *MockSentPacketHandler) SentPacketsAsRetransmission(arg0 []*ackhandler.Packet, arg1 protocol.PacketNumber) {
	m.ctrl.Call(m, "SentPacketsAsRetransmission", arg0, arg1)
}

// SentPacketsAsRetransmission indicates an expected call of SentPacketsAsRetransmission
func (mr *MockSentPacketHandlerMockRecorder) SentPacketsAsRetransmission(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SentPacketsAsRetransmission", reflect.TypeOf((*MockSentPacketHandler)(nil).SentPacketsAsRetransmission), arg0, arg1)
}

// SetHandshakeComplete mocks base method
func (m *MockSentPacketHandler) SetHandshakeComplete() {
	m.ctrl.Call(m, "SetHandshakeComplete")
}

// SetHandshakeComplete indicates an expected call of SetHandshakeComplete
func (mr *MockSentPacketHandlerMockRecorder) SetHandshakeComplete() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHandshakeComplete", reflect.TypeOf((*MockSentPacketHandler)(nil).SetHandshakeComplete))
}

// ShouldSendNumPackets mocks base method
func (m *MockSentPacketHandler) ShouldSendNumPackets() int {
	ret := m.ctrl.Call(m, "ShouldSendNumPackets")
	ret0, _ := ret[0].(int)
	return ret0
}

// ShouldSendNumPackets indicates an expected call of ShouldSendNumPackets
func (mr *MockSentPacketHandlerMockRecorder) ShouldSendNumPackets() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldSendNumPackets", reflect.TypeOf((*MockSentPacketHandler)(nil).ShouldSendNumPackets))
}

// TimeUntilSend mocks base method
func (m *MockSentPacketHandler) TimeUntilSend() time.Time {
	ret := m.ctrl.Call(m, "TimeUntilSend")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// TimeUntilSend indicates an expected call of TimeUntilSend
func (mr *MockSentPacketHandlerMockRecorder) TimeUntilSend() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeUntilSend", reflect.TypeOf((*MockSentPacketHandler)(nil).TimeUntilSend))
}
