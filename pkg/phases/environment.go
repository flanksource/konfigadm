package phases

import (
	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint
	"github.com/flanksource/konfigadm/pkg/utils"
)

var Environment Phase = environment{}

type environment struct{}

func (p environment) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}

	if len(sys.Environment) > 0 {
		files["/etc/environment"] = File{Content: utils.MapToIni(sys.Environment)}
	}
	return commands, files, nil
}
