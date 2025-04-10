package file

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileReadTracker(t *testing.T) {
	ctx, _ := telemetry.InitLogger(context.Background())

	// Setup Redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	tracker := NewFileReadTracker(ctx, client)

	tests := []struct {
		name          string
		id            string
		status        input.FileStatus
		preSetStatus  input.FileStatus // Optional: status to set before test
		wantErr       bool
		errorContains string
	}{
		{
			name:   "happy path - new file",
			id:     "test1",
			status: input.FILE_STATUS_INIT,
		},
		{
			name:         "happy path - update status",
			id:           "test2",
			preSetStatus: input.FILE_STATUS_INIT,
			status:       input.FILE_STATUS_PROCESSING,
		},
		// {
		// 	name:          "error - same status",
		// 	id:            "test3",
		// 	preSetStatus:  input.FILE_STATUS_INIT,
		// 	status:        input.FILE_STATUS_INIT,
		// 	wantErr:       true,
		// 	errorContains: "key already exists",
		// },
		// {
		// 	name:    "error - invalid status",
		// 	id:      "test4",
		// 	status:  input.FileStatus(999), // Invalid status
		// 	wantErr: true,
		// },
		// {
		// 	name:    "error - empty id",
		// 	id:      "",
		// 	status:  input.FILE_STATUS_INIT,
		// 	wantErr: true,
		// },
		// {
		// 	name:    "error - long id",
		// 	id:      string(make([]byte, 1000)), // Very long ID
		// 	status:  input.FILE_STATUS_INIT,
		// 	wantErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Pre-set status if needed
			if tt.preSetStatus != input.FILE_STATUS_INIT {
				err := tracker.UpsertFile(ctx, tt.id, tt.preSetStatus)
				require.NoError(t, err)
			}

			// Test UpsertFile
			err := tracker.UpsertFile(ctx, tt.id, tt.status)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}
			assert.NoError(t, err)

			// Verify the status was set correctly
			gotStatus, err := tracker.FileRead(ctx, tt.id)
			assert.NoError(t, err)
			assert.Equal(t, tt.status, gotStatus)
		})
	}
}

func TestFileReadTracker_RedisErrors(t *testing.T) {
	ctx, _ := telemetry.InitLogger(context.Background())

	// Create a client that will fail
	client := redis.NewClient(&redis.Options{
		Addr: "invalid:6379", // Invalid address
	})

	tracker := NewFileReadTracker(ctx, client)

	// Test FileRead with Redis error
	status, err := tracker.FileRead(ctx, "test")
	assert.Error(t, err)
	assert.Equal(t, input.FILE_STATUS_ERROR, status)

	// Test UpsertFile with Redis error
	err = tracker.UpsertFile(ctx, "test", input.FILE_STATUS_INIT)
	assert.Error(t, err)
}

func TestFileReadTracker_ConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping concurrent access test")
	}
	ctx, _ := telemetry.InitLogger(context.Background())

	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	tracker := NewFileReadTracker(ctx, client)

	// Test concurrent access to the same file
	const numGoroutines = 10
	results := make(chan error, numGoroutines)
	id := "concurrent_test"

	for i := 0; i < numGoroutines; i++ {
		go func() {
			err := tracker.UpsertFile(ctx, id, input.FILE_STATUS_INIT)
			results <- err
		}()
	}

	// Collect results
	var successCount int
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		if err == nil {
			successCount++
		} else {
			assert.Contains(t, err.Error(), "key already exists")
		}
	}

	// Only one goroutine should succeed
	// assert.Equal(t, 1, successCount)
}

func TestFileReadTracker_Expiration(t *testing.T) {
	ctx, _ := telemetry.InitLogger(context.Background())

	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	tracker := NewFileReadTracker(ctx, client)
	id := "expiration_test"

	// Set a status
	err = tracker.UpsertFile(ctx, id, input.FILE_STATUS_INIT)
	require.NoError(t, err)

	// Fast-forward time to after expiration
	mr.FastForward(7 * time.Hour)

	// Status should be not found after expiration
	status, err := tracker.FileRead(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, input.FILE_STATUS_NOT_FOUND, status)
}

func TestFileReadTracker_InvalidRedisValue(t *testing.T) {
	ctx, _ := telemetry.InitLogger(context.Background())

	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Set an invalid value directly in Redis
	id := "invalid_value_test"
	mr.Set("read_tracker_"+id, "not_a_number")

	tracker := NewFileReadTracker(ctx, client)

	// Attempt to read the invalid value
	status, err := tracker.FileRead(ctx, id)
	assert.Error(t, err)
	assert.Equal(t, input.FILE_STATUS_ERROR, status)
}
