package data

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Uuid     string `gorm:"uniqueIndex;size:20"`
	Email    string `gorm:"uniqueIndex;size:50"`
	Phone    string `gorm:"uniqueIndex;size:20"`
	Wechat   string `gorm:"uniqueIndex;size:100"`
	Qq       string `gorm:"uniqueIndex;size:100"`
	Gitee    int32  `gorm:"uniqueIndex;size:100"`
	Github   int32  `gorm:"uniqueIndex"`
	Password string `gorm:"size:500"`
}

//easyjson:json
type Profile struct {
	CreatedAt time.Time
	Updated   int64
	Uuid      string `gorm:"primaryKey;size:20"`
	Username  string `gorm:"index;not null;size:100"`
	Avatar    string `gorm:"size:200"`
	School    string `gorm:"size:50"`
	Company   string `gorm:"size:50"`
	Job       string `gorm:"size:50"`
	Homepage  string `gorm:"size:100"`
	Github    string `gorm:"size:100"`
	Gitee     string `gorm:"size:100"`
	Introduce string `gorm:"size:100"`
}

type ProfileUpdate struct {
	Profile
	Status int32 `gorm:"default:1"`
}

type AvatarReview struct {
	gorm.Model
	Uuid     string `gorm:"index;size:20"`
	JobId    string `gorm:"size:100"`
	Url      string `gorm:"size:1000"`
	Label    string `gorm:"size:100"`
	Result   int32
	Category string `gorm:"size:100"`
	SubLabel string `gorm:"size:100"`
	Score    int32
}

type CoverReview struct {
	gorm.Model
	Uuid     string `gorm:"index;size:20"`
	JobId    string `gorm:"size:100"`
	Url      string `gorm:"size:1000"`
	Label    string `gorm:"size:100"`
	Result   int32
	Category string `gorm:"size:100"`
	SubLabel string `gorm:"size:100"`
	Score    int32
}

type Follow struct {
	gorm.Model
	Follow   string `gorm:"uniqueIndex:idx_follow;size:20"`
	Followed string `gorm:"uniqueIndex:idx_follow;size:20"`
	Status   int32
}
