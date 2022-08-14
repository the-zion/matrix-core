package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/biz"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
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
	err := r.data.DB(ctx).Model(ach).Where("uuid = ? and agree > 0", uuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
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
	err := r.data.DB(ctx).Model(ach).Where("uuid = ? and collect > 0", uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to subtract achievement collect: uuid(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementAgreeToCache(ctx context.Context, uuid string) error {
	var script = redis.NewScript(`
					local uuid = KEYS[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						redis.call("HINCRBY", uuid, "agree", 1)
					end
					return 0
	`)
	keys := []string{uuid}
	_, err := script.Run(ctx, r.data.redisCli, keys).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set achievement agree to cache: uuid(%s)", uuid))
	}
	return nil
}

func (r *achievementRepo) CancelAchievementAgreeFromCache(ctx context.Context, uuid string) error {
	var script = redis.NewScript(`
					local uuid = KEYS[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", uuid, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", uuid, "agree", -1)
						end
					end
					return 0
	`)
	keys := []string{uuid}
	_, err := script.Run(ctx, r.data.redisCli, keys).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel achievement agree from cache: uuid(%s)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementViewToCache(ctx context.Context, uuid string) error {
	var script = redis.NewScript(`
					local uuid = KEYS[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						redis.call("HINCRBY", uuid, "view", 1)
					end
					return 0
	`)
	keys := []string{uuid}
	_, err := script.Run(ctx, r.data.redisCli, keys).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set achievement view to cache: uuid(%s)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementCollectToCache(ctx context.Context, uuid string) error {
	var script = redis.NewScript(`
					local uuid = KEYS[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						redis.call("HINCRBY", uuid, "collect", 1)
					end
					return 0
	`)
	keys := []string{uuid}
	_, err := script.Run(ctx, r.data.redisCli, keys).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set achievement collect to cache: uuid(%s)", uuid))
	}
	return nil
}

func (r *achievementRepo) CancelAchievementCollectFromCache(ctx context.Context, uuid string) error {
	var script = redis.NewScript(`
					local uuid = KEYS[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", uuid, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", uuid, "collect", -1)
						end
					end
					return 0
	`)
	keys := []string{uuid}
	_, err := script.Run(ctx, r.data.redisCli, keys).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel achievement collect from cache: uuid(%s)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementFollow(ctx context.Context, uuid string) error {
	ach := &Achievement{
		Uuid:   uuid,
		Follow: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"follow": gorm.Expr("follow + ?", 1)}),
	}).Create(ach).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add achievement follow: uuid(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementFollowed(ctx context.Context, uuid string) error {
	ach := &Achievement{
		Uuid:     uuid,
		Followed: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"followed": gorm.Expr("followed + ?", 1)}),
	}).Create(ach).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add achievement followed: uuid(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) SetAchievementFollowToCache(ctx context.Context, follow, followed string) error {
	var script = redis.NewScript(`
					local follow = KEYS[1]
					local exist = redis.call("EXISTS", follow)
					if exist == 1 then
						redis.call("HINCRBY", follow, "followed", 1)
					end

					local followed = KEYS[1]
					local exist = redis.call("EXISTS", followed)
					if exist == 1 then
						redis.call("HINCRBY", followed, "follow", 1)
					end
					return 0
	`)
	keys := []string{follow, followed}
	_, err := script.Run(ctx, r.data.redisCli, keys).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set achievement follow to cache: follow(%s), followed(%s)", follow, followed))
	}
	return nil
}

func (r *achievementRepo) CancelAchievementFollow(ctx context.Context, uuid string) error {
	ach := &Achievement{}
	err := r.data.DB(ctx).Model(ach).Where("uuid = ? and follow > 0", uuid).Update("follow", gorm.Expr("follow - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to subtract achievement follow: uuid(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) CancelAchievementFollowed(ctx context.Context, uuid string) error {
	ach := &Achievement{}
	err := r.data.DB(ctx).Model(ach).Where("uuid = ? and followed > 0", uuid).Update("followed", gorm.Expr("followed - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to subtract achievement followed: uuid(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) CancelAchievementFollowFromCache(ctx context.Context, follow, followed string) error {
	var script = redis.NewScript(`
					local follow = KEYS[1]
					local exist = redis.call("EXISTS", follow)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", follow, "followed"))
						if number > 0 then
  							redis.call("HINCRBY", follow, "followed", -1)
						end
					end

					local followed = KEYS[1]
					local exist = redis.call("EXISTS", followed)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", followed, "follow"))
						if number > 0 then
  							redis.call("HINCRBY", followed, "follow", -1)
						end
					end
					return 0
	`)
	keys := []string{follow, followed}
	_, err := script.Run(ctx, r.data.redisCli, keys).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel achievement follow to cache: follow(%s), followed(%s)", follow, followed))
	}
	return nil
}

func (r *achievementRepo) GetAchievementList(ctx context.Context, uuids []string) ([]*biz.Achievement, error) {
	achievementList := make([]*biz.Achievement, 0)
	exists, unExists, err := r.achievementListExist(ctx, uuids)
	if err != nil {
		return nil, err
	}

	g, _ := errgroup.WithContext(ctx)
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(exists) == 0 {
			return nil
		}
		return r.getAchievementListFromCache(ctx, exists, &achievementList)
	}))
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(unExists) == 0 {
			return nil
		}
		return r.getAchievementListFromDb(ctx, unExists, &achievementList)
	}))

	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return achievementList, nil
}

