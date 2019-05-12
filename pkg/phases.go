package pkg

import (
	"github.com/moshloop/configadm/pkg/apps"
	"github.com/moshloop/configadm/pkg/phases"
	"github.com/moshloop/configadm/pkg/types"
)

func init() {
	types.Dig.Provide(func() *[]types.Phase {
		return &[]types.Phase{
			phases.Context,
			apps.Kubernetes,
			apps.CRI,
			phases.Sysctl,
			phases.Environment,
			phases.Containers,
			phases.Packages,
			phases.Services,
			phases.Files,
			phases.Commands,
		}
	})
}
