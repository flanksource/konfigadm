package phases_test

import (
	"testing"
"os"
	"github.com/onsi/gomega"

	. "github.com/flanksource/konfigadm/pkg/types"
)

func TestArgs(t *testing.T) {
	os.Setenv("key1", "val2")
	cfg, g := NewFixture("env.yml", t).WithVars("env1=value1", "env2=value2").Build()
	cfg.ApplyPhases()
	g.Expect(cfg.Environment["env1"]).To(gomega.Equal("val: val2"))
}

func TestArgsFile(t *testing.T) {
}
func TestArgsRemotes(t *testing.T) {
}

func TestArgsPriority(t *testing.T) {
}

func TestArgRuntimeFlags(t *testing.T) {

}
