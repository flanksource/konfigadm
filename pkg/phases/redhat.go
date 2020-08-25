package phases

import (
	"strings"

	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
)

var (
	Redhat           = redhat{}
	Centos           = centos{}
	Fedora           = fedora{}
	RedhatEnterprise = rhel{}
	AmazonLinux      = amazonLinux{}
)

type redhat struct {
}

func (r redhat) GetPackageManager() types.PackageManager {
	return YumPackageManager{}
}

func (r redhat) GetTags() []types.Flag {
	return []types.Flag{types.REDHAT, types.REDHAT_LIKE}
}

func (r redhat) DetectAtRuntime() bool {
	id, ok := utils.IniToMap("/etc/os-release")["ID"]
	return ok && id == "rhel"
}

func (r redhat) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

type fedora struct {
}

func (r fedora) GetPackageManager() types.PackageManager {
	return DnfPackageManager{}
}

func (r fedora) GetTags() []types.Flag {
	return []types.Flag{types.FEDORA}
}

func (r fedora) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "fedoraproject")
}

func (r fedora) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

type centos struct {
}

func (c centos) GetPackageManager() types.PackageManager {
	return YumPackageManager{}
}

func (c centos) GetTags() []types.Flag {
	return []types.Flag{types.CENTOS, types.REDHAT_LIKE}
}

func (c centos) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "CentOS")
}

func (c centos) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

type rhel struct {
}

func (r rhel) GetPackageManager() types.PackageManager {
	return YumPackageManager{}
}

func (r rhel) GetTags() []types.Flag {
	return []types.Flag{types.RHEL, types.REDHAT_LIKE}
}

func (r rhel) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "RHEL")
}

func (r rhel) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

type amazonLinux struct {
}

func (a amazonLinux) GetPackageManager() types.PackageManager {
	return YumPackageManager{}
}

func (a amazonLinux) GetTags() []types.Flag {
	return []types.Flag{types.AMAZON_LINUX, types.REDHAT_LIKE}
}

func (a amazonLinux) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "Amazon Linux")
}

func (a amazonLinux) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}
