package server

import (
	"budgetbuddy/db"
	"budgetbuddy/models"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetExpenses(userId int) *models.Expenses {

	var todayExpenses []*models.Expense
	var monthExpenses []*models.Expense

	now := time.Now()
	today := now.Truncate(24 * time.Hour)

	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayOfNextMonth := firstDayOfMonth.AddDate(0, 1, 0)

	// Загрузка расходов с категориями
	db.DB.Preload("Category").Where("user_id = ? AND created_at >= ?", userId, today).Find(&todayExpenses)
	db.DB.Preload("Category").Where("user_id = ? AND created_at >= ? AND created_at < ?", userId, firstDayOfMonth, firstDayOfNextMonth).Find(&monthExpenses)

	var todayTotal float64
	for _, expense := range todayExpenses {
		todayTotal += expense.Price * float64(expense.Amount)
	}

	var monthTotal float64
	for _, expense := range monthExpenses {
		monthTotal += expense.Price * float64(expense.Amount)
	}

	expenses := &models.Expenses{
		TodayExpenses: todayExpenses,
		TodayTotal:    todayTotal,
		MonthExpenses: monthExpenses,
		MonthTotal:    monthTotal,
	}

	return expenses
}

func EditExpenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Получение данных из формы
	id := r.FormValue("id")
	title := r.FormValue("title")
	category := r.FormValue("category")
	category_id, _ := strconv.Atoi(category)
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	amount, _ := strconv.Atoi(r.FormValue("amount"))

	// Получение параметров фильтрации
	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Логика редактирования
	if err := db.DB.Model(&models.Expense{}).Where("id = ? AND user_id = ?", id, userId).Updates(models.Expense{
		Title:      title,
		CategoryID: category_id,
		Price:      price,
		Amount:     amount,
	}).Error; err != nil {
		http.Error(w, "Ошибка при обновлении расхода", http.StatusInternalServerError)
		return
	}

	// Построение URL для редиректа с параметрами фильтрации
	redirectURL := fmt.Sprintf("/filter?start_date=%s&end_date=%s", startDate, endDate)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func FormatDateForDisplay(date time.Time) string {
	months := map[string]string{
		"January":   "января",
		"February":  "февраля",
		"March":     "марта",
		"April":     "апреля",
		"May":       "мая",
		"June":      "июня",
		"July":      "июля",
		"August":    "августа",
		"September": "сентября",
		"October":   "октября",
		"November":  "ноября",
		"December":  "декабря",
	}

	day := date.Format("02")
	month := date.Format("January")
	return day + " " + months[month]
}
