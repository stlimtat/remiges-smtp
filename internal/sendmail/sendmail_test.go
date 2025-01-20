package sendmail

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/mjl-/adns"
	"github.com/mjl-/mox/dns"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLookupMX(t *testing.T) {
	type LookupMXResult struct {
		mxList     []*net.MX
		ADNSResult adns.Result
		err        error
	}
	var tests = []struct {
		name     string
		domain   dns.Domain
		cname    string
		mxResult LookupMXResult
		wantErr  bool
	}{
		{
			name: "happy",
			domain: dns.Domain{
				ASCII: "abc.com",
			},
			cname: "abc.com.",
			mxResult: LookupMXResult{
				mxList: []*net.MX{
					{Host: "host1.abc.com", Pref: uint16(0)},
					{Host: "host2.abc.com", Pref: uint16(0)},
				},
				ADNSResult: adns.Result{},
				err:        nil,
			},
			wantErr: false,
		},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.InitLogger(ctx)
			slogger := telemetry.GetSLogger(ctx)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			resolver := NewMockResolver(ctrl)
			resolver.EXPECT().
				LookupCNAME(gomock.Any(), tt.cname).
				DoAndReturn(func(_ context.Context, host string) (string, adns.Result, error) {
					assert.Equal(t, host, tt.cname)
					return tt.cname, tt.mxResult.ADNSResult, tt.mxResult.err
				}).AnyTimes()
			resolver.EXPECT().
				LookupMX(gomock.Any(), tt.cname).
				DoAndReturn(func(_ context.Context, domain string) ([]*net.MX, adns.Result, error) {
					assert.Equal(t, domain, tt.cname)
					return tt.mxResult.mxList, tt.mxResult.ADNSResult, tt.mxResult.err
				})

			m := NewMailSender(ctx, nil, resolver, slogger)
			_, err := m.LookupMX(ctx, tt.domain)
			if tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}

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

			m := NewMailSender(ctx, dialerFactory, nil, slogger)
			_, err := m.NewConn(ctx, tt.hosts)
			if tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}
