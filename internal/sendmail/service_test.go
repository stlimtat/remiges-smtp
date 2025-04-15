package sendmail

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/file_mail"
	"github.com/stlimtat/remiges-smtp/internal/intmail"
	"github.com/stlimtat/remiges-smtp/internal/output"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewSendMailService(t *testing.T) {
	tests := []struct {
		name        string
		concurrency int
		interval    time.Duration
		wantNil     bool
	}{
		{
			name:        "valid_configuration",
			concurrency: 5,
			interval:    time.Second,
			wantNil:     false,
		},
		{
			name:        "zero_concurrency",
			concurrency: 0,
			interval:    time.Second,
			wantNil:     false,
		},
		// {
		// 	name:        "zero_interval",
		// 	concurrency: 5,
		// 	interval:    0,
		// 	wantNil:     false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFileReader := file.NewMockIFileReader(ctrl)
			mockMailProcessor := intmail.NewMockIMailProcessor(ctrl)
			mockMailSender := NewMockIMailSender(ctrl)
			mockMailTransformer := file_mail.NewMockIMailTransformer(ctrl)
			mockOutput := output.NewMockIOutput(ctrl)

			service := NewSendMailService(
				context.Background(),
				tt.concurrency,
				mockFileReader,
				mockMailProcessor,
				mockMailSender,
				mockMailTransformer,
				mockOutput,
				tt.interval,
			)

			if tt.wantNil {
				assert.Nil(t, service)
			} else {
				assert.NotNil(t, service)
				assert.Equal(t, tt.concurrency, service.Concurrency)
				assert.Equal(t, tt.interval, service.PollInterval)
				assert.NotNil(t, service.ticker)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*file.MockIFileReader)
		contextTimeout time.Duration
		expectError    bool
	}{
		// {
		// 	name: "successful_run",
		// 	setupMocks: func(fr *file.MockIFileReader) {
		// 		fr.EXPECT().
		// 			RefreshList(gomock.Any()).
		// 			Return([]string{}, nil).
		// 			Times(1)
		// 	},
		// 	contextTimeout: 100 * time.Millisecond,
		// 	expectError:    false,
		// },
		// {
		// 	name: "refresh_list_error",
		// 	setupMocks: func(fr *file.MockIFileReader) {
		// 		fr.EXPECT().
		// 			RefreshList(gomock.Any()).
		// 			Return(nil, errors.New("refresh failed")).
		// 			Times(1)
		// 	},
		// 	contextTimeout: 100 * time.Millisecond,
		// 	expectError:    false,
		// },
		{
			name: "context_cancelled",
			setupMocks: func(fr *file.MockIFileReader) {
				// No mock setup needed as context will be cancelled
			},
			contextTimeout: 0,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFileReader := file.NewMockIFileReader(ctrl)
			mockMailProcessor := intmail.NewMockIMailProcessor(ctrl)
			mockMailSender := NewMockIMailSender(ctrl)
			mockMailTransformer := file_mail.NewMockIMailTransformer(ctrl)
			mockOutput := output.NewMockIOutput(ctrl)

			tt.setupMocks(mockFileReader)

			service := NewSendMailService(
				context.Background(),
				1,
				mockFileReader,
				mockMailProcessor,
				mockMailSender,
				mockMailTransformer,
				mockOutput,
				50*time.Millisecond,
			)

			ctx := context.Background()
			if tt.contextTimeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.contextTimeout)
				defer cancel()
			} else {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			err := service.Run(ctx)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadNextMail(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*file.MockIFileReader, *file_mail.MockIMailTransformer, *intmail.MockIMailProcessor, *MockIMailSender, *output.MockIOutput)
		expectError bool
		expectNil   bool
	}{
		{
			name: "successful_processing",
			setupMocks: func(fr *file.MockIFileReader, mt *file_mail.MockIMailTransformer, mp *intmail.MockIMailProcessor, ms *MockIMailSender, mo *output.MockIOutput) {
				fileInfo := &file.FileInfo{ID: "test-id"}
				mail := &pmail.Mail{}

				fr.EXPECT().
					ReadNextFile(gomock.Any()).
					Return(fileInfo, nil).
					Times(1)

				mt.EXPECT().
					Transform(gomock.Any(), fileInfo, gomock.Any()).
					Return(mail, nil).
					Times(1)

				mp.EXPECT().
					Process(gomock.Any(), mail).
					Return(mail, nil).
					Times(1)

				ms.EXPECT().
					SendMail(gomock.Any(), mail).
					Return(map[string][]pmail.Response{}, nil).
					Times(1)

				mo.EXPECT().
					Write(gomock.Any(), fileInfo, mail, map[string][]pmail.Response{}).
					Return(nil).
					Times(1)
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name: "no_file_available",
			setupMocks: func(fr *file.MockIFileReader, mt *file_mail.MockIMailTransformer, mp *intmail.MockIMailProcessor, ms *MockIMailSender, mo *output.MockIOutput) {
				fr.EXPECT().
					ReadNextFile(gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			expectError: false,
			expectNil:   true,
		},
		{
			name: "read_file_error",
			setupMocks: func(fr *file.MockIFileReader, mt *file_mail.MockIMailTransformer, mp *intmail.MockIMailProcessor, ms *MockIMailSender, mo *output.MockIOutput) {
				fr.EXPECT().
					ReadNextFile(gomock.Any()).
					Return(nil, errors.New("read failed")).
					Times(1)
			},
			expectError: true,
			expectNil:   true,
		},
		{
			name: "transform_error",
			setupMocks: func(fr *file.MockIFileReader, mt *file_mail.MockIMailTransformer, mp *intmail.MockIMailProcessor, ms *MockIMailSender, mo *output.MockIOutput) {
				fileInfo := &file.FileInfo{ID: "test-id"}

				fr.EXPECT().
					ReadNextFile(gomock.Any()).
					Return(fileInfo, nil).
					Times(1)

				mt.EXPECT().
					Transform(gomock.Any(), fileInfo, gomock.Any()).
					Return(nil, errors.New("transform failed")).
					Times(1)
			},
			expectError: true,
			expectNil:   true,
		},
		{
			name: "process_error",
			setupMocks: func(fr *file.MockIFileReader, mt *file_mail.MockIMailTransformer, mp *intmail.MockIMailProcessor, ms *MockIMailSender, mo *output.MockIOutput) {
				fileInfo := &file.FileInfo{ID: "test-id"}
				mail := &pmail.Mail{}

				fr.EXPECT().
					ReadNextFile(gomock.Any()).
					Return(fileInfo, nil).
					Times(1)

				mt.EXPECT().
					Transform(gomock.Any(), fileInfo, gomock.Any()).
					Return(mail, nil).
					Times(1)

				mp.EXPECT().
					Process(gomock.Any(), mail).
					Return(nil, errors.New("process failed")).
					Times(1)
			},
			expectError: true,
			expectNil:   true,
		},
		// {
		// 	name: "send_error",
		// 	setupMocks: func(fr *file.MockIFileReader, mt *file_mail.MockIMailTransformer, mp *intmail.MockIMailProcessor, ms *MockIMailSender) {
		// 		fileInfo := &file.FileInfo{ID: "test-id"}
		// 		mail := &pmail.Mail{}

		// 		fr.EXPECT().
		// 			ReadNextFile(gomock.Any()).
		// 			Return(fileInfo, nil).
		// 			Times(1)

		// 		mt.EXPECT().
		// 			Transform(gomock.Any(), fileInfo, gomock.Any()).
		// 			Return(mail, nil).
		// 			Times(1)

		// 		mp.EXPECT().
		// 			Process(gomock.Any(), mail).
		// 			Return(mail, nil).
		// 			Times(1)

		// 		ms.EXPECT().
		// 			SendMail(gomock.Any(), mail).
		// 			Return(nil, map[string]error{"test": errors.New("send failed")}).
		// 			Times(1)
		// 	},
		// 	expectError: true,
		// 	expectNil:   true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFileReader := file.NewMockIFileReader(ctrl)
			mockMailProcessor := intmail.NewMockIMailProcessor(ctrl)
			mockMailSender := NewMockIMailSender(ctrl)
			mockMailTransformer := file_mail.NewMockIMailTransformer(ctrl)
			mockOutput := output.NewMockIOutput(ctrl)

			tt.setupMocks(mockFileReader, mockMailTransformer, mockMailProcessor, mockMailSender, mockOutput)

			service := NewSendMailService(
				context.Background(),
				1,
				mockFileReader,
				mockMailProcessor,
				mockMailSender,
				mockMailTransformer,
				mockOutput,
				time.Second,
			)

			fileInfo, mail, err := service.ReadNextMail(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, mail)
			} else {
				assert.NoError(t, err)
				if tt.expectNil {
					assert.Nil(t, fileInfo)
					assert.Nil(t, mail)
				} else {
					assert.NotNil(t, fileInfo)
					assert.NotNil(t, mail)
					assert.Equal(t, input.FILE_STATUS_DELIVERED, fileInfo.Status)
				}
			}
		})
	}
}

