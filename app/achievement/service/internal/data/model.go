package data

import (
	"gorm.io/gorm"
)

type Achievement struct {
	gorm.Model
	Uuid     string `gorm:"primaryKey;size:36"`
	Agree    int32  `gorm:"default:0"`
	Collect  int32  `gorm:"default:0"`
	View     int32  `gorm:"default:0"`
	Follow   int32  `gorm:"default:0"`
	Followed int32  `gorm:"default:0"`
}
