package phases_test

import (
	"testing"

	. "github.com/flanksource/konfigadm/pkg/types"
	"github.com/onsi/gomega"
)

func init() {
	// log.SetLevel(log.TraceLevel)
}

func TestCommandRuntimeFlag(t *testing.T) {
}

func setupCommandFixture(t *testing.T, flag Flag) ([]Command, *gomega.WithT) {
	cfg, g := NewFixture("commands.yml", t).WithFlags(flag).Build()
	_, commands, _ := cfg.ApplyPhases()
	return commands, g
}

func TestCommand(t *testing.T) {
	commands, g := setupCommandFixture(t, DEBIAN)
	g.Expect(commands).To(MatchCommand("echo command"))
}

func TestPreCommand(t *testing.T) {
	commands, g := setupCommandFixture(t, DEBIAN)
	g.Expect(commands).To(MatchCommand("echo pre"))
}
func TestPostCommand(t *testing.T) {
	commands, g := setupCommandFixture(t, DEBIAN)
	g.Expect(commands).To(MatchCommand("echo post"))
}
func TestCommandInterpolation(t *testing.T) {
}
