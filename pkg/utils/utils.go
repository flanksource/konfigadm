package utils

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
)

func Exec(sh string) error {
	cmd := exec.Command("bash", "-c", sh)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
	}

	if !cmd.ProcessState.Success() {
		return errors.New("Failed")
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

func IsSlice(arg interface{}) bool {
	return reflect.ValueOf(arg).Kind() == reflect.Slice
}

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

func StructToIni(s interface{}) string {
	str := ""
	for k, v := range StructToMap(s) {
		str += k + "=" + ToString(v) + "\n"
	}
	return str
}

func MapToIni(Map map[string]string) string {
	str := ""
	for k, v := range Map {
		str += k + "=" + ToString(v) + "\n"
	}
	return str
}

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
