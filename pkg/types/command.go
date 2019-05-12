package types

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

//Command encapsulates a command and the tags for which it is applicable
type Command struct {
	Cmd   string
	Flags []Flag
}

func (c Command) String() string {
	return c.Cmd
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
			return fmt.Errorf("Unknown flag: %s", flag)
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
