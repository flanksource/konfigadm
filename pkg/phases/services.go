package phases

import (
	"fmt"
	"strings"

	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint
	"github.com/flanksource/konfigadm/pkg/utils"
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
		if svc.Extra.Install.WantedBy == "" && svc.Extra.Install.RequiredBy == "" {
			svc.Extra.Install.WantedBy = "multi-user.target"
		}
		files[filename] = File{Content: svc.Extra.ToUnitFile()}
		commands = append(commands, Command{Cmd: "systemctl enable " + name})
		commands = append(commands, Command{Cmd: "systemctl start " + name})
	}
	return commands, files, nil
}

func (p services) Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool {
	verify := true
	for name := range cfg.Services {
		verify = verify && VerifyService(name, results)

	}
	return verify
}

//VerifyService checks that the service is enabled and running
func VerifyService(name string, results *VerifyResults) bool {
	stdout, ok := utils.SafeExec("systemctl status %s | grep Active", name)
	stdout = strings.TrimSpace(strings.Split(stdout, "\n")[0])
	if !ok {
		results.Fail("%s is %s", name, stdout)

	} else if strings.Contains(stdout, "active (running)") {
		results.Pass("%s is  %s", name, stdout)
		return true
	} else {
		results.Fail("%s is %s", name, stdout)
	}
	return false
}
