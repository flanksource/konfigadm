package phases

import (
	"fmt"
	"strings"

	"github.com/flosch/pongo2"
)

var Context Phase = context{}

type context struct{}

func (p context) ApplyPhase(sys *SystemConfig, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}
	sys.Environment = ToStringMap(ctx.InterpolateMap(sys.Environment))
	sys.Files = ToStringMap(ctx.InterpolateMap(sys.Files))
	sys.Templates = ToStringMap(ctx.InterpolateMap(sys.Templates))

	return commands, files, nil
}

func (c SystemContext) Interpolate(s string) string {
	return InterpolateString(s, c.Vars)
}

func (c SystemContext) InterpolateMap(val map[string]string) map[string]interface{} {
	var out = make(map[string]interface{})
	for k, v := range val {
		out[k] = c.Interpolate(v)
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
