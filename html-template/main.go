package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Todo struct {
	Name string
	Done bool
}

func IsNotDone(todo Todo) bool {
	return !todo.Done
}

func main() {
	// Создаем шаблон и передаем туда функции
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"IsNotDone": IsNotDone,
	}).ParseFiles("html-template/index.html")

	if err != nil {
		log.Fatal("Can not expand template", err)
		return
	}

	todos := []Todo{
		{"Пресс качат", false},
		{"T) Бегит", false},
		{"Турник", false},
		{"Анжуманя", false},
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			// читаем из urlencoded запроса
			param := request.FormValue("id")
			//переобразуем строку в int
			index, _ := strconv.ParseInt(param, 10, 0)
			todos[index].Done = true
		}

		err := tmpl.Execute(writer, todos)
		if err != nil {
			//Вернем 500 и напишем ошибку
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8081", nil)

}
