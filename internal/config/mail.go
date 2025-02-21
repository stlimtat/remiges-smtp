package config

type MailProcessorConfig struct {
	Args  map[string]any `mapstructure:"args"`
	Index int            `mapstructure:"index"`
	Type  string         `mapstructure:"type"`
}

func DefaultMailProcessorConfigs() []MailProcessorConfig {
	return []MailProcessorConfig{
		{
			Args:  map[string]any{},
			Index: 0,
			Type:  "unixdos",
		},
		{
			Args:  map[string]any{},
			Index: 1,
			Type:  "body",
		},
		{
			Args:  map[string]any{},
			Index: 2,
			Type:  "bodyHeaders",
		},
		{
			Args:  map[string]any{},
			Index: 99,
			Type:  "mergeBody",
		},
	}
}
