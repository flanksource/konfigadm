package phases

import (
	"testing"

	"github.com/onsi/gomega"
)

func SetupFixtureWithArgs(name string, vars []string, t *testing.T) (*SystemConfig, *gomega.WithT) {
	cfg := NewSystemConfigFromFixture(name, vars, t)
	g := gomega.NewWithT(t)
	return cfg, g
}
func SetupFixture(name string, t *testing.T) (*SystemConfig, *gomega.WithT) {
	cfg := NewSystemConfigFromFixture(name, []string{}, t)
	g := gomega.NewWithT(t)
	return cfg, g
}

func NewSystemConfigFromFixture(name string, vars []string, t *testing.T) *SystemConfig {
	cfg, err := NewSystemConfig(
		vars,
		[]string{"../../fixtures/" + name},
	)

	if err != nil {
		t.Error(err)
	}
	return cfg
}

func TestMultipleConfigs(t *testing.T) {
}

func TestRemoteConfig(t *testing.T) {
}
