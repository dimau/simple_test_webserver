package main

import (
	"context"
	"encoding/json"
	"github.com/Dimau/simple_test_webserver/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
)

var ctx = context.TODO()
var mongoClient *mongo.Client
var err error

func init() {
	// Инициализируем клиента для работы с MongoDB
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Проверяем, что соединение с MongoDB успешно установлено
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	router := httprouter.New()
	router.POST("/book/", CreateBook)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func CreateBook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Объявляем переменную того типа, в который хотим преобразовать (декодировать) входящий JSON
	var book models.Book

	// Создаем декодер как обертку над stream http.Request.Body и декодируем JSON из запроса целиком
	// Полученную строку пытаемся преобразовать к типу указанной переменной и записать в нее (&book)
	err := json.NewDecoder(r.Body).Decode(&book)

	// Обработка ошибок, если декодирование JSON из запроса не удалось завершить корректно
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Присваивание уникального ID новой книге
	book.ID = primitive.NewObjectID()

	// Сохраняем полученную из запроса на создание книгу в MongoDB
	//_, err = mongoClient.Database("bookstore").Collection("books").InsertOne(ctx, bson.D{{"isbn", "4"}, {"title", "Jdfdfdfk"}})
	_, err = mongoClient.Database("bookstore").Collection("books").InsertOne(ctx, book)
	if err != nil {
		log.Fatalln("InsertOne problem", err.Error())
	}

	// Preparing and sending response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("{}"))
	if err != nil {
		log.Fatalln(err.Error())
	}
}