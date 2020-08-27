package phases

import (
	"fmt"
	"strings"

	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type YumPackageManager struct{}

func (p YumPackageManager) Install(pkg ...string) types.Commands {
	arg := strings.Join(pkg, " ")
	// Yum versions are specified using a -, not a =
	arg = strings.Replace(arg, "=", "-", -1)
	return types.NewCommand(fmt.Sprintf("yum install -y %s", arg))
}

func (p YumPackageManager) Update() types.Commands {
	return types.Commands{}
}
func (p YumPackageManager) Uninstall(pkg ...string) types.Commands {
	return types.NewCommand(fmt.Sprintf("yum remove -y %s", strings.Join(pkg, " ")))
}
func (p YumPackageManager) Mark(pkg ...string) types.Commands {
	return types.Commands{}
}

func (p YumPackageManager) CleanupCaches() types.Commands {
	return types.Commands{}
}

func (p YumPackageManager) GetInstalledVersion(pkg string) string {
	pkg = strings.Split(pkg, "=")[0]
	stdout, ok := utils.SafeExec("yum info " + pkg)
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

func (p YumPackageManager) AddRepo(url string, channel string, versionCodeName string, name string, gpgKey string, extraArgs map[string]string) types.Commands {
	repo := fmt.Sprintf(
		`[%s]
name=%s
enabled=1
`, name, name)

	if url != "" {
		repo += fmt.Sprintf("baseurl = %s\n", url)
	}

	if gpgKey != "" {
		repo += fmt.Sprintf(`gpgcheck=1
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
