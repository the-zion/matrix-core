package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/biz"
)

var _ biz.AchievementRepo = (*achievementRepo)(nil)

type achievementRepo struct {
	data *Data
	log  *log.Helper
}

func NewAchievementRepo(data *Data, logger log.Logger) biz.AchievementRepo {
	return &achievementRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "achievement/data/achievement")),
	}
}
