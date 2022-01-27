package main

import (
	"log"
	"net/http"
	"text/template"
)

// Инициализация контейнера с шаблонами HTML страниц
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

// Запуск HTTP сервера с одним обработчиком на все виды запросов
func main() {
	var handler httpHandler
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// Сам обработчик HTTP запросов
type httpHandler int

func (h httpHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// parsing and working with http request
	err := req.ParseForm() // before get form data from http.Request you need to call this method
	if err != nil {
		log.Fatalln(err.Error())
	}

	// executing template for HTML page
	err = tpl.ExecuteTemplate(rw, "tpl.gohtml", req.Form)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
