package os

import (
	"strings"

	"github.com/moshloop/konfigadm/pkg/utils"
)

var (
	Redhat      = redhat{}
	Centos      = centos{}
	RHEL        = rhel{}
	AmazonLinux = amazonLinux{}
)

type redhat struct {
}

func (r redhat) GetPackageManager() PackageManager {
	return YumPackageManager{}
}

func (r redhat) GetTags() []string {
	return []string{"redhat"}
}

func (r redhat) DetectAtRuntime() bool {
	return false
}

func (r redhat) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

type centos struct {
}

func (c centos) GetPackageManager() PackageManager {
	return YumPackageManager{}
}

func (c centos) GetTags() []string {
	return []string{"centos", "redhat"}
}

func (c centos) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "CentOS")
}

func (c centos) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

type rhel struct {
}

func (r rhel) GetPackageManager() PackageManager {
	return YumPackageManager{}
}

func (r rhel) GetTags() []string {
	return []string{"rhel", "redhat"}
}

func (r rhel) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "RHEL")
}

func (r rhel) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

type amazonLinux struct {
}

func (a amazonLinux) GetPackageManager() PackageManager {
	return YumPackageManager{}
}

func (a amazonLinux) GetTags() []string {
	return []string{"amazonLinux", "redhat"}
}

func (a amazonLinux) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "Amazon Linux")
}

func (a amazonLinux) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}
