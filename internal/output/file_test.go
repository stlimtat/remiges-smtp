package output

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mjl-/mox/smtpclient"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileOutput_Write(t *testing.T) {
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
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
					config.ConfigArgPath:         tmpDir,
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
					config.ConfigArgPath:         tmpDir,
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
					config.ConfigArgPath:         tmpDir,
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
					config.ConfigArgPath:         tmpDir,
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
					config.ConfigArgPath: tmpDir,
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
					config.ConfigArgPath: filepath.Join(tmpDir, "invalid"),
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
					config.ConfigArgPath:         tmpDir,
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
				&file.FileInfo{ID: msgID},
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
			generatedFile, err := os.Open(filePath)
			require.NoError(t, err)
			csvReader := csv.NewReader(generatedFile)
			content, err := csvReader.ReadAll()
			require.NoError(t, err)
			assert.Equal(t, []string{"msg_id", "status", "error"}, content[0])
			assert.Equal(t, []string{msgID, "250", "250 2.0.0 OK"}, content[1])
		})
	}
}

// func TestFileOutput_Write_EdgeCases(t *testing.T) {
// 	tmpDir := t.TempDir()
// 	defer os.RemoveAll(tmpDir)

// 	tests := []struct {
// 		name         string
// 		fileInfo     *file.FileInfo
// 		mail         *pmail.Mail
// 		responses    []pmail.Response
// 		wantWriteErr bool
// 	}{
// 		{
// 			name:         "nil file info",
// 			fileInfo:     nil,
// 			mail:         &pmail.Mail{MsgID: []byte("test")},
// 			responses:    []pmail.Response{{Response: smtpclient.Response{Code: 250, Line: "OK"}}},
// 			wantWriteErr: true,
// 		},
// 		{
// 			name:         "nil mail",
// 			fileInfo:     &file.FileInfo{ID: "test"},
// 			mail:         nil,
// 			responses:    []pmail.Response{{Response: smtpclient.Response{Code: 250, Line: "OK"}}},
// 			wantWriteErr: true,
// 		},
// 		{
// 			name:         "empty file ID",
// 			fileInfo:     &file.FileInfo{ID: ""},
// 			mail:         &pmail.Mail{MsgID: []byte("test")},
// 			responses:    []pmail.Response{{Response: smtpclient.Response{Code: 250, Line: "OK"}}},
// 			wantWriteErr: true,
// 		},
// 		{
// 			name:         "empty mail ID",
// 			fileInfo:     &file.FileInfo{ID: "test"},
// 			mail:         &pmail.Mail{MsgID: []byte("")},
// 			responses:    []pmail.Response{{Response: smtpclient.Response{Code: 250, Line: "OK"}}},
// 			wantWriteErr: true,
// 		},
// 		{
// 			name:         "nil responses",
// 			fileInfo:     &file.FileInfo{ID: "test"},
// 			mail:         &pmail.Mail{MsgID: []byte("test")},
// 			responses:    nil,
// 			wantWriteErr: true,
// 		},
// 		{
// 			name:         "empty responses",
// 			fileInfo:     &file.FileInfo{ID: "test"},
// 			mail:         &pmail.Mail{MsgID: []byte("test")},
// 			responses:    []pmail.Response{},
// 			wantWriteErr: true,
// 		},
// 		{
// 			name:     "multiple responses",
// 			fileInfo: &file.FileInfo{ID: "test"},
// 			mail:     &pmail.Mail{MsgID: []byte("test")},
// 			responses: []pmail.Response{
// 				{Response: smtpclient.Response{Code: 250, Line: "OK"}},
// 				{Response: smtpclient.Response{Code: 550, Line: "Error"}},
// 			},
// 			wantWriteErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx, _ := telemetry.InitLogger(context.Background())
// 			fo, err := NewFileOutput(ctx, config.OutputConfig{
// 				Type: config.ConfigOutputTypeFile,
// 				Args: map[string]any{
// 					config.ConfigArgPath:         tmpDir,
// 					config.ConfigArgFileNameType: config.ConfigArgFileNameTypeDate,
// 				},
// 			})
// 			require.NoError(t, err)

// 			err = fo.Write(ctx, tt.fileInfo, tt.mail, tt.responses)
// 			if tt.wantWriteErr {
// 				assert.Error(t, err)
// 				return
// 			}
// 			assert.NoError(t, err)

// 			// Verify file contents if write was successful
// 			if !tt.wantWriteErr {
// 				filePath := fo.GetFileName(ctx, tt.mail)
// 				generatedFile, err := os.Open(filePath)
// 				require.NoError(t, err)
// 				defer generatedFile.Close()

// 				csvReader := csv.NewReader(generatedFile)
// 				content, err := csvReader.ReadAll()
// 				require.NoError(t, err)

// 				// Check header
// 				assert.Equal(t, []string{"msg_id", "status", "error"}, content[0])

// 				// Check each response
// 				for i, resp := range tt.responses {
// 					assert.Equal(t, []string{
// 						string(tt.mail.MsgID),
// 						fmt.Sprintf("%d", resp.Code),
// 						resp.Line,
// 					}, content[i+1])
// 				}
// 			}
// 		})
// 	}
// }

func TestFileOutput_GetFileName(t *testing.T) {
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	now := time.Now()
	expectedDate := now.Format("2006-01-02")
	expectedHour := now.Format("2006-01-02-15")
	expectedQuarter := fmt.Sprintf("%s-%d", expectedHour, now.Minute()/15)
	msgID := "test-msg-id"

	tests := []struct {
		name           string
		fileNameType   string
		mail           *pmail.Mail
		expectedPrefix string
	}{
		{
			name:           "date format",
			fileNameType:   config.ConfigArgFileNameTypeDate,
			mail:           &pmail.Mail{MsgID: []byte(msgID)},
			expectedPrefix: expectedDate,
		},
		{
			name:           "hour format",
			fileNameType:   config.ConfigArgFileNameTypeHour,
			mail:           &pmail.Mail{MsgID: []byte(msgID)},
			expectedPrefix: expectedHour,
		},
		{
			name:           "quarter hour format",
			fileNameType:   config.ConfigArgFileNameTypeQuarterHour,
			mail:           &pmail.Mail{MsgID: []byte(msgID)},
			expectedPrefix: expectedQuarter,
		},
		{
			name:           "mail ID format",
			fileNameType:   config.ConfigArgFileNameTypeMailID,
			mail:           &pmail.Mail{MsgID: []byte(msgID)},
			expectedPrefix: msgID,
		},
		{
			name:           "invalid format falls back to date",
			fileNameType:   "invalid",
			mail:           &pmail.Mail{MsgID: []byte(msgID)},
			expectedPrefix: expectedDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			fo, err := NewFileOutput(ctx, config.OutputConfig{
				Type: config.ConfigOutputTypeFile,
				Args: map[string]any{
					config.ConfigArgPath:         tmpDir,
					config.ConfigArgFileNameType: tt.fileNameType,
				},
			})
			require.NoError(t, err)

			fileName := fo.GetFileName(ctx, tt.mail)
			assert.Contains(t, fileName, tt.expectedPrefix)
			// assert.Contains(t, fileName, DEFAULT_FILE_NAME)
			assert.Contains(t, fileName, tmpDir)
		})
	}
}
