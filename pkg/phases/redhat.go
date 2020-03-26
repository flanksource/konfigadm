package phases

import (
	"strings"

	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint, stylecheck
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

func (r redhat) GetPackageManager() PackageManager {
	return YumPackageManager{}
}

func (r redhat) GetTags() []Flag {
	return []Flag{REDHAT, REDHAT_LIKE}
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

func (r fedora) GetPackageManager() PackageManager {
	return DnfPackageManager{}
}

func (r fedora) GetTags() []Flag {
	return []Flag{FEDORA}
}

func (r fedora) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "fedoraproject")
}

func (r fedora) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}

type centos struct {
}

func (c centos) GetPackageManager() PackageManager {
	return YumPackageManager{}
}

func (c centos) GetTags() []Flag {
	return []Flag{CENTOS, REDHAT_LIKE}
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

func (r rhel) GetTags() []Flag {
	return []Flag{RHEL, REDHAT_LIKE}
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

func (a amazonLinux) GetTags() []Flag {
	return []Flag{AMAZON_LINUX, REDHAT_LIKE}
}

func (a amazonLinux) DetectAtRuntime() bool {
	return strings.Contains(utils.SafeRead("/etc/os-release"), "Amazon Linux")
}

func (a amazonLinux) GetVersionCodeName() string {
	return utils.IniToMap("/etc/os-release")["VERSION_CODENAME"]
}
