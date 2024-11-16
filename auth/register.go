package auth

import (
	db "budgetbuddy/db"
	"budgetbuddy/models"
	"budgetbuddy/server"
	"net/http"
)

type PageData struct {
	Error string
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		user := &models.User{
			Name:     r.FormValue("name"),
			Surname:  r.FormValue("surname"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		var existingUser models.User

		exists := db.DB.Where("email = ?", user.Email).First(&existingUser)
		if exists.Error == nil {
			data := PageData{Error: "Пользователь с этим email уже существует"}
			server.RenderTemplate(w, "registration", data)
			return
		}

		result := db.DB.Create(user)
		if result.Error != nil {
			http.Error(w, "User already exists or an error occurred", http.StatusConflict)
			return
		}

		http.Redirect(w, r, "/success", http.StatusSeeOther)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
