package data

import (
	"gorm.io/gorm"
	"time"
)

type CommentDraft struct {
	gorm.Model
	Uuid   string `gorm:"index;size:36"`
	Status int32  `gorm:"default:1"`
}

type Comment struct {
	CommentId    int32 `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	CreationId   int32          `gorm:"index:creation"`
	CreationType int32          `gorm:"index:creation"`
	Uuid         string         `gorm:"index;size:36"`
	Agree        int32          `gorm:"index;default:0"`
	Comment      int32          `gorm:"default:0"`
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
