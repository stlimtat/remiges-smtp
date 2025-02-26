package dns

import (
	"context"
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
			ctx, _ := telemetry.InitLogger(context.Background())
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

			r := NewResolver(ctx, resolver, slogger)
			_, err := r.LookupMX(ctx, tt.domain)
			if tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}
