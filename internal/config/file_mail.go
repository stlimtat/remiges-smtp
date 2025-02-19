package config

type FileMailConfig struct {
	Args  map[string]string `mapstructure:"args"`
	Index int               `mapstructure:"index"`
	Type  string            `mapstructure:"type"`
}

func DefaultFileMailConfigs() []FileMailConfig {
	return []FileMailConfig{
		{
			Args: map[string]string{
				"prefix": "H??",
			},
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
			Args: map[string]string{
				"default": "no subject",
			},
			Index: 3,
			Type:  "header_subject",
		},
		{
			Args:  map[string]string{},
			Index: 4,
			Type:  "header_contenttype",
		},
		{
			Args:  map[string]string{},
			Index: 5,
			Type:  "body",
		},
	}
}
