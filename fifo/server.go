package fifo

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

var tpl = `<!DOCTYPE HTML>
<html>
<head>
    <meta charset="utf-8">
    <title>{{.Title}}</title>
</head>
<body>
<h1>{{.Title}}</h1>
<p>{{.Content}}</p>
<p>{{.Date | dateFormat "2006-01-02 15:04:05"}}</p>
</body>
</html>`

type Page struct {
	Title, Content string
	Date           time.Time
}

func dateFormat(l string, d time.Time) string {
	return d.Format(l)
}

func parseTemplateString(name, tpl string) *template.Template {
	t := template.New(name)
	t.Funcs(funcMap)
	t = template.Must(t.Parse(tpl))
	//t = template.Must(t.ParseFiles("simple.html"))
	return t
}

var funcMap = template.FuncMap{
	"dateFormat": dateFormat,
}

func server() {
	http.HandleFunc("/", handlerDisplay)
	http.HandleFunc("/hi", handlerHello)
	http.ListenAndServe(":8080", nil)
}

func handlerHello(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Hello world!")
}

func handlerDisplay(w http.ResponseWriter, req *http.Request) {
	t := parseTemplateString("date", tpl)
	data := &Page{
		Title:   "An example html page",
		Content: "Have fun!!",
		Date:    time.Now(),
	}
	t.Execute(w, data)
}
