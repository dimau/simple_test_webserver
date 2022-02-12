package main

import (
	"encoding/json"
	"log"
)

// Объявление структуры, в которую будем демаршалить строку JSON
type city struct {
	Precision string  `json:"precision"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Address   string  `json:"address,omitempty"`
	City      string  `json:"city,omitempty"`
	State     string  `json:"state"`
	Zip       string  `json:"zip"`
	Country   string  `json:"country"`
}

type cities []city

func main() {
	// Исходная строка в JSON формате
	rcvd := `[{"precision":"zip","Latitude":37.7668,"Longitude":-122.3959,"Address":"",
"City":"SAN FRANCISCO","State":"CA","Zip":"94107","Country":"US"},
{"precision":"zip","Latitude":37.371991,"Longitude":-122.02602,"Address":"",
"City":"SUNNYVALE","State":"CA","Zip":"94085","Country":"US"}]`

	// Демаршалинг строки в JSON формате в структуру Go
	var data cities
	err := json.Unmarshal([]byte(rcvd), &data)
	if err != nil {
		log.Fatalln("Error unmarshalling", err.Error())
	}

	// Используем структуру, созданную из JSON строки
	log.Println(data)
}

