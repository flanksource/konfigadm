package cloudmeta

import (
	"fmt"

	. "github.com/moshloop/configadm/pkg/utils"
)

func init() {
	Register(TransformSysctls)
}

func TransformSysctls(sys *SystemConfig, ctx *SystemContext) (commands []string, files map[string]string, err error) {
	commands = []string{}
	files = make(map[string]string)
	if len(sys.Sysctls) > 0 {
		filename := fmt.Sprintf("/etc/sysctl.conf.d/100-%s.conf", sys.Context.Name)
		files[filename] = MapToIni(sys.Sysctls)
	}

	for k, v := range sys.Sysctls {
		commands = append(commands, fmt.Sprintf("sysctl -w %s %s", k, v))
	}
	return commands, files, nil
}
