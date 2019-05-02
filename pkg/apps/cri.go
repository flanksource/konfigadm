package apps

import . "github.com/moshloop/configadm/pkg/types"

var CRI Phase = cri{}

type cri struct{}

func (c cri) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {

	return []Command{}, Filesystem{}, nil
}
