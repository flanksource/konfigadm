package phases

import (
	"fmt"
	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
	"strings"
)

type TdnfPackageManager struct{}

func (p TdnfPackageManager) Install(pkg ...string) types.Commands {
	return types.NewCommand(fmt.Sprintf("tdnf install -y %s", strings.Join(pkg, " ")))
}

func (p TdnfPackageManager) Update() types.Commands {
	return types.NewCommand("tdnf makecache")
}

func (p TdnfPackageManager) Uninstall(pkg ...string) types.Commands {
	return types.NewCommand(fmt.Sprintf("tdnf remove -y %s", strings.Join(pkg, " ")))
}

func (p TdnfPackageManager) Mark(pkg ...string) types.Commands {
	return types.Commands{}
}

func (p TdnfPackageManager) CleanupCaches() types.Commands {
	return types.Commands{}
}

func (p TdnfPackageManager) GetInstalledVersion(pkg string) string {
	pkg = strings.Split(pkg, "=")[0]
	var version, release, arch string
	_, ok := utils.SafeExec("tdnf info " + pkg)
	if !ok {
		log.Debugf("No matching package available in db for: %s", pkg)
		return ""
	}

	stdout, ok := utils.SafeExec("tdnf info installed " + pkg)
	if !ok {
		log.Debugf("%s package available in db but not installed", pkg)
		return ""
	}

	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "Version") {
			version = strings.Split(line, ": ")[1]
		}
		if strings.HasPrefix(line, "Release") {
			release = fmt.Sprintf("-%s", strings.Split(line, ": ")[1])
		}
		if strings.HasPrefix(line, "Arch") {
			arch = fmt.Sprintf(".%s", strings.Split(line, ": ")[1])
		}
	}
	retval := fmt.Sprintf("%s%s%s", version, release, arch)
	if retval == "" {
		log.Debugf("Unable to find version info in " + stdout)
	}
	return retval
}

func (p TdnfPackageManager) GetKernelPackageNames(version string) (string, string) {
	// currently tdnf is only in photon, which does not ship rpm by default
	// Manually mapping %{dist} to os info for now
	var dist, arch string
	read, ok := utils.SafeExec("cat /etc/os-release")
	if !ok {
		return "", ""
	}
	distinfo := make(map[string]string)
	for _, line := range strings.Split(read, "\n") {
		linearr := strings.Split(line, "=")
		if len(linearr) != 2 {
			continue
		}
		distinfo[linearr[0]] = linearr[1]
	}
	if distinfo["ID"] != "photon" {
		return "", ""
	}
	if strings.HasPrefix(distinfo["VERSION_ID"], "3.") {
		dist = ".ph3"
	} else if strings.HasPrefix(distinfo["VERSION_ID"], "2.") {
		dist = ".ph2"
	} else {
		return "", ""
	}
	read, ok = utils.SafeExec("uname -p")
	if !ok {
		arch = "x86_64"
	} else {
		arch = strings.TrimSuffix(read, "\n")
	}
	kernelname := fmt.Sprintf("linux-%s%s.%s", version, dist, arch)
	headername := fmt.Sprintf("linux-devel-%s%s.%s", version, dist, arch)
	return kernelname, headername
}

func (p TdnfPackageManager) AddRepo(url string, channel string, versionCodeName string, name string, gpgKey string, extraArgs map[string]string) types.Commands {
	repo := fmt.Sprintf(
		`[%s]
name=%s
baseurl=%s
enabled=1
`, name, name, url)

	if gpgKey != "" {
		repo += fmt.Sprintf(`gpgcheck=1
repo_gpgcheck=1
gpgkey=%s
`, gpgKey)
	} else {
		repo += `
gpgcheck=0
`
	}

	for k, v := range extraArgs {
		repo += fmt.Sprintf("%s = %s\n", k, v)
	}

	return types.NewCommand(fmt.Sprintf(`cat <<EOF >/etc/yum.repos.d/%s.repo
%s
EOF`, name, repo))
}
