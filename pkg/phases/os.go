package phases

import (
	"fmt"

	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint
)

//OS provides an abstraction over different operating systems
type OS interface {

	// GetVersionCodeName returns the distributions version codename e.g. bionic, xenial, squeeze
	GetVersionCodeName() string

	//GetPackageManager returns the packagemanager used by the OS
	GetPackageManager() PackageManager

	//GetTags returns all the tags to which this OS applies
	GetTags() []Flag

	//DetectAtRuntime will detect if it is compatible with the current running OS
	DetectAtRuntime() bool
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
}

var OperatingSystems = map[string]OS{
	"ubuntu":      Ubuntu,
	"debian":      Debian,
	"redhat":      Redhat,
	"amazonLinux": AmazonLinux,
	"centos":      Centos,
	"fedora":      Fedora,
}

//BaseOperatingSystems is the list of base distributions that are supported, which is currently only debian and redhat
var BaseOperatingSystems = OperatingSystemList{
	Debian,
	Redhat,
	Fedora,
}

func GetOSForTag(tags ...Flag) (OS, error) {
	for _, t := range tags {
		for _, os := range SupportedOperatingSystems {
			for _, tag := range os.GetTags() {
				if tag.Name == t.Name {
					return os, nil
				}
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
