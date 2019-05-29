package os

//OS provides an abstraction over different operating systems
type OS interface {

	// GetVersionCodeName returns the distributions version codename e.g. bionic, xenial, squeeze
	GetVersionCodeName() string

	//GetPackageManager returns the packagemanager used by the OS
	GetPackageManager() PackageManager

	//GetTags returns all the tags to which this OS applies
	/*TODO GetTag doesn't return a Tag directly as it would create an import cycle**/
	GetTags() []string

	//DetectAtRuntime will detect if it is compatible with the current running OS
	DetectAtRuntime() bool
}

type OperatingSystemList []OS

//SupportedOperatingSystems is a list of all supported OS's, used primarily for detecting runtime flags
var SupportedOperatingSystems = OperatingSystemList{
	Debian,
	Redhat,
	Ubuntu,
	RHEL,
	Centos,
	AmazonLinux,
}

//BaseOperatingSystems is the list of base distributions that are supported, which is currently only debian and redhat
var BaseOperatingSystems = OperatingSystemList{
	Debian,
	Redhat,
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
