package config

type FileMailConfig struct {
	Args  map[string]any `mapstructure:"args"`
	Index int            `mapstructure:"index"`
	Type  string         `mapstructure:"type"`
}

func DefaultFileMailConfigs() []FileMailConfig {
	return []FileMailConfig{
		{
			Args: map[string]any{
				"prefix": "H??",
			},
			Index: 0,
			Type:  "headers",
		},
		{
			Args:  map[string]any{},
			Index: 1,
			Type:  "header_from",
		},
		{
			Args:  map[string]any{},
			Index: 2,
			Type:  "header_to",
		},
		{
			Args: map[string]any{
				"default": "no subject",
			},
			Index: 3,
			Type:  "header_subject",
		},
		{
			Args:  map[string]any{},
			Index: 4,
			Type:  "header_contenttype",
		},
		{
			Args:  map[string]any{},
			Index: 5,
			Type:  "header_msgid",
		},
		{
			Args:  map[string]any{},
			Index: 6,
			Type:  "body",
		},
	}
}
