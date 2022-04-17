package data

import "gorm.io/gorm"

type Achievement struct {
	gorm.Model
	UserId   int64 `gorm:"uniqueIndex"`
	Follow   int64
	Followed int64
	Collect  int64
	View     int64
}

type Follower struct {
	gorm.Model
	Follow   int64 `gorm:"index"`
	Followed int64 `gorm:"index"`
}
