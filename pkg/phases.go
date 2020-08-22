package pkg

import (
	"github.com/flanksource/konfigadm/pkg/apps"
	"github.com/flanksource/konfigadm/pkg/phases"
	"github.com/flanksource/konfigadm/pkg/types"
	log "github.com/sirupsen/logrus"
)

func init() {
	if err := types.Dig.Provide(func() *[]types.Phase {
		return &[]types.Phase{
			apps.Kubernetes,
			apps.CRI,
			phases.Sysctl,
			phases.Environment,
			phases.Containers,
			phases.Packages,
			phases.Services,
			phases.TrustedCA,
			phases.Files,
			phases.CommandsPhase,
			phases.Users,
			apps.Cleanup,
		}
	}); err != nil {
		log.Errorf("Failed to provide dependencies: %s", err)
	}
}
