package phases

import (
	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
)

var Photon = photon{}

type photon struct {
}

func (p photon) GetPackageManager() types.PackageManager {
	return TdnfPackageManager{}
}

func (p photon) GetTags() []types.Flag {
	return []types.Flag{types.PHOTON}
}

func (p photon) DetectAtRuntime() bool {
	id, ok := utils.IniToMap("/etc/os-release")["ID"]
	return ok && id == "photon"
}

func (p photon) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}
