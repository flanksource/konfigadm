package phases

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
)

var (
	flagProcessors   = make([]FlagProcessor, 0)
	DEBIAN           = Flag{Name: "debian", AlsoMatches: []Flag{UBUNTU}}
	REDHAT           = Flag{Name: "redhat", AlsoMatches: []Flag{RHEL, CENTOS, AMAZON_LINUX}}
	AMAZON_LINUX     = Flag{Name: "amazonLinux"}
	RHEL             = Flag{Name: "rhel"}
	CENTOS           = Flag{Name: "centos"}
	UBUNTU           = Flag{Name: "ubuntu"}
	AWS              = Flag{Name: "aws"}
	VMWARE           = Flag{Name: "vmware"}
	NOT_DEBIAN       = Flag{Name: "!debian", Negates: []Flag{DEBIAN, UBUNTU}}
	NOT_REDHAT       = Flag{Name: "!redhat", Negates: []Flag{REDHAT, CENTOS, AMAZON_LINUX}}
	NOT_CENTOS       = Flag{Name: "!centos", Negates: []Flag{CENTOS}}
	NOT_RHEL         = Flag{Name: "!rhel", Negates: []Flag{RHEL, REDHAT}}
	NOT_UBUNTU       = Flag{Name: "!ubuntu", Negates: []Flag{DEBIAN}}
	NOT_AWS          = Flag{Name: "!aws", Negates: []Flag{AWS}}
	NOT_VMWARE       = Flag{Name: "!vmware", Negates: []Flag{VMWARE}}
	NOT_AMAZON_LINUX = Flag{Name: "!amazonLinux", Negates: []Flag{AMAZON_LINUX}}
	FLAG_MAP         = make(map[string]Flag)
	FLAGS            = []Flag{DEBIAN, REDHAT, AMAZON_LINUX, CENTOS, RHEL, UBUNTU, AWS, VMWARE, NOT_DEBIAN, NOT_REDHAT, NOT_CENTOS, NOT_RHEL, NOT_UBUNTU, NOT_AWS, NOT_VMWARE, NOT_AMAZON_LINUX}
)

func init() {
	for _, flag := range FLAGS {
		name := flag.Name
		FLAG_MAP[name] = flag
		if !strings.HasPrefix(name, "!") && !strings.HasPrefix(name, "+") {
			FLAG_MAP[fmt.Sprintf("+%s", name)] = flag
		}
	}

}

func RegisterFlagProcessor(fn FlagProcessor) {
	log.Printf("Registering Flag Processor %+v\n", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
	flagProcessors = append(flagProcessors, fn)
}

//MinifyWithFlags will evaluate all flags, removing items that do not match flags
func (sys *SystemConfig) MinifyWithFlags(flags ...Flag) {
	for _, processor := range flagProcessors {
		processor(sys, flags...)
	}
}

type FlagProcessor func(cfg *SystemConfig, flags ...Flag)

type Flag struct {
	Name        string
	Negates     []Flag
	AlsoMatches []Flag
}

func (f Flag) String() string {
	return f.Name
}

func (f *Flag) Matches(other Flag) bool {
	if f.Name == other.Name {
		return true
	}
	for _, flag2 := range other.AlsoMatches {
		if f.Matches(flag2) {
			return true
		}
	}

	if len(other.Negates) > 0 {
		for _, flag2 := range other.Negates {
			if f.Matches(flag2) {
				return false
			}
		}
		return true
	}
	return false
}

//MatchAll returns true if all constraints match at least one flag AND none of the constraints negates any flag
func MatchAll(flags []Flag, constraints []Flag) bool {
outer:
	for _, constraint := range constraints {
		for _, flag := range flags {
			if flag.Matches(constraint) {
				continue outer
			}
		}
		return false
	}
	return true
}
