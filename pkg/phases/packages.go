package phases

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint, stylecheck
	"github.com/flanksource/konfigadm/pkg/utils"
)

var Packages AllPhases = packages{}

type packages struct{}

func (p packages) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	commands := Commands{}
	files := Filesystem{}

	for _, repo := range *sys.PackageRepos {
		var os OS
		var err error
		if os, err = GetOSForTag(repo.Flags...); err != nil {
			return nil, nil, err
		}

		log.Tracef("Adding %s\n", repo)
		if repo.URL != "" || repo.ExtraArgs["mirrorlist"] != "" {
			_commands := os.GetPackageManager().
				AddRepo(repo.URL, repo.Channel, repo.VersionCodeName, repo.Name, repo.GPGKey, repo.ExtraArgs)
			commands.Append(_commands.WithTags(repo.Flags...))
		}
	}
	addPackageCommands(sys, &commands)
	_commands := commands.Merge()
	return _commands, files, nil
}

type packageOperations struct {
	install   []string
	uninstall []string
	mark      []string
	tags      []Flag
}

func appendStrings(slice []string, s string) []string {
	var newSlice []string
	if slice != nil {
		newSlice = slice
	}
	newSlice = append(newSlice, s)
	return newSlice
}

func getKeyFromTags(tags ...Flag) string {
	return fmt.Sprintf("%s", tags)
}

func addPackageCommands(sys *Config, commands *Commands) {
	// package installation can have 2 scenarios:
	// 1) tags specified and we know the package manager
	// 2) tags not specified, so we need to add tagged commands for each base operating system

	// track operations by tag group
	// TODO merge compatible tags, e.g. ubuntu and debian-like tags can be included in the same command
	var managers = make(map[string]packageOperations)

	// handle case 1) tags specified
	for _, p := range *sys.Packages {
		if len(p.Flags) == 0 {
			continue
		}
		var ops packageOperations
		var ok bool
		if ops, ok = managers[getKeyFromTags(p.Flags...)]; !ok {
			ops = packageOperations{tags: p.Flags}
		}
		if p.Uninstall {
			ops.uninstall = appendStrings(ops.uninstall, p.Name)
		} else {
			ops.install = appendStrings(ops.install, p.VersionedName())
		}
		if p.Mark {
			ops.mark = appendStrings(ops.mark, p.Name)
		}

		managers[getKeyFromTags(p.Flags...)] = ops
	}

	// handle case 2) tags not specified
	for _, os := range BaseOperatingSystems {
		for _, p := range *sys.Packages {
			if len(p.Flags) > 0 {
				continue
			}
			var ops packageOperations
			var ok bool
			if ops, ok = managers[getKeyFromTags(os.GetTags()...)]; !ok {
				ops = packageOperations{tags: os.GetTags()}
			}
			if p.Uninstall {
				ops.uninstall = appendStrings(ops.uninstall, p.Name)
			} else {
				ops.install = appendStrings(ops.install, p.VersionedName())
			}
			if p.Mark {
				ops.mark = appendStrings(ops.mark, p.Name)
			}
			managers[getKeyFromTags(os.GetTags()...)] = ops
		}
	}

	// iterate over all tag/op combinations and emit commands
	for _, ops := range managers {
		os, _ := GetOSForTag(ops.tags...)
		commands.Append(os.GetPackageManager().Update().WithTags(ops.tags...))

		if ops.install != nil && len(ops.install) > 0 {
			commands = commands.Append(os.GetPackageManager().Install(ops.install...).WithTags(ops.tags...))
		}
		if ops.uninstall != nil && len(ops.uninstall) > 0 {
			commands = commands.Append(os.GetPackageManager().Uninstall(ops.uninstall...).WithTags(ops.tags...))
		}
		if ops.mark != nil && len(ops.mark) > 0 {
			commands = commands.Append(os.GetPackageManager().Mark(ops.mark...).WithTags(ops.tags...))
		}
	}
}

func (p packages) ProcessFlags(sys *Config, flags ...Flag) {
	minified := []Package{}
	for _, pkg := range *sys.Packages {
		if MatchAll(flags, pkg.Flags) {
			minified = append(minified, pkg)
		}
	}
	sys.Packages = &minified

	minifiedRepos := []PackageRepo{}
	for _, repo := range *sys.PackageRepos {
		if MatchesAny(flags, repo.Flags) {
			minifiedRepos = append(minifiedRepos, repo)
		}
	}
	sys.PackageRepos = &minifiedRepos
}

func (p packages) Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool {
	verify := true
	var os OS
	var err error
	if os, err = GetOSForTag(flags...); err != nil {
		results.Fail("Unable to find OS for tags %s", flags)
		return false
	}

	for _, p := range *cfg.Packages {
		if !MatchesAny(flags, p.Flags) {
			continue
		}
		expandedVersion, _ := utils.SafeExec("echo %s", p.Version)
		expandedVersion = strings.Replace(expandedVersion, "\n", "", -1)
		log.Tracef("Verifying package: %s, version: %s => %s", p.Name, p.Version, expandedVersion)
		installed := os.GetPackageManager().GetInstalledVersion(p.Name)
		if p.Uninstall {
			if installed == "" {
				results.Pass("%s is not installed", p)
			} else {
				results.Fail("%s should not be installed", p)
				verify = false
			}
		} else if p.Version == "" && installed != "" {
			results.Pass("%s is installed with any version: %s", p, installed)
		} else if p.Version == "" && installed == "" {
			results.Fail("%s is not installed, any version required", p)
			verify = false
		} else if strings.HasPrefix(expandedVersion, installed) || strings.HasPrefix(installed, expandedVersion) {
			results.Pass("%s is installed with expected version: %s", p, installed)
		} else {
			results.Fail("%s is installed, but expected %s, got %s", p.Name, expandedVersion, installed)
			verify = false
		}
	}

	return verify
}
