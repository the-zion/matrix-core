package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/biz"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *achievementRepo) SetAchievementAgree(ctx context.Context, uuid string) error {
	ach := &Achievement{
		Uuid:  uuid,
		Agree: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"agree": gorm.Expr("agree + ?", 1)}),
	}).Create(ach).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add achievement agree: c(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) CancelAchievementAgree(ctx context.Context, uuid string) error {
	ach := &Achievement{}
	err := r.data.DB(ctx).Model(ach).Where("uuid = ?", uuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to subtract achievement agree: uuid(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementView(ctx context.Context, uuid string) error {
	ach := &Achievement{
		Uuid: uuid,
		View: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"view": gorm.Expr("view + ?", 1)}),
	}).Create(ach).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add achievement view: uuid(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementCollect(ctx context.Context, uuid string) error {
	ach := &Achievement{
		Uuid:    uuid,
		Collect: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"collect": gorm.Expr("collect + ?", 1)}),
	}).Create(ach).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add achievement collect: uuid(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) CancelAchievementCollect(ctx context.Context, uuid string) error {
	ach := &Achievement{}
	err := r.data.DB(ctx).Model(ach).Where("uuid = ?", uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to subtract achievement collect: uuid(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementAgreeToCache(ctx context.Context, uuid string) error {
	_, err := r.data.redisCli.HIncrBy(ctx, uuid, "agree", 1).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set achievement agree to cache: uuid(%s)", uuid))
	}
	return nil
}

func (r *achievementRepo) CancelAchievementAgreeFromCache(ctx context.Context, uuid string) error {
	_, err := r.data.redisCli.HIncrBy(ctx, uuid, "agree", -1).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel achievement agree from cache: uuid(%s)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementViewToCache(ctx context.Context, uuid string) error {
	_, err := r.data.redisCli.HIncrBy(ctx, uuid, "view", 1).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set achievement view to cache: uuid(%s)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementCollectToCache(ctx context.Context, uuid string) error {
	_, err := r.data.redisCli.HIncrBy(ctx, uuid, "collect", 1).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set achievement collect to cache: uuid(%s)", uuid))
	}
	return nil
}

func (r *achievementRepo) CancelAchievementCollectFromCache(ctx context.Context, uuid string) error {
	_, err := r.data.redisCli.HIncrBy(ctx, uuid, "collect", -1).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel achievement collect from cache: uuid(%s)", uuid))
	}
	return nil
}
