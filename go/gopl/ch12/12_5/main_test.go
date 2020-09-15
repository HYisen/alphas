package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestEquality(t *testing.T) {
	tests := []struct {
		name string
		item interface{}
	}{
		{name: "num", item: 8964},
		{name: "string", item: "hello"},
		{name: "map", item: map[string]int{"a": 1, "b": 2}},
		{name: "array", item: [...]int{1, 2, 4, 8}},
		{name: "slice", item: []int{3, 2, 1}},
		{name: "struct", item: struct {
			Cost int
			Name string
		}{
			Cost: 0,
			Name: "free",
		}}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				i := recover()
				if i != nil {
					t.Fatalf("face unexpect error %v", i)
				}
			}()

			// through marshaller
			old := test.item
			ref := canonicalEncode(old)
			tst := Encode(old)
			if ref != tst {
				t.Fatalf("expect:|%s|\nactual:|%s|\n", ref, tst)
			}

			// through unmarshaler
			switch reflect.ValueOf(old).Kind() {
			case reflect.Map, reflect.Struct: // there is little difference between struct and map in json
			default:
				oldStr := fmt.Sprintf("%v", old)
				neo := old
				err := json.Unmarshal([]byte(tst), &neo)
				if err != nil {
					panic(err)
				}
				neoStr := fmt.Sprintf("%v", neo)
				if neoStr != oldStr {
					t.Fatalf("expect:|%v|\nactual:|%v|\n", oldStr, neoStr)
				}
			}
		})
	}
}

func canonicalEncode(i interface{}) string {
	var sb strings.Builder
	encoder := json.NewEncoder(&sb)
	err := encoder.Encode(i)
	if err != nil {
		panic(err)
	}
	return sb.String()
}
