package phases

import (
	"fmt"
	"strings"

	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type DnfPackageManager struct{}

func (p DnfPackageManager) Install(pkg ...string) types.Commands {
	return types.NewCommand(fmt.Sprintf("dnf install -y %s", strings.Join(pkg, " ")))
}

func (p DnfPackageManager) Update() types.Commands {
	return types.Commands{}
}
func (p DnfPackageManager) Uninstall(pkg ...string) types.Commands {
	return types.NewCommand(fmt.Sprintf("dnf remove -y %s", strings.Join(pkg, " ")))
}
func (p DnfPackageManager) Mark(pkg ...string) types.Commands {
	return types.Commands{}
}

func (p DnfPackageManager) CleanupCaches() types.Commands {
	return types.Commands{}
}

func (p DnfPackageManager) GetInstalledVersion(pkg string) string {
	pkg = strings.Split(pkg, "=")[0]
	stdout, ok := utils.SafeExec("dnf info " + pkg)
	if !ok {
		log.Debugf("Failed installation check for %s -> %s", pkg, stdout)
		return ""
	}

	if !strings.Contains(stdout, "Installed Packages") {
		log.Debugf("%s is in db, but is not installed", pkg)
		return ""
	}

	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "Version") {
			return strings.Split(line, ": ")[1]
		}
	}
	log.Debugf("Unable to find version info in " + stdout)
	return ""
}

func (p DnfPackageManager) AddRepo(url string, channel string, versionCodeName string, name string, gpgKey string, extraArgs map[string]string) types.Commands {
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
	}

	for k, v := range extraArgs {
		repo += fmt.Sprintf("%s = %s\n", k, v)
	}

	return types.NewCommand(fmt.Sprintf(`cat <<EOF >/etc/yum.repos.d/%s.repo
%s
EOF`, name, repo))
}
