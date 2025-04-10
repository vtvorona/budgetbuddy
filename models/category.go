package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID     int `gorm:"primaryKey"`
	UserID int
	Name   string
}
