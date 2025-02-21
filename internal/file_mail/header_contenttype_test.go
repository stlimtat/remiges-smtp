package file_mail

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderContentTypeTransformer(t *testing.T) {
	tests := []struct {
		name            string
		cfg             config.FileMailConfig
		headers         map[string][]byte
		wantContentType []byte
		wantErr         bool
	}{
		{
			name: "happy",
			cfg: config.FileMailConfig{
				Type: HeaderContentTypeTransformerType,
				Args: map[string]any{},
			},
			headers: map[string][]byte{
				input.HeaderContentTypeKey: []byte("text/plain"),
			},
			wantContentType: []byte("text/plain"),
			wantErr:         false,
		},
		{
			name: "no content type",
			cfg: config.FileMailConfig{
				Type: HeaderContentTypeTransformerType,
				Args: map[string]any{},
			},
			headers:         map[string][]byte{},
			wantContentType: make([]byte, 0),
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			transformer := &HeaderContentTypeTransformer{}
			err := transformer.Init(ctx, tt.cfg)
			require.NoError(t, err)

			got, err := transformer.Transform(
				ctx,
				&file.FileInfo{},
				&mail.Mail{
					Metadata: tt.headers,
				},
			)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantContentType, got.ContentType)
		})
	}
}
