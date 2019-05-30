package types

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

var (
	flagProcessors   = make([]FlagProcessor, 0)
	DEBIAN           = Flag{Name: "debian"}
	DEBIAN_LIKE      = Flag{Name: "debian-like"}
	REDHAT           = Flag{Name: "redhat"}
	REDHAT_LIKE      = Flag{Name: "redhat-like"}
	AMAZON_LINUX     = Flag{Name: "amazonLinux"}
	RHEL             = Flag{Name: "rhel"}
	CENTOS           = Flag{Name: "centos"}
	UBUNTU           = Flag{Name: "ubuntu"}
	AWS              = Flag{Name: "aws"}
	VMWARE           = Flag{Name: "vmware"}
	NOT_DEBIAN       = Flag{Name: "!debian", Negates: []Flag{DEBIAN}}
	NOT_REDHAT       = Flag{Name: "!redhat", Negates: []Flag{REDHAT}}
	NOT_DEBIAN_LIKE  = Flag{Name: "!debian", Negates: []Flag{DEBIAN_LIKE}}
	NOT_REDHAT_LIKE  = Flag{Name: "!redhat", Negates: []Flag{REDHAT_LIKE}}
	NOT_CENTOS       = Flag{Name: "!centos", Negates: []Flag{CENTOS}}
	NOT_RHEL         = Flag{Name: "!rhel", Negates: []Flag{RHEL}}
	NOT_UBUNTU       = Flag{Name: "!ubuntu", Negates: []Flag{UBUNTU}}
	NOT_AWS          = Flag{Name: "!aws", Negates: []Flag{AWS}}
	NOT_VMWARE       = Flag{Name: "!vmware", Negates: []Flag{VMWARE}}
	NOT_AMAZON_LINUX = Flag{Name: "!amazonLinux", Negates: []Flag{AMAZON_LINUX}}
	FLAG_MAP         = make(map[string]Flag)
	FLAGS            = []Flag{DEBIAN, DEBIAN_LIKE, REDHAT, REDHAT_LIKE, AMAZON_LINUX, CENTOS, RHEL, UBUNTU, AWS, VMWARE, NOT_DEBIAN_LIKE, NOT_REDHAT_LIKE, NOT_DEBIAN, NOT_REDHAT, NOT_CENTOS, NOT_RHEL, NOT_UBUNTU, NOT_AWS, NOT_VMWARE, NOT_AMAZON_LINUX}
)

type Flag struct {
	Name    string
	Negates []Flag
}

func GetTag(name string) *Flag {
	for _, tag := range FLAGS {
		if strings.ToLower(tag.Name) == strings.ToLower(name) {
			return &tag
		}
	}
	return nil
}

func init() {
	for _, flag := range FLAGS {
		name := flag.Name
		FLAG_MAP[name] = flag
		if !strings.HasPrefix(name, "!") && !strings.HasPrefix(name, "+") {
			FLAG_MAP[fmt.Sprintf("+%s", name)] = flag
		}
	}

}

func (f Flag) String() string {
	return f.Name
}

func (f *Flag) Matches(other Flag) bool {
	if f.Name == other.Name {
		return true
	}

	if len(other.Negates) > 0 {
		for _, flag2 := range other.Negates {
			if f.Matches(flag2) {
				return false
			}
		}
		return true
	}
	return false
}

//MatchAll returns true if all constraints match at least one flag AND none of the constraints negates any flag
func MatchAll(flags []Flag, constraints []Flag) bool {
	if len(constraints) == 0 {
		return true
	}
outer:
	for _, constraint := range constraints {
		for _, flag := range flags {
			if flag.Matches(constraint) {
				continue outer
			}
		}
		log.Debugf("%s don't match any constraints %s\n", flags, constraints)
		return false
	}
	return true
}

func MatchesAny(flags []Flag, constraints []Flag) bool {
	for _, constraint := range constraints {
		for _, flag := range flags {
			if constraint.Matches(flag) {
				return true
			}
		}
	}
	return false
}

func NegatesAny(flags []Flag, constraints []Flag) bool {
	for _, constraint := range constraints {
		for _, negate := range constraint.Negates {
			for _, flag := range flags {
				if negate.Matches(flag) {
					return true
				}
			}
		}
	}
	return false
}

func Marshall(flags []Flag) string {
	if len(flags) == 0 {
		return ""
	}
	s := ""
	for _, flag := range flags {
		s += flag.String() + " "
	}
	return strings.TrimSpace("#" + s)
}

//MarshalYAML ads tags as comments
func (t Flag) MarshalYAML() (interface{}, error) {
	return t.Name, nil
}

//UnmarshalYAML decodes comments into tags and parses modifiers for packages
func (t *Flag) UnmarshalYAML(node *yaml.Node) error {
	tag := GetTag(node.Value)
	t.Name = tag.Name
	t.Negates = tag.Negates
	log.Tracef("Unmarshal %s into %s\n", node.Value, tag)
	return nil
}

func FilterFlags(commands []Command, flags ...Flag) []Command {
	minified := []Command{}
	for _, cmd := range commands {
		if NegatesAny(flags, cmd.Flags) {
			continue
		}
		if len(cmd.Flags) == 0 || MatchesAny(flags, cmd.Flags) {
			minified = append(minified, cmd)
		} else {
			log.Debugf("%s with tags %s does not match any constraints %s\n", cmd, cmd.Flags, flags)
		}
	}
	return minified
}

func FilterFilesystemByFlags(files Filesystem, flags ...Flag) Filesystem {
	var filtered = make(Filesystem)
	for path, file := range files {
		if NegatesAny(flags, file.Flags) {
			continue
		}
		if len(file.Flags) == 0 || MatchesAny(flags, file.Flags) {
			filtered[path] = file
		} else {
			log.Debugf("%s with tags %s does not match any constraints %s\n", path, file.Flags, flags)
		}
	}
	return filtered
}
