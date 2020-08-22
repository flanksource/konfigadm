package phases

import (
	"fmt"
	"os"
	"strings"

	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint

	. "github.com/flanksource/konfigadm/pkg/utils" // nolint: golint
)

var Sysctl Phase = sysctl{}

type sysctl struct{}

func (p sysctl) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}

	if len(sys.Sysctls) > 0 {
		filename := fmt.Sprintf("/etc/sysctl.d/100-%s.conf", sys.Context.Name)
		files[filename] = File{Content: MapToIni(sys.Sysctls)}
	}

	for k, v := range sys.Sysctls {
		// make sysctl application errors warnings
		commands = append(commands, Command{Cmd: fmt.Sprintf("sysctl -w %s=%s || true", k, v)})
	}
	return commands, files, nil
}

func (p sysctl) Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool {
	verify := true
	for k, v := range cfg.Sysctls {
		if os.Getenv("container") != "" {
			results.Skip("sysctl[%s]: cannot test inside a container", k)
			continue
		}
		value := SafeRead("/proc/sys" + strings.Replace(k, ".", "/", -1))
		if value == v {
			results.Pass("sysctl[%s]: %s", k, v)
		} else {
			results.Fail("sysctl[%s]: %s != %s", k, value, v)
			verify = false
		}
	}
	return verify
}
