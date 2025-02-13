package config

type FileMailConfig struct {
	Args  map[string]string `mapstructure:"args"`
	Index int               `mapstructure:"index"`
	Type  string            `mapstructure:"type"`
}

func DefaultFileMailConfigs() []FileMailConfig {
	return []FileMailConfig{
		{
			Args:  map[string]string{},
			Index: 0,
			Type:  "headers",
		},
		{
			Args:  map[string]string{},
			Index: 1,
			Type:  "header_from",
		},
		{
			Args:  map[string]string{},
			Index: 2,
			Type:  "header_to",
		},
		{
			Args:  map[string]string{},
			Index: 3,
			Type:  "header_subject",
		},
	}
}
