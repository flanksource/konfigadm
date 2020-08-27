package phases

import (
	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
)

var Environment types.Phase = environment{}

type environment struct{}

func (p environment) ApplyPhase(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	var commands []types.Command
	files := types.Filesystem{}

	if len(sys.Environment) > 0 {
		files["/etc/environment"] = types.File{Content: utils.MapToIni(sys.Environment)}
	}
	return commands, files, nil
}
