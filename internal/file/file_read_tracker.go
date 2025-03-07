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

type FileReadTracker struct {
	redisClient *redis.Client
}

func NewFileReadTracker(
	_ context.Context,
	redisClient *redis.Client,
) *FileReadTracker {
	return &FileReadTracker{redisClient: redisClient}
}

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
