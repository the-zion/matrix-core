package data

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Uuid     string `gorm:"uniqueIndex;size:200"`
	Email    string `gorm:"uniqueIndex;size:200"`
	Phone    string `gorm:"uniqueIndex;size:200"`
	Wechat   string `gorm:"uniqueIndex;size:500"`
	Github   string `gorm:"uniqueIndex;size:500"`
	Password string `gorm:"size:500"`
}

type Profile struct {
	gorm.Model
	Uuid      string `gorm:"uniqueIndex;size:200"`
	Username  string `gorm:"uniqueIndex;size:100"`
	Avatar    string `gorm:"size:200"`
	School    string `gorm:"size:100"`
	Company   string `gorm:"size:100"`
	Job       string `gorm:"size:100"`
	Homepage  string `gorm:"size:200"`
	Introduce string `gorm:"size:500"`
}

type ProfileUpdate struct {
	Profile
	Status int64 `gorm:"default:1"`
}
