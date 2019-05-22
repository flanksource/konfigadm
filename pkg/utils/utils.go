package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	Reset        = "\x1b[0m"
	Red          = "\x1b[31m"
	LightRed     = "\x1b[31;1m"
	Green        = "\x1b[32m"
	LightGreen   = "\x1b[32;1m"
	LightBlue    = "\x1b[34;1m"
	Magenta      = "\x1b[35m"
	LightMagenta = "\x1b[35;1m"
	Cyan         = "\x1b[36m"
	LightCyan    = "\x1b[36;1m"
	White        = "\x1b[37;1m"
	Bold         = "\x1b[1m"
	BoldOff      = "\x1b[22m"
)

//SafeExec executes the sh script and returns the stdout and stderr, errors will result in a nil return only.
func SafeExec(sh string, args ...interface{}) (string, bool) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(sh, args...))
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Debugf("Failed to exec %s, %s %s\n", sh, data, err)
		return "", false
	}

	if !cmd.ProcessState.Success() {
		log.Debugf("Command did not succeed %s\n", sh)
		return "", false
	}
	return string(data), true

}

//Exec runs the sh script and forwards stderr/stdout to the console
func Exec(sh string, args ...interface{}) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(sh, args...))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s failed with %s", sh, err)
	}

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("%s failed to run", sh)
	}
	return nil
}

func takeSliceArg(arg interface{}) (out []interface{}, ok bool) {
	val := reflect.ValueOf(arg)
	if val.Kind() != reflect.Slice {
		return nil, false
	}

	c := val.Len()
	out = make([]interface{}, c)
	for i := 0; i < val.Len(); i++ {
		out[i] = val.Index(i).Interface()
	}
	return out, true
}

//IsSlice returns true if the argument is a slice
func IsSlice(arg interface{}) bool {
	return reflect.ValueOf(arg).Kind() == reflect.Slice
}

//ToString takes an object and tries to convert it to a string
func ToString(i interface{}) string {
	if slice, ok := takeSliceArg(i); ok {
		s := ""
		for _, v := range slice {
			if s != "" {
				s += ", "
			}
			s += ToString(v)
		}
		return s

	}
	switch v := i.(type) {
	case fmt.Stringer:
		return v.String()
	case string:
		return v
	case interface{}:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		// panic(fmt.Sprintf("I don't know about type %T!\n", v))
	}
	return ""
}

//StructToMap takes an object and returns all it's field in a map
func StructToMap(s interface{}) map[string]interface{} {
	values := make(map[string]interface{})
	value := reflect.ValueOf(s)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.CanInterface() {
			v := field.Interface()
			if v != nil && v != "" {
				values[value.Type().Field(i).Name] = v
			}
		}
	}
	return values
}

//StructToIni takes an object and serializes it's fields in INI format
func StructToIni(s interface{}) string {
	str := ""
	for k, v := range StructToMap(s) {
		str += k + "=" + ToString(v) + "\n"
	}
	return str
}

//MapToIni takes a map and converts it into an INI formatted string
func MapToIni(Map map[string]string) string {
	str := ""
	for k, v := range Map {
		str += k + "=" + ToString(v) + "\n"
	}
	return str
}

//IniToMap takes the path to an INI formatted file and transforms it into a map
func IniToMap(path string) map[string]string {
	result := make(map[string]string)
	ini := SafeRead(path)
	for _, line := range strings.Split(ini, "\n") {
		values := strings.Split(line, "=")
		if len(values) == 2 {
			result[values[0]] = values[1]
		}
	}
	return result
}

//GzipFile takes the path to a file and returns a Gzip comppressed byte slic
func GzipFile(path string) ([]byte, error) {
	var buf bytes.Buffer

	w := gzip.NewWriter(&buf)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(contents)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	result := buf.Bytes()
	return result, nil
}

//SafeRead reads a path and returns the text contents or nil,
func SafeRead(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

//ReplaceAllInSlice runs strings.Replace on all elements in a slice and returns the result
func ReplaceAllInSlice(a []string, find string, replacement string) (replaced []string) {
	for _, s := range a {
		replaced = append(replaced, strings.Replace(s, find, replacement, -1))
	}
	return
}

//SplitAllInSlice runs strings.Split on all elements in a slice and returns the results at the given index
func SplitAllInSlice(a []string, split string, index int) (replaced []string) {
	for _, s := range a {
		replaced = append(replaced, strings.Split(s, split)[index])
	}
	return
}
