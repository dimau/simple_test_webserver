package main

import (
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"text/template"
)

var tpl *template.Template // инициализация контейнера с шаблонами HTML страниц
var UserDB map[string]User
var SessionDB map[string]string

type User struct {
	FirstName string
	LastName string
	Username string
	Password string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	UserDB = map[string]User{}
	SessionDB = map[string]string{}
}

func main() {
	http.HandleFunc("/", indexPageHandler)
	http.HandleFunc("/signup", signUpPageHandler)
	http.HandleFunc("/secretpage", onlyForRegisteredUsersPageHandler)
	http.HandleFunc("/favicon.ico", notFoundPageHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err.Error())
	}
}

func indexPageHandler(w http.ResponseWriter, req *http.Request) {
	user := getUser(w, req)

	// executing template for HTML page
	if err := tpl.ExecuteTemplate(w, "index.gohtml", user); err != nil {
		log.Fatalln(err.Error())
	}
}

func signUpPageHandler(w http.ResponseWriter, r *http.Request) {
	// Отправляем на главную страницу, если пользователь уже авторизован
	if isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Обработка данных из заполненной формы
	if r.Method == http.MethodPost {

		// Получаем поля из формы
		if err := r.ParseForm(); err != nil {
			log.Fatalln(err.Error())
		}
		user := User{
			FirstName: r.FormValue("first_name"),
			LastName:  r.FormValue("last_name"),
			Username:  r.FormValue("email"),
			Password:  r.FormValue("password"),
		}

		// Проверяем есть ли такой пользователь в базе
		if _, ok := UserDB[user.Username]; ok {
			http.Error(w, "This username is already exist", http.StatusForbidden)
			return
		}

		// Сохраняем нового пользователя в базу
		UserDB[user.Username] = user

		// Создаем сессию
		sID := uuid.New().String()
		c := &http.Cookie{
			Name: "session",
			Value: sID,
		}
		SessionDB[sID] = user.Username
		http.SetCookie(w, c)

		// Редирект на главную после успешной регистрации
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Выполнение шаблона и отправка страницы с формой
	if err := tpl.ExecuteTemplate(w, "signup.gohtml", nil); err != nil {
		log.Fatalln(err.Error())
	}

}

func onlyForRegisteredUsersPageHandler(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if _, err := io.WriteString(w, "<!DOCTYPE html><html><head><title>Secret Page</title></head><body>This is a secret page</body></html>"); err != nil {
		log.Fatalln(err.Error())
	}
}

func notFoundPageHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	_, err := io.WriteString(w, "<!DOCTYPE html><html><head><title>Not Found Page</title></head><body><p>This page is not found</p></body></html>")
	if err != nil {
		log.Fatalln(err.Error())
	}
}
