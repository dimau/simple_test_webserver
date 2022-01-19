package main

import (
	"log"
	"os"
	"text/template"
)

type country struct {
	Name    string
	Capital string
}

func main() {
	tpl, err := template.ParseGlob("templates/*.gohtml")
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
