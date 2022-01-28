package main

import (
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

// Инициализация контейнера с шаблонами HTML страниц
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(notFoundPageHandler)
	router.GET("/", indexPageHandler)
	router.GET("/about", aboutPageHandler)
	router.GET("/main", mainPageHandler)
	router.POST("/form_handler", mainPageHandler)
	router.GET("/form_handler", mainPageHandler)
	router.GET("/user/:name/:surname", userPageHandler)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func indexPageHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	io.WriteString(w, "<!DOCTYPE html><html><head><title>Index page</title></head><body><p>It's index page</p></body></html>")
}

func aboutPageHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	io.WriteString(w, "<!DOCTYPE html><html><head><title>About page</title></head><body><p>It's about page</p></body></html>")
}

func userPageHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	userName := params.ByName("name")
	userSurname := params.ByName("surname")
	io.WriteString(w, "<!DOCTYPE html><html><head><title>User page</title></head><body><p>User name: "+userName+" Surname: "+userSurname+"</p></body></html>")
}

func mainPageHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {

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

func notFoundPageHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	io.WriteString(w, "<!DOCTYPE html><html><head><title>Not Found Page</title></head><body><p>This page is not found</p></body></html>")
}
