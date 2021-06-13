package phases

import (
	"strconv"
	"strings"

	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
	"gopkg.in/ini.v1"
)

var Photon = photon{}

type photon struct {
}

func (p photon) GetName() string {
	return "photon"
}

func (p photon) GetPackageManager() types.PackageManager {
	return TdnfPackageManager{}
}

func (p photon) GetTags() []types.Flag {
	osrelease, _ := ini.Load("/etc/os-release")
	majorVersionID, _ := strconv.Atoi(strings.Split(osrelease.Section("").Key("VERSION_ID").String(), ".")[0])
	if majorVersionID == 2 {
		return []types.Flag{types.PHOTON2, types.PHOTON}
	} else if majorVersionID == 3 {
		return []types.Flag{types.PHOTON3, types.PHOTON}
	}
	return []types.Flag{types.PHOTON}
}

func (p photon) DetectAtRuntime() bool {
	id, ok := utils.IniToMap("/etc/os-release")["ID"]
	return ok && id == "photon"
}

func (p photon) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

func (p photon) GetKernelPackageNames(version string) (string, string) {
	return TdnfPackageManager{}.GetKernelPackageNames(version)
}

func (p photon) UpdateDefaultKernel(version string) types.Commands {
	return GrubbyManager{}.UpdateDefaultKernel(version)
}

func (p photon) ReadDefaultKernel() string {
	readkernel, _ := GrubbyManager{}.ReadDefaultKernel()
	return readkernel
}
