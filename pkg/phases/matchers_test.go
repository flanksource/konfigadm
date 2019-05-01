package phases

import (
	"fmt"
	"strings"

	"github.com/onsi/gomega/types"
)

func MatchCommand(expected interface{}) types.GomegaMatcher {
	return &CommandMatcher{
		expected: expected,
	}
}

type CommandMatcher struct {
	expected   interface{}
	commands   []Command
	filesystem Filesystem
	err        error
}

func (matcher *CommandMatcher) Match(actual interface{}) (success bool, err error) {
	sys, ok := actual.(*Config)
	if !ok {
		return false, fmt.Errorf("CommandMatcher matcher expects a SystemConfig")
	}
	_, commands, err := sys.ApplyPhases()
	return strings.Contains(commands, matcher.expected.(string)), nil

}

func (matcher *CommandMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected \t%#v to contain command: \t%#v", actual, matcher.expected)
}

func (matcher *CommandMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected \t%#v to NOT contain command: \t%#v", actual, matcher.expected)
}

func ContainPackage(expected interface{}) types.GomegaMatcher {
	return &PackageMatcher{
		expected: expected,
	}
}

type PackageMatcher struct {
	expected   interface{}
	commands   []Command
	filesystem Filesystem
	err        error
}

func (matcher *PackageMatcher) Match(actual interface{}) (success bool, err error) {
	sys, ok := actual.(*Config)

	if !ok {
		return false, fmt.Errorf("PackageMatcher matcher expects a SystemConfig")
	}
	_, _, err = sys.ApplyPhases()
	if err != nil {
		return false, fmt.Errorf("Error applying phases: %s", err)
	}
	for _, p := range sys.Packages {
		if p.Name == matcher.expected {
			return true, nil
		}
	}
	return false, nil
}

func (matcher *PackageMatcher) FailureMessage(actual interface{}) (message string) {
	cfg, _ := actual.(*Config)
	return fmt.Sprintf("Expected %s to contain package: %#v", cfg.Packages, matcher.expected)
}

func (matcher *PackageMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	cfg, _ := actual.(*Config)
	return fmt.Sprintf("Expected %s to NOT contain Package: \t%#v", cfg.Packages, matcher.expected)
}
