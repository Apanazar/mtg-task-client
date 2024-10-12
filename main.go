package main

import (
	"log"
	"net/http"

	"mtg-client/handlers"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	handlers.RegisterRoutes()

	log.Println("[INFO] клиент запущен на http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
