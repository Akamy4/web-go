package main

import (
	"fmt"
	"net/http"

	"github.com/common-nighthawk/go-figure"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		text := request.URL.Path[1:] // Убираем слеш в начале пути
		if text == "" {
			text = "Hello"
		}

		// Генерация ASCII-арта
		asciiArt := figure.NewFigure(text, "", true)
		result := asciiArt.String()

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(result))
	})

	fmt.Println("Server is running on :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
