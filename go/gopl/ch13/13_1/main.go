package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	fmt.Println(Equal(999999, 1000000))
	fmt.Println(Equal(1000001, 1000000))
	fmt.Println(Equal(1e10, 1e10+5))
	fmt.Println(Equal(1e3, 1e3+5))
}

type record struct {
	l, r unsafe.Pointer
	t    reflect.Type
}

func withinSimilarRate(rate float64) bool {
	return rate < 1e-6
}

func equal(l, r reflect.Value, dejavu map[record]bool) bool {
	if !l.IsValid() || !r.IsValid() {
		return l.IsValid() == r.IsValid()
	}

	if l.Type() != r.Type() {
		return false
	}

	if l.CanAddr() && r.CanAddr() {
		pl := unsafe.Pointer(l.UnsafeAddr())
		pr := unsafe.Pointer(r.UnsafeAddr())
		if pl == pr {
			return true
		}
		r := record{
			l: pl,
			r: pr,
			t: l.Type(), // l&r's identical type shall be confirmed by a previous check
		}
		if _, ok := dejavu[r]; ok {
			return true
		}
		dejavu[r] = true
	}

	switch l.Kind() {
	case reflect.Int16, reflect.Int8:
		return l.Int() == r.Int()
	case reflect.Int, reflect.Int64, reflect.Int32:
		// Whoever care edge overflow situations! I'm just kidding.
		big, small := l.Int(), r.Int()
		if big == small {
			return true
		}
		if big < small {
			big, small = small, big
		}
		diff := big - small
		return withinSimilarRate(float64(diff) / float64(big))
	case reflect.Uint16, reflect.Uint8:
		// How can I merge the similar code? There are only type difference.
		big, small := l.Uint(), r.Uint()
		if big == small {
			return true
		}
		if big < small {
			big, small = small, big
		}
		diff := big - small
		return withinSimilarRate(float64(diff) / float64(big))
	case reflect.Uint, reflect.Uint64, reflect.Uint32:
		return l.Uint() == r.Uint()
	case reflect.Float32, reflect.Float64:
		big, small := l.Float(), r.Float()
		if big < small {
			big, small = small, big
		}
		return withinSimilarRate((big - small) / big)
	case reflect.Bool:
		return l.Bool() == r.Bool()
	case reflect.String:
		return l.String() == r.String()
	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return l.Pointer() == r.Pointer()
	case reflect.Ptr, reflect.Interface:
		return equal(l.Elem(), r.Elem(), dejavu)
	case reflect.Array, reflect.Slice:
		if l.Len() != r.Len() {
			return false
		}
		for i := 0; i < l.Len(); i++ {
			if !equal(l.Index(i), r.Index(i), dejavu) {
				return false
			}
		}
		return true
	default:
		panic(fmt.Errorf("unsupported kind %v", l.Kind()))
	}
}

func Equal(l, r interface{}) bool {
	dejavu := make(map[record]bool)
	return equal(reflect.ValueOf(l), reflect.ValueOf(r), dejavu)
}
