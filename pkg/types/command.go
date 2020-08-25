package types

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/flanksource/yaml.v3"
)

//Command encapsulates a command and the tags for which it is applicable
type Command struct {
	Cmd   string
	Flags []Flag
}

type Commands struct {
	dependencies *[]Command
	commands     *[]Command
}

func NewCommand(cmd string) Commands {
	return Commands{
		commands: &[]Command{Command{Cmd: cmd}},
	}
}

func (c *Commands) AddAll(cmd ...Command) *Commands {
	if c.commands == nil {
		c.commands = &[]Command{}
	}
	commandsSlice := *c.commands
	commandsSlice = append(commandsSlice, cmd...)
	c.commands = &commandsSlice
	return c
}

func (c *Commands) Add(commands ...string) *Commands {
	if c.commands == nil {
		c.commands = &[]Command{}
	}
	commandsSlice := *c.commands
	for _, command := range commands {
		commandsSlice = append(commandsSlice, Command{Cmd: command})
	}
	c.commands = &commandsSlice
	return c
}

func (c *Commands) AddDependency(commands ...string) *Commands {
	if c.dependencies == nil {
		c.dependencies = &[]Command{}
	}
	commandsSlice := *c.dependencies
	for _, command := range commands {
		commandsSlice = append(commandsSlice, Command{Cmd: command})
	}
	c.dependencies = &commandsSlice
	return c
}

func contains(commands []Command, other Command) bool {
	for _, cmd := range commands {
		if cmd.Cmd == other.Cmd && MatchesAny(cmd.Flags, other.Flags) {
			return true
		}
	}
	return false
}

func (c Commands) GetCommands() []Command {
	if c.dependencies == nil && c.commands == nil {
		return []Command{}
	}
	if c.dependencies == nil {
		return *c.commands
	}
	if c.commands == nil {
		return *c.dependencies
	}
	return append(*c.dependencies, *c.commands...)
}

func (c *Commands) Append(c2 Commands) *Commands {
	var cmdSlice []Command
	var depSlice []Command
	if c.commands == nil {
		c.commands = &[]Command{}
	}

	if c2.commands != nil {
		cmdSlice = append(*c.commands, *c2.commands...)
		c.commands = &cmdSlice
	}
	if c.dependencies == nil {
		c.dependencies = &[]Command{}
	}
	if c2.dependencies != nil {
		depSlice = append(*c.dependencies, *c2.dependencies...)
		c.dependencies = &depSlice
	}
	return c
}

func (c *Commands) Merge() []Command {
	commands := []Command{}
	if c.dependencies != nil {
		for _, cmd := range *c.dependencies {
			if contains(commands, cmd) {
				continue
			}
			commands = append(commands, cmd)
		}
	}
	if c.commands != nil {
		commands = append(commands, *c.commands...)
	}
	return commands
}

func (c Commands) WithTags(tags ...Flag) Commands {
	new := Commands{commands: &[]Command{}, dependencies: &[]Command{}}
	if c.commands != nil {
		commands := *new.commands
		for _, command := range *c.commands {
			command.Flags = tags
			commands = append(commands, command)
		}
		new.commands = &commands
	}

	if c.dependencies != nil {
		dependencies := *new.dependencies
		for _, command := range *c.dependencies {
			command.Flags = tags
			dependencies = append(dependencies, command)
		}
		new.dependencies = &dependencies
	}
	return new
}

func (c Command) String() string {
	return c.Cmd + fmt.Sprintf("%s", c.Flags)
}

func (cfg *Config) AddCommand(cmd string, flags ...*Flag) *Config {
	command := Command{Cmd: cmd}
	for _, flag := range flags {
		if flag != nil {
			command.Flags = append(command.Flags, *flag)
		}
	}
	cfg.Commands = append(cfg.Commands, command)
	return cfg
}

//UnmarshalYAML decodes comments into tags
func (c *Command) UnmarshalYAML(node *yaml.Node) error {
	c.Cmd = node.Value
	comment := node.LineComment
	if !strings.Contains(comment, "#") {
		return nil
	}
	comment = comment[1:]
	for _, flag := range strings.Split(comment, " ") {
		if FLAG, ok := FLAG_MAP[flag]; ok {
			c.Flags = append(c.Flags, FLAG)
		} else {
			log.Debugf("Ignoring flags: %s on line: %s\n", comment, node.Value)
			return nil
		}
	}
	return nil
}

//MarshalYAML ads tags as comments
func (c Command) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:        yaml.ScalarNode,
		Tag:         "!!str",
		LineComment: Marshall(c.Flags),
		Value:       c.Cmd,
	}, nil
}

//FindCmd returns a list of commands starting with prefix
func (cfg *Config) FindCmd(prefix string) []*Command {
	cmds := []*Command{}

	for _, cmd := range cfg.PreCommands {
		if strings.HasPrefix(cmd.Cmd, prefix) {
			cmds = append(cmds, &cmd)
		}
	}
	return cmds
}
