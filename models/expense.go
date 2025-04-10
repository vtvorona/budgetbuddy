package models

import "gorm.io/gorm"

type Expense struct {
	gorm.Model
	ID         int `gorm:"primaryKey"`
	UserID     int `gorm:"index"`
	Title      string
	Price      float64
	Amount     int
	Total      float64
	CategoryID int
	Category   Category
}
