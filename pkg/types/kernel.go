package types

import (
	"fmt"
	"strings"

	"gopkg.in/flanksource/yaml.v3"
)

//Package includes the package name, modifiers (mark, uninstall) and runtime tags
type KernelInput struct {
	Version string
	Flags   []Flag
}

func (p KernelInput) String() string {
	return p.Version + fmt.Sprintf("%s", p.Flags)
}

//MarshalYAML adds tags as comments
func (p KernelInput) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:        yaml.ScalarNode,
		Tag:         "!!str",
		LineComment: Marshall(p.Flags),
		Value:       p.Version,
	}, nil
}

//UnmarshalYAML decodes comments into tags and parses modifiers for packages
func (p *KernelInput) UnmarshalYAML(node *yaml.Node) error {
	p.Version = node.Value

	comment := node.LineComment
	if !strings.Contains(comment, "#") {
		return nil
	}
	comment = comment[1:]
	for _, flag := range strings.Split(comment, " ") {
		if FLAG, ok := FLAG_MAP[flag]; ok {
			p.Flags = append(p.Flags, FLAG)
		} else {
			return fmt.Errorf("Unknown flag: %s", flag)
		}
	}
	return nil
}
