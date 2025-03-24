package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"io"
	"log"
	"net/http"
	"os"
)

type Joke struct {
	Value string `json:"value"`
}

var buttons = [][]tgbotapi.KeyboardButton{
	{tgbotapi.KeyboardButton{Text: "Get Joke"}},
}

// WebhookURL При старте приложени, оно скажет телеграму ходить с обновлениями в этот URL
const WebhookURL = "https://636c-195-7-13-201.ngrok-free.app"

func getJoke() string {
	c := http.Client{}

	resp, err := c.Get("https://api.chucknorris.io/jokes/random")

	if err != nil {
		return "Jokes not unenviable"
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Println("Raw Response:", string(body))
	joke := Joke{}

	err = json.Unmarshal(body, &joke)

	if err != nil {
		return "Joke error"
	}

	return joke.Value
}
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, skipping...")
	}

	port := os.Getenv("PORT")

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_API"))

	fmt.Println(bot, port)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Println("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))

	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/")
	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for update := range updates {
		var message tgbotapi.MessageConfig

		log.Println("Received text: ", update.Message.Text)

		switch update.Message.Text {
		case "Get Joke":
			message = tgbotapi.NewMessage(update.Message.Chat.ID, getJoke())
		default:
			message = tgbotapi.NewMessage(update.Message.Chat.ID, `Press "Get Joke" to get joke`)
		}

		message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons...)

		bot.Send(message)
	}
}