func TestProcessFileLoop(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*file.MockIFileReader)
		contextTimeout time.Duration
	}{
		{
			name: "process_files_until_context_done",
			setupMocks: func(fr *file.MockIFileReader) {
				fr.EXPECT().
					ReadNextFile(gomock.Any()).
					Return(nil, nil).
					AnyTimes()
			},
			contextTimeout: 100 * time.Millisecond,
		},
		{
			name: "immediate_context_cancel",
			setupMocks: func(fr *file.MockIFileReader) {
				// No mock setup needed as context will be cancelled immediately
			},
			contextTimeout: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFileReader := file.NewMockIFileReader(ctrl)
			mockMailProcessor := intmail.NewMockIMailProcessor(ctrl)
			mockMailSender := NewMockIMailSender(ctrl)
			mockMailTransformer := file_mail.NewMockIMailTransformer(ctrl)
			mockOutput := output.NewMockIOutput(ctrl)

			tt.setupMocks(mockFileReader)

			service := NewSendMailService(
				context.Background(),
				1,
				mockFileReader,
				mockMailProcessor,
				mockMailSender,
				mockMailTransformer,
				mockOutput,
				50*time.Millisecond,
			)

			ctx := context.Background()
			if tt.contextTimeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.contextTimeout)
				defer cancel()
			} else {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			var wg sync.WaitGroup
			wg.Add(1)
			service.ProcessFileLoop(ctx, &wg)
			wg.Wait()
		})
	}
}
