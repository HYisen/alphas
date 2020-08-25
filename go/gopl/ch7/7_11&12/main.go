package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, _ *http.Request) {
	for item, price := range db {
		_, _ = fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

var tmpl = template.Must(template.New("list").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<table>
    <tr style='text-align: left'>
        <th>
            Item
        </th>
        <th>
            Price
        </th>
    </tr>
    {{range $item,$price:=.}}
        <tr>
            <td>{{$item}}</td>
            <td>{{$price}}</td>
        </tr>
    {{end}}
</table>

</body>
</html>`))

func (db database) web(w http.ResponseWriter, _ *http.Request) {
	_ = tmpl.Execute(w, db)
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	_, _ = fmt.Fprintf(w, "%s\n", price)
}

func (db database) put(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	num, _ := strconv.ParseFloat(req.URL.Query().Get("price"), 32)
	price := dollars(num)
	_, ok := db[item]
	db[item] = price
	if ok {
		// old one exists.
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func (db database) remove(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	delete(db, item)
	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	db := database{"shows": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/list.html", db.web)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.put)
	http.HandleFunc("/delete", db.remove)
	http.HandleFunc("/update", db.put)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
