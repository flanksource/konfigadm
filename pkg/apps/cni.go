package apps

import (
	"github.com/flanksource/konfigadm/pkg/types"
)

var CNI types.Phase = cni{}

type cni struct{}

func (k cni) ApplyPhase(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {

	return []types.Command{}, types.Filesystem{}, nil
}
