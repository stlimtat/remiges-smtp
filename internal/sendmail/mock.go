// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go
//
// Generated by this command:
//
//	mockgen -destination=mock.go -package=sendmail -source=interface.go
//

// Package sendmail is a generated GoMock package.
package sendmail

import (
	context "context"
	net "net"
	reflect "reflect"

	dns "github.com/mjl-/mox/dns"
	smtpclient "github.com/mjl-/mox/smtpclient"
	mail "github.com/stlimtat/remiges-smtp/internal/mail"
	gomock "go.uber.org/mock/gomock"
)

// MockINetDialerFactory is a mock of INetDialerFactory interface.
type MockINetDialerFactory struct {
	ctrl     *gomock.Controller
	recorder *MockINetDialerFactoryMockRecorder
	isgomock struct{}
}

// MockINetDialerFactoryMockRecorder is the mock recorder for MockINetDialerFactory.
type MockINetDialerFactoryMockRecorder struct {
	mock *MockINetDialerFactory
}

// NewMockINetDialerFactory creates a new mock instance.
func NewMockINetDialerFactory(ctrl *gomock.Controller) *MockINetDialerFactory {
	mock := &MockINetDialerFactory{ctrl: ctrl}
	mock.recorder = &MockINetDialerFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockINetDialerFactory) EXPECT() *MockINetDialerFactoryMockRecorder {
	return m.recorder
}

// NewDialer mocks base method.
func (m *MockINetDialerFactory) NewDialer() smtpclient.Dialer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewDialer")
	ret0, _ := ret[0].(smtpclient.Dialer)
	return ret0
}

// NewDialer indicates an expected call of NewDialer.
func (mr *MockINetDialerFactoryMockRecorder) NewDialer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewDialer", reflect.TypeOf((*MockINetDialerFactory)(nil).NewDialer))
}

// MockIMailSender is a mock of IMailSender interface.
type MockIMailSender struct {
	ctrl     *gomock.Controller
	recorder *MockIMailSenderMockRecorder
	isgomock struct{}
}

// MockIMailSenderMockRecorder is the mock recorder for MockIMailSender.
type MockIMailSenderMockRecorder struct {
	mock *MockIMailSender
}

// NewMockIMailSender creates a new mock instance.
func NewMockIMailSender(ctrl *gomock.Controller) *MockIMailSender {
	mock := &MockIMailSender{ctrl: ctrl}
	mock.recorder = &MockIMailSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIMailSender) EXPECT() *MockIMailSenderMockRecorder {
	return m.recorder
}

// LookupMX mocks base method.
func (m *MockIMailSender) LookupMX(ctx context.Context, domain dns.Domain) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupMX", ctx, domain)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupMX indicates an expected call of LookupMX.
func (mr *MockIMailSenderMockRecorder) LookupMX(ctx, domain any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupMX", reflect.TypeOf((*MockIMailSender)(nil).LookupMX), ctx, domain)
}

// NewConn mocks base method.
func (m *MockIMailSender) NewConn(ctx context.Context, hosts []string) (net.Conn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewConn", ctx, hosts)
	ret0, _ := ret[0].(net.Conn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewConn indicates an expected call of NewConn.
func (mr *MockIMailSenderMockRecorder) NewConn(ctx, hosts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewConn", reflect.TypeOf((*MockIMailSender)(nil).NewConn), ctx, hosts)
}

// SendMail mocks base method.
func (m *MockIMailSender) SendMail(ctx context.Context, conn net.Conn, mail *mail.Mail) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMail", ctx, conn, mail)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMail indicates an expected call of SendMail.
func (mr *MockIMailSenderMockRecorder) SendMail(ctx, conn, mail any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMail", reflect.TypeOf((*MockIMailSender)(nil).SendMail), ctx, conn, mail)
}
