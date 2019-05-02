package apps

import (
	. "github.com/moshloop/configadm/pkg/types"
)

var CNI Phase = cni{}

type cni struct{}

func (k cni) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {

	return []Command{}, Filesystem{}, nil
}
