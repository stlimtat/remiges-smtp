package output

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mjl-/mox/smtpclient"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestFileTrackerOutput_Write(t *testing.T) {
	tests := []struct {
		name           string
		fileTrackerErr error
		wantErr        bool
	}{
		{
			name:           "happy path",
			fileTrackerErr: nil,
			wantErr:        false,
		},
		{
			name:           "file tracker error",
			fileTrackerErr: errors.New("tracker error"),
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx, _ := telemetry.InitLogger(context.Background())
			msgID := uuid.New().String()

			// Create mock file tracker
			mockTracker := file.NewMockIFileReadTracker(ctrl)
			mockTracker.EXPECT().
				UpsertFile(ctx, msgID, input.FILE_STATUS_DONE).
				Return(tt.fileTrackerErr)

			// Create file tracker output
			fto, err := NewFileTrackerOutput(
				ctx,
				config.OutputConfig{
					Type: config.ConfigOutputTypeFileTracker,
				},
				mockTracker,
			)
			require.NoError(t, err)

			// Test Write method
			err = fto.Write(
				ctx,
				&file.FileInfo{ID: msgID},
				&pmail.Mail{
					MsgID: []byte(msgID),
				},
				[]pmail.Response{
					{
						Response: smtpclient.Response{
							Code: 250,
							Line: "250 2.0.0 OK",
						},
					},
				},
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

// func TestFileTrackerOutput_Write_EdgeCases(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fileInfo *file.FileInfo
// 		mail     *pmail.Mail
// 		wantErr  bool
// 	}{
// 		{
// 			name:     "nil file info",
// 			fileInfo: nil,
// 			mail:     &pmail.Mail{MsgID: []byte("test")},
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "nil mail",
// 			fileInfo: &file.FileInfo{ID: "test"},
// 			mail:     nil,
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "empty file ID",
// 			fileInfo: &file.FileInfo{ID: ""},
// 			mail:     &pmail.Mail{MsgID: []byte("test")},
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "empty mail ID",
// 			fileInfo: &file.FileInfo{ID: "test"},
// 			mail:     &pmail.Mail{MsgID: []byte("")},
// 			wantErr:  true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			ctx, _ := telemetry.InitLogger(context.Background())

// 			// Create mock file tracker
// 			mockTracker := file.NewMockIFileReadTracker(ctrl)
// 			if !tt.wantErr {
// 				mockTracker.EXPECT().
// 					UpsertFile(gomock.Any(), gomock.Any(), gomock.Any()).
// 					Return(nil)
// 			}

// 			// Create file tracker output
// 			fto, err := NewFileTrackerOutput(
// 				ctx,
// 				config.OutputConfig{
// 					Type: config.ConfigOutputTypeFileTracker,
// 				},
// 				mockTracker,
// 			)
// 			require.NoError(t, err)

// 			// Test Write method
// 			err = fto.Write(
// 				ctx,
// 				tt.fileInfo,
// 				tt.mail,
// 				[]pmail.Response{
// 					{
// 						Response: smtpclient.Response{
// 							Code: 250,
// 							Line: "250 2.0.0 OK",
// 						},
// 					},
// 				},
// 			)

// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}
// 			assert.NoError(t, err)
// 		})
// 	}
// }