func (r *achievementRepo) achievementListExist(ctx context.Context, uuids []string) ([]string, []string, error) {
	exists := make([]string, 0)
	unExists := make([]string, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range uuids {
			pipe.Exists(ctx, item)
		}
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, fmt.Sprintf("fail to check if achievement exist from cache: uuids(%v)", uuids))
	}

	for index, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		if exist == 1 {
			exists = append(exists, uuids[index])
		} else {
			unExists = append(unExists, uuids[index])
		}
	}
	return exists, unExists, nil
}

func (r *achievementRepo) getAchievementListFromCache(ctx context.Context, exists []string, achievementList *[]*biz.Achievement) error {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, uuid := range exists {
			pipe.HMGet(ctx, uuid, "agree", "collect", "view", "follow", "followed")
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get achievement list from cache: ids(%v)", exists))
	}

	for index, item := range cmd {
		val := []int32{0, 0, 0, 0, 0}
		for _index, count := range item.(*redis.SliceCmd).Val() {
			if count == nil {
				break
			}
			num, err := strconv.ParseInt(count.(string), 10, 32)
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
			}
			val[_index] = int32(num)
		}
		*achievementList = append(*achievementList, &biz.Achievement{
			Uuid:     exists[index],
			Agree:    val[0],
			Collect:  val[1],
			View:     val[2],
			Follow:   val[3],
			Followed: val[4],
		})
	}
	return nil
}

func (r *achievementRepo) getAchievementListFromDb(ctx context.Context, unExists []string, achievementList *[]*biz.Achievement) error {
	list := make([]*Achievement, 0)
	err := r.data.db.WithContext(ctx).Where("uuid IN ?", unExists).Find(&list).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get achievement list from db: uuids(%v)", unExists))
	}

	for _, item := range list {
		*achievementList = append(*achievementList, &biz.Achievement{
			Uuid:     item.Uuid,
			View:     item.View,
			Agree:    item.Agree,
			Collect:  item.Collect,
			Follow:   item.Follow,
			Followed: item.Followed,
		})
	}

	if len(list) != 0 {
		go r.setAchievementListToCache(list)
	}
	return nil
}

func (r *achievementRepo) setAchievementListToCache(achievementList []*Achievement) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, item := range achievementList {
			key := item.Uuid
			pipe.HSetNX(context.Background(), key, "agree", item.Agree)
			pipe.HSetNX(context.Background(), key, "collect", item.Collect)
			pipe.HSetNX(context.Background(), key, "view", item.View)
			pipe.HSetNX(context.Background(), key, "follow", item.Follow)
			pipe.HSetNX(context.Background(), key, "followed", item.Followed)
			pipe.Expire(context.Background(), key, time.Hour*8)
		}
		return nil
	})

	if err != nil {
		r.log.Errorf("fail to set achievement to cache, err(%s)", err.Error())
	}
}

func (r *achievementRepo) GetUserAchievement(ctx context.Context, uuid string) (*biz.Achievement, error) {
	var achievement *biz.Achievement
	key := uuid
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		achievement, err = r.getAchievementFromCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return achievement, nil
	}

	achievement, err = r.getAchievementFromDB(ctx, uuid)
	if err != nil {
		return nil, err
	}

	go r.setAchievementToCache(key, achievement)

	return achievement, nil
}

func (r *achievementRepo) getAchievementFromCache(ctx context.Context, key string) (*biz.Achievement, error) {
	achievement, err := r.data.redisCli.HMGet(ctx, key, "agree", "collect", "view", "follow", "followed").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get achievement form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0, 0}
	for _index, count := range achievement {
		if count == nil {
			break
		}
		num, err := strconv.ParseInt(count.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
		}
		val[_index] = int32(num)
	}
	return &biz.Achievement{
		Agree:    val[0],
		Collect:  val[1],
		View:     val[2],
		Follow:   val[3],
		Followed: val[4],
	}, nil
}

func (r *achievementRepo) getAchievementFromDB(ctx context.Context, uuid string) (*biz.Achievement, error) {
	ach := &Achievement{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).First(ach).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get achievement from db: uuid(%s)", uuid))
	}
	return &biz.Achievement{
		Uuid:     ach.Uuid,
		Agree:    ach.Agree,
		Collect:  ach.Collect,
		View:     ach.View,
		Follow:   ach.Follow,
		Followed: ach.Followed,
	}, nil
}

func (r *achievementRepo) setAchievementToCache(key string, achievement *biz.Achievement) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		pipe.HSetNX(context.Background(), key, "agree", achievement.Agree)
		pipe.HSetNX(context.Background(), key, "collect", achievement.Collect)
		pipe.HSetNX(context.Background(), key, "view", achievement.View)
		pipe.HSetNX(context.Background(), key, "follow", achievement.Follow)
		pipe.HSetNX(context.Background(), key, "followed", achievement.Followed)
		pipe.Expire(context.Background(), key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set achievement to cache, err(%s)", err.Error())
	}
}

func (r *achievementRepo) AddAchievementScore(ctx context.Context, uuid string, score int32) error {
	ach := &Achievement{
		Uuid:  uuid,
		Score: score,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"score": gorm.Expr("score + ?", score)}),
	}).Create(ach).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add achievement score: c(%v)", uuid))
	}
	return nil
}

func (r *achievementRepo) AddAchievementScoreToCache(ctx context.Context, uuid string, score int32) error {
	var script = redis.NewScript(`
					local uuid = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						redis.call("HINCRBY", uuid, "score", value)
					end
					return 0
	`)
	keys := []string{uuid}
	values := []interface{}{score}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()

	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add achievement score to cache: uuid(%s), score(%v)", uuid, score))
	}
	return nil
}
