package phases

import . "github.com/moshloop/konfigadm/pkg/types"

var CommandsPhase AllPhases = command{}

type command struct{}

func (p command) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}
	commands = append(commands, sys.PreCommands...)
	commands = append(commands, sys.Commands...)
	commands = append(commands, sys.PostCommands...)

	return commands, files, nil
}
func (p command) ProcessFlags(sys *Config, flags ...Flag) {
	sys.PreCommands = FilterFlags(sys.PreCommands, flags...)
	sys.Commands = FilterFlags(sys.Commands, flags...)
	sys.PostCommands = FilterFlags(sys.PostCommands, flags...)
}
