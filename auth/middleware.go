package auth

import (
	"context"
	"fmt"
	"net/http"
)

// Определяем тип для ключа в контексте (чтобы избежать коллизий)
type contextKey string

// Константа для ключа
const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil {
			if err == http.ErrNoCookie {
				// Если куки нет, перенаправляем на страницу логина
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Достаём user_id из куки
		userID := cookie.Value

		// Можно добавить проверку или загрузку пользователя из базы
		fmt.Println("Пользователь ID:", userID)

		// Передаём user_id в контекст
		ctx := context.WithValue(r.Context(), UserIDKey, cookie.Value)

		// Передаём управление следующему обработчику
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
