package sendmail

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

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
				NewDialer().Return(dialer)

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
