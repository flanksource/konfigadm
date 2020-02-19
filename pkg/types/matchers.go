package types

import (
	"fmt"
	"strings"

	"github.com/onsi/gomega/types"
)

func MatchCommand(expected interface{}) types.GomegaMatcher {
	return &CommandMatcher{
		Expected: expected,
	}
}

type CommandMatcher struct {
	Expected   interface{}
	Commands   []Command
	Filesystem Filesystem
	Err        error
}

func (matcher *CommandMatcher) Match(actual interface{}) (success bool, err error) {
	switch v := actual.(type) {
	case *Config:
		_, commands, err := v.ApplyPhases()
		if err != nil {
			return false, err
		}
		for _, cmd := range commands {
			if strings.Contains(cmd.Cmd, matcher.Expected.(string)) {
				return true, nil
			}
		}
	case []Command:
		for _, cmd := range v {
			if strings.Contains(cmd.Cmd, matcher.Expected.(string)) {
				return true, nil
			}
		}
	default:
		return false, fmt.Errorf("CommandMatcher matcher expects a SystemConfig or []Config")
	}
	return false, nil
}

func (matcher *CommandMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected \t%#v to contain command: \t%#v", actual, matcher.Expected)
}

func (matcher *CommandMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected \t%#v to NOT contain command: \t%#v", actual, matcher.Expected)
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
	for _, p := range *sys.Packages {
		if p.Name == matcher.expected {
			return true, nil
		}
	}
	return false, nil
}

func (matcher *PackageMatcher) FailureMessage(actual interface{}) (message string) {
	cfg, _ := actual.(*Config)
	return fmt.Sprintf("Expected %s to contain package: %#v", *cfg.Packages, matcher.expected)
}

func (matcher *PackageMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	cfg, _ := actual.(*Config)
	return fmt.Sprintf("Expected %s to NOT contain Package: \t%#v", *cfg.Packages, matcher.expected)
}
