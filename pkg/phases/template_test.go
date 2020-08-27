package phases_test

import (
	. "github.com/flanksource/konfigadm/pkg/types"
	"github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestArgs(t *testing.T) {
	os.Setenv("key1", "val2")
	cfg, g := NewFixture("env.yml", t).WithVars("env1=value1", "env2=value2").Build()
	if _, _, err := cfg.ApplyPhases(); err != nil {
		log.Errorf("Failed to apply phases: %s", err)
	}
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
