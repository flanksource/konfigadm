package phases

import (
	"fmt"

	. "github.com/moshloop/configadm/pkg/utils"
)

var Sysctl Phase = sysctl{}

type sysctl struct{}

func (p sysctl) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}

	if len(sys.Sysctls) > 0 {
		filename := fmt.Sprintf("/etc/sysctl.conf.d/100-%s.conf", sys.Context.Name)
		files[filename] = File{Content: MapToIni(sys.Sysctls)}
	}

	for k, v := range sys.Sysctls {
		commands = append(commands, Command{Cmd: fmt.Sprintf("sysctl -w %s %s", k, v)})
	}
	return commands, files, nil
}
