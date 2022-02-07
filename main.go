package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Соединение с БД объявляем доступным в рамках всего пакета
var db *sql.DB
var err error

func main() {
	// С помощью "sql.Open" мы инициализируем пул для хранения соединений с базой данных "db" с типом "*sql.DB"
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	db, err = sql.Open("pgx", "postgres://evoniser_user:fdjf34HFn345Uj5n@localhost:5432/evoniser_db")
	if err != nil {
		log.Fatalln(err.Error())
	}
	// It is idiomatic to defer db.Close() if the sql.DB should not have a lifetime beyond the scope of the function
	defer db.Close()

	// Проверяем, что БД нам отвечает
	// The first actual connection to the underlying datastore will be established lazily, when it’s needed for the first time.
	// If you want to check right away that the database is available and accessible (for example, check that you
	// can establish a network connection and log in), use db.Ping() to do that
	if err = db.Ping(); err != nil {
		log.Fatalln(err.Error())
	}

	http.HandleFunc("/", indexPageHandler)
	http.HandleFunc("/select", selectPageHandler)
	http.HandleFunc("/createtable", createTablePageHandler)
	http.HandleFunc("/insert", insertPageHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err.Error())
	}
}

func selectPageHandler(w http.ResponseWriter, r *http.Request) {

	// Выполнение SQL запроса к БД, результат записываем в "rows"
	rows, err := db.Query(`SELECT task_id, summary FROM tasks`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Важно выполнять Close для значений с типом *sql.Rows после завершения работы,
	// чтобы освобождать соответствующее соединение с базой данных, иначе может быть перерасход памяти
	defer rows.Close()

	// Объявляем переменные, в которые будем считывать значения столбцов из строк, хранящихся в rows
	var (
		task_id, summary string
		oneResult []string  // Будем упаковывать каждую пару task_id и summary в такой массив
		allResults [][]string // Будем упаковывать в итоге все результаты выполнения запроса в одну переменную
	)

	// Итерируем по строкам результата с помощью "rows.Next()"
	for rows.Next() {
		// We read the columns in each row into variables with rows.Scan()
		err = rows.Scan(&task_id, &summary)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Собираем результаты в один массив для дальнейшей передачи в HTML шаблон для отображения
		oneResult = []string{task_id, summary}
		allResults = append(allResults, oneResult)
	}

	// Выход из цикла мог случиться не потому что обработаны все строки, а потому что при итерировании по ним произошла ошибка
	// Поэтому всегда следует проверять наличие ошибки после выхода из цикла "for rows.Next()"
	err = rows.Err()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Вырожденный способ собрать строку и отправить ее в качестве ответа в браузер
	_, _ = io.WriteString(w, gatherAllResultsIntoString(allResults))
}

func gatherAllResultsIntoString(allResults [][]string) string {
	allResultsString := ""
	for _, i := range allResults {
		for _, j := range i {
			allResultsString += j + " "
		}
		allResultsString += "\n"
	}
	return allResultsString
}

func createTablePageHandler(w http.ResponseWriter, r *http.Request) {

	// Создаем шаблон для выполнения (особенно повторяющихся) операций в БД - statement
	stmt, err := db.Prepare(`CREATE TABLE customer (name VARCHAR(20));`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Выполняем запрос из шаблона stmt
	result, err := stmt.Exec()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Возвращает количество строк, которые были изменены после выполнения SQL запроса
	n, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Вырожденный способ собрать строку и отправить ее в качестве ответа в браузер
	_, _ = io.WriteString(w, "Table is successfully created: " + strconv.FormatInt(n, 10))
}

func insertPageHandler(w http.ResponseWriter, r *http.Request) {

	// Создаем шаблон для выполнения (особенно повторяющихся) операций в БД - statement
	stmt, err := db.Prepare("INSERT INTO customer (name, surname) VALUES ($1, $2);")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Выполняем запрос из шаблона stmt с подстановкой переменных
	res, err := stmt.Exec("Richard", "Stevenson")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Возвращает количество строк, которые были изменены после выполнения SQL запроса
	n, err := res.RowsAffected()
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Вырожденный способ собрать строку и отправить ее в качестве ответа в браузер
	_, _ = io.WriteString(w, strconv.FormatInt(n, 10))
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "It's working!")
}