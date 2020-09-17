package phases

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
)

var Kernel types.AllPhases = kernel{}

type kernel struct{}

func (p kernel) ApplyPhase(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	commands := types.Commands{}
	files := types.Filesystem{}
	var os OS
	var err error
	if len(*sys.Kernel) == 0 {
		return []types.Command{}, files, nil
	}
	// Ensure repo caches are up to date before trying to install kernel
	for _, os := range SupportedOperatingSystems {
		commands = *commands.Append(os.GetPackageManager().Update().WithTags(os.GetTags()...))
	}
	for _, kern := range *sys.Kernel {
		var oslist []OS
		if len(kern.Flags) > 0 {
			// Kernel version is tagged for a specific OS (family)
			// Add specific kernel and headers package
			for _, flag := range kern.Flags {
				if os, err = GetOSForTag(flag); err != nil {
					continue
				}
				oslist = append(oslist, os)
			}
		} else {
			// Input is not tagged
			// Add for all supported OSes
			oslist = SupportedOperatingSystems
		}
		for _, os := range oslist {
			commandtags := os.GetTags()
			// rpm based systems will fail to install kernels in containers
			if types.MatchesAny(commandtags, []types.Flag{types.FEDORA}) {
				commandtags = append(commandtags, types.NOT_CONTAINER)
			}
			kernelname, kernelheadername := os.GetKernelPackageNames(kern.Version)
			commands = *commands.Append(os.GetPackageManager().Install(kernelname, kernelheadername).WithTags(commandtags...))
			commands = *commands.Append(os.GetPackageManager().Mark(kernelname, kernelheadername).WithTags(commandtags...))
			commands = *commands.Append(os.UpdateDefaultKernel(kern.Version).WithTags(append(os.GetTags(), types.NOT_CONTAINER)...))

		}
	}
	return commands.Merge(), files, nil
}

func (p kernel) Verify(cfg *types.Config, results *types.VerifyResults, flags ...types.Flag) bool {
	verify := true
	var os OS
	var err error
	if os, err = GetOSForTag(flags...); err != nil {
		results.Fail("Unable to find OS for tags %s", flags)
		return false
	}
	if len(*cfg.Kernel) == 0 {
		return verify
	}
	incontainer := types.MatchesAny(flags, []types.Flag{types.CONTAINER})
	// dockerskip := types.MatchesAny(flags, []types.Flag{types.FEDORA})
	dockerskip := false
	// An improperly configured input (eg different tags for both redhatlike and
	// centos then run on centos) will fail to verify since only one of the
	// two kernels can be the default
	for _, kern := range *cfg.Kernel {
		kernelpackage, kernelheaderpackage := os.GetKernelPackageNames(kern.Version)
		// Certain containers will skip kernel install, skip verification
		if !(incontainer && dockerskip) {
			installed := os.GetPackageManager().GetInstalledVersion(kernelpackage)
			if installed != "" {
				results.Pass("%s is installed", kernelpackage)
			} else {
				results.Fail("%s is not installed", kernelpackage)
				verify = verify && false
			}
			installed = os.GetPackageManager().GetInstalledVersion(kernelheaderpackage)
			if installed != "" {
				results.Pass("%s is installed", kernelheaderpackage)
			} else {
				results.Fail("%s is not installed", kernelheaderpackage)
				verify = verify && false
			}
		} else {
			results.Skip("Cannot test inside a container")
		}
		// Ignore grub config in containers
		if !incontainer {
			if os.ReadDefaultKernel() == kern.Version {
				results.Pass("%s is the default kernel", kern.Version)
			} else {
				results.Fail("%s is not the default kernel", kern.Version)
				verify = verify && false
			}
		} else {
			results.Skip("Cannot test inside a container")
		}
	}

	return verify
}

func (p kernel) ProcessFlags(sys *types.Config, flags ...types.Flag) {
	minified := []types.KernelInput{}
	for _, pkg := range *sys.Kernel {
		if types.MatchesAny(flags, pkg.Flags) {
			minified = append(minified, pkg)
		}
	}
	sys.Kernel = &minified
}

// Interface and corresponding structs for reading/writing grub config with grub.conf or grubby
type GrubConfigManager interface {
	ReadDefaultKernel() (string, bool)
	UpdateDefaultKernel(version string) types.Commands
}

// Generally used by yum systems
type GrubbyManager struct{}

func (p GrubbyManager) ReadDefaultKernel() (string, bool) {
	grubbyout, ok := utils.SafeExec("grubby --default-kernel")
	if ok {
		// Strip dist and arch flags
		for i := 0; i < 2; i++ {
			grubbyout = grubbyout[:strings.LastIndex(grubbyout, ".")]
		}
		re := regexp.MustCompile("/boot/vmlinuz-(.*)")
		match := re.FindStringSubmatch(grubbyout)
		if len(match) > 1 {
			return match[1], ok
		}
		return "", ok
	}
	return "", ok
}

func (p GrubbyManager) UpdateDefaultKernel(version string) types.Commands {
	commands := types.Commands{}
	commands.Add(fmt.Sprintf("grubby --set-default=/boot/vmlinuz-%s$(rpm --eval %%{dist}).$(uname -p)", version))
	return commands
}

// Generally used by apt systems
type GrubConfManager struct{}

func (p GrubConfManager) ReadDefaultKernel() (string, bool) {
	return utils.SafeExec("cat /etc/default/grub | grep '^GRUB_DEFAULT=' | cut -d '=' -f 2")
}

func (p GrubConfManager) UpdateDefaultKernel(version string) types.Commands {
	commands := types.Commands{}
	commands.Add(fmt.Sprintf("sed 's/^GRUB_DEFAULT=.*/GRUB_DEFAULT=\"%s\"/' -i /etc/default/grub", version))
	commands.Add("update-grub")
	return commands
}
