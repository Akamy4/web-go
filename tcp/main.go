package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	//Bind на порт
	listener, err := net.Listen("tcp", ":5001")
	if err != nil {
		fmt.Println("Error binding port:", err)
		return
	}
	defer listener.Close()

	for {
		// Ждем пока не придет клиент
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Can not connect:", err)
			continue
		}

		fmt.Println("Connected")

		// Reader для чтения информации из сокета
		bufReader := bufio.NewReader(conn)
		fmt.Println("Start reading")

		go func(conn net.Conn) {
			defer conn.Close()
			for {
				// Считывание побайтово
				rByte, err := bufReader.ReadByte()
				if err != nil {
					fmt.Println("Can not read:", err)
					break
				}

				fmt.Print(string(rByte))
			}
		}(conn)
	}
}
