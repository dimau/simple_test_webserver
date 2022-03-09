package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Dimau/simple_test_webserver/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
)

var db *mongo.Client
var err error

func main() {
	// Database connection initialization
	db, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Defer a call to Disconnect after instantiating your MongoDB client
	// To Do in main function of the application (time of life = time of life of application)
	defer func() {
		if err = db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Checking that connection with MongoDB is working properly
	if err = db.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatalln(err.Error())
	}

	router := httprouter.New()
	router.POST("/book/", CreateBook)
	router.GET("/initialize_test_data/", CreateTestDataAboutTea)
	router.GET("/add_one_test_document", AddNewDocument)
	router.GET("/find", Find)
	router.GET("/find_one", FindOne)
	router.GET("/delete_many_documents", DeleteTestDocuments)
	router.GET("/delete_one_document", DeleteTestDocument)
	router.GET("/update_one_document", UpdateOneDocument)
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
	_, err = db.Database("bookstore").Collection("books").InsertOne(context.TODO(), book)
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

func CreateTestDataAboutTea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	collection := db.Database("tea").Collection("ratings")
	docs := []interface{}{
		bson.D{{"type", "Masala"}, {"rating", 10}},
		bson.D{{"type", "Matcha"}, {"rating", 7}},
		bson.D{{"type", "Assam"}, {"rating", 4}},
		bson.D{{"type", "Oolong"}, {"rating", 9}},
		bson.D{{"type", "Chrysanthemum"}, {"rating", 5}},
		bson.D{{"type", "Earl Grey"}, {"rating", 8}},
		bson.D{{"type", "Jasmine"}, {"rating", 3}},
		bson.D{{"type", "English Breakfast"}, {"rating", 6}},
		bson.D{{"type", "White Peony"}, {"rating", 4}},
	}
	result, err := collection.InsertMany(context.TODO(), docs)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Number of documents inserted: %d\n", len(result.InsertedIDs))
}

func AddNewDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	collection := db.Database("myDB").Collection("favorite_books")

	doc := bson.D{
		{"title", "Invisible Cities"},
		{"author", "Italo Calvino"},
		{"year_published", 1974},
	}

	result, _ := collection.InsertOne(context.TODO(), doc)
	w.Write([]byte(fmt.Sprintf("Inserted document with id %v\n", result.InsertedID)))
}

func Find(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Выбираем конкретную базу данных и коллекцию
	coll := db.Database("tea").Collection("ratings")

	// Собираем условия фильтрации документов
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"rating", bson.D{{"$gt", 5}}}},
				bson.D{{"rating", bson.D{{"$lt", 10}}}},
			},
		},
	}

	// Собираем сортировку
	sort := bson.D{{"rating", -1}}

	// Выбираем какие именно поля в документах нам нужно вернуть
	projection := bson.D{
		{"type", 1}, {"rating", 1}, {"_id", 0},
	}

	// Собираем все в options к запросу
	opts := options.Find().SetSort(sort).SetProjection(projection)

	// Выполняем запрос, получаем назад структуру типа "cursor"
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		panic(err)
	}
	// Отложенный вызов закрытия "cursor", чтобы не забыть освободить ресурсы
	defer cursor.Close(context.TODO())

	// В цикле перебираем все результаты из "cursor" и обрабатываем каждый
	for cursor.Next(context.TODO()) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

}

func FindOne(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Выбираем конкретную базу данных и коллекцию
	coll := db.Database("tea").Collection("ratings")

	// Собираем пустой фильтр - все документы коллекции подходят под запрос
	filter := bson.D{}

	// Собираем сортировку - по возрастанию по полю "rating"
	sort := bson.D{{"rating", 1}}

	// Выбираем какие именно поля в документах нам нужно вернуть
	projection := bson.D{{"type", 1}, {"rating", 1}, {"_id", 0}}

	// Собираем все в options к запросу
	opts := options.FindOne().SetSort(sort).SetProjection(projection)

	// Выполняем запрос к MongoDB и складываем единственный результат в переменную "result"
	var result bson.D
	err := coll.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	// Выполняем действия с полученным результатом из MongoDB
	fmt.Println(result)
}

func DeleteTestDocuments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Выбираем конкретную базу данных и коллекцию
	coll := db.Database("tea").Collection("ratings")

	// Собираем фильтр документов, которые будем удалять
	filter := bson.D{
		{"rating", bson.D{{"$gt", 8}}},
	}

	// Выполняем удаление в БД всех документов, соответствующих условиям фильтрации
	result, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatalln(err)
	}

	// Выводим количество удаленных документов
	fmt.Printf("Number of documents deleted: %v\n", result.DeletedCount)
}

func DeleteTestDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Выбираем конкретную базу данных и коллекцию
	coll := db.Database("tea").Collection("ratings")

	// Собираем фильтр для выборки нужного документа
	objectID, _ := primitive.ObjectIDFromHex("622887304092ea3375e1c3ab")
	filter := bson.D{
		{"_id", objectID},
	}

	// Выполняем удаление в БД документа, соответствующего условиям фильтрации
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatalln(err)
	}

	// Выводим количество удаленных документов, в этом случае может быть = 0 или 1
	fmt.Printf("Number of documents deleted: %v\n", result.DeletedCount)
}

func UpdateOneDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Выбираем конкретную базу данных и коллекцию
	coll := db.Database("tea").Collection("ratings")

	// Собираем фильтр для выборки нужного документа - по ID
	objectID, _ := primitive.ObjectIDFromHex("622887304092ea3375e1c3b2")
	filter := bson.D{{"_id", objectID}}

	// Какие поля и на какие значения нужно изменить
	update := bson.D{
		{"$set", bson.D{
			{"type", "Black Peony"},
			{"rating", 10},
		}},
		//{"$inc", bson.D{
		//	{"bonus", 2000},
		//}},
	}

	// Производим изменений документа в MongoDB
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatalln(err)
	}

	// Смотрим на результаты проведения изменений
	fmt.Printf("Documents matched: %v\n", result.MatchedCount)
	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)

}
