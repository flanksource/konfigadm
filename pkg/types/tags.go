package types

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/flanksource/yaml.v3"
)

var (
	CONTAINER        = Flag{Name: "container"}
	DEBIAN           = Flag{Name: "debian"}
	DEBIAN9          = Flag{Name: "debian9"}
	DEBIAN10         = Flag{Name: "debian10"}
	DEBIAN_LIKE      = Flag{Name: "debian-like"} // nolint: golint
	REDHAT           = Flag{Name: "redhat"}
	FEDORA           = Flag{Name: "fedora"}
	FEDORA30         = Flag{Name: "fedora30"}
	FEDORA31         = Flag{Name: "fedora31"}
	FEDORA32         = Flag{Name: "fedora32"}
	REDHAT_LIKE      = Flag{Name: "redhat-like"} // nolint: golint
	AMAZON_LINUX     = Flag{Name: "amazonLinux"} // nolint: golint
	RHEL             = Flag{Name: "rhel"}
	RHEL7            = Flag{Name: "rhel7"}
	RHEL8            = Flag{Name: "rhel8"}
	CENTOS           = Flag{Name: "centos"}
	CENTOS7          = Flag{Name: "centos7"}
	CENTOS8          = Flag{Name: "centos8"}
	UBUNTU           = Flag{Name: "ubuntu"}
	UBUNTU18         = Flag{Name: "ubuntu18"}
	UBUNTU20         = Flag{Name: "ubuntu20"}
	UBUNTU16         = Flag{Name: "ubuntu16"}
	AWS              = Flag{Name: "aws"}
	PHOTON           = Flag{Name: "photon"}
	PHOTON2          = Flag{Name: "photon2"}
	PHOTON3          = Flag{Name: "photon3"}
	VMWARE           = Flag{Name: "vmware"}
	NOT_PHOTON       = Flag{Name: "!photon", Negates: []Flag{PHOTON}}            // nolint: golint
	NOT_PHOTON2      = Flag{Name: "!photon2", Negates: []Flag{PHOTON2}}          // nolint: golint
	NOT_PHOTON3      = Flag{Name: "!photon3", Negates: []Flag{PHOTON3}}          // nolint: golint
	NOT_CONTAINER    = Flag{Name: "!container", Negates: []Flag{CONTAINER}}      // nolint: golint
	NOT_FEDORA       = Flag{Name: "!fedora", Negates: []Flag{FEDORA}}            // nolint: golint
	NOT_FEDORA30     = Flag{Name: "!fedora30", Negates: []Flag{FEDORA30}}        // nolint: golint
	NOT_FEDORA31     = Flag{Name: "!fedora31", Negates: []Flag{FEDORA31}}        // nolint: golint
	NOT_FEDORA32     = Flag{Name: "!fedora32", Negates: []Flag{FEDORA32}}        // nolint: golint
	NOT_DEBIAN       = Flag{Name: "!debian", Negates: []Flag{DEBIAN}}            // nolint: golint
	NOT_DEBIAN9      = Flag{Name: "!debian9", Negates: []Flag{DEBIAN9}}          // nolint: golint
	NOT_DEBIAN10     = Flag{Name: "!debian10", Negates: []Flag{DEBIAN10}}        // nolint: golint
	NOT_REDHAT       = Flag{Name: "!redhat", Negates: []Flag{REDHAT}}            // nolint: golint
	NOT_DEBIAN_LIKE  = Flag{Name: "!debian", Negates: []Flag{DEBIAN_LIKE}}       // nolint: golint
	NOT_REDHAT_LIKE  = Flag{Name: "!redhat", Negates: []Flag{REDHAT_LIKE}}       // nolint: golint
	NOT_CENTOS       = Flag{Name: "!centos", Negates: []Flag{CENTOS}}            // nolint: golint
	NOT_CENTOS7      = Flag{Name: "!centos7", Negates: []Flag{CENTOS7}}          // nolint: golint
	NOT_CENTOS8      = Flag{Name: "!centos8", Negates: []Flag{CENTOS8}}          // nolint: golint
	NOT_RHEL         = Flag{Name: "!rhel", Negates: []Flag{RHEL}}                // nolint: golint
	NOT_RHEL7        = Flag{Name: "!rhel7", Negates: []Flag{RHEL7}}              // nolint: golint
	NOT_RHEL8        = Flag{Name: "!rhel8", Negates: []Flag{RHEL8}}              // nolint: golint
	NOT_UBUNTU       = Flag{Name: "!ubuntu", Negates: []Flag{UBUNTU}}            // nolint: golint
	NOT_UBUNTU18     = Flag{Name: "!ubuntu18", Negates: []Flag{UBUNTU18}}        // nolint: golint
	NOT_UBUNTU20     = Flag{Name: "!ubuntu20", Negates: []Flag{UBUNTU20}}        // nolint: golint
	NOT_UBUNTU16     = Flag{Name: "!ubuntu16", Negates: []Flag{UBUNTU16}}        // nolint: golint
	NOT_AWS          = Flag{Name: "!aws", Negates: []Flag{AWS}}                  // nolint: golint
	NOT_VMWARE       = Flag{Name: "!vmware", Negates: []Flag{VMWARE}}            // nolint: golint
	NOT_AMAZON_LINUX = Flag{Name: "!amazonLinux", Negates: []Flag{AMAZON_LINUX}} // nolint: golint
	FLAG_MAP         = make(map[string]Flag)                                     // nolint: golint
	FLAGS            = []Flag{CONTAINER, DEBIAN, DEBIAN9, DEBIAN10, DEBIAN_LIKE, REDHAT, FEDORA, FEDORA30, FEDORA31, FEDORA32, REDHAT_LIKE, AMAZON_LINUX, CENTOS, CENTOS7, CENTOS8, RHEL, RHEL7, RHEL8, UBUNTU, UBUNTU16, UBUNTU18, UBUNTU20, PHOTON, PHOTON2, PHOTON3, AWS, VMWARE, NOT_CONTAINER, NOT_FEDORA, NOT_FEDORA30, NOT_FEDORA31, NOT_FEDORA32, NOT_PHOTON, NOT_PHOTON2, NOT_PHOTON3, NOT_DEBIAN_LIKE, NOT_REDHAT_LIKE, NOT_DEBIAN, NOT_DEBIAN9, NOT_DEBIAN10, NOT_REDHAT, NOT_CENTOS, NOT_CENTOS7, NOT_CENTOS8, NOT_RHEL, NOT_RHEL7, NOT_RHEL8, NOT_UBUNTU, NOT_AWS, NOT_VMWARE, NOT_AMAZON_LINUX, NOT_UBUNTU16, NOT_UBUNTU18, NOT_UBUNTU20}
)

type Flag struct {
	Name    string
	Negates []Flag
}

func GetTag(name string) *Flag {
	for _, tag := range FLAGS {
		if strings.EqualFold(tag.Name, name) {
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
	if len(constraints) == 0 {
		return true
	}
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
func (f Flag) MarshalYAML() (interface{}, error) {
	return f.Name, nil
}

//UnmarshalYAML decodes comments into tags and parses modifiers for packages
func (f *Flag) UnmarshalYAML(node *yaml.Node) error {
	tag := GetTag(node.Value)
	f.Name = tag.Name
	f.Negates = tag.Negates
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
			log.Tracef("%s with tags %s does not match any constraints %s", cmd, cmd.Flags, flags)
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
			log.Tracef("%s with tags %s does not match any constraints %s", path, file.Flags, flags)
		}
	}
	return filtered
}
