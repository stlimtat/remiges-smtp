package config

type MailProcessorConfig struct {
	Args  map[string]string `mapstructure:"args"`
	Index int               `mapstructure:"index"`
	Type  string            `mapstructure:"type"`
}

func DefaultMailProcessorConfigs() []MailProcessorConfig {
	return []MailProcessorConfig{
		{
			Args:  map[string]string{},
			Index: 0,
			Type:  "unixdos",
		},
		{
			Args:  map[string]string{},
			Index: 1,
			Type:  "bodyHeaders",
		},
		{
			Args:  map[string]string{},
			Index: 99,
			Type:  "mergeBody",
		},
	}
}
