package main

import (
	"alphas/go/gopl/ch7/7_16/eval"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl = template.Must(template.New("calculator").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Calculator</title>
</head>
<body>
<p id="screen">{{.Screen}}</p>
<p id="result">={{.Result}}</p>
<button onclick="addChar(0)">0</button>
<button onclick="addChar(1)">1</button>
<button onclick="addChar(2)">2</button>
<button onclick="addChar(3)">3</button>
<button onclick="addChar(4)">4</button>
<button onclick="addChar(5)">5</button>
<button onclick="addChar(6)">6</button>
<button onclick="addChar(7)">7</button>
<button onclick="addChar(8)">8</button>
<button onclick="addChar(9)">9</button>
<button onclick="addChar('.')">.</button>
<button onclick="addChar('+')">+</button>
<button onclick="addChar('-')">-</button>
<button onclick="addChar('*')">*</button>
<button onclick="addChar('/')">/</button>
<button onclick="clr()">C</button>
<button onclick="calc()">=</button>
<script>
    const addChar = neo => {
        setScreen(getScreen() + neo.toString())
    }

    const getScreen = () => {
        return document.getElementById("screen").innerText;
    }

    const setScreen = neo => {
        document.getElementById("screen").innerText = neo;
    }

    const clr = () => {
        setScreen('')
    }

    const calc = () => {
        const queryParams = new URLSearchParams(window.location.search);

        queryParams.set("expr", getScreen());

        history.pushState(null, null, "?" + queryParams.toString());

        location.reload()
        return false
    };
</script>
</body>
</html>`))

func main() {
	fmt.Println(os.Getwd())
	http.HandleFunc("/calculator", func(writer http.ResponseWriter, request *http.Request) {
		var result string
		screen := request.URL.Query().Get("expr")
		fmt.Println(screen)
		if len(screen) != 0 {
			expr, err := eval.Parse(screen)
			if err != nil {
				log.Println(err)
			}
			result = fmt.Sprintf("%.3f", expr.Eval(eval.Env{}))
		}
		err := tmpl.Execute(writer, struct {
			Result string
			Screen string
		}{
			Result: result,
			Screen: screen,
		})
		if err != nil {
			log.Println(err)
		}
		writer.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
