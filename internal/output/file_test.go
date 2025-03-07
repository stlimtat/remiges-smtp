package output

import (
	"context"
	"encoding/csv"
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
	tempDir, err := os.MkdirTemp("", "remiges-smtp-output-file-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	tests := []struct {
		name         string
		cfg          config.OutputConfig
		wantInitErr  bool
		wantWriteErr bool
	}{
		{
			name: "happy - date",
			cfg: config.OutputConfig{
				Type: config.ConfigOutputTypeFile,
				Args: map[string]any{
					config.ConfigArgPath:         tempDir,
					config.ConfigArgFileNameType: config.ConfigArgFileNameTypeDate,
				},
			},
			wantInitErr:  false,
			wantWriteErr: false,
		},
		{
			name: "happy - mail_id",
			cfg: config.OutputConfig{
				Type: config.ConfigOutputTypeFile,
				Args: map[string]any{
					config.ConfigArgPath:         tempDir,
					config.ConfigArgFileNameType: config.ConfigArgFileNameTypeMailID,
				},
			},
			wantInitErr:  false,
			wantWriteErr: false,
		},
		{
			name: "happy - hour",
			cfg: config.OutputConfig{
				Type: config.ConfigOutputTypeFile,
				Args: map[string]any{
					config.ConfigArgPath:         tempDir,
					config.ConfigArgFileNameType: config.ConfigArgFileNameTypeHour,
				},
			},
			wantInitErr:  false,
			wantWriteErr: false,
		},
		{
			name: "happy - quarter_hour",
			cfg: config.OutputConfig{
				Type: config.ConfigOutputTypeFile,
				Args: map[string]any{
					config.ConfigArgPath:         tempDir,
					config.ConfigArgFileNameType: config.ConfigArgFileNameTypeQuarterHour,
				},
			},
			wantInitErr:  false,
			wantWriteErr: false,
		},
		{
			name: "happy - date by default",
			cfg: config.OutputConfig{
				Type: config.ConfigOutputTypeFile,
				Args: map[string]any{
					config.ConfigArgPath: tempDir,
				},
			},
			wantInitErr:  false,
			wantWriteErr: false,
		},
		{
			name: "error - invalid path",
			cfg: config.OutputConfig{
				Type: config.ConfigOutputTypeFile,
				Args: map[string]any{
					config.ConfigArgPath: filepath.Join(tempDir, "invalid"),
				},
			},
			wantInitErr:  true,
			wantWriteErr: false,
		},
		{
			name: "happy - invalid file name type",
			cfg: config.OutputConfig{
				Type: config.ConfigOutputTypeFile,
				Args: map[string]any{
					config.ConfigArgPath:         tempDir,
					config.ConfigArgFileNameType: "invalid",
				},
			},
			wantInitErr:  false,
			wantWriteErr: false,
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
			mail := &pmail.Mail{
				MsgID: []byte(msgID),
			}

			filePath := fo.GetFileName(ctx, mail)
			err = fo.Write(
				ctx,
				mail,
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
