package phases

import (
	"fmt"
	"strings"

	. "github.com/moshloop/configadm/pkg/types"
	"github.com/moshloop/configadm/pkg/utils"
)

var Services Phase = services{}

type services struct{}

func (p services) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}

	for name, svc := range sys.Services {
		filename := fmt.Sprintf("/etc/systemd/system/%s.service", name)
		svc.Extra.Service.ExecStart = svc.ExecStart
		svc.Extra.Unit.Description = name
		files[filename] = File{Content: svc.Extra.ToUnitFile()}
		commands = append(commands, Command{Cmd: "systemctl enable " + name})
		commands = append(commands, Command{Cmd: "systemctl start " + name})
	}
	return commands, files, nil
}

func (p services) Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool {
	verify := true
	for name, _ := range cfg.Services {
		stdout, ok := utils.SafeExec("systemctl status %s | grep Active", name)
		if !ok {
			results.Fail("%s is not running %s", name, stdout)
			verify = false
		} else if strings.Contains(stdout, "active (running)") {
			results.Pass("%s is  running %s", name, stdout)
		} else {
			results.Fail("%s is not running %s", name, stdout)
		}
	}
	return verify
}
