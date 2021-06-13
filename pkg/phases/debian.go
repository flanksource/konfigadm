package phases

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
)

var (
	Ubuntu = ubuntu{}
	Debian = debian{}
)

type ubuntu struct {
}

func (u ubuntu) GetName() string {
	return "ubuntu"
}

func (u ubuntu) GetPackageManager() types.PackageManager {
	return AptPackageManager{}
}

func (u ubuntu) GetTags() []types.Flag {
	osrelease, _ := ini.Load("/etc/os-release")
	majorVersionID, _ := strconv.Atoi(strings.Split(osrelease.Section("").Key("VERSION_ID").String(), ".")[0])
	if majorVersionID == 20 {
		return []types.Flag{types.UBUNTU, types.UBUNTU20, types.DEBIAN_LIKE}
	} else if majorVersionID == 18 {
		return []types.Flag{types.UBUNTU, types.UBUNTU18, types.DEBIAN_LIKE}
	} else if majorVersionID == 16 {
		return []types.Flag{types.UBUNTU, types.UBUNTU16, types.DEBIAN_LIKE}
	}
	return []types.Flag{types.UBUNTU, types.DEBIAN_LIKE}
}

func (u ubuntu) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "Ubuntu")
}

func (u ubuntu) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

func (u ubuntu) GetKernelPackageNames(version string) (string, string) {
	return AptPackageManager{}.GetKernelPackageNames(fmt.Sprintf("%s-generic", version))
}

func (u ubuntu) UpdateDefaultKernel(version string) types.Commands {
	return GrubConfManager{}.UpdateDefaultKernel(fmt.Sprintf("Advanced options for Ubuntu>Ubuntu, with Linux %s-generic", version))
}

func (u ubuntu) ReadDefaultKernel() string {
	readkernel, ok := GrubbyManager{}.ReadDefaultKernel()
	if ok {
		re := regexp.MustCompile("Advanced options for Ubuntu>Ubuntu, with Linux (.*?)-generic")
		match := re.FindStringSubmatch(readkernel)
		if len(match) > 1 {
			return match[1]
		}
		return ""
	}
	return ""
}

type debian struct {
}

func (d debian) GetName() string {
	return "debian"
}

func (d debian) GetPackageManager() types.PackageManager {
	return AptPackageManager{}
}

func (d debian) GetTags() []types.Flag {
	osrelease, err := ini.Load("/etc/os-release")
	if err != nil {
		return []types.Flag{}
	}
	majorVersionID, _ := strconv.Atoi(strings.Split(osrelease.Section("").Key("VERSION_ID").String(), ".")[0])
	if majorVersionID == 9 {
		return []types.Flag{types.DEBIAN, types.DEBIAN9, types.DEBIAN_LIKE}
	} else if majorVersionID == 10 {
		return []types.Flag{types.DEBIAN, types.DEBIAN10, types.DEBIAN_LIKE}
	}
	return []types.Flag{types.DEBIAN, types.DEBIAN_LIKE}
}

func (d debian) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "Debian")
}

func (d debian) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

func (d debian) GetKernelPackageNames(version string) (string, string) {
	read, ok := utils.SafeExec("uname -p")
	if !ok {
		return AptPackageManager{}.GetKernelPackageNames(version)
	}
	arch := strings.TrimSuffix(read, "\n")
	if arch == "unknown" {
		arch = "amd64"
	}
	return AptPackageManager{}.GetKernelPackageNames(fmt.Sprintf("%s-%s", version, arch))
}

func (d debian) UpdateDefaultKernel(version string) types.Commands {
	return GrubConfManager{}.UpdateDefaultKernel(fmt.Sprintf("Debian GNU/Linux, with Linux %s", version))
}

func (d debian) ReadDefaultKernel() string {
	readkernel, ok := GrubbyManager{}.ReadDefaultKernel()
	if ok {
		re := regexp.MustCompile("Debian GNU/Linux, with Linux (.*?)")
		match := re.FindStringSubmatch(readkernel)
		if len(match) > 1 {
			return match[1]
		}
		return ""
	}
	return ""
}

func (d debian) Cleanup() string {
	return `
	unset HISTFILE
	find /var/log -type f | while read f; do echo -ne '' > "${f}"; done;
	UBUNTU_VERSION=$(lsb_release -sr)
if [[ ${UBUNTU_VERSION} == 16.04 ]] || [[ ${UBUNTU_VERSION} == 16.10 ]]; then
    # Modified version of
    # https://github.com/cbednarski/packer-ubuntu/blob/master/scripts-1604/vm_cleanup.sh#L9-L15
    # Instead of eth0 the interface is now called ens5 to mach the PCI
    # slot, so we need to change the networking scripts to enable the
    # correct interface.
    #
    # NOTE: After the machine is rebooted Packer will not be able to reconnect
    # (Vagrant will be able to) so make sure this is done in your final
    # provisioner.
    sed -i "s/ens3/ens5/g" /etc/network/interfaces
fi
	`
}
