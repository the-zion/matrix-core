package data

import "gorm.io/gorm"

type SystemNotification struct {
	gorm.Model
	ContentId        int32
	NotificationType string
	Title            string `gorm:"size:100"`
	Uid              string
	Uuid             string `gorm:"index;size:36"`
	Label            string `gorm:"size:100"`
	Result           int32
	Section          string
	Text             string
	Comment          string `gorm:"size:100"`
}
