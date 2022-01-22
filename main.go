package main

import (
	"log"
	"os"
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

func main() {
	tpl, err := template.New("").Funcs(fm).ParseGlob("templates/*.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	russia := country{Name: "Russia", Capital: "Moscow"}
	germany := country{Name: "Germany", Capital: "Berlin"}
	countries := []country{russia, germany}
	err = tpl.ExecuteTemplate(os.Stdout, "tpl.gohtml", countries)

	if err != nil {
		log.Fatalln(err)
	}
}
