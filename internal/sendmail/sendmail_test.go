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
					assert.Equal(t, TCP_NETWORK, network)
					hosts2 := []string{}
					for _, host := range tt.hosts {
						host = fmt.Sprintf("%s:%s", host, DEFAULT_SMTP_PORT_STR)
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
