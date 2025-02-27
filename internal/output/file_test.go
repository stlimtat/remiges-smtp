package output

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/mjl-/mox/smtpclient"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileOutput_Write(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.OutputConfig
		wantInitErr  bool
		wantWriteErr bool
	}{
		{
			name: "happy",
			cfg: config.OutputConfig{
				Type: config.ConfigOutputTypeFile,
				Args: map[string]any{
					config.ConfigArgPath: "/tmp",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			fo, err := NewFileOutput(ctx, tt.cfg)
			if tt.wantInitErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			msgID := uuid.New().String()

			err = fo.Write(
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
			// check that the file was created
			filePath := filepath.Join(tt.cfg.Args[config.ConfigArgPath].(string), fmt.Sprintf(DEFAULT_FILE_NAME, msgID))
			_, err = os.Stat(filePath)
			require.NoError(t, err)
			// check that the file contains the correct data
			file, err := os.Open(filePath)
			require.NoError(t, err)
			csvReader := csv.NewReader(file)
			content, err := csvReader.ReadAll()
			require.NoError(t, err)
			assert.Equal(t, []string{"msg_id", "status", "error"}, content[0])
			assert.Equal(t, []string{msgID, "250", "250 2.0.0 OK"}, content[1])
		})
	}
}
