package sendmail

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/file_mail"
	"github.com/stlimtat/remiges-smtp/internal/intmail"
	"github.com/stlimtat/remiges-smtp/internal/output"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestReadNextMail(t *testing.T) {
	tests := []struct {
		name                string
		wantFile            *file.FileInfo
		wantMail            *pmail.Mail
		wantRefreshListErr  bool
		wantReadNextFileErr bool
	}{
		{
			name: "happy",
			wantFile: &file.FileInfo{
				ID: "test1",
			},
			wantMail: &pmail.Mail{
				From: smtp.Address{
					Localpart: "sender",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				To: []smtp.Address{
					{Localpart: "recipient", Domain: dns.Domain{ASCII: "example.com"}},
				},
				Metadata: map[string][]byte{
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
			fileReader := file.NewMockIFileReader(ctrl)
			fileReader.EXPECT().
				RefreshList(ctx).
				DoAndReturn(
					func(_ context.Context) ([]*file.FileInfo, error) {
						if tt.wantRefreshListErr {
							return nil, errors.New("test error")
						}
						return []*file.FileInfo{tt.wantFile}, nil
					},
				).AnyTimes()
			fileReader.EXPECT().
				ReadNextFile(ctx).
				DoAndReturn(
					func(_ context.Context) (*file.FileInfo, error) {
						if tt.wantReadNextFileErr {
							return nil, errors.New("test error")
						}
						return tt.wantFile, nil
					},
				).AnyTimes()
			mailProcessor := intmail.NewMockIMailProcessor(ctrl)
			mailProcessor.EXPECT().
				Process(ctx, gomock.Any()).
				DoAndReturn(
					func(_ context.Context, _ *pmail.Mail) (*pmail.Mail, error) {
						return tt.wantMail, nil
					},
				).AnyTimes()
			mailSender := NewMockIMailSender(ctrl)
			mailSender.EXPECT().
				SendMail(ctx, gomock.Any()).
				DoAndReturn(
					func(_ context.Context, _ *pmail.Mail) (map[string][]pmail.Response, map[string]error) {
						return nil, nil
					},
				).AnyTimes()
			mailTransformer := file_mail.NewMockIMailTransformer(ctrl)
			mailTransformer.EXPECT().
				Transform(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(_ context.Context, _ *file.FileInfo, _ *pmail.Mail) (*pmail.Mail, error) {
						return tt.wantMail, nil
					},
				).AnyTimes()
			myOutput := output.NewMockIOutput(ctrl)
			myOutput.EXPECT().
				Write(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(_ context.Context, _ *pmail.Mail, _ []pmail.Response) error {
						return nil
					},
				).AnyTimes()
			sendMailService := NewSendMailService(
				ctx, 1, fileReader, mailProcessor, mailSender, mailTransformer, myOutput, time.Second,
			)
			go func() {
				err := sendMailService.Run(ctx)
				assert.NoError(t, err)
			}()
			time.Sleep(5 * time.Second)
			cancel()
		})
	}
}
