package phases

import (
	"fmt"
	"strings"

	. "github.com/moshloop/konfigadm/pkg/types"

	"github.com/flosch/pongo2"
)

var Context Phase = context{}

type context struct{}

func (p context) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}
	sys.Environment = ToStringMap(InterpolateMap(ctx, sys.Environment))
	sys.Files = ToStringMap(InterpolateMap(ctx, sys.Files))
	sys.Templates = ToStringMap(InterpolateMap(ctx, sys.Templates))

	return commands, files, nil
}

func Interpolate(c *SystemContext, s string) string {
	return InterpolateString(s, c.Vars)
}

func InterpolateMap(c *SystemContext, val map[string]string) map[string]interface{} {
	var out = make(map[string]interface{})
	for k, v := range val {
		out[k] = Interpolate(c, v)
	}
	return out
}

func ToGenericMap(m map[string]string) map[string]interface{} {
	var out = map[string]interface{}{}
	for k, v := range m {
		out[k] = v
	}
	return out
}

func ToStringMap(m map[string]interface{}) map[string]string {
	var out = make(map[string]string)
	for k, v := range m {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}

func ConvertSyntaxFromJinjaToPongo(template string) string {
	// jinja used filter(arg), pongo uses filter:arg
	template = strings.Replace(template, "(", ":", -1)
	template = strings.Replace(template, ")", "", -1)
	return template
}

func InterpolateString(template string, vars map[string]interface{}) string {
	if strings.Contains(template, "lookup(") {
		// log.Warningf("ansible lookups not supported %s", template)
		return template
	}

	template = ConvertSyntaxFromJinjaToPongo(template)
	tpl, err := pongo2.FromString(template)
	if err != nil {
		// log.Debugf("Error parsing: %s: %v", template, err)
		return template
	}
	out, err := tpl.Execute(vars)
	if err != nil {
		// log.Debugf("Error parsing: %s: %v", template, err)
		return template
	}
	//log.Errorf("%s => %s", template, out)
	return out
}
