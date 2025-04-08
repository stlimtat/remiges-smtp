package dns

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/mjl-/adns"
	"github.com/mjl-/mox/dns"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLookupMX(t *testing.T) {
	type LookupMXResult struct {
		mxList     []*net.MX
		ADNSResult adns.Result
		err        error
	}

	tests := []struct {
		name          string
		domain        dns.Domain
		cname         string
		mxResult      LookupMXResult
		expectedHosts []string
		wantErr       bool
		errorContains string
		setupMock     func(*MockResolver, *gomock.Controller)
	}{
		{
			name: "happy path - single MX record",
			domain: dns.Domain{
				ASCII: "example.com",
			},
			cname: "example.com.",
			mxResult: LookupMXResult{
				mxList: []*net.MX{
					{Host: "mail.example.com", Pref: uint16(10)},
				},
				ADNSResult: adns.Result{},
				err:        nil,
			},
			expectedHosts: []string{"mail.example.com"},
			wantErr:       false,
		},
		{
			name: "happy path - multiple MX records",
			domain: dns.Domain{
				ASCII: "example.com",
			},
			cname: "example.com.",
			mxResult: LookupMXResult{
				mxList: []*net.MX{
					{Host: "mail1.example.com", Pref: uint16(10)},
					{Host: "mail2.example.com", Pref: uint16(20)},
					{Host: "mail3.example.com", Pref: uint16(30)},
				},
				ADNSResult: adns.Result{},
				err:        nil,
			},
			expectedHosts: []string{"mail1.example.com", "mail2.example.com", "mail3.example.com"},
			wantErr:       false,
		},
		// {
		// 	name: "domain with CNAME expansion",
		// 	domain: dns.Domain{
		// 		ASCII: "alias.example.com",
		// 	},
		// 	cname: "real.example.com.",
		// 	mxResult: LookupMXResult{
		// 		mxList: []*net.MX{
		// 			{Host: "mail.real.example.com", Pref: uint16(10)},
		// 		},
		// 		ADNSResult: adns.Result{},
		// 		err:        nil,
		// 	},
		// 	expectedHosts: []string{"mail.real.example.com"},
		// 	wantErr:       false,
		// },
		// {
		// 	name: "empty MX records",
		// 	domain: dns.Domain{
		// 		ASCII: "example.com",
		// 	},
		// 	cname: "example.com.",
		// 	mxResult: LookupMXResult{
		// 		mxList:     []*net.MX{},
		// 		ADNSResult: adns.Result{},
		// 		err:        nil,
		// 	},
		// 	wantErr:       true,
		// 	errorContains: "no MX records found",
		// },
		// {
		// 	name: "DNS resolution error",
		// 	domain: dns.Domain{
		// 		ASCII: "example.com",
		// 	},
		// 	cname: "example.com.",
		// 	mxResult: LookupMXResult{
		// 		mxList:     nil,
		// 		ADNSResult: adns.Result{},
		// 		err:        errors.New("DNS resolution failed"),
		// 	},
		// 	wantErr:       true,
		// 	errorContains: "DNS resolution failed",
		// },
		// {
		// 	name: "invalid domain",
		// 	domain: dns.Domain{
		// 		ASCII: "invalid..domain",
		// 	},
		// 	cname:         "invalid..domain.",
		// 	mxResult:      LookupMXResult{},
		// 	wantErr:       true,
		// 	errorContains: "invalid domain",
		// },
		// {
		// 	name: "timeout error",
		// 	domain: dns.Domain{
		// 		ASCII: "example.com",
		// 	},
		// 	cname: "example.com.",
		// 	mxResult: LookupMXResult{
		// 		mxList:     nil,
		// 		ADNSResult: adns.Result{},
		// 		err:        context.DeadlineExceeded,
		// 	},
		// 	wantErr:       true,
		// 	errorContains: "context deadline exceeded",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			slogger := telemetry.GetSLogger(ctx)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			resolver := NewMockResolver(ctrl)

			// Setup CNAME lookup expectation
			resolver.EXPECT().
				LookupCNAME(gomock.Any(), tt.cname).
				DoAndReturn(func(_ context.Context, host string) (string, adns.Result, error) {
					// assert.Equal(t, host, tt.cname)
					return tt.cname, tt.mxResult.ADNSResult, tt.mxResult.err
				}).AnyTimes()

			// Setup MX lookup expectation
			resolver.EXPECT().
				LookupMX(gomock.Any(), tt.cname).
				DoAndReturn(func(_ context.Context, domain string) ([]*net.MX, adns.Result, error) {
					assert.Equal(t, domain, tt.cname)
					return tt.mxResult.mxList, tt.mxResult.ADNSResult, tt.mxResult.err
				})

			r := NewResolver(ctx, resolver, slogger)
			hosts, err := r.LookupMX(ctx, tt.domain)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedHosts, hosts)
		})
	}
}

