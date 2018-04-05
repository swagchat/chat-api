package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

var builderStr *StringBuilder

type StringBuilder struct {
	str *strings.Builder
}

func NewStringBuilder() *StringBuilder {
	return &StringBuilder{
		str: &strings.Builder{},
	}
}

func (b *StringBuilder) PrintStruct(name string, v interface{}) string {
	r := reflect.ValueOf(v)
	if r.IsValid() {
		switch r.Kind() {
		case reflect.String:
			if v.(string) != "" {
				b.str.WriteString(fmt.Sprintf("%s:%s\t", name, v))
			}
		case reflect.Bool:
			b.str.WriteString(fmt.Sprintf("%s:%t\t", name, v))
		case reflect.Int:
			b.str.WriteString(fmt.Sprintf("%s:%d\t", name, v))
		case reflect.Ptr, reflect.Slice:
			if !r.IsNil() {
				fields := structs.Fields(v)
				var printName string
				for _, f := range fields {
					if name == "config" {
						printName = f.Name()
					} else {
						printName = fmt.Sprintf("%s.%s", name, f.Name())
					}
					b.PrintStruct(printName, f.Value())
				}
			}
		case reflect.Struct:
			fields := structs.Fields(v)
			var printName string
			for _, f := range fields {
				if name == "config" {
					printName = f.Name()
				} else {
					printName = fmt.Sprintf("%s.%s", name, f.Name())
				}
				b.PrintStruct(printName, f.Value())
			}
		}
	}

	return b.str.String()
}
