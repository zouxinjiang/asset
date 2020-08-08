package cobra_args_parser

import (
	"github.com/spf13/cobra"
	"reflect"
	"strconv"
	"strings"
)

type tag string

func (c tag) GetName() string {
	items := strings.Split(string(c), ";")
	for _, item := range items {
		kv := strings.SplitN(item, ":", 2)
		if len(kv) == 2 && strings.ToLower(kv[0]) == "name" {
			return kv[1]
		}
	}
	return ""
}

func (c tag) GetDefault() string {
	items := strings.Split(string(c), ";")
	for _, item := range items {
		kv := strings.SplitN(item, ":", 2)
		if len(kv) == 2 && strings.ToLower(kv[0]) == "default" {
			return kv[1]
		}
	}
	return ""
}

func (c tag) GetUsage() string {
	items := strings.Split(string(c), ";")
	for _, item := range items {
		kv := strings.SplitN(item, ":", 2)
		if len(kv) == 2 && strings.ToLower(kv[0]) == "usage" {
			return kv[1]
		}
	}
	return ""
}

func InitArgs(arg interface{}, cmd *cobra.Command) {
	typ := reflect.TypeOf(arg)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		var (
			name  = strings.ToLower(f.Name)
			value = ""
			usage = name
		)
		tg := tag(f.Tag.Get("cmd"))
		value = tg.GetDefault()
		usage = tg.GetUsage()
		tmp := tg.GetName()
		if tmp != "" {
			name = tmp
		}
		switch f.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			fallthrough
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			v, _ := strconv.ParseInt(value, 10, 64)
			cmd.Flags().Int64(name, v, usage)
		case reflect.String:
			cmd.Flags().String(name, value, usage)
		}
	}
}

func ParseArgs(parg interface{}, cmd *cobra.Command) {
	typ := reflect.TypeOf(parg)
	if typ.Kind() != reflect.Ptr {
		return
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		return
	}
	val := reflect.ValueOf(parg)
	if !val.IsValid() {
		return
	}
	val = val.Elem()
	for i := 0; i < val.NumField(); i++ {
		tg := tag(typ.Field(i).Tag.Get("cmd"))
		var name = ""
		if name = tg.GetName(); name == "" {
			name = strings.ToLower(typ.Field(i).Name)
		}
		switch typ.Field(i).Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			argVal, _ := cmd.Flags().GetInt64(name)
			f := val.Field(i)
			if f.CanSet() {
				f.SetInt(argVal)
			}
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			argVal, _ := cmd.Flags().GetUint64(name)
			f := val.Field(i)
			if f.CanSet() {
				f.SetUint(argVal)
			}
		case reflect.String:
			argVal, _ := cmd.Flags().GetString(name)
			f := val.Field(i)
			if f.CanSet() {
				f.SetString(argVal)
			}
		}
	}
}
