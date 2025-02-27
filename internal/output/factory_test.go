package output

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mjl-/mox/smtpclient"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestNewOutputs(t *testing.T) {
	tests := []struct {
		name         string
		cfgs         []config.OutputConfig
		wantLen      int
		wantInitErr  bool
		wantWriteErr bool
	}{
		{
			name: "happy - file",
			cfgs: []config.OutputConfig{
				{
					Type: config.ConfigOutputTypeFile,
					Args: map[string]any{
						config.ConfigArgPath: "/tmp",
					},
				},
			},
			wantLen:      1,
			wantInitErr:  false,
			wantWriteErr: false,
		},
		{
			name: "alternate - file path does not exist",
			cfgs: []config.OutputConfig{
				{
					Type: config.ConfigOutputTypeFile,
					Args: map[string]any{
						config.ConfigArgPath: "/tmp/does-not-exist",
					},
				},
			},
			wantLen:      1,
			wantInitErr:  true,
			wantWriteErr: false,
		},
		{
			name:         "alternate - configs empty",
			cfgs:         []config.OutputConfig{},
			wantLen:      0,
			wantInitErr:  true,
			wantWriteErr: false,
		},
		{
			name: "alternate - invalid config",
			cfgs: []config.OutputConfig{
				{
					Type: "invalid",
				},
			},
			wantLen:      0,
			wantInitErr:  true,
			wantWriteErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx, _ := telemetry.InitLogger(context.Background())

			msgID := uuid.New().String()

			factory := OutputFactory{}
			got, err := factory.NewOutputs(ctx, tt.cfgs)
			if tt.wantInitErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, tt.wantLen)

			// Replacing the outputs with our mock
			myOutput := NewMockIOutput(ctrl)
			myOutput.EXPECT().
				Write(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, _ *pmail.Mail, _ []pmail.Response) error {
					if tt.wantWriteErr {
						err = errors.New("test error")
					}
					return err
				})
			factory.Outputs = []IOutput{myOutput}

			err = factory.Write(
				ctx,
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
			if tt.wantWriteErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
