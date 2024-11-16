package auth

import (
	db "budgetbuddy/db"
	"budgetbuddy/models"
	"budgetbuddy/server"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var user models.User

		// Ищем пользователя в базе данных
		result := db.DB.Where("email = ? AND password = ?", email, password).First(&user)
		if result.Error != nil {
			// Если ошибка, значит пользователь не найден
			if result.Error == gorm.ErrRecordNotFound {
				data := PageData{Error: "Неверный email или пароль"}
				server.RenderTemplate(w, "login", data)
				return
			}
		} else {
			cookie := &http.Cookie{
				Name:     "user_id",                      // Имя куки
				Value:    fmt.Sprintf("%d", user.ID),     // Значение куки
				Expires:  time.Now().Add(24 * time.Hour), // Время жизни куки
				HttpOnly: true,                           // Доступно только через HTTP, не JavaScript
				Path:     "/",                            // Доступно на всём сайте
			}

			fmt.Println(cookie.Value)

			http.SetCookie(w, cookie)
			// Логика успешного входа (например, установка сессии, редирект и т.д.)
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}
}
