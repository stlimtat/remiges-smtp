package file_mail

import (
	"bytes"
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyTransformer(t *testing.T) {
	tests := []struct {
		name             string
		cfg              config.FileMailConfig
		body             []byte
		wantBody         []byte
		wantInitErr      bool
		wantTransformErr bool
	}{
		{
			name:             "happy",
			cfg:              config.FileMailConfig{},
			body:             []byte("test body"),
			wantBody:         []byte("test body"),
			wantInitErr:      false,
			wantTransformErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			transformer := &BodyTransformer{}
			err := transformer.Init(ctx, tt.cfg)
			if tt.wantInitErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			fileInfo := &file.FileInfo{
				DfReader: bytes.NewReader(tt.body),
			}
			gotMail, err := transformer.Transform(ctx, fileInfo, &mail.Mail{})
			if tt.wantTransformErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantBody, gotMail.Body)
		})
	}
}
