package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

func firstThree(s string) string {
	s = strings.TrimSpace(s)
	s = s[:3]
	return s
}

var fm = template.FuncMap{
	"uc": strings.ToUpper,
	"ft": firstThree,
}

type country struct {
	Name    string
	Capital string
}

type dataForTemplate struct {
	Countries    []country
	HeaderString string
}

type httpHandler int

func (h httpHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// prepairing template container
	tpl, err := template.New("").Funcs(fm).ParseGlob("templates/*.gohtml")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// prepairing data for the HTML template
	russia := country{Name: "Russia", Capital: "Moscow"}
	germany := country{Name: "Germany", Capital: "Berlin"}
	countries := []country{russia, germany}
	data := dataForTemplate{
		Countries:    countries,
		HeaderString: "This is a header string",
	}

	err = tpl.ExecuteTemplate(rw, "tpl.gohtml", data)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {

	// running HTTP server
	var handler httpHandler
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalln(err.Error())
	}

}
