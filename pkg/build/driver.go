package build

import "github.com/moshloop/konfigadm/pkg/types"

type Driver interface {
	Build(image string, cfg *types.Config)
}
