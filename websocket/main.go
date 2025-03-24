package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

// Определяем структуру для JSON-ответа из API шуток
type Joke struct {
	Value string `json:"value"` // В поле "value" будет храниться текст шутки
}

// Настройка WebSocket соединения
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // Размер буфера для чтения данных из WebSocket
	WriteBufferSize: 1024, // Размер буфера для записи данных в WebSocket
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем соединения с любого источника (для тестов нормально, но на проде лучше ограничить)
	},
}

// Bus — это шина событий, в которой происходит регистрация клиентов и рассылка сообщений
type Bus struct {
	register  chan *websocket.Conn     // Канал для регистрации новых клиентов
	broadcast chan []byte              // Канал для отправки сообщений всем клиентам
	clients   map[*websocket.Conn]bool // Мапа для хранения подключенных клиентов
}

// Метод запуска шины событий (постоянный цикл обработки событий)
func (b *Bus) Run() {
	for {
		select {
		// Отправляем сообщение всем подключенным клиентам
		case message := <-b.broadcast:
			for client := range b.clients {
				// Пытаемся открыть соединение для записи сообщения
				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					// Если клиент отключился — удаляем его из списка
					delete(b.clients, client)
					continue
				}
				// Отправляем сообщение клиенту
				w.Write(message)
			}
		// Регистрируем нового клиента
		case client := <-b.register:
			log.Println("User registered") // Логируем подключение
			b.clients[client] = true
		}
	}
}

// Конструктор для создания нового экземпляра Bus
func NewBus() *Bus {
	return &Bus{
		register:  make(chan *websocket.Conn),     // Канал для регистрации новых клиентов
		broadcast: make(chan []byte),              // Канал для отправки сообщений
		clients:   make(map[*websocket.Conn]bool), // Мапа для хранения клиентов
	}
}

// runJoker отправляет новую шутку всем клиентам каждые 5 секунд
func runJoker(b *Bus) {
	for {
		// Ждём 5 секунд
		<-time.After(5 * time.Second)
		log.Println("Joke runJoker") // Логируем событие

		// Получаем шутку из API и отправляем её всем клиентам
		b.broadcast <- []byte(getJoke())
	}
}

// getJoke делает HTTP-запрос к API и возвращает текст шутки
func getJoke() string {
	c := http.Client{}

	// Делаем GET-запрос к API шуток
	resp, err := c.Get("https://api.chucknorris.io/jokes/random")
	if err != nil {
		return "Jokes not unenviable" // Если ошибка — возвращаем сообщение об ошибке
	}

	defer resp.Body.Close() // Закрываем тело ответа после завершения

	// Читаем тело ответа
	body, _ := io.ReadAll(resp.Body)
	log.Println("Raw Response:", string(body))

	// Распарсиваем JSON в структуру Joke
	joke := Joke{}
	err = json.Unmarshal(body, &joke)
	if err != nil {
		return "Joke error" // Если ошибка при парсинге — возвращаем сообщение об ошибке
	}

	return joke.Value // Возвращаем текст шутки
}

func main() {
	// Создаём новую шину событий
	bus := NewBus()

	// Запускаем обработку событий (горутиной)
	go bus.Run()

	// Запускаем отправку шуток каждые 5 секунд (горутиной)
	go runJoker(bus)

	// Регистрируем обработчик WebSocket-соединений
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// Пробуем поднять WebSocket-соединение
		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Fatal(err) // Если ошибка — логируем и выходим
		}

		// Регистрируем клиента в Bus
		bus.register <- ws
	})

	// Запускаем HTTP-сервер на порту 8081
	http.ListenAndServe(":8081", nil)
}
