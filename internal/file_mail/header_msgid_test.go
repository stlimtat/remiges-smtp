package file_mail

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderMsgIDTransformer(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config.FileMailConfig
		headers     map[string][]byte
		wantMsgID   []byte
		wantMsgUuid bool
		wantErr     bool
	}{
		{
			name: "happy - default",
			cfg: config.FileMailConfig{
				Args: map[string]any{
					HeaderConfigArgType:    config.ConfigTypeDefaultStr,
					HeaderConfigArgDefault: "default msgid",
				},
			},
			headers: map[string][]byte{
				input.HeaderMsgIDKey: []byte("test msgid"),
			},
			wantMsgID:   []byte("default msgid"),
			wantMsgUuid: false,
			wantErr:     false,
		},
		{
			name: "happy - header",
			cfg: config.FileMailConfig{
				Args: map[string]any{
					HeaderConfigArgType: config.ConfigTypeHeadersStr,
				},
			},
			headers: map[string][]byte{
				input.HeaderMsgIDKey: []byte("test msgid"),
			},
			wantMsgID:   []byte("test msgid"),
			wantMsgUuid: false,
			wantErr:     false,
		},
		{
			name: "happy - uuid",
			cfg: config.FileMailConfig{
				Args: map[string]any{
					HeaderConfigArgType: HeaderMsgIDConfigArgUuid,
				},
			},
			headers: map[string][]byte{
				input.HeaderMsgIDKey: []byte("test msgid"),
			},
			wantMsgID:   []byte{},
			wantMsgUuid: true,
			wantErr:     false,
		},
		{
			name: "alternate - header not found",
			cfg: config.FileMailConfig{
				Args: map[string]any{
					HeaderConfigArgType: config.ConfigTypeHeadersStr,
				},
			},
			headers:     map[string][]byte{},
			wantMsgID:   []byte{},
			wantMsgUuid: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			transformer := &HeaderMsgIDTransformer{}
			err := transformer.Init(ctx, tt.cfg)
			require.NoError(t, err)
			gotMail, err := transformer.Transform(ctx, &file.FileInfo{}, &mail.Mail{
				Metadata: tt.headers,
			})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantMsgUuid {
				assert.NotNil(t, gotMail.MsgID)
				assert.Len(t, gotMail.MsgID, 16)
				assert.NotEqual(t, uuid.Nil, gotMail.MsgID)
				_, err = uuid.FromBytes(gotMail.MsgID)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.wantMsgID, gotMail.MsgID)
			}
		})
	}
}
