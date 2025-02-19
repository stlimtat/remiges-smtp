package file_mail

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestReadNextMail(t *testing.T) {
	tests := []struct {
		name                string
		wantFile            *file.FileInfo
		wantMail            *mail.Mail
		wantRefreshListErr  bool
		wantReadNextFileErr bool
	}{
		{
			name: "happy",
			wantFile: &file.FileInfo{
				ID: "test1",
			},
			wantMail: &mail.Mail{
				From: smtp.Address{
					Localpart: "sender",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				To: []smtp.Address{
					{Localpart: "recipient", Domain: dns.Domain{ASCII: "example.com"}},
				},
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
			fr := file.NewMockIFileReader(ctrl)
			fr.EXPECT().
				RefreshList(ctx).
				DoAndReturn(
					func(_ context.Context) ([]*file.FileInfo, error) {
						if tt.wantRefreshListErr {
							return nil, errors.New("test error")
						}
						return []*file.FileInfo{tt.wantFile}, nil
					},
				).AnyTimes()
			fr.EXPECT().
				ReadNextFile(ctx).
				DoAndReturn(
					func(_ context.Context) (*file.FileInfo, error) {
						if tt.wantReadNextFileErr {
							return nil, errors.New("test error")
						}
						return tt.wantFile, nil
					},
				).AnyTimes()
			mt := NewMockIMailTransformer(ctrl)
			mt.EXPECT().
				Transform(ctx, tt.wantFile, gomock.Any()).
				DoAndReturn(
					func(_ context.Context, _ *file.FileInfo, _ *mail.Mail) (*mail.Mail, error) {
						return tt.wantMail, nil
					},
				).AnyTimes()
			fs := NewFileMailService(ctx, 1, fr, mt, time.Second)
			go func() {
				err := fs.Run(ctx)
				assert.NoError(t, err)
			}()
			time.Sleep(5 * time.Second)
			cancel()
		})
	}
}
