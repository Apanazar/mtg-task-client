package handlers

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// Регистрируем маршруты приложения
func RegisterRoutes() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/send", sendHandler)
}
