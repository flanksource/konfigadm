package types

import (
	"log"
	"strings"
)

type ConfigBuilder struct {
	configs []string
	vars    []string
	flags   []Flag
}

func (f *ConfigBuilder) WithVars(vars ...string) *ConfigBuilder {
	f.vars = vars
	return f
}

func (f *ConfigBuilder) WithFlags(flags ...Flag) *ConfigBuilder {
	f.flags = flags
	return f
}

func (builder *ConfigBuilder) Build() (*Config, error) {
	cfg := &Config{}
	cfg.Init()
	cfg.Context.Flags = builder.flags
	for _, config := range builder.configs {
		c, err := newConfig(config)
		if err != nil {
			log.Fatalf("Error parsing %s: %s", config, err)
		}
		cfg.ImportConfig(*c)
	}

	for _, v := range builder.vars {
		if strings.Contains(v, "=") {
			cfg.Context.Vars[strings.Split(v, "=")[0]] = strings.Split(v, "=")[1]
		}
	}

	return cfg, nil
}

func NewConfig(configs ...string) *ConfigBuilder {
	return &ConfigBuilder{
		configs: configs,
	}
}
