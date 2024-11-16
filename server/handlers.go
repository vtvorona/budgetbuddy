package server

import (
	"budgetbuddy/auth"
	"budgetbuddy/db"
	"budgetbuddy/models"
	"net/http"
)

// Структура для хранения данных шаблонов
type PageData struct {
	Menu string
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	RenderTemplate(writer, "home")
}

func RegistHandler(writer http.ResponseWriter, request *http.Request) {
	RenderTemplate(writer, "registration")
}

func SuccessHandler(writer http.ResponseWriter, request *http.Request) {
	RenderTemplate(writer, "success")
}

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	RenderTemplate(writer, "login")
}

func DashboardHandler(writer http.ResponseWriter, request *http.Request) {
	RenderTemplate(writer, "dashboard")
}

func AddExpense(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Не удалось определить ID пользователя", http.StatusInternalServerError)
		return
	}

	// Получаем данные формы
	title := r.FormValue("title")
	amount := r.FormValue("amount")
	category := r.FormValue("category")

	// Создаём структуру Expense
	expense := models.Expense{
		UserID:   userID,
		Title:    title,
		Amount:   amount,
		Category: category,
	}

	// Сохраняем в базу данных
	if err := db.DB.Create(&expense).Error; err != nil {
		http.Error(w, "Ошибка при добавлении расхода", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
