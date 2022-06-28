package data

import "gorm.io/gorm"

type ArticleDraft struct {
	gorm.Model
	Uuid   string `gorm:"index;size:36"`
	Title  string `gorm:"size:50"`
	Status int32  `gorm:"default:1"`
}
