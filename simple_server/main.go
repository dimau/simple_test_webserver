package main

import (
	"io"
	"log"
	"net/http"
)

// Объявляем первую функцию для обработки HTTP запросов
func indexPageHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "<!DOCTYPE html><html><head><title>Index page</title></head><body><p>It's index page</p></body></html>")
}

// Объявляем вторую функцию для обработки HTTP запросов
func aboutPageHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "<!DOCTYPE html><html><head><title>About page</title></head><body><p>It's about page</p></body></html>")
}

func main() {

	// Регистрируем функции для обработки HTTP запросов с соответствующими им шаблонами URL
	// Так как вызываем HandleFunc как функцию пакета http (а не метод конкретного ServeMux),
	//   то указанная функция-обработчик будет привязана к DefaultServeMux
	http.HandleFunc("/", indexPageHandler)
	http.HandleFunc("/about", aboutPageHandler)

	// Запуск HTTP сервера на порту 8080 localhost
	// Так как в качестве handler (второй аргумент) указан nil,
	//   то для роутинга запросов будет использоваться DefaultServeMux
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
