package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

const SecretKey = "ef97e31910863ca504cfef2a6fd778f3682da302305791689b8b9a11b494dbd69cdcd62e4a19b95193da979d50943040f653709f66a289024dcad09e914b3035"

func GenerateToken(id int) (string, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(signedToken string, SecretKey string) (int, error) {
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		id := int(claims["id"].(float64))
		return id, nil
	}
	return 0, errors.New("Invalid token")
}

func ExtractIdFromToken(tokenString string, SecretKey string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неподдерживаемый метод подписи")
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга токена: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["id"].(float64); ok {
			return int(id), nil
		} else {
			return 0, fmt.Errorf("id не найдено в токене")
		}
	} else {
		return 0, fmt.Errorf("невалидный токен")
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем токен
		token := cookie.Value
		userId, err := ValidateToken(token, SecretKey)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/" {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)

		if r.URL.Path == "/expense/create" || r.URL.Path == "/categories/add" {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
