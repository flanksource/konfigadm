package phases

import (
	"github.com/flanksource/konfigadm/pkg/types"
)

var CommandsPhase types.AllPhases = command{}

type command struct{}

func (p command) ApplyPhase(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	var commands []types.Command
	files := types.Filesystem{}
	commands = append(commands, sys.PreCommands...)
	commands = append(commands, sys.Commands...)
	commands = append(commands, sys.PostCommands...)
	sys.PreCommands = nil
	sys.Commands = nil
	sys.PostCommands = nil

	return commands, files, nil
}
func (p command) ProcessFlags(sys *types.Config, flags ...types.Flag) {
	sys.PreCommands = types.FilterFlags(sys.PreCommands, flags...)
	sys.Commands = types.FilterFlags(sys.Commands, flags...)
	sys.PostCommands = types.FilterFlags(sys.PostCommands, flags...)
}
