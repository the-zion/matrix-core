package data

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ArticleId int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Uuid      string `gorm:"index;size:36"`
	Status    int32  `gorm:"default:1"`
}

type ArticleStatistic struct {
	ArticleId int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Uuid      string `gorm:"index;size:36"`
	Agree     int32  `gorm:"default:0"`
	View      int32  `gorm:"default:0"`
	Collect   int32  `gorm:"default:0"`
	Comment   int32  `gorm:"default:0"`
}

type ArticleDraft struct {
	gorm.Model
	Uuid   string `gorm:"index;size:36"`
	Status int32  `gorm:"default:1"`
}
