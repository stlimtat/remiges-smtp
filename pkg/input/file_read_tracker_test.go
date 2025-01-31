package input

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileReadTracker(t *testing.T) {
	tests := []struct {
		name       string
		getErr     error
		setErr     error
		want       FileStatus
		wantGetErr bool
		wantSetErr bool
		wantErr    bool
	}{
		{"happy", nil, nil, FILE_STATUS_INIT, false, false, false},
		{"set_err", nil, fmt.Errorf("set_err"), FILE_STATUS_INIT, false, true, false},
		{"get_err", fmt.Errorf("get_err"), nil, FILE_STATUS_INIT, true, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.InitLogger(ctx)

			client, mock := redismock.NewClientMock()

			expectSet := mock.ExpectSet(
				"read_tracker_123",
				FILE_STATUS_INIT,
				6*time.Hour,
			)
			if tt.setErr != nil {
				expectSet.SetErr(tt.setErr)
			} else {
				expectSet.SetVal("OK")
			}
			expectGet := mock.ExpectGet("read_tracker_123")
			if tt.getErr != nil {
				expectGet.SetErr(tt.getErr)
			} else {
				expectGet.
					SetVal(strconv.Itoa(int(FILE_STATUS_INIT)))
			}

			frt := NewFileReadTracker(ctx, client)
			err := frt.UpsertFile(ctx, "123", FILE_STATUS_INIT)
			if tt.wantSetErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			got, err := frt.FileRead(ctx, "123")
			if tt.wantGetErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
