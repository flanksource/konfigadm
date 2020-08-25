package phases

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type AptPackageManager struct {
}

func (p AptPackageManager) Install(pkg ...string) types.Commands {
	return types.NewCommand("DEBIAN_FRONTEND=noninteractive apt-get install -y --allow-downgrades --allow-change-held-packages --no-install-recommends " + strings.Join(utils.ReplaceAllInSlice(pkg, "==", "="), " "))
}

func (p AptPackageManager) Uninstall(pkg ...string) types.Commands {
	return types.NewCommand("apt-get purge -y " + strings.Join(utils.SplitAllInSlice(pkg, "=", 0), " "))
}

func (p AptPackageManager) Mark(pkg ...string) types.Commands {
	return types.NewCommand("apt-mark hold " + strings.Join(utils.SplitAllInSlice(pkg, "=", 0), " "))
}

func (p AptPackageManager) ListInstalled() string {
	return "dpkg --get-selections"
}

func (p AptPackageManager) Update() types.Commands {
	return types.NewCommand("apt-get update")
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

	log.Tracef("package: %s, status: %s,  version: %s", pkg, status, version)
	if status != "installed" {
		log.Debugf("%s is in db, but is not installed: %s", pkg, status)
		return ""
	}
	return version
}

func (p AptPackageManager) AddRepo(uri string, channel string, versionCodeName string, name string, gpgKey string, extraArgs map[string]string) types.Commands {
	cmds := &types.Commands{}
	if channel == "" {
		channel = "main"
	}

	cmds = cmds.AddDependency("(test \"$(ls /var/lib/apt/lists)\" =  \"\"  && apt-get update) || true")

	if versionCodeName == "" {
		cmds = cmds.AddDependency("which lsb_release 2>&1 > /dev/null || apt-get install -y lsb-release")
		versionCodeName = "$(lsb_release -cs)"
	}
	if name == "" {
		_uri, _ := url.Parse(uri)
		name = filepath.Base(_uri.Path)
	}

	if strings.HasPrefix(uri, "https://") {
		cmds = cmds.AddDependency("test -e /usr/share/doc/apt-transport-https || apt-get install -y apt-transport-https")
	}
	cmds = cmds.
		AddDependency("which curl 2>&1 > /dev/null || apt-get install -y curl")

	if gpgKey != "" {
		cmds = cmds.
			AddDependency("which gpg2 2>&1 > /dev/null || apt-get install -y gnupg")
		cmds = cmds.Add(fmt.Sprintf("curl -skL \"%s\" | apt-key add -", gpgKey))
	}

	return *cmds.Add(fmt.Sprintf("echo deb [arch=amd64] %s %s %s > /etc/apt/sources.list.d/%s.list", uri, versionCodeName, channel, name))
}

func (p AptPackageManager) CleanupCaches() types.Commands {
	set := types.Commands{}
	return *set.Add("apt-get -y autoremove --purge", "apt-get -y clean", "apt-get -y autoclean")
}
