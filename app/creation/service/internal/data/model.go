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
	Agree     int32          `gorm:"index;type:int unsigned;default:0"`
	View      int32          `gorm:"type:int unsigned;default:0"`
	Collect   int32          `gorm:"type:int unsigned;default:0"`
	Comment   int32          `gorm:"type:int unsigned;default:0"`
	Auth      int32          `gorm:"default:1"`
}

type TalkStatistic struct {
	TalkId    int32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Uuid      string         `gorm:"index;size:36"`
	Agree     int32          `gorm:"index;type:int unsigned;default:0"`
	View      int32          `gorm:"type:int unsigned;default:0"`
	Collect   int32          `gorm:"type:int unsigned;default:0"`
	Comment   int32          `gorm:"type:int unsigned;default:0"`
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

type ArticleAgree struct {
	gorm.Model
	ArticleId int32  `gorm:"uniqueIndex:idx_unique"`
	Uuid      string `gorm:"uniqueIndex:idx_unique;size:36"`
	Status    int32  `gorm:"default:1"`
}

type ArticleCollect struct {
	gorm.Model
	ArticleId int32  `gorm:"uniqueIndex:idx_unique"`
	Uuid      string `gorm:"uniqueIndex:idx_unique;size:36"`
	Status    int32  `gorm:"default:1"`
}

type TalkAgree struct {
	gorm.Model
	TalkId int32  `gorm:"uniqueIndex:idx_unique"`
	Uuid   string `gorm:"uniqueIndex:idx_unique;size:36"`
	Status int32  `gorm:"default:1"`
}

type TalkCollect struct {
	gorm.Model
	TalkId int32  `gorm:"uniqueIndex:idx_unique"`
	Uuid   string `gorm:"uniqueIndex:idx_unique;size:36"`
	Status int32  `gorm:"default:1"`
}

type ColumnAgree struct {
	gorm.Model
	ColumnId int32  `gorm:"uniqueIndex:idx_unique"`
	Uuid     string `gorm:"uniqueIndex:idx_unique;size:36"`
	Status   int32  `gorm:"default:1"`
}

type ColumnCollect struct {
	gorm.Model
	ColumnId int32  `gorm:"uniqueIndex:idx_unique"`
	Uuid     string `gorm:"uniqueIndex:idx_unique;size:36"`
	Status   int32  `gorm:"default:1"`
}

type Collections struct {
	gorm.Model
	Uuid      string `gorm:"index;size:36"`
	Name      string `gorm:"size:50"`
	Introduce string `gorm:"size:100"`
	Auth      int32  `gorm:"default:1"`
	Article   int32  `gorm:"type:int unsigned;default:0"`
	Column    int32  `gorm:"type:int unsigned;default:0"`
	Talk      int32  `gorm:"type:int unsigned;default:0"`
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
	Agree     int32          `gorm:"type:int unsigned;index;default:0"`
	View      int32          `gorm:"type:int unsigned;default:0"`
	Collect   int32          `gorm:"type:int unsigned;default:0"`
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

type Subscribe struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ColumnId  int32  `gorm:"primarykey"`
	Uuid      string `gorm:"primarykey;size:36"`
	AuthorId  string `gorm:"size:36"`
	Status    int32  `gorm:"default:1"`
}

type Record struct {
	gorm.Model
	Uuid       string `gorm:"size:36"`
	CreationId int32
	Mode       int32 `gorm:"default:1"`
	Operation  string
	Ip         string
}

type CreationUser struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Uuid        string `gorm:"primarykey;size:36"`
	Article     int32  `gorm:"type:int unsigned;default:0"`
	Column      int32  `gorm:"type:int unsigned;default:0"`
	Talk        int32  `gorm:"type:int unsigned;default:0"`
	Collections int32  `gorm:"type:int unsigned;default:0"`
	Collect     int32  `gorm:"type:int unsigned;default:0"`
	Subscribe   int32  `gorm:"type:int unsigned;default:0"`
}

type CreationUserVisitor struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Uuid        string `gorm:"primarykey;size:36"`
	Article     int32  `gorm:"type:int unsigned;default:0"`
	Column      int32  `gorm:"type:int unsigned;default:0"`
	Talk        int32  `gorm:"type:int unsigned;default:0"`
	Collections int32  `gorm:"type:int unsigned;default:0"`
}
