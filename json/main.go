package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Value string `json:"value"`
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		c := http.Client{}

		resp, err := c.Get("https://api.chucknorris.io/jokes/random")

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		log.Println("Raw Response:", string(body))
		joke := Response{}

		err = json.Unmarshal(body, &joke)

		if err != nil {
			log.Fatal(err)
		}

		jsonJoke, _ := json.Marshal(joke)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(jsonJoke)
	})

	log.Println("Server is running on http://127.0.0.1:8081")
	http.ListenAndServe(":8081", nil)
}
