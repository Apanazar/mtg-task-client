package models

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var (
	dataChannel = make(chan string)
	once        sync.Once
	fileMutex   sync.Mutex
)

type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

// Запуск горутины для записи в файл
func StartFileWriter() {
	once.Do(func() {
		go writeDataToFile()
	})
}

// Читаем данные от сервера и передаём их для записи в файл
func ReadDataFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		encodedData, err := reader.ReadString('\n')
		if err != nil {
			log.Println("[ERR] ошибка чтения от сервера:", err)
			break
		}
		encodedData = encodedData[:len(encodedData)-1] // Удаляем символ новой строки

		// Декодируем данные из Base64
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			log.Println("[ERR] ошибка декодирования Base64:", err)
			continue
		}

		// Декодируем JSON
		var item Item
		err = json.Unmarshal(decodedData, &item)
		if err != nil {
			log.Println("[ERR] ошибка декодирования JSON:", err)
			continue
		}

		dataLine := fmt.Sprintf("%+v\n", item)
		dataChannel <- dataLine
	}
}

// PassDataToFile отправляет данные в горутину записи в файл
func PassDataToFile(data string) {
	dataChannel <- data
}

func writeDataToFile() {
	// Открытие или создание файла для добавления данных
	file, err := os.OpenFile("client_data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("[ERR] ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	for dataLine := range dataChannel {
		fileMutex.Lock()
		_, err := file.WriteString(dataLine)
		fileMutex.Unlock()
		if err != nil {
			log.Println("[ERR] ошибка записи в файл:", err)
			return
		}
	}
}
