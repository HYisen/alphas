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
	type Movie struct {
		Fuck            bool
		Title, Subtitle string
		Year            int
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}
	strangelove := Movie{
		Fuck:     false,
		Title:    "Dr. Strangelove",
		Subtitle: "",
		Year:     0,
		Actor:    nil,
		Oscars:   nil,
	}
	fmt.Println(Encode(strangelove))

	fmt.Print("\n" + strings.Repeat("-", 60) + "\n\n")

	strangelove = Movie{
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
	fmt.Println(Encode(strangelove))
}

func Encode(v interface{}) string {
	var buf bytes.Buffer
	if err := encode(0, false, &buf, reflect.ValueOf(v)); err != nil {
		log.Fatalln(err)
	}
	return buf.String()
}

func encode(indent int, skipFirstIndent bool, buf *bytes.Buffer, v reflect.Value) error {
	if !skipFirstIndent {
		buf.WriteString(strings.Repeat(" ", indent))
	}
	switch v.Kind() {
	case reflect.Invalid:
		break
	case reflect.Int:
		if v.Int() == 0 {
			break
		}
		_, _ = fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint:
		if v.Uint() == 0 {
			break
		}
		buf.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.String:
		if v.Len() == 0 {
			break
		}
		_, _ = fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(0, false, buf, v.Elem())
	case reflect.Array, reflect.Slice:
		if v.Len() == 0 {
			break
		}
		err := loopWrapper(indent, buf, v, reflect.Value.Len, reflect.Value.Index, nil, nil)
		if err != nil {
			return err
		}
	case reflect.Struct:
		err := loopWrapper(0, buf, v, reflect.Value.NumField, reflect.Value.Field, genFieldNamePrefix, genProvider(")"))
		if err != nil {
			return err
		}
	case reflect.Map:
		if len(v.MapKeys()) == 0 {
			break
		}
		limbs, err := mapSlice(v.MapKeys(), genMapper(v))
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(buf, "(%s)", strings.Join(limbs, "\n"+strings.Repeat(" ", indent+1)))
	case reflect.Bool:
		if v.Bool() {
			buf.WriteRune('t')
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

		if err := encode(0, false, &buf, key); err != nil {
			return "", err
		}

		buf.WriteRune(' ')

		if err := encode(0, false, &buf, v.MapIndex(key)); err != nil {
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
func loopWrapper(indent int, buf *bytes.Buffer, v reflect.Value,
	lengthFunc func(v reflect.Value) int, extractFunc func(v reflect.Value, i int) reflect.Value,
	prefixGeneratorNullable, suffixGeneratorNullable func(v reflect.Value, i int) string) error {

	buf.WriteRune('(')
	for i := 0; i < lengthFunc(v); i++ {
		var b bytes.Buffer

		if i > 0 {
			b.WriteString("\n" + strings.Repeat(" ", indent+1))
		}
		prefixLength := 0
		if prefixGeneratorNullable != nil {
			prefix := prefixGeneratorNullable(v, i)
			prefixLength = len(prefix)
			b.WriteString(prefix)
		}

		var nb bytes.Buffer
		if err := encode(prefixLength+1, true, &nb, extractFunc(v, i)); err != nil {
			return err
		}
		if nb.Len() == 0 {
			continue
		}
		b.WriteString(nb.String())

		if suffixGeneratorNullable != nil {
			b.WriteString(suffixGeneratorNullable(v, i))
		}

		buf.WriteString(b.String())
	}
	buf.WriteRune(')')
	return nil
}
