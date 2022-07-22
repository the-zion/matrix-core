package data

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ArticleId int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Uuid      string         `gorm:"index;size:36"`
	Status    int32          `gorm:"default:1"`
	Auth      int32          `gorm:"default:1"`
}

type Talk struct {
	TalkId    int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Uuid      string         `gorm:"index;size:36"`
	Status    int32          `gorm:"default:1"`
	Auth      int32          `gorm:"default:1"`
}

type ArticleStatistic struct {
	ArticleId int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Uuid      string         `gorm:"index;size:36"`
	Agree     int32          `gorm:"index;default:0"`
	View      int32          `gorm:"default:0"`
	Collect   int32          `gorm:"default:0"`
	Comment   int32          `gorm:"default:0"`
	Auth      int32          `gorm:"default:1"`
}

type TalkStatistic struct {
	TalkId    int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Uuid      string         `gorm:"index;size:36"`
	Agree     int32          `gorm:"index;default:0"`
	View      int32          `gorm:"default:0"`
	Collect   int32          `gorm:"default:0"`
	Comment   int32          `gorm:"default:0"`
	Auth      int32          `gorm:"default:1"`
}

type ArticleDraft struct {
	gorm.Model
	Uuid   string `gorm:"index;size:36"`
	Status int32  `gorm:"default:1"`
}

type TalkDraft struct {
	gorm.Model
	Uuid   string `gorm:"index;size:36"`
	Status int32  `gorm:"default:1"`
}

type Collections struct {
	gorm.Model
	Uuid      string `gorm:"index;size:36"`
	Name      string `gorm:"size:50"`
	Introduce string `gorm:"size:100"`
	Auth      int32  `gorm:"default:1"`
}

type Collect struct {
	CollectionsId int32 `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Uuid          string `gorm:"size:36"`
	CreationsId   int32  `gorm:"primarykey"`
	Mode          int32  `gorm:"primarykey;default:1"`
	Status        int32  `gorm:"default:1"`
}

type Column struct {
	ColumnId  int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Uuid      string         `gorm:"index;size:36"`
	Status    int32          `gorm:"default:1"`
	Auth      int32          `gorm:"default:1"`
}

type ColumnDraft struct {
	gorm.Model
	Uuid   string `gorm:"index;size:36"`
	Status int32  `gorm:"default:1"`
}

type ColumnStatistic struct {
	ColumnId  int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Uuid      string         `gorm:"index;size:36"`
	Agree     int32          `gorm:"index;default:0"`
	View      int32          `gorm:"default:0"`
	Collect   int32          `gorm:"default:0"`
	Auth      int32          `gorm:"default:1"`
}

type ColumnInclusion struct {
	ColumnId  int32 `gorm:"primarykey"`
	ArticleId int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Uuid      string `gorm:"size:36"`
	Status    int32  `gorm:"default:1"`
}
