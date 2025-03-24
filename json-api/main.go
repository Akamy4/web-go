package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Todo struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
}

func main() {
	todos := []Todo{
		{"Атжуманя", false},
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// здесь нужно вернуть статический файл, который будет общаться с API из браузера
		//открываем файл
		fileContents, err := ioutil.ReadFile("json-api/index.html")

		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		// и выводим содержимое файла
		writer.Write(fileContents)
	})

	http.HandleFunc("/todos/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("request ", request.URL.Path)
		defer request.Body.Close()

		// разные методы обрабатывающие по-разному
		switch request.Method {
		case http.MethodGet:
			// преобразуем структуру в json
			productJson, _ := json.Marshal(todos)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			writer.Write(productJson)

		case http.MethodPost:
			decoder := json.NewDecoder(request.Body)
			todo := Todo{}
			// преобразуем json запрос в структуру
			err := decoder.Decode(&todo)

			if err != nil {
				log.Println(err)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			todos = append(todos, todo)
		case http.MethodPut:
			id := request.URL.Path[len("/todos/"):]
			index, _ := strconv.ParseInt(id, 10, 0)
			todos[index].Done = true

		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8081", nil)
}
