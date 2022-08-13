package data

type Achievement struct {
	Uuid     string `gorm:"primaryKey;size:36"`
	Score    int32  `gorm:"type:int unsigned;default:0"`
	Agree    int32  `gorm:"type:int unsigned;default:0"`
	Collect  int32  `gorm:"type:int unsigned;default:0"`
	View     int32  `gorm:"type:int unsigned;default:0"`
	Follow   int32  `gorm:"type:int unsigned;default:0"`
	Followed int32  `gorm:"type:int unsigned;default:0"`
}
