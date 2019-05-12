package os

import (
	"strings"

	"github.com/moshloop/configadm/pkg/utils"
)

var (
	tmpFolders = []string{
		"/root/.bash_history",
		"/home/${SSH_USER}/.bash_history",
		"/dev/.udev/",
		"/lib/udev/rules.d/75-persistent-net-generator.rules",
		"/var/lib/dhcp3/*",
		"/var/lib/dhcp/*",
		"/tmp/*",
	}
	tmpFiles = []string{
		"/var/log/lastlog",
		"/var/log/wtmp",
		"/var/log/btmp",
		"/etc/machine-id",
	}
	Ubuntu = ubuntu{}
	Debian = debian{}
)

type ubuntu struct {
}

func (u ubuntu) GetPackageManager() PackageManager {
	return AptPackageManager{}
}

func (u ubuntu) GetTag() string {
	return "ubuntu"
}

func (u ubuntu) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "Ubuntu")
}

func (u ubuntu) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

type debian struct {
}

func (d debian) GetPackageManager() PackageManager {
	return AptPackageManager{}
}

func (d debian) GetTag() string {
	return "debian"
}

func (d debian) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "Debian")
}

func (d debian) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
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
