package file

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileReadTracker(t *testing.T) {
	tests := []struct {
		name        string
		getErr      error
		setErr      error
		want        input.FileStatus
		wantGetErr  bool
		wantSetErr  bool
		wantSet2Err bool
		wantErr     bool
	}{
		{"happy", nil, nil, input.FILE_STATUS_INIT, false, false, true, false},
		{"set_err", nil, fmt.Errorf("set_err"), input.FILE_STATUS_INIT, false, true, true, false},
		{"get_err", fmt.Errorf("get_err"), nil, input.FILE_STATUS_INIT, true, false, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.InitLogger(ctx)

			client, mock := redismock.NewClientMock()

			expectGet := mock.ExpectGet("read_tracker_123")
			expectGet.SetErr(redis.Nil)
			expectSet := mock.ExpectSet(
				"read_tracker_123",
				int(tt.want),
				6*time.Hour,
			)
			if tt.setErr != nil {
				expectSet.SetErr(tt.setErr)
			} else {
				expectSet.SetVal("OK")
			}

			frt := NewFileReadTracker(ctx, client)
			err := frt.UpsertFile(ctx, "123", tt.want)
			if tt.wantSetErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			expectGet = mock.ExpectGet("read_tracker_123")
			if tt.wantGetErr {
				expectGet.SetErr(tt.getErr)
			} else {
				expectGet.SetVal(strconv.Itoa(int(tt.want)))
			}
			got, err := frt.FileRead(ctx, "123")
			if tt.wantGetErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)

			expectGet = mock.ExpectGet("read_tracker_123")
			expectGet.SetVal(strconv.Itoa(int(tt.want)))
			err = frt.UpsertFile(ctx, "123", tt.want)
			if !tt.wantSet2Err {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
		})
	}
}
