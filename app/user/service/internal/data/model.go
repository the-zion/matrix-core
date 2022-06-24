package data

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uuid     string `gorm:"uniqueIndex;size:36"`
	Email    string `gorm:"uniqueIndex;size:50"`
	Phone    string `gorm:"uniqueIndex;size:20"`
	Wechat   string `gorm:"uniqueIndex;size:100"`
	Github   string `gorm:"uniqueIndex;size:100"`
	Password string `gorm:"size:500"`
}

type Profile struct {
	Updated   int64
	Uuid      string `gorm:"primaryKey;size:36"`
	Username  string `gorm:"uniqueIndex;not null;size:20"`
	Avatar    string `gorm:"size:200"`
	School    string `gorm:"size:50"`
	Company   string `gorm:"size:50"`
	Job       string `gorm:"size:50"`
	Homepage  string `gorm:"size:100"`
	Introduce string `gorm:"size:100"`
}

type ProfileUpdate struct {
	Profile
	Status int32 `gorm:"default:1"`
}

type ProfileUpdateRetry struct {
	Profile
}
