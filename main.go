package main

import (
	"github.com/Dimau/simple_test_webserver/controllers"
	"github.com/julienschmidt/httprouter"
	"html/template"
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
	uc := controllers.NewUserController()
	router.GET("/user/:id", uc.GetUser)
	router.POST("/user/", uc.CreateUser)
	router.DELETE("/user/:id", uc.DeleteUser)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
