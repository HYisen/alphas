package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	pack, err := Pack(&struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}{
		Labels:     []string{"a", "b"},
		MaxResults: 5,
		Exact:      false,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pack)
}

func Pack(ptr interface{}) (string, error) {
	var sb strings.Builder

	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		if i != 0 {
			sb.WriteString("&")
		}

		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}

		str, err := gen(name, v.Field(i))
		if err != nil {
			log.Println(err)
			continue
		}
		sb.WriteString(str)
	}

	return sb.String(), nil
}

func gen(name string, v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		var limbs []string
		for i := 0; i < v.Len(); i++ {
			s, err := gen(name, v.Index(i))
			if err != nil {
				log.Println(err)
				continue
			}
			limbs = append(limbs, s)
		}
		return strings.Join(limbs, "&"), nil
	case reflect.String:
		s := v.String()
		if !strings.Contains(s, "@") {
			return "", fmt.Errorf("str %s is not email", s)
		}
		return name + "=" + s, nil
	case reflect.Int:
		return name + "=" + strconv.FormatInt(v.Int(), 10), nil
	case reflect.Bool:
		return name + "=" + strconv.FormatBool(v.Bool()), nil
	default:
		return "", fmt.Errorf("unsupported kind %v", v.Kind())
	}
}

func search(resp http.ResponseWriter, req *http.Request) {
	var data struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	data.MaxResults = 10
	if err := Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(resp, "Search: %+v\n", data)
}

func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		if num, err := strconv.ParseInt(value, 10, 64); err != nil {
			return err
		} else {
			v.SetInt(num)
		}
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("unsupported king %s", v.Type())
	}
	return nil
}
