package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
)

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next() {
	lex.token = lex.scan.Scan()
}

func (lex *lexer) text() string {
	return lex.scan.TokenText()
}

func (lex *lexer) consume(want rune) {
	if lex.token != want {
		panic(fmt.Errorf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

func read(lex *lexer, v reflect.Value) {
	switch lex.token {
	case scanner.Ident:
		if lex.text() == "nil" {
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text())
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text())
		v.SetInt(int64(i))
		lex.next()
		return
	case '(':
		lex.next()
		readList(lex, v)
		lex.next()
		return
	}
	panic(fmt.Errorf("unexpected toke %q", lex.text()))
}

func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		for i := 0; listShallEnd(lex); i++ {
			read(lex, v.Index(i))
		}
	case reflect.Slice:
		for !listShallEnd(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}
	case reflect.Struct:
		for !listShallEnd(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Errorf("got token %q, want a field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}
	case reflect.Map:
		v.Set(reflect.MakeMap(v.Type()))
		for !listShallEnd(lex) {
			lex.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}
	default:
		panic(fmt.Errorf("cannot decode list into %v", v.Type()))
	}
}

func listShallEnd(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

func Unmarshal(data []byte, out interface{}) error {
	decoder := NewDecoder(bytes.NewReader(data))
	return decoder.decode(out)
}

type Decoder struct {
	r io.Reader
}

func (d *Decoder) decode(v interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(d.r)
	lex.next()
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(v).Elem())
	return nil
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func main() {
	var num int
	err := Unmarshal([]byte("12"), &num)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(num)

	var s string
	decoder := NewDecoder(strings.NewReader("\"hello\""))
	err = decoder.decode(&s)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(s)
}
