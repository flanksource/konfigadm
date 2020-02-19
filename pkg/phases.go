package pkg

import (
	"github.com/flanksource/konfigadm/pkg/apps"
	"github.com/flanksource/konfigadm/pkg/phases"
	"github.com/flanksource/konfigadm/pkg/types"
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
			phases.CommandsPhase,
			apps.Cleanup,
		}
	})
}
