package os

import (
	"fmt"
	"strings"

	"github.com/moshloop/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type Dpkg struct {
	Package            string `json:"Package"`
	Status             string `json:"Status"`
	Priority           string `json:"Priority"`
	Section            string `json:"Section"`
	InstalledSize      int    `json:"Installed-Size"`
	Maintainer         string `json:"Maintainer"`
	Architecture       string `json:"Architecture"`
	MultiArch          string `json:"Multi-Arch"`
	Version            string `json:"Version"`
	Depends            string `json:"Depends"`
	Suggests           string `json:"Suggests"`
	Description        string `json:"Description"`
	Homepage           string `json:"Homepage"`
	OriginalMaintainer string `json:"Original-Maintainer"`
}

//IsInstalled returns true if the metadata lists the package as installed
func (pkg Dpkg) IsInstalled() bool {
	return pkg.Status == "install ok installed"
}

type AptPackageManager struct {
}

func (p AptPackageManager) Install(pkg ...string) string {
	return "apt-get install -y --no-install-recommends " + strings.Join(utils.ReplaceAllInSlice(pkg, "==", "="), " ")
}

func (p AptPackageManager) Uninstall(pkg ...string) string {
	return "apt-get purge -y " + strings.Join(utils.SplitAllInSlice(pkg, "=", 0), " ")
}

func (p AptPackageManager) Mark(pkg ...string) string {
	return "apt-get mark -y " + strings.Join(utils.SplitAllInSlice(pkg, "=", 0), " ")
}

func (p AptPackageManager) ListInstalled() string {
	return "dpkg --get-selections"
}

func (p AptPackageManager) Update() string {
	return "apt-get update"
}

func (p AptPackageManager) GetInstalledVersion(pkg string) string {
	pkg = strings.Split(pkg, "=")[0]
	stdout, ok := utils.SafeExec("dpkg-query -W -f='${db:Status-Status}\t${Version}' " + pkg)
	if !ok {
		log.Debugf("Failed installation check for %s -> %s", pkg, stdout)
		return ""
	}
	status := strings.Split(stdout, "\t")[0]
	version := strings.Split(stdout, "\t")[1]

	if status != "installed" {
		log.Debugf("%s is in db, but is not installed: %s", pkg, status)
		return ""
	}
	return version
}

func (p AptPackageManager) AddKey(url string) string {
	return fmt.Sprintf("curl -fsSKL %s | sudo apt-key add -", url)
}

func (p AptPackageManager) AddRepo(url string, channel string, versionCodeName string) string {
	if channel == "" {
		channel = "main"
	}
	return fmt.Sprintf("add-apt-repository \"deb [arch=amd64] %s %s %s\"", url, versionCodeName, channel)
}

func (p AptPackageManager) CleanupCaches() string {
	return `apt-get -y autoremove --purge
apt-get -y clean
apt-get -y autoclean`
}

func (p AptPackageManager) Setup() string {
	return "[[ $(which add-apt-repository 2> /dev/null) && $(which curl 2> /dev/null) && $(which sudo 2> /dev/null) ]] || (apt-get update;  apt-get install -y --no-install-recommends --ignore-missing sudo apt-transport-https ca-certificates curl gnupg2 software-properties-common)"
}
