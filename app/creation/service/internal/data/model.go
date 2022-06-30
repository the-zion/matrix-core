package data

import (
	"gorm.io/gorm"
)

type ArticleDraft struct {
	gorm.Model
	Updated int64
	Uuid    string `gorm:"index;size:36"`
	Status  int32  `gorm:"default:1"`
}

type ArticleDraftRetry struct {
	ArticleDraft
}
