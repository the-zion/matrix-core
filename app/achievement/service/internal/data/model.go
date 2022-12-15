package data

type Achievement struct {
	Uuid     string `gorm:"primaryKey;size:20"`
	Score    int32  `gorm:"type:int unsigned;default:0"`
	Agree    int32  `gorm:"type:int unsigned;default:0"`
	Collect  int32  `gorm:"type:int unsigned;default:0"`
	View     int32  `gorm:"type:int unsigned;default:0"`
	Follow   int32  `gorm:"type:int unsigned;default:0"`
	Followed int32  `gorm:"type:int unsigned;default:0"`
}

type Active struct {
	Uuid  string `gorm:"primaryKey;size:20"`
	Agree int32  `gorm:"type:int unsigned;default:0"`
}

type Medal struct {
	Uuid      string `gorm:"primaryKey;size:20"`
	Creation1 int32  `gorm:"type:int unsigned;default:0"`
	Creation2 int32  `gorm:"type:int unsigned;default:0"`
	Creation3 int32  `gorm:"type:int unsigned;default:0"`
	Creation4 int32  `gorm:"type:int unsigned;default:0"`
	Creation5 int32  `gorm:"type:int unsigned;default:0"`
	Creation6 int32  `gorm:"type:int unsigned;default:0"`
	Creation7 int32  `gorm:"type:int unsigned;default:0"`
	Agree1    int32  `gorm:"type:int unsigned;default:0"`
	Agree2    int32  `gorm:"type:int unsigned;default:0"`
	Agree3    int32  `gorm:"type:int unsigned;default:0"`
	Agree4    int32  `gorm:"type:int unsigned;default:0"`
	Agree5    int32  `gorm:"type:int unsigned;default:0"`
	Agree6    int32  `gorm:"type:int unsigned;default:0"`
	View1     int32  `gorm:"type:int unsigned;default:0"`
	View2     int32  `gorm:"type:int unsigned;default:0"`
	View3     int32  `gorm:"type:int unsigned;default:0"`
	Comment1  int32  `gorm:"type:int unsigned;default:0"`
	Comment2  int32  `gorm:"type:int unsigned;default:0"`
	Comment3  int32  `gorm:"type:int unsigned;default:0"`
	Collect1  int32  `gorm:"type:int unsigned;default:0"`
	Collect2  int32  `gorm:"type:int unsigned;default:0"`
	Collect3  int32  `gorm:"type:int unsigned;default:0"`
}
