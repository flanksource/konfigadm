package phases

import (
	"strings"
	"testing"
)

func TestCommandRuntimeFlag(t *testing.T) {
}

func TestPreCommand(t *testing.T) {
}
func TestPostCommand(t *testing.T) {
}
func TestCommand(t *testing.T) {
}
func TestCommandInterpolation(t *testing.T) {
}

func (cfg *SystemConfig) FindCmd(prefix string) []*Command {
	cmds := []*Command{}

	for _, cmd := range cfg.PreCommands {
		if strings.HasPrefix(cmd.Cmd, prefix) {
			cmds = append(cmds, &cmd)
		}
	}
	return cmds
}
