package phases

var Commands AllPhases = command{}

type command struct{}

func (p command) ApplyPhase(sys *SystemConfig, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}
	commands = append(commands, sys.PreCommands...)
	commands = append(commands, sys.Commands...)
	commands = append(commands, sys.PostCommands...)

	return commands, files, nil
}
func (p command) ProcessFlags(sys *SystemConfig, flags ...Flag) {
	sys.PreCommands = filter(sys.PreCommands, flags...)
	sys.Commands = filter(sys.Commands, flags...)
	sys.PostCommands = filter(sys.PostCommands, flags...)
}

func filter(commands []Command, flags ...Flag) []Command {
	minified := []Command{}
	for _, cmd := range commands {
		if MatchAll(flags, cmd.Flags) {
			minified = append(minified, cmd)
		}
	}
	return minified
}
