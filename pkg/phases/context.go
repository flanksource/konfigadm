package phases

import (
	"os"
	"strings"

	"github.com/flosch/pongo2"
	log "github.com/sirupsen/logrus"

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

func interpolateSlice(c *SystemContext, val []string) []string {
	var out []string
	for _, v := range val {
		out = append(out, interpolate(c, v))
	}
	return out
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
	if len(template) == 0 {
		return ""
	}
	if strings.Contains(template, "lookup(") {
		log.Tracef("ansible lookups not supported %s", template)
		return template
	}

	if strings.HasPrefix(template, "$") && os.Getenv(template[1:]) != "" {
		return os.Getenv(template[1:])
	}

	template = ConvertSyntaxFromJinjaToPongo(template)
	tpl, err := pongo2.FromString(template)
	if err != nil {
		log.Tracef("error parsing template %s: %v", template, err)
		return template
	}
	out, err := tpl.Execute(vars)
	if err != nil {
		log.Tracef("error executing template %s: %v", template, err)
		return template
	}
	return out
}
