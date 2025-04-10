package server

import (
	"budgetbuddy/db"
	"budgetbuddy/models"
	"budgetbuddy/render"
	"budgetbuddy/runes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	data := models.Data{
		IsAuth: false,
	}
	render.RenderTemplate(writer, "home", data)
}

func RegistHandler(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "registration")
}

func SuccessHandler(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "success")
}

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "login")
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var user models.User
	if err := db.DB.Preload("Categories").First(&user, userId).Error; err != nil {
		fmt.Println("Ошибка получения пользователя:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	expenses := GetExpenses(userId)

	data := models.Data{
		IsAuth:   true,
		User:     user,
		Expenses: *expenses,
	}

	render.RenderTemplate(w, "dashboard", data)
}

func CreateExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusSeeOther)
			return
		}

		userId, _ := r.Context().Value("userId").(int)

		title := r.FormValue("title")
		title = runes.CapitalizeFirstLetter(title)

		categoryId := r.FormValue("category")

		// Загрузка категории по ID
		var category models.Category
		if err := db.DB.First(&category, categoryId).Error; err != nil {
			http.Error(w, "Category not found", http.StatusBadRequest)
			return
		}

		price := r.FormValue("price")
		priceFloat64, _ := strconv.ParseFloat(price, 64)
		amount := r.FormValue("amount")
		amountInt, _ := strconv.Atoi(amount)
		total := priceFloat64 * float64(amountInt)

		expense := &models.Expense{
			Title:      title,
			CategoryID: category.ID,
			Price:      priceFloat64,
			Amount:     amountInt,
			UserID:     userId,
			Total:      total,
		}

		fmt.Println(expense.Total)

		result := db.DB.Create(expense)
		if result.Error != nil {
			http.Error(w, "An error occurred while creating the expense", http.StatusConflict)
			return
		}

		fmt.Println("Expense created successfully")
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получение ID расхода из формы
	expenseID := r.FormValue("id")
	if expenseID == "" {
		http.Error(w, "Expense ID is required", http.StatusBadRequest)
		return
	}

	// Конвертация ID из строки в число
	id, err := strconv.Atoi(expenseID)
	if err != nil {
		http.Error(w, "Invalid Expense ID", http.StatusBadRequest)
		return
	}

	// Получение userId из контекста
	userIdVal := r.Context().Value("userId")
	userId, ok := userIdVal.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Удаление записи
	result := db.DB.Where("user_id = ? AND id = ?", userId, id).Delete(&models.Expense{})
	if result.Error != nil {
		http.Error(w, "Failed to delete expense", http.StatusInternalServerError)
		return
	}

	// Получение заголовка Referer
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/filter" // Перенаправить на фильтр по умолчанию, если Referer отсутствует
	}

	// Перенаправление обратно на страницу, с которой был сделан запрос
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func UserSettings(w http.ResponseWriter, r *http.Request) {
	// Получение userId из контекста
	userIdVal := r.Context().Value("userId")
	userId, ok := userIdVal.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := db.DB.Preload("Categories").First(&user, userId).Error; err != nil {
		fmt.Println("Ошибка получения пользователя:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Подготовка данных для рендеринга
	data := models.Data{
		User:   user,
		IsAuth: true,
	}

	// Отображение шаблона
	render.RenderTemplate(w, "settings", data)
}

func AddCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	categoryName := r.FormValue("category")
	if categoryName == "" {
		http.Error(w, "Категория не может быть пустой", http.StatusBadRequest)
		return
	}

	categoryName = runes.CapitalizeFirstLetter(categoryName)

	userIdVal := r.Context().Value("userId")
	userId, ok := userIdVal.(int)
	if !ok {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	category := models.Category{
		Name:   categoryName,
		UserID: userId,
	}

	if err := db.DB.Create(&category).Error; err != nil {
		fmt.Println("Ошибка сохранения категории:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	fmt.Println("Категория сохранена:", category.Name)
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	categoryId := r.FormValue("categoryId")
	if categoryId == "" {
		http.Error(w, "Идентификатор категории отсутствует", http.StatusBadRequest)
		return
	}

	userIdVal := r.Context().Value("userId")
	userId, ok := userIdVal.(int)
	if !ok {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	// Удаление категории
	if err := db.DB.Where("id = ? AND user_id = ?", categoryId, userId).Delete(&models.Category{}).Error; err != nil {
		fmt.Println("Ошибка удаления категории:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	fmt.Println("Категория удалена, ID:", categoryId)
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func FilterExpensesHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметров даты
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// Проверка на валидность параметров
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Неверный формат начальной даты", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Неверный формат конечной даты", http.StatusBadRequest)
		return
	}

	// Получение userId из контекста
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Фильтрация расходов
	var expenses []models.Expense
	if err := db.DB.Where("user_id = ? AND created_at >= ? AND created_at <= ?", userId, startDate, endDate).Find(&expenses).Error; err != nil {
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	// Возврат данных в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func RenderFilteredExpenses(w http.ResponseWriter, r *http.Request) {
	// Получение параметров фильтрации
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// Если параметры пустые, устанавливаем значения по умолчанию
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Неверный формат начальной даты", http.StatusBadRequest)
			return
		}
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Неверный формат конечной даты", http.StatusBadRequest)
			return
		}
	} else {
		// Установить дату начала как "давно", а конца - "сейчас"
		startDate = time.Time{} // Zero time для начала всех времён
		endDate = time.Now()
	}

	// Получение userId из контекста
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Фильтрация расходов
	var expenses []models.Expense
	db.DB.Where("user_id = ? AND created_at >= ? AND created_at <= ?", userId, startDate, endDate).Find(&expenses)

	// Группировка расходов по дням с подсчётом суммы
	groupedExpenses := GroupExpensesByDayWithTotals(expenses)

	// Получение категорий пользователя
	var categories []models.Category
	db.DB.Where("user_id = ?", userId).Find(&categories)

	// Передача данных в шаблон
	data := map[string]interface{}{
		"ExpensesByDay": groupedExpenses,
		"User": map[string]interface{}{
			"Categories": categories,
		},
		"Filter": map[string]string{
			"StartDate": startDateStr,
			"EndDate":   endDateStr,
		},
	}

	render.RenderTemplate(w, "filterdate", data)
}

func GroupExpensesByDay(expenses []models.Expense) map[string][]models.Expense {
	grouped := make(map[string][]models.Expense)
	for _, expense := range expenses {
		day := FormatDateForDisplay(expense.CreatedAt) // Форматируем дату
		grouped[day] = append(grouped[day], expense)
	}
	return grouped
}

func GroupExpensesByDayWithTotals(expenses []models.Expense) map[string]map[string]interface{} {
	grouped := make(map[string]map[string]interface{})

	for _, expense := range expenses {
		day := expense.CreatedAt.Format("2006-01-02")
		if _, exists := grouped[day]; !exists {
			grouped[day] = map[string]interface{}{
				"Expenses": []models.Expense{},
				"Total":    0.0,
			}
		}

		// Добавляем расход в список
		grouped[day]["Expenses"] = append(grouped[day]["Expenses"].([]models.Expense), expense)
		// Увеличиваем общую сумму
		grouped[day]["Total"] = grouped[day]["Total"].(float64) + expense.Price*float64(expense.Amount)
	}

	return grouped
}
