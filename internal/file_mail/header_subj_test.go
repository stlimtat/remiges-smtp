package file_mail

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderSubjectTransformer(t *testing.T) {
	tests := []struct {
		name             string
		cfg              config.FileMailConfig
		headers          map[string][]byte
		wantSubject      []byte
		wantInitErr      bool
		wantTransformErr bool
	}{
		{
			name: "happy",
			cfg:  config.FileMailConfig{},
			headers: map[string][]byte{
				"Subject": []byte("test subject"),
			},
			wantSubject:      []byte("test subject"),
			wantInitErr:      false,
			wantTransformErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())

			transformer := &HeaderSubjectTransformer{}
			err := transformer.Init(ctx, tt.cfg)
			if tt.wantInitErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			gotMail, err := transformer.Transform(ctx, nil, &mail.Mail{
				Metadata: tt.headers,
			})
			if tt.wantTransformErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantSubject, gotMail.Subject)
		})
	}
}
