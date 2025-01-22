package config

type MailProcessorConfig struct {
	Args  map[string]string `json:"args" mapstructure:"args"`
	Index int               `json:"index" mapstructure:"index"`
	Type  string            `json:"type" mapstructure:"type"`
}

func DefaultMailProcessorConfigs() []MailProcessorConfig {
	return []MailProcessorConfig{
		{
			Args:  map[string]string{},
			Index: 0,
			Type:  "unixdosProcessor",
		},
	}
}
