package main

import (
	"budgetbuddy/auth"
	db "budgetbuddy/db"
	"budgetbuddy/server"
	"log"
	"net/http"
)

func main() {
	server.LoadTemplates() // Загружаем шаблоны при запуске
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	db.InitDB()

	// Открытые маршруты (без авторизации)
	http.HandleFunc("/", server.HomeHandler)
	http.HandleFunc("/registration", server.RegistHandler)
	http.HandleFunc("/register", auth.RegisterUserHandler)
	http.HandleFunc("/login", server.LoginHandler)
	http.HandleFunc("/auth", auth.LoginUserHandler)

	// Защищённые маршруты (с AuthMiddleware)
	http.Handle("/dashboard", auth.AuthMiddleware(http.HandlerFunc(server.DashboardHandler)))
	http.Handle("/success", auth.AuthMiddleware(http.HandlerFunc(server.SuccessHandler)))

	log.Println("Сервер запущен на http://localhost:5555")
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
