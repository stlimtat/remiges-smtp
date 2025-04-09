package sendmail

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/mjl-/mox/smtpclient"
	"github.com/stlimtat/remiges-smtp/internal/dns"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewMailSender(t *testing.T) {
	tests := []struct {
		name     string
		debug    bool
		wantNil  bool
		wantInit bool
	}{
		{
			name:     "normal_initialization",
			debug:    false,
			wantNil:  false,
			wantInit: true,
		},
		{
			name:     "debug_mode",
			debug:    true,
			wantNil:  false,
			wantInit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDialerFactory := NewMockINetDialerFactory(ctrl)
			mockResolver := dns.NewMockIResolver(ctrl)
			logger := slog.Default()

			sender := NewMailSender(context.Background(), tt.debug, mockDialerFactory, mockResolver, logger)

			if tt.wantNil {
				assert.Nil(t, sender)
			} else {
				assert.NotNil(t, sender)
				if tt.wantInit {
					assert.NotNil(t, sender.CachedMX)
					assert.Equal(t, tt.debug, sender.Debug)
					assert.NotNil(t, sender.metrics)
					assert.Equal(t, 3, sender.maxRetries)
					assert.Equal(t, 5*time.Second, sender.retryDelay)
				}
			}
		})
	}
}

// func TestSendMail(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		mail        *pmail.Mail
// 		setupMocks  func(*MockINetDialerFactory, *dns.MockIResolver)
// 		expectError bool
// 	}{
// 		{
// 			name: "invalid_mail",
// 			mail: &pmail.Mail{}, // Empty mail, will fail validation
// 			setupMocks: func(df *MockINetDialerFactory, r *dns.MockIResolver) {
// 				// No mock setup needed as validation will fail
// 			},
// 			expectError: true,
// 		},
// 		{
// 			name: "mx_lookup_failure",
// 			mail: &pmail.Mail{
// 				From:      smtp.Address{Localpart: "sender", Domain: moxDns.Domain{ASCII: "example.com"}},
// 				To:        []smtp.Address{{Localpart: "recipient", Domain: moxDns.Domain{ASCII: "example.com"}}},
// 				MsgID:     []byte("test-message-id"),
// 				Subject:   []byte("Test Subject"),
// 				FinalBody: []byte("Test Body"),
// 			},
// 			setupMocks: func(df *MockINetDialerFactory, r *dns.MockIResolver) {
// 				r.EXPECT().
// 					LookupMX(gomock.Any(), "example.com").
// 					Return(nil, errors.New("mx lookup failed")).
// 					Times(1)
// 			},
// 			expectError: true,
// 		},
// 		{
// 			name: "connection_failure",
// 			mail: &pmail.Mail{
// 				From:      smtp.Address{Localpart: "sender", Domain: moxDns.Domain{ASCII: "example.com"}},
// 				To:        []smtp.Address{{Localpart: "recipient", Domain: moxDns.Domain{ASCII: "example.com"}}},
// 				MsgID:     []byte("test-message-id"),
// 				Subject:   []byte("Test Subject"),
// 				FinalBody: []byte("Test Body"),
// 			},
// 			setupMocks: func(df *MockINetDialerFactory, r *dns.MockIResolver) {
// 				r.EXPECT().
// 					LookupMX(gomock.Any(), "example.com").
// 					Return([]string{"mx.example.com"}, nil).
// 					Times(1)

// 				df.EXPECT().
// 					NewDialer(gomock.Any()).
// 					Return(nil, errors.New("connection failed")).
// 					Times(1)
// 			},
// 			expectError: true,
// 		},
// 		{
// 			name: "debug_mode_success",
// 			mail: &pmail.Mail{
// 				From:      smtp.Address{Localpart: "sender", Domain: moxDns.Domain{ASCII: "example.com"}},
// 				To:        []smtp.Address{{Localpart: "recipient", Domain: moxDns.Domain{ASCII: "example.com"}}},
// 				MsgID:     []byte("test-message-id"),
// 				Subject:   []byte("Test Subject"),
// 				FinalBody: []byte("Test Body"),
// 			},
// 			setupMocks: func(df *MockINetDialerFactory, r *dns.MockIResolver) {
// 				// No mock setup needed as debug mode skips actual delivery
// 			},
// 			expectError: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockDialerFactory := NewMockINetDialerFactory(ctrl)
// 			mockResolver := dns.NewMockIResolver(ctrl)
// 			logger := slog.Default()

// 			tt.setupMocks(mockDialerFactory, mockResolver)

// 			sender := NewMailSender(context.Background(), tt.name == "debug_mode_success", mockDialerFactory, mockResolver, logger)
// 			responses, errs := sender.SendMail(context.Background(), tt.mail)

// 			if tt.expectError {
// 				assert.NotNil(t, errs)
// 			} else {
// 				assert.Nil(t, errs)
// 				if tt.name == "debug_mode_success" {
// 					assert.Empty(t, responses)
// 				}
// 			}
// 		})
// 	}
// }

// func TestDeliverToRecipient(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		setupMocks  func(*MockINetDialerFactory, *dns.MockIResolver)
// 		expectError bool
// 		retries     int
// 	}{
// 		{
// 			name: "success_first_try",
// 			setupMocks: func(df *MockINetDialerFactory, r *dns.MockIResolver) {
// 				r.EXPECT().
// 					LookupMX(gomock.Any(), "example.com").
// 					Return([]string{"mx.example.com"}, nil).
// 					Times(1)

// 				mockDialer := NewMockDialer(gomock.NewController(t))
// 				mockDialer.EXPECT().
// 					DialContext(gomock.Any(), TCPNetwork, gomock.Any()).
// 					Return(&net.TCPConn{}, nil).
// 					Times(1)

