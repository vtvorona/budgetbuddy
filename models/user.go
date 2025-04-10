package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            int `gorm:"primaryKey"`
	Name          string
	Surname       string
	Email         string `gorm:"unique"`
	Password      string
	MonthlyBudget float64
	Categories    []*Category
}
