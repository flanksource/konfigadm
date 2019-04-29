package phases

import (
	"fmt"
	"strings"
)

var Packages AllPhases = packages{}

type packages struct{}

func (p packages) ApplyPhase(sys *SystemConfig, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}
	install := []string{}
	uninstall := []string{}
	mark := []string{}
	for _, p := range sys.Packages {
		if p.Uninstall {
			uninstall = append(uninstall, p.Name)
		} else {
			install = append(install, p.Name)
		}
	}

	if len(install) > 0 {
		commands = append(commands, Command{
			Cmd:   fmt.Sprintf("apt-get install -y %s", strings.Join(install, " ")),
			Flags: []Flag{DEBIAN},
		})
		commands = append(commands, Command{
			Cmd:   fmt.Sprintf("yum install -y %s", strings.Join(install, " ")),
			Flags: []Flag{REDHAT},
		})
	}
	if len(uninstall) > 0 {
		commands = append(commands, Command{
			Cmd:   fmt.Sprintf("apt-get remove -y %s", strings.Join(uninstall, " ")),
			Flags: []Flag{DEBIAN},
		})
		commands = append(commands, Command{
			Cmd:   fmt.Sprintf("yum remove -y %s", strings.Join(uninstall, " ")),
			Flags: []Flag{REDHAT},
		})
	}

	if len(mark) > 0 {
		commands = append(commands, Command{
			Cmd:   fmt.Sprintf("apt-get mark %s", strings.Join(mark, " ")),
			Flags: []Flag{DEBIAN},
		})
	}
	return commands, files, nil
}
func (p packages) ProcessFlags(sys *SystemConfig, flags ...Flag) {
	minified := []Package{}
	for _, pkg := range sys.Packages {
		if MatchAll(flags, pkg.Flags) {
			minified = append(minified, pkg)
		}
	}
	sys.Packages = minified
}
