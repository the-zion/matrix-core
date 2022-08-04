package data

import (
	"gorm.io/gorm"
)

type CommentDraft struct {
	gorm.Model
	Uuid   string `gorm:"index;size:36"`
	Status int32  `gorm:"default:1"`
}

type Comment struct {
	gorm.Model
	CreationId   int32
	CreationType int32
	Uuid         string `gorm:"index;size:36"`
	Agree        int32  `gorm:"index;default:0"`
}

type SubComment struct {
	gorm.Model
	CommentId int32 `gorm:"index"`
	ParentId  int32
	Agree     int32
}

type Record struct {
	gorm.Model
	Uuid     string `gorm:"size:36"`
	CommonId int32
	Ip       string
}
