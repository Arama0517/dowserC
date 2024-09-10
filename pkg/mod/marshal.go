package mod

import (
	"fmt"
	"reflect"
	"strings"
)

// Marshal 将数据结构编码为 Paradox 脚本
func Marshal(src interface{}) ([]byte, error) {
	var sb strings.Builder

	v := reflect.ValueOf(src)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct but got %s", v.Kind())
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("paradox")
		if tag == "" {
			continue
		}

		switch value.Kind() {
		case reflect.String:
			sb.WriteString(fmt.Sprintf("%s=\"%s\"\n", tag, value.String()))
		case reflect.Slice:
			if value.Type().Elem().Kind() == reflect.String {
				if value.Len() > 0 {
					sb.WriteString(fmt.Sprintf("%s = {\n", tag))
					for j := 0; j < value.Len(); j++ {
						sb.WriteString(fmt.Sprintf("\t\"%s\"\n", value.Index(j).String())) // 使用制表符缩进
					}
					sb.WriteString("}\n")
				}
			}
		default:
			return nil, fmt.Errorf("unsupported field type: %s", value.Kind())
		}
	}

	return []byte(sb.String()), nil
}
