package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string   `gorm:"not null;size: 100" json:"title" validate:"required,min=3,max=100"`
	Description string   `gorm:"not null;size: 1000" json:"description" validate:"required,min=1,max=1000"`
	Author      string   `gorm:"not null;size: 50" json:"author" validate:"required,min=3,max=50"`
	Price       float64  `gorm:"not null;size: 20" json:"price" validate:"required"`
	CategoryID  uint     `gorm:"not null;size: 100" json:"categoryId" validate:"required,min=0,max=10000"`
	Category    Category `gorm:"foreignKey:CategoryID" json:"category"`
}
