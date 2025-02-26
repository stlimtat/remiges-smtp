package config

import "context"

const (
	ConfigOutputTypeFile string = "file"
	ConfigArgPath        string = "path"
)

type OutputConfig struct {
	Args map[string]any `mapstructure:"args,omitempty"`
	Type string         `mapstructure:"type"`
}

func DefaultOutputConfig(
	_ context.Context,
) []OutputConfig {
	result := []OutputConfig{
		{
			Type: ConfigOutputTypeFile,
			Args: map[string]any{
				ConfigArgPath: "~/",
			},
		},
	}
	return result
}
