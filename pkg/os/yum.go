package os

type YumPackageManager struct{}

func (p YumPackageManager) Install(pkg ...string) string {
	return ""
}

func (p YumPackageManager) Update() string {
	return ""
}
func (p YumPackageManager) Uninstall(pkg ...string) string {
	return ""
}
func (p YumPackageManager) Mark(pkg ...string) string {
	return ""
}
func (p YumPackageManager) ListInstalled() string {
	return ""
}
func (p YumPackageManager) CleanupCaches() string {
	return ""
}

func (p YumPackageManager) GetInstalledVersion(pkg string) string {
	return ""
}

func (p YumPackageManager) AddKey(url string) string {
	return ""
}

func (p YumPackageManager) AddRepo(url string, channel string, versionCodeName string) string {
	return ""
}

func (p YumPackageManager) Setup() string {
	return ""
}
