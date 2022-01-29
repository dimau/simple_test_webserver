package main

import (
	"io"
	"log"
	"net/http"
)

func indexPageHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "<!DOCTYPE html><html><head><title>Index page</title></head><body><img src='/resources/Niceguys.jpg'></body></html>")
}

func main() {

	http.HandleFunc("/", indexPageHandler)
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("../assets"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
