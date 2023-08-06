package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null;size:50" json:"username" validate:"required,min=3,max=50"`
	Email    string `gorm:"uniqueIndex;not null;size:50" json:"email"`
	Password string `gorm:"column:password;not null;size:1:100" json:"password" validate:"required,min=3,max=50"`
	Names    string `json:"names"`
}
