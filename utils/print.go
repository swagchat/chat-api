package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

// func PrintStruct(v interface{}) {
// 	fields := structs.Fields(v)
// 	for _, f := range fields {
// 		r := reflect.ValueOf(f.Value())
// 		if r.IsValid() {
// 			switch r.Kind() {
// 			case reflect.String:
// 				fmt.Printf("%s: %s, ", f.Name(), f.Value())
// 			case reflect.Bool:
// 				fmt.Printf("%s: %t, ", f.Name(), f.Value())
// 			case reflect.Ptr:
// 				PrintStruct(f.Value())
// 			default:
// 				fmt.Println(r.Interface())
// 			}
// 		}
// 	}
// }

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
				b.str.WriteString(fmt.Sprintf("%s: %s,", name, v))
			}
		case reflect.Bool:
			b.str.WriteString(fmt.Sprintf("%s: %t, ", name, v))
		case reflect.Ptr:
			if !r.IsNil() {
				fields := structs.Fields(v)
				for _, f := range fields {
					b.PrintStruct(f.Name(), f.Value())
				}
			}
		}
	}

	return b.str.String()
}
