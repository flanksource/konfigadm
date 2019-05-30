package os

import (
	"strings"
)

type PackageManager interface {
	Install(pkg ...string) string
	Uninstall(pkg ...string) string
	Mark(pkg ...string) string
	AddRepo(url string, channel string, versionCodeName string, name string, gpgKey string) string
	GetInstalledVersion(pkg string) string
	CleanupCaches() string
	Update() string
	Setup() string
}

func CompareVersions(version string, compareTo string) bool {
	if strings.Contains(compareTo, "==") {
		compareTo = strings.Split(compareTo, "==")[1]
	} else if strings.Contains(compareTo, "=") {
		compareTo = strings.Split(compareTo, "=")[1]
	}
	return version == compareTo
}
