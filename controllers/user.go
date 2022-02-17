package controllers

import (
	"context"
	"encoding/json"
	"github.com/Dimau/simple_test_webserver/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"net/http"
)

type UserController struct{
	mongoClient *mongo.Client
	ctx context.Context
}

func NewUserController(mc *mongo.Client, ctx context.Context) *UserController {
	return &UserController{
		mongoClient: mc,
		ctx: ctx,
	}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.User{
		Name:   "James Bond",
		Gender: "Male",
		Age:    32,
		//Id:     p.ByName("id"),
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

	// Create bson ID
	u.Id = primitive.NewObjectID()

	// Store the user to MongoDB
	collection := uc.mongoClient.Database("simple_test_webserver").Collection("users")
	_, err = collection.InsertOne(uc.ctx, u)
	if err != nil {
		log.Fatalln("InsertOne problem", err.Error())
	}

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
