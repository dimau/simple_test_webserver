package main

import (
	"encoding/json"
	"github.com/Dimau/simple_test_webserver/models"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"log"
	"net/http"
)

var db *template.Template
var err error

func init()  {
	db, err = template.ParseGlob("templates/*.gohtml")
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/user/:id", gerUser)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := db.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func gerUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.User{
		Name:   "James Bond",
		Gender: "Male",
		Age:    32,
		Id:     p.ByName("id"),
	}

	// Marshal into JSON
	uj, err := json.Marshal(u)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Write content type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = io.WriteString(w, string(uj))
	if err != nil {
		log.Fatalln(err.Error())
	}
}