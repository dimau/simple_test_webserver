package main

import (
	"io"
	"log"
	"net/http"
	"text/template"
)

// Инициализация контейнера с шаблонами HTML страниц
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	http.HandleFunc("/", indexPageHandler)
	http.HandleFunc("/about", aboutPageHandler)
	http.HandleFunc("/main", mainPageHandler)
	http.HandleFunc("/form_handler", mainPageHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func indexPageHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "<!DOCTYPE html><html><head><title>Index page</title></head><body><p>It's index page</p></body></html>")
}

func aboutPageHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "<!DOCTYPE html><html><head><title>About page</title></head><body><p>It's about page</p></body></html>")
}

func mainPageHandler(w http.ResponseWriter, req *http.Request) {

	// parsing and working with http request
	err := req.ParseForm() // before get form data from http.Request you need to call this method
	if err != nil {
		log.Fatalln(err.Error())
	}

	// executing template for HTML page
	err = tpl.ExecuteTemplate(w, "tpl.gohtml", req.Form)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
