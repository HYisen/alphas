package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type node struct {
	id   int
	next *node
}

func main() {
	a := &node{
		id: 0,
	}
	b := &node{
		id: 1,
	}
	c := &node{
		id: 2,
	}

	fmt.Println(reflect.ValueOf(a).Elem().CanAddr())
	fmt.Println(reflect.ValueOf(b).CanAddr())
	fmt.Println(reflect.ValueOf(b).Type())
	fmt.Println(Cyclic(a))
	a.next = b
	b.next = a
	fmt.Println(Cyclic(a))
	b.next = c
	fmt.Println(Cyclic(a))
	c.next = a
	fmt.Println(Cyclic(a))
}

func Cyclic(v interface{}) bool {
	dejavu := make(map[unsafe.Pointer]bool)
	return cyclic(reflect.ValueOf(v), dejavu)
}

func cyclic(v reflect.Value, dejavu map[unsafe.Pointer]bool) bool {
	if !v.IsValid() {
		return false
	}

	if v.CanAddr() {
		ptr := unsafe.Pointer(v.UnsafeAddr())
		if _, ok := dejavu[ptr]; ok {
			return true
		}
		dejavu[ptr] = true
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool:
	case reflect.String:
	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
	case reflect.Ptr, reflect.Interface:
		return cyclic(v.Elem(), dejavu)
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if cyclic(key, dejavu) || cyclic(v.MapIndex(key), dejavu) {
				return true
			}
		}
		return false
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if cyclic(v.Index(i), dejavu) {
				return true
			}
		}
		return false
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if cyclic(v.Field(i), dejavu) {
				return true
			}
		}
		return false
	default:
		panic(fmt.Errorf("unsupported kind %v", v.Kind()))
	}
	return false
}
