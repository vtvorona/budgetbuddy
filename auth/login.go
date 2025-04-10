package auth

import (
	db "budgetbuddy/db"
	"budgetbuddy/jwt"
	"budgetbuddy/models"
	"budgetbuddy/render"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		password := r.FormValue("password")

		var user models.User
		result := db.DB.Where("email = ? AND password = ?", email, password).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				data := models.Data{Error: "Неверный email или пароль", IsAuth: false}
				render.RenderTemplate(w, "login", data)
				return
			}
		} else {
			token, _ := jwt.GenerateToken(int(user.ID))
			cookie := &http.Cookie{
				Name:     "auth_token",
				Value:    token,
				HttpOnly: true,
				Secure:   false,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
			}

			http.SetCookie(w, cookie)

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}
}

func GetUser(id int) (models.User, error) {
	var user models.User
	result := db.DB.Where("ID = ?", id).Find(&user)
	return user, result.Error
}
