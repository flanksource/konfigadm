package types

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type Package struct {
	Name      string
	Mark      bool
	Uninstall bool
	Flags     []Flag
}

func (p Package) String() string {
	return p.Name
}

func (p Package) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:        yaml.ScalarNode,
		Tag:         "!!str",
		LineComment: Marshall(p.Flags),
		Value:       p.Name,
	}, nil
}

func (p *Package) UnmarshalYAML(node *yaml.Node) error {
	p.Name = node.Value
	if strings.HasPrefix(node.Value, "!") {
		p.Name = node.Value[1:]
		p.Uninstall = true
	}
	if strings.HasPrefix(node.Value, "=") {
		p.Name = node.Value[1:]
		p.Mark = true
	}
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

type PackageRepo struct {
	URL    string
	GPGKey string
	Flags  []Flag
}

func (p PackageRepo) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:        yaml.ScalarNode,
		Tag:         "!!str",
		LineComment: Marshall(p.Flags),
		Value:       p.URL,
	}, nil
}
