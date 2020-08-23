package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

type multiSort struct {
	t    []*Track
	keys []string
}

func (m multiSort) Len() int {
	return len(m.t)
}

func (m multiSort) Less(i, j int) bool {
	l := m.t[i]
	r := m.t[j]
	for _, key := range m.keys {
		fmt.Println("check " + key)
		switch key {
		case "Title":
			if l.Title != r.Title {
				return l.Title < r.Title
			}
		case "Artist":
			if l.Artist != r.Artist {
				return l.Artist < r.Artist
			}
		case "Album":
			if l.Album != r.Album {
				return l.Album < r.Album
			}
		case "Year":
			if l.Year != r.Year {
				return l.Year < r.Year
			}
		case "Length":
			if l.Length != r.Length {
				return l.Length < r.Length
			}
		default:
			panic("bad key " + key)
		}
	}
	return false
}

func (m multiSort) Swap(i, j int) {
	m.t[i], m.t[j] = m.t[j], m.t[i]
}

func main() {
	printTracks(tracks)

	http.HandleFunc("/tracks", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}

var tmpl = template.Must(template.New("tracks").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<table>
    <tr style='text-align: left'>
        <th>
            <button onclick="update('Title')">Title</button>
        </th>
        <th>
            <button onclick="update('Artist')">Artist</button>
        </th>
        <th>
            <button onclick="update('Album')">Album</button>
        </th>
        <th>
            <button onclick="update('Year')">Year</button>
        </th>
    </tr>
    {{range .Items}}
        <tr>
            <td>{{.Title}}</td>
            <td>{{.Artist}}</td>
            <td>{{.Album}}</td>
            <td>{{.Year}}</td>
        </tr>
    {{end}}
</table>

<script>
    const update = (neoKey) => {
        let keys = []
        const queryParams = new URLSearchParams(window.location.search);
        const param = queryParams.get("key");
        if (param !== null) {
            // I forgot to add ... at the beginning. Mismatch of types hurts me.
            keys.push(...param.split("0"))
        }

        const goodOlds = keys.filter(value => value !== neoKey)

        let neo = [neoKey].concat(goodOlds);
        queryParams.set("key", neo.join("0"));

        history.pushState(null, null, "?" + queryParams.toString());
        
        location.reload()
        return false
    }
</script>
</body>
</html>`))

func handler(w http.ResponseWriter, r *http.Request) {
	var keys []string
	param := r.URL.Query().Get("key")
	if len(param) > 0 {
		keys = strings.Split(param, "0")
	}

	sort.Sort(multiSort{
		t:    tracks,
		keys: keys,
	})

	err := tmpl.Execute(w, struct {
		Items []*Track
	}{tracks})
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		w.WriteHeader(500)
	}
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	_, _ = fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	_, _ = fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		_, _ = fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	_ = tw.Flush() // calculate column widths and print table
	fmt.Println()
}
