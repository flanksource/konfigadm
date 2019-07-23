package phases

import (
	"strings"

	"github.com/flosch/pongo2"
	. "github.com/moshloop/konfigadm/pkg/types"
	. "github.com/moshloop/konfigadm/pkg/utils"
)

var Context Phase = context{}

type context struct{}

func (p context) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}
	sys.Environment = ToStringMap(interpolateMap(ctx, sys.Environment))
	sys.Files = ToStringMap(interpolateMap(ctx, sys.Files))
	sys.Templates = ToStringMap(interpolateMap(ctx, sys.Templates))

	return commands, files, nil
}

func interpolate(c *SystemContext, s string) string {
	return interpolateString(s, c.Vars)
}

func interpolateMap(c *SystemContext, val map[string]string) map[string]interface{} {
	var out = make(map[string]interface{})
	for k, v := range val {
		out[k] = interpolate(c, v)
	}
	return out
}

func ConvertSyntaxFromJinjaToPongo(template string) string {
	// jinja used filter(arg), pongo uses filter:arg
	template = strings.Replace(template, "(", ":", -1)
	template = strings.Replace(template, ")", "", -1)
	return template
}

func interpolateString(template string, vars map[string]interface{}) string {
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
