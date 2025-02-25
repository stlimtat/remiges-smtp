package config

import "context"

type OutputConfig struct {
	Args map[string]any `mapstructure:"args,omitempty"`
	Type string         `mapstructure:"type"`
}

func DefaultOutputConfig(
	_ context.Context,
) OutputConfig {
	result := OutputConfig{
		Type: "file",
		Args: map[string]any{
			"path": "output.log",
		},
	}
	return result
}
