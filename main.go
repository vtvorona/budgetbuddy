package main

import (
	"budgetbuddy/auth"
	db "budgetbuddy/db"
	"budgetbuddy/jwt"
	"budgetbuddy/render"
	"budgetbuddy/server"
	"log"
	"net/http"
)

func main() {
	render.LoadTemplates()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	db.InitDB()

	http.Handle("/", jwt.AuthMiddleware(http.HandlerFunc(server.HomeHandler)))
	http.HandleFunc("/registration", server.RegistHandler)
	http.HandleFunc("/register", auth.RegisterUserHandler)
	http.HandleFunc("/login", server.LoginHandler)
	http.HandleFunc("/auth", auth.LoginUserHandler)
	http.HandleFunc("/success", server.SuccessHandler)
	http.Handle("/dashboard", jwt.AuthMiddleware(http.HandlerFunc(server.DashboardHandler)))
	http.Handle("/expense/create", jwt.AuthMiddleware(http.HandlerFunc(server.CreateExpense)))
	http.Handle("/expense/delete", jwt.AuthMiddleware(http.HandlerFunc(server.DeleteExpense)))
	http.Handle("/user", jwt.AuthMiddleware(http.HandlerFunc(server.UserSettings)))
	http.Handle("/categories/add", jwt.AuthMiddleware(http.HandlerFunc(server.AddCategory)))
	http.Handle("/categories/delete", jwt.AuthMiddleware(http.HandlerFunc(server.DeleteCategory)))
	http.Handle("/filter", jwt.AuthMiddleware(http.HandlerFunc(server.RenderFilteredExpenses)))
	http.Handle("/filter/get", jwt.AuthMiddleware(http.HandlerFunc(server.FilterExpensesHandler)))
	http.Handle("/expense/edit", jwt.AuthMiddleware(http.HandlerFunc(server.EditExpenseHandler)))

	log.Println("Сервер запущен на http://localhost:5555")
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
