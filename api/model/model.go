package model

import "github.com/jinzhu/gorm"

// Model model
type Model struct {
	db *gorm.DB
}

// New new model
func New(db *gorm.DB) *Model {
	return &Model{db}
}
