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
	CommentId      int32 `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	CreationId     int32          `gorm:"index:creation"`
	CreationType   int32          `gorm:"index:creation"`
	CreationAuthor string         `gorm:"index;size:36"`
	Uuid           string         `gorm:"index;size:36"`
	Agree          int32          `gorm:"index;type:int unsigned;default:0"`
	Comment        int32          `gorm:"type:int unsigned;default:0"`
}

type CommentAgree struct {
	gorm.Model
	CommentId int32  `gorm:"uniqueIndex:idx_unique"`
	Uuid      string `gorm:"uniqueIndex:idx_unique;size:36"`
	Status    int32  `gorm:"default:1"`
}

type SubComment struct {
	CommentId      int32 `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	CreationId     int32          `gorm:"index:creation"`
	CreationType   int32          `gorm:"index:creation"`
	CreationAuthor string         `gorm:"size:36"`
	RootId         int32          `gorm:"index"`
	RootUser       string         `gorm:"index;size:36"`
	ParentId       int32          `gorm:"index"`
	Uuid           string         `gorm:"index;size:36"`
	Reply          string         `gorm:"index;size:36"`
	Agree          int32          `gorm:"type:int unsigned;default:0"`
}

type CommentUser struct {
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Uuid              string `gorm:"primarykey;size:36"`
	Comment           int32  `gorm:"type:int unsigned;default:0"`
	ArticleReply      int32  `gorm:"type:int unsigned;default:0"`
	ArticleReplySub   int32  `gorm:"type:int unsigned;default:0"`
	TalkReply         int32  `gorm:"type:int unsigned;default:0"`
	TalkReplySub      int32  `gorm:"type:int unsigned;default:0"`
	ArticleReplied    int32  `gorm:"type:int unsigned;default:0"`
	ArticleRepliedSub int32  `gorm:"type:int unsigned;default:0"`
	TalkReplied       int32  `gorm:"type:int unsigned;default:0"`
	TalkRepliedSub    int32  `gorm:"type:int unsigned;default:0"`
}

type CommentContentReview struct {
	gorm.Model
	CommentId int32
	Comment   string
	Kind      string
	Uuid      string `gorm:"index;size:36"`
	JobId     string `gorm:"size:100"`
	Label     string `gorm:"size:100"`
	Result    int32
	Section   string
}

type Record struct {
	gorm.Model
	Uuid     string `gorm:"size:36"`
	CommonId int32
	Ip       string
}
