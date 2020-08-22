package phases

import (
	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint
)

var CommandsPhase AllPhases = command{}

type command struct{}

func (p command) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}
	commands = append(commands, sys.PreCommands...)
	commands = append(commands, sys.Commands...)
	commands = append(commands, sys.PostCommands...)
	sys.PreCommands = nil
	sys.Commands = nil
	sys.PostCommands = nil

	return commands, files, nil
}
func (p command) ProcessFlags(sys *Config, flags ...Flag) {
	sys.PreCommands = FilterFlags(sys.PreCommands, flags...)
	sys.Commands = FilterFlags(sys.Commands, flags...)
	sys.PostCommands = FilterFlags(sys.PostCommands, flags...)
}
