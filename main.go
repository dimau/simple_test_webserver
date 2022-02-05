package main

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"text/template"
)

var tpl *template.Template // инициализация контейнера с шаблонами HTML страниц
var userDB map[string]user
var sessionDB map[string]string

type user struct {
	FirstName string
	LastName string
	Username string
	Password []byte
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	userDB = map[string]user{}
	sessionDB = map[string]string{}
}

func main() {
	http.HandleFunc("/", indexPageHandler)
	http.HandleFunc("/signup", signUpPageHandler)
	http.HandleFunc("/login", logInPageHandler)
	http.HandleFunc("/logout", logoutPageHandler)
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

func logInPageHandler(w http.ResponseWriter, r *http.Request) {
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
		username := r.FormValue("email")
		password := []byte(r.FormValue("password"))

		// Находим пользователя с таким именем в базе и проверяем пароль
		u, ok := userDB[username]
		if !ok {
			http.Error(w, "Username does not exist", http.StatusNotFound)
			return
		}
		if err:= bcrypt.CompareHashAndPassword(u.Password, password); err != nil {
			http.Error(w, "Password is not correct", http.StatusForbidden)
			return
		}

		// Создаем сессию
		sID := uuid.New().String()
		c := &http.Cookie{
			Name:  "session",
			Value: sID,
		}
		sessionDB[sID] = username
		http.SetCookie(w, c)

		// Редирект на главную после успешной регистрации
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Выполнение шаблона и отправка страницы с формой
	if err := tpl.ExecuteTemplate(w, "login.gohtml", nil); err != nil {
		log.Fatalln(err.Error())
	}
}

func logoutPageHandler(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Удаляем информацию о сессии из БД на сервере
	c, _ := r.Cookie("session")
	delete(sessionDB, c.Value)

	// Очищаем куку с id сессии на стороне клиента
	c = &http.Cookie{
		Name: "session",
		Value: "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// Редирект на главную страницу сайта
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
		user := user{
			FirstName: r.FormValue("first_name"),
			LastName:  r.FormValue("last_name"),
			Username:  r.FormValue("email"),
			Password:  []byte(r.FormValue("password")),
		}

		// Проверяем есть ли такой пользователь в базе
		if _, ok := userDB[user.Username]; ok {
			http.Error(w, "This username is already exist", http.StatusForbidden)
			return
		}

		// Шифруем пароль и сохраняем нового пользователя в базу
		encryptedPass, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		user.Password = encryptedPass
		userDB[user.Username] = user

		// Создаем сессию
		sID := uuid.New().String()
		c := &http.Cookie{
			Name: "session",
			Value: sID,
		}
		sessionDB[sID] = user.Username
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

func notFoundPageHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(404)
	_, err := io.WriteString(w, "<!DOCTYPE html><html><head><title>Not Found Page</title></head><body><p>This page is not found</p></body></html>")
	if err != nil {
		log.Fatalln(err.Error())
	}
}
