package input

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestReadNextMail(t *testing.T) {
	tests := []struct {
		name                string
		wantFile            *FileInfo
		wantMail            *mail.Mail
		wantRefreshListErr  bool
		wantReadNextFileErr bool
	}{
		{
			name: "happy",
			wantFile: &FileInfo{
				ID: "test1",
			},
			wantMail: &mail.Mail{
				From: "sender@example.com",
				To:   "recipient@example.com",
				Headers: map[string][]byte{
					"Test1": []byte("test1"),
				},
				Body: []byte("test1"),
			},
			wantRefreshListErr:  false,
			wantReadNextFileErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			ctx, _ = telemetry.InitLogger(ctx)
			fr := NewMockIFileReader(ctrl)
			fr.EXPECT().
				RefreshList(ctx).
				DoAndReturn(
					func(_ context.Context) ([]*FileInfo, error) {
						if tt.wantRefreshListErr {
							return nil, errors.New("test error")
						}
						return []*FileInfo{tt.wantFile}, nil
					},
				).AnyTimes()
			fr.EXPECT().
				ReadNextFile(ctx).
				DoAndReturn(
					func(_ context.Context) (*FileInfo, error) {
						if tt.wantReadNextFileErr {
							return nil, errors.New("test error")
						}
						return tt.wantFile, nil
					},
				).AnyTimes()
			mt := NewMockIMailTransformer(ctrl)
			mt.EXPECT().
				Transform(ctx, tt.wantFile).
				DoAndReturn(
					func(_ context.Context, _ *FileInfo) (*mail.Mail, error) {
						return tt.wantMail, nil
					},
				).AnyTimes()
			fs := NewFileService(ctx, 1, fr, mt, time.Second)
			go func() {
				err := fs.Run(ctx)
				assert.NoError(t, err)
			}()
			time.Sleep(5 * time.Second)
			cancel()
		})
	}
}
