package cli

// func TestNewGenericSvc(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		setup       func(t *testing.T) (context.Context, func())
// 		expectError bool
// 	}{
// 		{
// 			name: "valid configuration",
// 			setup: func(t *testing.T) (context.Context, func()) {
// 				ctx := context.Background()
// 				ctx, _ = telemetry.InitLogger(ctx)

// 				// Setup Redis configuration
// 				viper.Set("read-file.redis-addr", "localhost:6379")
// 				viper.Set("read-file.in-path", "/tmp")
// 				viper.Set("read-file.file-mails", []config.FileMailConfig{})
// 				viper.Set("outputs", []config.OutputConfig{})
// 				viper.Set("mail-processors", []config.MailProcessorConfig{})
// 				viper.Set("debug", false)
// 				viper.Set("read-file.concurrency", 1)
// 				viper.Set("read-file.poll-interval", time.Second)

// 				cfg := config.SendMailConfig{
// 					ReadFileConfig: config.ReadFileConfig{
// 						RedisAddr:    "localhost:6379",
// 						InPath:       "/tmp",
// 						FileMails:    []config.FileMailConfig{},
// 						Concurrency:  1,
// 						PollInterval: time.Second,
// 					},
// 					Outputs:        []config.OutputConfig{},
// 					MailProcessors: []config.MailProcessorConfig{},
// 					Debug:          false,
// 				}
// 				ctx = config.SetContextConfig(ctx, cfg)

// 				cmd := &cobra.Command{}
// 				cmd.SetContext(ctx)

// 				cleanup := func() {
// 					viper.Reset()
// 				}

// 				return ctx, cleanup
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "invalid Redis configuration",
// 			setup: func(t *testing.T) (context.Context, func()) {
// 				ctx := context.Background()
// 				ctx, _ = telemetry.InitLogger(ctx)

// 				// Setup invalid Redis configuration
// 				viper.Set("read-file.redis-addr", "invalid:address")
// 				viper.Set("read-file.in-path", "/tmp")
// 				viper.Set("read-file.file-mails", []config.FileMailConfig{})
// 				viper.Set("outputs", []config.OutputConfig{})
// 				viper.Set("mail-processors", []config.MailProcessorConfig{})
// 				viper.Set("debug", false)
// 				viper.Set("read-file.concurrency", 1)
// 				viper.Set("read-file.poll-interval", time.Second)

// 				cfg := config.SendMailConfig{
// 					ReadFileConfig: config.ReadFileConfig{
// 						RedisAddr:    "invalid:address",
// 						InPath:       "/tmp",
// 						FileMails:    []config.FileMailConfig{},
// 						Concurrency:  1,
// 						PollInterval: time.Second,
// 					},
// 					Outputs:        []config.OutputConfig{},
// 					MailProcessors: []config.MailProcessorConfig{},
// 					Debug:          false,
// 				}
// 				ctx = config.SetContextConfig(ctx, cfg)

// 				cmd := &cobra.Command{}
// 				cmd.SetContext(ctx)

// 				cleanup := func() {
// 					viper.Reset()
// 				}

// 				return ctx, cleanup
// 			},
// 			expectError: true,
// 		},
// 		{
// 			name: "invalid file path",
// 			setup: func(t *testing.T) (context.Context, func()) {
// 				ctx := context.Background()
// 				ctx, _ = telemetry.InitLogger(ctx)

// 				// Setup invalid file path
// 				viper.Set("read-file.redis-addr", "localhost:6379")
// 				viper.Set("read-file.in-path", "/nonexistent/path")
// 				viper.Set("read-file.file-mails", []config.FileMailConfig{})
// 				viper.Set("outputs", []config.OutputConfig{})
// 				viper.Set("mail-processors", []config.MailProcessorConfig{})
// 				viper.Set("debug", false)
// 				viper.Set("read-file.concurrency", 1)
// 				viper.Set("read-file.poll-interval", time.Second)