// 				df.EXPECT().
// 					NewDialer(gomock.Any()).
// 					Return(mockDialer, nil).
// 					Times(1)
// 			},
// 			expectError: false,
// 			retries:     0,
// 		},
// 		{
// 			name: "retry_success",
// 			setupMocks: func(df *MockINetDialerFactory, r *dns.MockIResolver) {
// 				r.EXPECT().
// 					LookupMX(gomock.Any(), "example.com").
// 					Return([]string{"mx.example.com"}, nil).
// 					Times(2)

// 				mockDialer := NewMockDialer(gomock.NewController(t))
// 				mockDialer.EXPECT().
// 					DialContext(gomock.Any(), TCPNetwork, gomock.Any()).
// 					Return(nil, errors.New("first attempt failed")).
// 					Times(1)
// 				mockDialer.EXPECT().
// 					DialContext(gomock.Any(), TCPNetwork, gomock.Any()).
// 					Return(&net.TCPConn{}, nil).
// 					Times(1)

// 				df.EXPECT().
// 					NewDialer(gomock.Any()).
// 					Return(mockDialer, nil).
// 					Times(2)
// 			},
// 			expectError: false,
// 			retries:     1,
// 		},
// 		{
// 			name: "max_retries_exceeded",
// 			setupMocks: func(df *MockINetDialerFactory, r *dns.MockIResolver) {
// 				r.EXPECT().
// 					LookupMX(gomock.Any(), "example.com").
// 					Return([]string{"mx.example.com"}, nil).
// 					Times(3)

// 				mockDialer := NewMockDialer(gomock.NewController(t))
// 				mockDialer.EXPECT().
// 					DialContext(gomock.Any(), TCPNetwork, gomock.Any()).
// 					Return(nil, errors.New("connection failed")).
// 					Times(3)

// 				df.EXPECT().
// 					NewDialer(gomock.Any()).
// 					Return(mockDialer, nil).
// 					Times(3)
// 			},
// 			expectError: true,
// 			retries:     3,
// 		},
// 		{
// 			name: "context_cancelled",
// 			setupMocks: func(df *MockINetDialerFactory, r *dns.MockIResolver) {
// 				// No mock setup needed as context will be cancelled
// 			},
// 			expectError: true,
// 			retries:     0,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockDialerFactory := NewMockINetDialerFactory(ctrl)
// 			mockResolver := dns.NewMockIResolver(ctrl)
// 			logger := slog.Default()

// 			tt.setupMocks(mockDialerFactory, mockResolver)

// 			sender := NewMailSender(context.Background(), false, mockDialerFactory, mockResolver, logger)
// 			sender.maxRetries = 3
// 			sender.retryDelay = 10 * time.Millisecond

// 			ctx := context.Background()
// 			if tt.name == "context_cancelled" {
// 				var cancel context.CancelFunc
// 				ctx, cancel = context.WithCancel(ctx)
// 				cancel()
// 			}

// 			mail := &pmail.Mail{
// 				From:      smtp.Address{Localpart: "sender", Domain: moxDns.Domain{ASCII: "example.com"}},
// 				To:        []smtp.Address{{Localpart: "recipient", Domain: moxDns.Domain{ASCII: "example.com"}}},
// 				MsgID:     []byte("test-message-id"),
// 				Subject:   []byte("Test Subject"),
// 				FinalBody: []byte("Test Body"),
// 			}

// 			result := sender.deliverToRecipient(ctx, mail, mail.To[0])

// 			if tt.expectError {
// 				assert.NotNil(t, result.err)
// 			} else {
// 				assert.Nil(t, result.err)
// 			}
// 		})
// 	}
// }

func TestNewConn(t *testing.T) {
	var tests = []struct {
		name    string
		hosts   []string
		wantErr bool
	}{
		{"happy", []string{"aspmx.l.google.com"}, false},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.InitLogger(ctx)
			slogger := telemetry.GetSLogger(ctx)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dialer := NewMockDialer(ctrl)
			dialer.EXPECT().
				DialContext(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, network, address string) (net.Conn, error) {
					assert.Equal(t, TCPNetwork, network)
					hosts2 := []string{}
					for _, host := range tt.hosts {
						host = fmt.Sprintf("%s:%s", host, DefaultSMTPPort)
						hosts2 = append(hosts2, host)
					}
					assert.Contains(t, hosts2, address)
					return &net.TCPConn{}, nil
				})
			dialerFactory := NewMockINetDialerFactory(ctrl)
			dialerFactory.EXPECT().
				NewDialer(gomock.Any()).
				DoAndReturn(func(_ context.Context) (smtpclient.Dialer, error) {
					return dialer, nil
				})

			m := NewMailSender(ctx, false, dialerFactory, nil, slogger)
			_, err := m.NewConn(ctx, tt.hosts)
			if tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}

// Add more comprehensive test cases
// func TestMailSender_SendMail(t *testing.T) {
// 	tests := []struct {
// 			name    string
// 			mail    *pmail.Mail
// 			mocks   func(*MockResolver, *MockDialerFactory)
// 			want    map[string][]pmail.Response
// 			wantErr bool
// 	}{
// 			// Add test cases for various scenarios
// 			{
// 					name: "successful delivery",
// 					mail: &pmail.Mail{...},
// 					mocks: func(r *MockResolver, d *MockDialerFactory) {
// 							// Setup mock expectations
// 					},
// 					want: map[string][]pmail.Response{...},
// 					wantErr: false,
// 			},
// 			// Add more test cases
// 	}

// 	for _, tt := range tests {
// 			t.Run(tt.name, func(t *testing.T) {
// 					// Test implementation
// 					assert.True(t, true)
// 			})
// 	}
// }
