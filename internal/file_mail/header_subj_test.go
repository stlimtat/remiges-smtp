package file_mail

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/mail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderSubjectTransformer(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config.FileMailConfig
		headers     map[string][]byte
		wantSubject []byte
		wantErr     bool
	}{
		{
			name: "happy",
			cfg:  config.FileMailConfig{},
			headers: map[string][]byte{
				input.HeaderSubjectKey: []byte("test subject"),
			},
			wantSubject: []byte("test subject"),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())

			transformer := &HeaderSubjectTransformer{}
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
			assert.Equal(t, tt.wantSubject, gotMail.Subject)
		})
	}
}
