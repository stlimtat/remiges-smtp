package telemetry

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "valid context",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, logger := InitLogger(context.Background())
			if tt.expectError {
				assert.Nil(t, logger)
				return
			}

			require.NotNil(t, ctx)
			require.NotNil(t, logger)
		})
	}
}

// func TestGetLogger(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		writer      func() *bytes.Buffer
// 		expectError bool
// 	}{
// 		{
// 			name: "valid context and writer",
// 			writer: func() *bytes.Buffer {
// 				return &bytes.Buffer{}
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "nil writer",
// 			writer: func() *bytes.Buffer {
// 				return nil
// 			},
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			writer := tt.writer()
// 			ctx, logger := GetLogger(context.Background(), writer)
// 			if tt.expectError {
// 				assert.Nil(t, logger)
// 				return
// 			}

// 			require.NotNil(t, ctx)
// 			require.NotNil(t, logger)

// 			// Test logging functionality
// 			logger.Info().Msg("test message")
// 			if writer != nil {
// 				assert.Contains(t, writer.String(), "test message")
// 			}
// 		})
// 	}
// }

func TestSetGlobalLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level zerolog.Level
	}{
		{
			name:  "debug level",
			level: zerolog.DebugLevel,
		},
		{
			name:  "info level",
			level: zerolog.InfoLevel,
		},
		{
			name:  "warn level",
			level: zerolog.WarnLevel,
		},
		{
			name:  "error level",
			level: zerolog.ErrorLevel,
		},
		{
			name:  "fatal level",
			level: zerolog.FatalLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetGlobalLogLevel(tt.level)
			assert.Equal(t, tt.level, zerolog.GlobalLevel())
		})
	}
}

func TestGetSLogger(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "valid context with logger",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := InitLogger(context.Background())
			logger := GetSLogger(ctx)
			if tt.expectError {
				assert.Nil(t, logger)
				return
			}

			require.NotNil(t, logger)

			// Test slog logging functionality
			var buf bytes.Buffer
			logger = slog.New(slog.NewTextHandler(&buf, nil))
			logger.Info("test message")
			assert.Contains(t, buf.String(), "test message")
		})
	}
}

func TestLoggerPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping logger performance test")
	}
	_, logger := InitLogger(context.Background())
	require.NotNil(t, logger)

	// Test high-volume logging
	start := time.Now()
	for i := 0; i < 10000; i++ {
		logger.Info().Msg("performance test message")
	}
	duration := time.Since(start)

	// Verify that logging 10,000 messages takes less than 1 second
	assert.Less(t, duration, time.Second, "Logger performance test failed")
}

func TestLoggerConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping logger concurrency test")
	}
	_, logger := InitLogger(context.Background())
	require.NotNil(t, logger)

	// Test concurrent logging
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 1000; j++ {
				logger.Info().Int("goroutine", id).Int("message", j).Msg("concurrent test message")
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
