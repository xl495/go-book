package model

import "gorm.io/gorm"

// Category 分类
type Category struct {
	gorm.Model
	Name        string `gorm:"not null;size: 50" json:"name" validate:"required,min=1,max=50"`
	Description string `gorm:"not null;size: 100" json:"description" validate:"min=3,max=100"`
	Books       []Book `gorm:"foreignKey:CategoryID" json:"books"`
}

// GetBooks 获取与该分类关联的书籍
func (category *Category) GetBooks() []Book {
	return category.Books
}
