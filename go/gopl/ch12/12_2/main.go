package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type node struct {
	next *node
	name string
}

func main() {
	a := node{
		name: "one",
	}
	b := node{
		name: "two",
	}
	a.next = &b
	b.next = &a

	Display("rec", a)
}

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x), 0)
}

func display(path string, v reflect.Value, cnt int) {
	cnt++
	if cnt > 4 {
		_, _ = fmt.Fprintf(os.Stderr, "too much recurse, stop at %s", path)
		return
	}
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i), cnt)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			display(fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name), v.Field(i), cnt)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s = %s]", path, formatAtom(key), key), v.MapIndex(key), cnt)
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem(), cnt)
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem(), cnt)
		}
	default:
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'E', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'E', -1, 64)
	case reflect.Complex64:
		return strconv.FormatComplex(v.Complex(), 'E', -1, 64)
	case reflect.Complex128:
		return strconv.FormatComplex(v.Complex(), 'E', -1, 128)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Array, reflect.Struct, reflect.Interface:
		return v.Type().String() + " value"
	default:
		panic("unknown type " + v.Type().String())
	}
}
