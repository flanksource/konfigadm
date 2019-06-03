package phases

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	. "github.com/moshloop/konfigadm/pkg/types"
	"github.com/moshloop/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type AptPackageManager struct {
}

func (p AptPackageManager) Install(pkg ...string) Commands {
	return NewCommand("apt-get install -y --no-install-recommends " + strings.Join(utils.ReplaceAllInSlice(pkg, "==", "="), " "))
}

func (p AptPackageManager) Uninstall(pkg ...string) Commands {
	return NewCommand("apt-get purge -y " + strings.Join(utils.SplitAllInSlice(pkg, "=", 0), " "))
}

func (p AptPackageManager) Mark(pkg ...string) Commands {
	cmds := NewCommand("apt-get mark -y " + strings.Join(utils.SplitAllInSlice(pkg, "=", 0), " "))
	return cmds.AddDependency("apt-get install -y apt-mark")
}

func (p AptPackageManager) ListInstalled() string {
	return "dpkg --get-selections"
}

func (p AptPackageManager) Update() Commands {
	return NewCommand("apt-get update")
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

func (p AptPackageManager) AddRepo(uri string, channel string, versionCodeName string, name string, gpgKey string) Commands {
	cmds := Commands{}
	if channel == "" {
		channel = "main"
	}

	if versionCodeName == "" {
		versionCodeName = "$(lsb_release -cs)"
	}
	if name == "" {
		_uri, _ := url.Parse(uri)
		name = filepath.Base(_uri.Path)
	}

	if strings.HasPrefix(uri, "https://") {
		cmds = cmds.AddDependency("[ -e /usr/share/doc/apt-transport-https ] || apt-get install -y apt-transport-https")
	}
	cmds = cmds.
		AddDependency("which sudo 2>&1 > /dev/null || apt-get install -y sudo").
		AddDependency("which curl 2>&1 > /dev/null || apt-get install -y curl").
		AddDependency("which add-apt-repository 2>&1 /dev/null  || apt-get install -y software-properties-common")

	if gpgKey != "" {
		cmds = cmds.Add(fmt.Sprintf("curl -fsSKL %s | sudo apt-key add -", uri))
	}

	return cmds.Add(fmt.Sprintf("add-apt-repository \"deb [arch=amd64] %s %s %s\"", uri, versionCodeName, channel))
}

func (p AptPackageManager) CleanupCaches() Commands {
	set := Commands{}
	return set.Add("apt-get -y autoremove --purge", "apt-get -y clean", "apt-get -y autoclean")
}
