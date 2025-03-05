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
		{
			Args: map[string]any{
				"domain-str": "stlim.net",
				"dkim": map[string]any{
					"selectors": map[string]any{
						"key001": map[string]any{
							"domain":           "key001",
							"algorithm":        "rsa",
							"hash":             "sha256",
							"headers":          []string{"from", "to", "subject"},
							"private-key-file": "./config/key001.pem",
						},
					},
					"sign": []string{"key001"},
				},
			},
			Index: 100,
			Type:  "dkim",
		},
	}
}
