package data

import (
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var _ biz.ProfileRepo = (*profileRepo)(nil)

type profileRepo struct {
	data *Data
	log  *log.Helper
}

func NewProfileRepo(data *Data, logger log.Logger) biz.ProfileRepo {
	return &profileRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "user/data/profile")),
	}
}

type Profile struct {
	gorm.Model
	UserId          int64  `gorm:"uniqueIndex"`
	Username        string `gorm:"uniqueIndex;size:200"`
	Sex             string `gorm:"size:100"`
	Introduce       string `gorm:"size:200"`
	Address         string `gorm:"size:100"`
	Industry        string `gorm:"size:100"`
	PersonalProfile string `gorm:"size:500"`
	Tag             string `gorm:"size:100"`
	Background      string `gorm:"size:500"`
	Image           string `gorm:"size:500"`
}
