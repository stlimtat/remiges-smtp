package sendmail

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/proxy"
)

func TestNewDefaultDialerFactory(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.DialerConfig
	}{
		{
			name: "with_direct_connection",
			cfg: config.DialerConfig{
				Timeout: 30 * time.Second,
			},
		},
		{
			name: "with_socks5_proxy",
			cfg: config.DialerConfig{
				Timeout: 30 * time.Second,
				Socks5:  "localhost:1080",
				Auth: proxy.Auth{
					User:     "testuser",
					Password: "testpass",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			factory := NewDefaultDialerFactory(ctx, tt.cfg)
			require.NotNil(t, factory)
			assert.Equal(t, tt.cfg, factory.cfg)
		})
	}
}

func TestNewDialer(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config.DialerConfig
		expectError bool
	}{
		{
			name: "direct_connection",
			cfg: config.DialerConfig{
				Timeout: 30 * time.Second,
			},
			expectError: false,
		},
		// {
		// 	name: "invalid_socks5_proxy",
		// 	cfg: config.DialerConfig{
		// 		Timeout: 30 * time.Second,
		// 		Socks5:  "invalid:proxy:1080", // Invalid proxy address
		// 		Auth: proxy.Auth{
		// 			User:     "testuser",
		// 			Password: "testpass",
		// 		},
		// 	},
		// 	expectError: true,
		// },
		{
			name: "zero_timeout",
			cfg: config.DialerConfig{
				Timeout: 0,
			},
			expectError: false,
		},
		{
			name: "negative_timeout",
			cfg: config.DialerConfig{
				Timeout: -1 * time.Second,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			factory := NewDefaultDialerFactory(ctx, tt.cfg)
			dialer, err := factory.NewDialer(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, dialer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dialer)

				// Verify the dialer type based on configuration
				if tt.cfg.Socks5 != "" {
					_, ok := dialer.(proxy.ContextDialer)
					assert.True(t, ok, "Expected SOCKS5 proxy dialer")
				} else {
					_, ok := dialer.(*net.Dialer)
					assert.True(t, ok, "Expected direct TCP dialer")
				}
			}
		})
	}
}

func TestDialerConnection(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config.DialerConfig
		target      string
		expectError bool
	}{
		{
			name: "invalid_target_address",
			cfg: config.DialerConfig{
				Timeout: 1 * time.Second,
			},
			target:      "invalid:address:25",
			expectError: true,
		},
		{
			name: "connection_timeout",
			cfg: config.DialerConfig{
				Timeout: 1 * time.Millisecond, // Very short timeout
			},
			target:      "example.com:25",
			expectError: true,
		},
		{
			name: "non_existent_host",
			cfg: config.DialerConfig{
				Timeout: 1 * time.Second,
			},
			target:      "non.existent.host:25",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			factory := NewDefaultDialerFactory(ctx, tt.cfg)
			dialer, err := factory.NewDialer(ctx)
			require.NoError(t, err)
			require.NotNil(t, dialer)

			conn, err := dialer.DialContext(ctx, TCPNetwork, tt.target)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, conn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, conn)
				conn.Close()
			}
		})
	}
}
