package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func main() {
	type Movie struct {
		Fuck            bool
		Title, Subtitle string
		Year            int
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}
	strangelove := Movie{
		Fuck:     true,
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	var sb strings.Builder
	encoder := json.NewEncoder(&sb)
	err := encoder.Encode(strangelove)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(sb.String())

	fmt.Println(Encode(strangelove))
}

func Encode(v interface{}) string {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		log.Fatalln(err)
	}
	buf.WriteString("\n")
	return buf.String()
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("null")
	case reflect.Int:
		_, _ = fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint:
		buf.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.String:
		_, _ = fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem())
	case reflect.Array, reflect.Slice:
		err := loopWrapper(
			buf,
			v,
			"[",
			"]",
			reflect.Value.Len,
			reflect.Value.Index,
			nil,
			nil,
		)
		if err != nil {
			return err
		}
	case reflect.Struct:
		err := loopWrapper(
			buf,
			v,
			"{",
			"}",
			reflect.Value.NumField,
			reflect.Value.Field,
			genFieldNamePrefix,
			nil,
		)
		if err != nil {
			return err
		}
	case reflect.Map:
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			var lsb, rsb bytes.Buffer
			_ = encode(&lsb, keys[i])
			_ = encode(&rsb, keys[j])
			//fmt.Printf("%t=[%d]%s|[%d]%s\n", lsb.String() < rsb.String(), i, lsb.String(), j, rsb.String())
			return lsb.String() < rsb.String()
		})
		limbs, err := mapSlice(keys, genMapper(v))
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(buf, "{%s}", strings.Join(limbs, ","))
	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("true")
		} else {
			buf.WriteString("nil")
		}
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func genMapper(v reflect.Value) func(key reflect.Value) (string, error) {
	return func(key reflect.Value) (string, error) {
		var buf bytes.Buffer

		if err := encode(&buf, key); err != nil {
			return "", err
		}

		buf.WriteRune(':')

		if err := encode(&buf, v.MapIndex(key)); err != nil {
			return "", err
		}

		return buf.String(), nil
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
	return "\"" + v.Type().Field(i).Name + "\":"
}

// String.format("(%s)",
// IntStream
// .range(0, lengthFunc(value))
// .mapToObj(i->prefixGeneratorNullable?(v,i)+extractFunc(v,i)+suffixGeneratorNullable?(v,i)})
// .collect(Collectors.joining(" ")))
func loopWrapper(
	buf *bytes.Buffer,
	v reflect.Value,
	head, tail string,
	lengthFunc func(v reflect.Value) int,
	extractFunc func(v reflect.Value, i int) reflect.Value,
	prefixGeneratorNullable, suffixGeneratorNullable func(v reflect.Value, i int) string,
) error {

	buf.WriteString(head)
	for i := 0; i < lengthFunc(v); i++ {
		if i > 0 {
			buf.WriteRune(',')
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
	buf.WriteString(tail)
	return nil
}