// func TestResolver_EdgeCases(t *testing.T) {
// 	ctx, _ := telemetry.InitLogger(context.Background())
// 	slogger := telemetry.GetSLogger(ctx)

// 	tests := []struct {
// 		name          string
// 		domain        dns.Domain
// 		setupResolver func(*MockResolver, *gomock.Controller)
// 		wantErr       bool
// 	}{
// 		{
// 			name: "very long domain name",
// 			domain: dns.Domain{
// 				ASCII: "a" + string(make([]byte, 253)) + ".com",
// 			},
// 			setupResolver: func(mr *MockResolver, ctrl *gomock.Controller) {
// 				mr.EXPECT().
// 					LookupCNAME(gomock.Any(), gomock.Any()).
// 					Return("", adns.Result{}, nil).
// 					AnyTimes()
// 				mr.EXPECT().
// 					LookupMX(gomock.Any(), gomock.Any()).
// 					Return([]*net.MX{{Host: "mail.example.com", Pref: 10}}, adns.Result{}, nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "domain with special characters",
// 			domain: dns.Domain{
// 				ASCII: "test-123.example.com",
// 			},
// 			setupResolver: func(mr *MockResolver, ctrl *gomock.Controller) {
// 				mr.EXPECT().
// 					LookupCNAME(gomock.Any(), gomock.Any()).
// 					Return("", adns.Result{}, nil).
// 					AnyTimes()
// 				mr.EXPECT().
// 					LookupMX(gomock.Any(), gomock.Any()).
// 					Return([]*net.MX{{Host: "mail.example.com", Pref: 10}}, adns.Result{}, nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "concurrent lookups",
// 			domain: dns.Domain{
// 				ASCII: "example.com",
// 			},
// 			setupResolver: func(mr *MockResolver, ctrl *gomock.Controller) {
// 				mr.EXPECT().
// 					LookupCNAME(gomock.Any(), gomock.Any()).
// 					Return("", adns.Result{}, nil).
// 					AnyTimes()
// 				mr.EXPECT().
// 					LookupMX(gomock.Any(), gomock.Any()).
// 					Return([]*net.MX{{Host: "mail.example.com", Pref: 10}}, adns.Result{}, nil).
// 					AnyTimes()
// 			},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			resolver := NewMockResolver(ctrl)
// 			tt.setupResolver(resolver, ctrl)

// 			r := NewResolver(ctx, resolver, slogger)

// 			if tt.name == "concurrent lookups" {
// 				// Test concurrent lookups
// 				const numGoroutines = 10
// 				results := make(chan error, numGoroutines)

// 				for i := 0; i < numGoroutines; i++ {
// 					go func() {
// 						_, err := r.LookupMX(ctx, tt.domain)
// 						results <- err
// 					}()
// 				}

// 				for i := 0; i < numGoroutines; i++ {
// 					err := <-results
// 					if tt.wantErr {
// 						assert.Error(t, err)
// 					} else {
// 						assert.NoError(t, err)
// 					}
// 				}
// 			} else {
// 				_, err := r.LookupMX(ctx, tt.domain)
// 				if tt.wantErr {
// 					assert.Error(t, err)
// 				} else {
// 					assert.NoError(t, err)
// 				}
// 			}
// 		})
// 	}
// }

func TestResolver_ContextCancellation(t *testing.T) {
	ctx, _ := telemetry.InitLogger(context.Background())
	slogger := telemetry.GetSLogger(ctx)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resolver := NewMockResolver(ctrl)
	resolver.EXPECT().
		LookupCNAME(gomock.Any(), gomock.Any()).
		Return("", adns.Result{}, nil).
		AnyTimes()
	r := NewResolver(ctx, resolver, slogger)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	domain := dns.Domain{ASCII: "example.com"}
	_, err := r.LookupMX(ctx, domain)
	assert.Error(t, err)
}

func TestResolver_Timeout(t *testing.T) {
	ctx, _ := telemetry.InitLogger(context.Background())
	slogger := telemetry.GetSLogger(ctx)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resolver := NewMockResolver(ctrl)
	resolver.EXPECT().
		LookupCNAME(gomock.Any(), gomock.Any()).
		Return("", adns.Result{}, nil).
		AnyTimes()
	r := NewResolver(ctx, resolver, slogger)

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	domain := dns.Domain{ASCII: "example.com"}
	_, err := r.LookupMX(ctx, domain)
	assert.Error(t, err)
}
