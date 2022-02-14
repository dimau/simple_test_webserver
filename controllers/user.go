package controllers

import (
	"encoding/json"
	"github.com/Dimau/simple_test_webserver/models"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

// В итоге мы разместим здесь connection к БД
type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Объявляем переменную того типа, в который хотим преобразовать (декодировать) JSON
	var u models.User

	// Создаем декодер как обертку над stream http.Request.Body и декодируем его целиком
	// Полученную строку пытаемся преобразовать к типу указанной переменной и записать в нее (&u)
	err := json.NewDecoder(r.Body).Decode(&u)

	// Обработка ошибок, если декодирование не удалось завершить корректно
	if err != nil {
		log.Fatalln(err.Error())
	}

	u.Id = "007"

	// Prepare json response
	uj, err := json.Marshal(u)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Preparing and sending response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(uj)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//TODO: write code to delete user
	w.WriteHeader(http.StatusOK)
}
