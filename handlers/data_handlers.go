package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"mtg-client/models"
)

type Request struct {
	ThreadCount int `json:"threadCount"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

// Обработчик для получения данных от сервера
func startHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.ThreadCount < 1 {
		log.Println("[ERR] указано неверное количество потоков")
		http.Error(w, "Неверное количество потоков", http.StatusBadRequest)
		return
	}

	for i := 0; i < req.ThreadCount; i++ {
		go receiveData(i)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Запрос данных запущен в %d поток(ах).", req.ThreadCount),
	})
}

// Обработчик для отправки данных на сервер
func sendHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.ThreadCount < 1 {
		log.Println("[ERR] указано неверное количество потоков")
		http.Error(w, "Неверное количество потоков", http.StatusBadRequest)
		return
	}

	for i := 0; i < req.ThreadCount; i++ {
		go sendRandomData(i)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Отправка данных запущена в %d поток(ах).", req.ThreadCount),
	})
}

func receiveData(threadID int) {
	// Подключение к серверу
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Printf("[ERR] поток %d: ошибка подключения к серверу: %v\n", threadID, err)
		return
	}
	defer conn.Close()

	// Отправка запроса на получение данных
	_, err = conn.Write([]byte("GET_DATA\n"))
	if err != nil {
		log.Printf("поток %d: ошибка отправки запроса: %v\n", threadID, err)
		return
	}

	models.StartFileWriter()
	models.ReadDataFromServer(conn)
}

func sendRandomData(threadID int) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Printf("поток %d: ошибка подключения к серверу: %v\n", threadID, err)
		return
	}
	defer conn.Close()

	rand.Seed(time.Now().UnixNano() + int64(threadID))

	for {
		size := rand.Intn(9001) + 1000 // Размер от 1000 до 10000 байт
		data := make([]byte, size)
		for i := range data {
			data[i] = '0' + byte(rand.Intn(10)) // Случайная цифра от '0' до '9'
		}
		encodedData := base64.StdEncoding.EncodeToString(data)

		_, err = conn.Write([]byte(encodedData + "\n"))
		if err != nil {
			log.Printf("поток %d: ошибка отправки данных: %v\n", threadID, err)
			return
		}

		// Передача данных для записи в файл
		dataLine := string(data) + "\n"
		models.PassDataToFile(dataLine)

		time.Sleep(500 * time.Millisecond)
	}
}
