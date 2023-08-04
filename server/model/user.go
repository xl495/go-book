package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null;size:1:50" json:"username"`
	Email    string `gorm:"uniqueIndex;not null;size:1:50" json:"email"`
	Password string `gorm:"not null;size:1:100" json:"password"`
	Names    string `json:"names;size:1:100"`
}
