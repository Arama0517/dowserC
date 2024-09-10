package mod

import (
	"bufio"
	"reflect"
	"regexp"
	"strings"
)

var regexPattern = regexp.MustCompile(`^(\w+)\s*=\s*(.*)`)

// Unmarshal 解析 Paradox 脚本
func Unmarshal(value []byte, dest interface{}) error {
	sv := string(value)
	v := reflect.ValueOf(dest).Elem()
	tagFieldMap := getTagFieldMap(v)

	scanner := bufio.NewScanner(strings.NewReader(sv))
	var currentTag string
	insideBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if insideBlock {
			if line == "}" {
				insideBlock = false
				currentTag = ""
				continue
			}
			processBlockLine(line, currentTag, tagFieldMap)
			continue
		}

		if match := regexPattern.FindStringSubmatch(line); len(match) == 3 {
			key, value := match[1], strings.Trim(match[2], " \t\n\"")
			field, ok := tagFieldMap[key]
			if !ok {
				continue
			}
			if field.Kind() == reflect.String {
				field.SetString(value)
			} else if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.String {
				currentTag = key
				insideBlock = true
				if field.IsNil() {
					field.Set(reflect.MakeSlice(field.Type(), 0, 0))
				}
			}
		}
	}

	return scanner.Err()
}

func getTagFieldMap(v reflect.Value) map[string]reflect.Value {
	t := v.Type()
	tagFieldMap := make(map[string]reflect.Value)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("paradox")
		if tag != "" {
			tagFieldMap[tag] = v.Field(i)
		}
	}
	return tagFieldMap
}

func processBlockLine(line, currentTag string, tagFieldMap map[string]reflect.Value) {
	if field, ok := tagFieldMap[currentTag]; ok && field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.String {
		field.Set(reflect.Append(field, reflect.ValueOf(strings.Trim(line, " \t\n\""))))
	}
}
