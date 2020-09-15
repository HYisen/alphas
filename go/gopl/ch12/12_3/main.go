package main

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2

	s := struct {
		id    int
		name  string
		score float64
		good  bool
		vec   complex128
	}{17, "alex", 2.14, true, complex(-0.5, 2)}

	fmt.Println(Encode(m))
	fmt.Println(Encode(s))
}

func Encode(v interface{}) string {
	var buf bytes.Buffer
	fmt.Println(reflect.ValueOf(v))
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		log.Fatalln(err)
	}
	return buf.String()
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Int:
		_, _ = fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint:
		buf.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.String:
		_, _ = fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem())
	case reflect.Array, reflect.Slice:
		err := loopWrapper(buf, v, reflect.Value.Len, reflect.Value.Index, nil, nil)
		if err != nil {
			return err
		}
	case reflect.Struct:
		err := loopWrapper(buf, v, reflect.Value.NumField, reflect.Value.Field, genFieldNamePrefix, genProvider(")"))
		if err != nil {
			return err
		}
	case reflect.Map:
		limbs, err := mapSlice(v.MapKeys(), genMapper(v))
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(buf, "(%s)", strings.Join(limbs, " "))
	case reflect.Bool:
		if v.Bool() {
			buf.WriteRune('t')
		} else {
			buf.WriteString("nil")
		}
	case reflect.Float64:
		buf.WriteString(strconv.FormatFloat(v.Float(), 'b', -1, 64))
	case reflect.Complex128:
		c := v.Complex()
		_, _ = fmt.Fprintf(buf, "#(%.1f %.1f)", real(c), imag(c))
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func genProvider(item string) func(v reflect.Value, i int) string {
	return func(v reflect.Value, i int) string { return item }
}

func genMapper(v reflect.Value) func(key reflect.Value) (string, error) {
	return func(key reflect.Value) (string, error) {
		var buf bytes.Buffer

		if err := encode(&buf, key); err != nil {
			return "", err
		}

		buf.WriteRune(' ')

		if err := encode(&buf, v.MapIndex(key)); err != nil {
			return "", err
		}

		return "(" + buf.String() + ")", nil
	}
}

func mapSlice(source []reflect.Value, mapper func(value reflect.Value) (string, error)) ([]string, error) {
	var ret []string
	for _, val := range source {
		if str, err := mapper(val); err != nil {
			return nil, err
		} else {
			ret = append(ret, str)
		}
	}
	return ret, nil
}

func genFieldNamePrefix(v reflect.Value, i int) string {
	return "(" + v.Type().Field(i).Name + " "
}

// String.format("(%s)",
// IntStream
// .range(0, lengthFunc(value))
// .mapToObj(i->prefixGeneratorNullable?(v,i)+extractFunc(v,i)+suffixGeneratorNullable?(v,i)})
// .collect(Collectors.joining(" ")))
func loopWrapper(buf *bytes.Buffer, v reflect.Value,
	lengthFunc func(v reflect.Value) int, extractFunc func(v reflect.Value, i int) reflect.Value,
	prefixGeneratorNullable, suffixGeneratorNullable func(v reflect.Value, i int) string) error {

	buf.WriteRune('(')
	for i := 0; i < lengthFunc(v); i++ {
		if i > 0 {
			buf.WriteRune(' ')
		}
		if prefixGeneratorNullable != nil {
			buf.WriteString(prefixGeneratorNullable(v, i))
		}
		if err := encode(buf, extractFunc(v, i)); err != nil {
			return err
		}
		if suffixGeneratorNullable != nil {
			buf.WriteString(suffixGeneratorNullable(v, i))
		}
	}
	buf.WriteRune(')')
	return nil
}
