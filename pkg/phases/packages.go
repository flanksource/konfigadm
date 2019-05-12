package phases

import (
	operatingSystem "github.com/moshloop/configadm/pkg/os"

	. "github.com/moshloop/configadm/pkg/types"
)

var Packages AllPhases = packages{}

type packages struct{}

func (p packages) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}

	for _, _os := range operatingSystem.BaseOperatingSystems {

		if len(*sys.PackageRepos) > 0 || len(*sys.Packages) > 0 {
			commands = append(commands, Command{
				Cmd:   _os.GetPackageManager().Setup(),
				Flags: GetTags(_os),
			})
		}

		for _, repo := range *sys.PackageRepos {
			if repo.GPGKey != "" {
				commands = append(commands, Command{
					Cmd:   _os.GetPackageManager().AddKey(repo.GPGKey),
					Flags: GetTags(_os),
				})
			}
			if repo.URL != "" {
				commands = append(commands, Command{
					Cmd:   _os.GetPackageManager().AddRepo(repo.URL, repo.Channel, _os.GetVersionCodeName()),
					Flags: GetTags(_os),
				})
			}
			commands = append(commands, Command{
				Cmd:   _os.GetPackageManager().Update(),
				Flags: GetTags(_os),
			})
		}

		install := []string{}
		uninstall := []string{}
		mark := []string{}
		for _, p := range *sys.Packages {
			if p.Uninstall {
				uninstall = append(uninstall, p.Name)
			} else {
				install = append(install, p.Name)
			}
		}

		if len(install) > 0 {
			// update package repos before installing
			commands = append(commands, Command{
				Cmd:   _os.GetPackageManager().Update(),
				Flags: GetTags(_os),
			})
			commands = append(commands, Command{
				Cmd:   _os.GetPackageManager().Install(install...),
				Flags: GetTags(_os),
			})

		}
		if len(uninstall) > 0 {
			commands = append(commands, Command{
				Cmd:   _os.GetPackageManager().Uninstall(install...),
				Flags: GetTags(_os),
			})

		}

		if len(mark) > 0 {
			commands = append(commands, Command{
				Cmd:   _os.GetPackageManager().Mark(install...),
				Flags: GetTags(_os),
			})
		}
	}

	return commands, files, nil
}
func (p packages) ProcessFlags(sys *Config, flags ...Flag) {
	minified := []Package{}
	for _, pkg := range *sys.Packages {
		if MatchAll(flags, pkg.Flags) {
			minified = append(minified, pkg)
		}
	}
	sys.Packages = &minified
}

func (p packages) Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool {
	verify := true
	os := cfg.Context.OS
	for _, p := range *cfg.Packages {
		installed := os.GetPackageManager().GetInstalledVersion(p.Name)
		if p.Uninstall {
			if installed == "" {
				results.Pass("%s is not installed", p)
			} else {
				results.Fail("%s-%s should not be installed", p, installed)
				verify = false
			}
		} else if p.Version == "" && installed != "" {
			results.Pass("%s-%s is installed", p, installed)
		} else if p.Version == "" && installed == "" {
			results.Fail("%s is not installed, any version required", p)
			verify = false
		} else if installed == p.Version {
			results.Pass("%s-%s is installed", p, installed)
		} else {
			results.Fail("%s-%s is installed, but not the correct version: %s", p.Name, installed, p.Version)
			verify = false
		}
	}

	return verify
}
