package models

import "time"

// Expense представляет траты пользователя
type Expense struct {
	ID          int `json:"id"`      // Уникальный идентификатор
	UserID      int `json:"user_id"` // ID пользователя, которому принадлежит расход
	Title       string
	Amount      float64   `json:"amount"`      // Сумма расхода
	Category    string    `json:"category"`    // Категория расхода (например, "еда", "транспорт")
	Description string    `json:"description"` // Описание расхода
	Date        time.Time `json:"date"`        // Дата расхода
	CreatedAt   time.Time `json:"created_at"`  // Дата создания записи
	UpdatedAt   time.Time `json:"updated_at"`  // Дата обновления записи
}
