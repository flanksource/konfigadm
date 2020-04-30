package types

import (
	"strings"

	"github.com/flanksource/commons/logger"
	"gopkg.in/flanksource/yaml.v3"
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
		if config == "" {
			continue
		}
		c, err := newConfig(config)
		if err != nil {
			logger.Fatalf("Error parsing %s: %s", config, err)
		} else {
			data, _ := yaml.Marshal(c)
			logger.Tracef("\n%s\n", string(data))
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
