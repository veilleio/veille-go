package veille

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// fillPath replaces {name} segments and returns the parameters left for the query string or body.
func fillPath(template string, params []param) (string, []param, error) {
	rest := make([]param, 0, len(params))
	used := map[string]bool{}
	path := template
	for strings.Contains(path, "{") {
		start := strings.Index(path, "{")
		end := strings.Index(path[start:], "}")
		if end < 0 {
			break
		}
		name := path[start+1 : start+end]
		value, found := "", false
		for _, p := range params {
			if p.name == name && !isEmpty(p.value) {
				text, err := wireString(name, p.value)
				if err != nil {
					return "", nil, err
				}
				value, found = text, true
				break
			}
		}
		if !found {
			return "", nil, &SerializationError{Param: name, Message: fmt.Sprintf("The %q path parameter is required.", name)}
		}
		used[name] = true
		path = path[:start] + url.PathEscape(value) + path[start+end+1:]
	}
	for _, p := range params {
		if !used[p.name] && !isEmpty(p.value) {
			rest = append(rest, p)
		}
	}
	return path, rest, nil
}

func encodeQuery(params []param) (string, error) {
	parts := make([]string, 0, len(params))
	for _, p := range params {
		value := reflect.ValueOf(p.value)
		if value.Kind() == reflect.Slice && value.Type().Elem().Kind() != reflect.Uint8 {
			for i := 0; i < value.Len(); i++ {
				text, err := wireString(p.name, value.Index(i).Interface())
				if err != nil {
					return "", err
				}
				parts = append(parts, url.QueryEscape(p.name)+"="+url.QueryEscape(text))
			}
			continue
		}
		text, err := wireString(p.name, p.value)
		if err != nil {
			return "", err
		}
		parts = append(parts, url.QueryEscape(p.name)+"="+url.QueryEscape(text))
	}
	return strings.ReplaceAll(strings.Join(parts, "&"), "+", "%20"), nil
}

func wireString(name string, value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case bool:
		return strconv.FormatBool(v), nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return "", &SerializationError{Param: name, Message: fmt.Sprintf("The %q parameter must be a finite number.", name)}
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return wireString(name, float64(v))
	case fmt.Stringer:
		return v.String(), nil
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", &SerializationError{Param: name, Message: fmt.Sprintf("The %q parameter could not be serialized.", name), Err: err}
	}
	return string(encoded), nil
}

func isEmpty(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.IsNil() || v.Len() == 0
	}
	return false
}

// String returns a pointer to v, for optional parameters.
func String(v string) *string { return &v }

// Float returns a pointer to v, for optional number parameters.
func Float(v float64) *float64 { return &v }

// Int returns a pointer to v, for optional integer parameters.
func Int(v int64) *int64 { return &v }

// Bool returns a pointer to v, for optional boolean parameters.
func Bool(v bool) *bool { return &v }