// 				cfg := config.SendMailConfig{
// 					ReadFileConfig: config.ReadFileConfig{
// 						RedisAddr:    "localhost:6379",
// 						InPath:       "/nonexistent/path",
// 						FileMails:    []config.FileMailConfig{},
// 						Concurrency:  1,
// 						PollInterval: time.Second,
// 					},
// 					Outputs:        []config.OutputConfig{},
// 					MailProcessors: []config.MailProcessorConfig{},
// 					Debug:          false,
// 				}
// 				ctx = config.SetContextConfig(ctx, cfg)

// 				cmd := &cobra.Command{}
// 				cmd.SetContext(ctx)

// 				cleanup := func() {
// 					viper.Reset()
// 				}

// 				return ctx, cleanup
// 			},
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx, cleanup := tt.setup(t)
// 			defer cleanup()

// 			cmd := &cobra.Command{}
// 			cmd.SetContext(ctx)

// 			svc := newGenericSvc(cmd, nil)

// 			if tt.expectError {
// 				// Verify that critical components are nil
// 				assert.Nil(t, svc.RedisClient)
// 				assert.Nil(t, svc.FileReader)
// 				assert.Nil(t, svc.MailProcessor)
// 			} else {
// 				// Verify that all components are initialized
// 				assert.NotNil(t, svc.RedisClient)
// 				assert.NotNil(t, svc.FileReader)
// 				assert.NotNil(t, svc.MailProcessor)
// 				assert.NotNil(t, svc.MailSender)
// 				assert.NotNil(t, svc.MailTransformerFactory)
// 				assert.NotNil(t, svc.MyOutput)
// 				assert.NotNil(t, svc.MyResolver)
// 				assert.NotNil(t, svc.SendMailService)
// 				assert.NotNil(t, svc.Slogger)
// 			}
// 		})
// 	}
// }

// func TestGenericSvc_Integration(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	tests := []struct {
// 		name        string
// 		setup       func(t *testing.T) (*GenericSvc, func())
// 		expectError bool
// 	}{
// 		{
// 			name: "full integration test",
// 			setup: func(t *testing.T) (*GenericSvc, func()) {
// 				// Create temporary directory for test files
// 				tmpDir := t.TempDir()
// 				defer os.RemoveAll(tmpDir)

// 				ctx := context.Background()
// 				ctx, _ = telemetry.InitLogger(ctx)

// 				// Setup configuration
// 				viper.Set("read-file.redis-addr", "localhost:6379")
// 				viper.Set("read-file.in-path", tmpDir)
// 				viper.Set("read-file.file-mails", []config.FileMailConfig{})
// 				viper.Set("outputs", []config.OutputConfig{})
// 				viper.Set("mail-processors", []config.MailProcessorConfig{})
// 				viper.Set("debug", false)
// 				viper.Set("read-file.concurrency", 1)
// 				viper.Set("read-file.poll-interval", time.Second)

// 				cfg := config.SendMailConfig{
// 					ReadFileConfig: config.ReadFileConfig{
// 						RedisAddr:    "localhost:6379",
// 						InPath:       tempDir,
// 						FileMails:    []config.FileMailConfig{},
// 						Concurrency:  1,
// 						PollInterval: time.Second,
// 					},
// 					Outputs:        []config.OutputConfig{},
// 					MailProcessors: []config.MailProcessorConfig{},
// 					Debug:          false,
// 				}
// 				ctx = config.SetContextConfig(ctx, cfg)

// 				cmd := &cobra.Command{}
// 				cmd.SetContext(ctx)

// 				svc := newGenericSvc(cmd, nil)

// 				cleanup := func() {
// 					os.RemoveAll(tempDir)
// 					viper.Reset()
// 				}

// 				return svc, cleanup
// 			},
// 			expectError: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svc, cleanup := tt.setup(t)
// 			defer cleanup()

// 			// Test Redis connection
// 			_, err := svc.RedisClient.Ping(context.Background()).Result()
// 			if tt.expectError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			// Test file operations
// 			if !tt.expectError {
// 				testFile := filepath.Join(svc.Cfg.ReadFileConfig.InPath, "test.txt")
// 				err := os.WriteFile(testFile, []byte("test"), 0644)
// 				assert.NoError(t, err)

// 				// Test file reader
// 				_, err = svc.FileReader.RefreshList(context.Background())
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }
