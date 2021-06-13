package phases

import (
	"fmt"
	"strings"

	"github.com/flanksource/konfigadm/pkg/types"
)

//OS provides an abstraction over different operating systems
type OS interface {

	// GetVersionCodeName returns the distributions version codename e.g. bionic, xenial, squeeze
	GetVersionCodeName() string

	//GetPackageManager returns the packagemanager used by the OS
	GetPackageManager() types.PackageManager

	//GetTags returns all the tags to which this OS applies
	GetTags() []types.Flag

	GetName() string

	//DetectAtRuntime will detect if it is compatible with the current running OS
	DetectAtRuntime() bool

	// Returns the names of the kernel and kernel header packages
	GetKernelPackageNames(version string) (string, string)

	// Calls the version specific grub config logic
	UpdateDefaultKernel(version string) types.Commands

	// Check whether defaul kernel matches specific version
	ReadDefaultKernel() string
}

type OperatingSystemList []OS

//SupportedOperatingSystems is a list of all supported OS's, used primarily for detecting runtime flags
var SupportedOperatingSystems = OperatingSystemList{
	Debian,
	Redhat,
	Ubuntu,
	AmazonLinux,
	RedhatEnterprise,
	Centos,
	Fedora,
	Photon,
}

var OperatingSystems = map[string]OS{
	"ubuntu":      Ubuntu,
	"debian":      Debian,
	"redhat":      Redhat,
	"amazonLinux": AmazonLinux,
	"centos":      Centos,
	"fedora":      Fedora,
	"photon":      Photon,
}

//BaseOperatingSystems is the list of base distributions that are supported, which is currently only debian and redhat
var BaseOperatingSystems = OperatingSystemList{
	Debian,
	Redhat,
	Fedora,
	Photon,
}

func GetOSForTag(tags ...types.Flag) (OS, error) {
	for _, t := range tags {
		for _, os := range SupportedOperatingSystems {
			if strings.HasPrefix(t.Name, os.GetName()) {
				return os, nil
			}
		}
	}
	return nil, fmt.Errorf("Unable to find OS for %s", tags)
}

//Detect returns a list of all compatible operating systems at runtime
func (l OperatingSystemList) Detect() []OS {
	var list []OS
	for _, os := range l {
		if os.DetectAtRuntime() {
			list = append(list, os)
		}
	}
	return list

}
