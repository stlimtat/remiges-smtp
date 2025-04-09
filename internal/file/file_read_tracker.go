// Package file provides file reading and tracking functionality for mail processing.
// It includes implementations for tracking file processing states, reading file contents,
// and managing file processing workflows.
package file

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/pkg/input"
)

// FileReadTracker manages the state of file processing using Redis as a persistent store.
// It tracks which files have been read and their current processing status, ensuring
// that files are not processed multiple times and maintaining the processing state
// across application restarts.
//
// The tracker uses Redis keys with a 6-hour TTL to prevent stale entries from accumulating.
// Each file is identified by a unique ID, and its status is stored as an integer
// representing the FileStatus enum.
type FileReadTracker struct {
	redisClient *redis.Client
}

// NewFileReadTracker creates a new instance of FileReadTracker with the provided Redis client.
// The Redis client is used to persist file processing states with a 6-hour TTL.
//
// Parameters:
//   - ctx: Context for initialization (currently unused but reserved for future use)
//   - redisClient: The Redis client to use for state persistence
//
// Returns:
//   - *FileReadTracker: A new tracker instance configured with the provided Redis client
func NewFileReadTracker(
	_ context.Context,
	redisClient *redis.Client,
) *FileReadTracker {
	return &FileReadTracker{redisClient: redisClient}
}

// FileRead retrieves the current processing status of a file by its ID.
// It queries Redis for the file's status and returns the appropriate FileStatus.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - id: The unique identifier of the file to check
//
// Returns:
//   - input.FileStatus: The current processing status of the file
//   - error: Non-nil if the Redis operation fails
//
// Possible return values:
//   - FILE_STATUS_NOT_FOUND: If the file hasn't been tracked yet
//   - FILE_STATUS_ERROR: If Redis operations fail
//   - Other FileStatus values: The actual processing status of the file
func (f *FileReadTracker) FileRead(
	ctx context.Context, id string,
) (input.FileStatus, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Str("id", id).
		Msg("FileRead")
	getResult := f.redisClient.Get(ctx, "read_tracker_"+id)
	if getResult.Err() != nil {
		if errors.Is(getResult.Err(), redis.Nil) {
			return input.FILE_STATUS_NOT_FOUND, nil
		}
		return input.FILE_STATUS_ERROR, getResult.Err()
	}
	getResultInt, err := strconv.ParseInt(getResult.Val(), 10, 8)
	if err != nil {
		return input.FILE_STATUS_ERROR, err
	}
	return input.FileStatus(getResultInt), nil
}

// UpsertFile updates or inserts a file's processing status in Redis.
// It first checks if the status has changed before performing the update
// to prevent unnecessary Redis operations.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - id: The unique identifier of the file
//   - status: The new processing status to set
//
// Returns:
//   - error: Non-nil if:
//   - The status hasn't changed (returns "key already exists")
//   - Redis operations fail
//   - The provided status is invalid
//
// The status is stored with a 6-hour TTL to prevent stale entries from accumulating.
func (f *FileReadTracker) UpsertFile(
	ctx context.Context,
	id string,
	status input.FileStatus,
) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Str("id", id).
		Int("status", int(status)).
		Msg("UpsertFile")

	gotStatus, err := f.FileRead(ctx, id)
	if err != nil {
		logger.Error().Err(err).Msg("UpsertFile: FileRead")
		return err
	}
	if gotStatus == status {
		logger.Debug().Msg("UpsertFile: key already exists")
		return fmt.Errorf("key already exists")
	}
	setResult := f.redisClient.Set(
		ctx,
		"read_tracker_"+id,
		int(status),
		6*time.Hour,
	)
	if setResult.Err() != nil {
		logger.Error().Err(setResult.Err()).Msg("UpsertFile: setResult")
		return setResult.Err()
	}
	return nil
}
