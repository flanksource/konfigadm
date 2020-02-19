package phases_test

import (
	"testing"

	"github.com/onsi/gomega"

	. "github.com/flanksource/konfigadm/pkg/types"
)

func TestArgs(t *testing.T) {
	cfg, g := NewFixture("env.yml", t).WithVars("env1=value1", "env2=value2").Build()
	cfg.ApplyPhases()
	g.Expect(cfg.Environment["env1"]).To(gomega.Equal("val: value1"))
}

func TestArgsFile(t *testing.T) {
}
func TestArgsRemotes(t *testing.T) {
}

func TestArgsPriority(t *testing.T) {
}

func TestArgRuntimeFlags(t *testing.T) {

}
