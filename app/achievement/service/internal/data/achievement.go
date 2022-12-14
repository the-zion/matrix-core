package data

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
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

func (r *achievementRepo) SetActiveAgree(ctx context.Context, uuid string) error {
	ac := &Active{
		Uuid:  uuid,
		Agree: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"agree": gorm.Expr("agree + ?", 1)}),
	}).Create(ac).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add active agree: c(%v)", uuid))
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

func (r *achievementRepo) CancelActiveAgree(ctx context.Context, userUuid string) error {
	active := &Active{}
	err := r.data.DB(ctx).Model(active).Where("uuid = ? and agree > 0", userUuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to subtract achievement agree: userUuid(%v)", userUuid))
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

func (r *achievementRepo) SetUserMedalToCache(ctx context.Context, medal, uuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("HSET", key, change, 1)
					end
					return 0
	`)
	keys := []string{"medal_" + uuid}
	values := []interface{}{medal}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user medal to cache: medal(%s), uuid(%s)", medal, uuid))
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

					local followed = KEYS[2]
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

					local followed = KEYS[2]
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

func (r *achievementRepo) CancelUserMedalFromCache(ctx context.Context, medal, uuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
					local change = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("HSET", key, change, 2)
					end
					return 0
	`)
	keys := []string{"medal_" + uuid}
	values := []interface{}{medal}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user medal from cache: medal(%s), uuid(%s)", medal, uuid))
	}
	return nil
}

func (r *achievementRepo) CancelUserMedal(ctx context.Context, medal, uuid string) error {
	m := &Medal{}
	err := r.data.DB(ctx).Model(m).Where(fmt.Sprintf("uuid = ? and %s = ?", medal), uuid, 1).Update(medal, 2).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user medal from db: medal(%s), uuid(%s)", medal, uuid))
	}
	return nil
}

func (r *achievementRepo) GetAchievementList(ctx context.Context, uuids []string) ([]*biz.Achievement, error) {
	exists, unExists, err := r.achievementListExist(ctx, uuids)
	if err != nil {
		return nil, err
	}

	achievementList := make([]*biz.Achievement, 0, cap(exists))

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
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range uuids {
			pipe.Exists(ctx, item)
		}
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, fmt.Sprintf("fail to check if achievement exist from cache: uuids(%v)", uuids))
	}

	exists := make([]string, 0, len(cmd))
	unExists := make([]string, 0, len(cmd))
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
	list := make([]*Achievement, 0, cap(unExists))
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
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setAchievementListToCache(list)
		})()
	}
	return nil
}

func (r *achievementRepo) setAchievementListToCache(achievementList []*Achievement) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range achievementList {
			key := item.Uuid
			pipe.HSetNX(ctx, key, "agree", item.Agree)
			pipe.HSetNX(ctx, key, "collect", item.Collect)
			pipe.HSetNX(ctx, key, "view", item.View)
			pipe.HSetNX(ctx, key, "follow", item.Follow)
			pipe.HSetNX(ctx, key, "followed", item.Followed)
			pipe.Expire(ctx, key, time.Minute*30)
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

	newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
	go r.data.Recover(newCtx, func(ctx context.Context) {
		r.setAchievementToCache(key, achievement)
	})()

	return achievement, nil
}

func (r *achievementRepo) getAchievementFromCache(ctx context.Context, key string) (*biz.Achievement, error) {
	achievement, err := r.data.redisCli.HMGet(ctx, key, "agree", "collect", "view", "follow", "followed", "score").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get achievement form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0, 0, 0}
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
		Score:    val[5],
	}, nil
}

func (r *achievementRepo) getAchievementFromDB(ctx context.Context, uuid string) (*biz.Achievement, error) {
	ach := &Achievement{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).First(ach).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get achievement from db: uuid(%s)", uuid))
	}
	return &biz.Achievement{
		Uuid:     ach.Uuid,
		Agree:    ach.Agree,
		Collect:  ach.Collect,
		View:     ach.View,
		Follow:   ach.Follow,
		Followed: ach.Followed,
		Score:    ach.Score,
	}, nil
}

func (r *achievementRepo) setAchievementToCache(key string, achievement *biz.Achievement) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSetNX(ctx, key, "agree", achievement.Agree)
		pipe.HSetNX(ctx, key, "collect", achievement.Collect)
		pipe.HSetNX(ctx, key, "view", achievement.View)
		pipe.HSetNX(ctx, key, "follow", achievement.Follow)
		pipe.HSetNX(ctx, key, "followed", achievement.Followed)
		pipe.HSetNX(ctx, key, "followed", achievement.Followed)
		pipe.HSetNX(ctx, key, "score", achievement.Score)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set achievement to cache, err(%s)", err.Error())
	}
}

func (r *achievementRepo) GetUserMedal(ctx context.Context, uuid string) (*biz.Medal, error) {
	var medal *biz.Medal
	key := "medal_" + uuid
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		medal, err = r.getMedalFromCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return medal, nil
	}

	medal, err = r.getMedalFromDB(ctx, uuid)
	if err != nil {
		return nil, err
	}

	newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
	go r.data.Recover(newCtx, func(ctx context.Context) {
		r.setMedalToCache(key, medal)
	})()

	return medal, nil
}

func (r *achievementRepo) getMedalFromCache(ctx context.Context, key string) (*biz.Medal, error) {
	medal, err := r.data.redisCli.HMGet(ctx, key,
		"creation1", "creation2", "creation3", "creation4", "creation5", "creation6", "creation7",
		"agree1", "agree2", "agree3", "agree4", "agree5", "agree6",
		"view1", "view2", "view3",
		"comment1", "comment2", "comment3",
		"collect1", "collect2", "collect3").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get medal form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for _index, count := range medal {
		if count == nil {
			break
		}
		num, err := strconv.ParseInt(count.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
		}
		val[_index] = int32(num)
	}
	return &biz.Medal{
		Creation1: val[0],
		Creation2: val[1],
		Creation3: val[2],
		Creation4: val[3],
		Creation5: val[4],
		Creation6: val[5],
		Creation7: val[6],
		Agree1:    val[7],
		Agree2:    val[8],
		Agree3:    val[9],
		Agree4:    val[10],
		Agree5:    val[11],
		Agree6:    val[12],
		View1:     val[13],
		View2:     val[14],
		View3:     val[15],
		Comment1:  val[16],
		Comment2:  val[17],
		Comment3:  val[18],
		Collect1:  val[19],
		Collect2:  val[20],
		Collect3:  val[21],
	}, nil
}

func (r *achievementRepo) getMedalFromDB(ctx context.Context, uuid string) (*biz.Medal, error) {
	medal := &Medal{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).First(medal).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get achievement from db: uuid(%s)", uuid))
	}
	return &biz.Medal{
		Creation1: medal.Creation1,
		Creation2: medal.Creation2,
		Creation3: medal.Creation3,
		Creation4: medal.Creation4,
		Creation5: medal.Creation5,
		Creation6: medal.Creation6,
		Creation7: medal.Creation7,
		Agree1:    medal.Agree1,
		Agree2:    medal.Agree2,
		Agree3:    medal.Agree3,
		Agree4:    medal.Agree4,
		Agree5:    medal.Agree5,
		Agree6:    medal.Agree6,
		View1:     medal.View1,
		View2:     medal.View2,
		View3:     medal.View3,
		Comment1:  medal.Comment1,
		Comment2:  medal.Comment2,
		Comment3:  medal.Comment3,
		Collect1:  medal.Collect1,
		Collect2:  medal.Collect2,
		Collect3:  medal.Collect3,
	}, nil
}

func (r *achievementRepo) setMedalToCache(key string, medal *biz.Medal) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSetNX(ctx, key, "creation1", medal.Creation1)
		pipe.HSetNX(ctx, key, "creation2", medal.Creation2)
		pipe.HSetNX(ctx, key, "creation3", medal.Creation3)
		pipe.HSetNX(ctx, key, "creation4", medal.Creation4)
		pipe.HSetNX(ctx, key, "creation5", medal.Creation5)
		pipe.HSetNX(ctx, key, "creation6", medal.Creation6)
		pipe.HSetNX(ctx, key, "creation7", medal.Creation7)
		pipe.HSetNX(ctx, key, "agree1", medal.Agree1)
		pipe.HSetNX(ctx, key, "agree2", medal.Agree2)
		pipe.HSetNX(ctx, key, "agree3", medal.Agree3)
		pipe.HSetNX(ctx, key, "agree4", medal.Agree4)
		pipe.HSetNX(ctx, key, "agree5", medal.Agree5)
		pipe.HSetNX(ctx, key, "agree6", medal.Agree6)
		pipe.HSetNX(ctx, key, "view1", medal.View1)
		pipe.HSetNX(ctx, key, "view2", medal.View2)
		pipe.HSetNX(ctx, key, "view3", medal.View3)
		pipe.HSetNX(ctx, key, "comment1", medal.Comment1)
		pipe.HSetNX(ctx, key, "comment2", medal.Comment2)
		pipe.HSetNX(ctx, key, "comment3", medal.Comment3)
		pipe.HSetNX(ctx, key, "collect1", medal.Collect1)
		pipe.HSetNX(ctx, key, "collect2", medal.Collect2)
		pipe.HSetNX(ctx, key, "collect3", medal.Collect3)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set medal to cache, err(%s)", err.Error())
	}
}

func (r *achievementRepo) GetUserActive(ctx context.Context, uuid string) (*biz.Active, error) {
	var active *biz.Active
	key := "active_" + uuid
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		active, err = r.getActiveFromCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return active, nil
	}

	active, err = r.getActiveFromDB(ctx, uuid)
	if err != nil {
		return nil, err
	}

	newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
	go r.data.Recover(newCtx, func(ctx context.Context) {
		r.setActiveToCache(key, active)
	})()

	return active, nil
}

func (r *achievementRepo) getActiveFromCache(ctx context.Context, key string) (*biz.Active, error) {
	active, err := r.data.redisCli.HMGet(ctx, key, "agree").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get active form cache: key(%s)", key))
	}
	val := []int32{0}
	for _index, count := range active {
		if count == nil {
			break
		}
		num, err := strconv.ParseInt(count.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
		}
		val[_index] = int32(num)
	}
	return &biz.Active{
		Agree: val[0],
	}, nil
}

func (r *achievementRepo) getActiveFromDB(ctx context.Context, uuid string) (*biz.Active, error) {
	active := &Active{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).First(active).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get active from db: uuid(%s)", uuid))
	}
	return &biz.Active{
		Agree: active.Agree,
	}, nil
}

func (r *achievementRepo) setActiveToCache(key string, active *biz.Active) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSetNX(ctx, key, "agree", active.Agree)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set active to cache, err(%s)", err.Error())
	}
}

func (r *achievementRepo) SendMedalToMq(ctx context.Context, medal, uuid, mode string) error {
	commentMap := &biz.CommentMap{
		Uuid:  uuid,
		Medal: medal,
		Mode:  mode,
	}

	data, err := commentMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "achievement",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.achievementMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send medal to mq: %v", medal))
	}
	return nil
}

func (r *achievementRepo) SetUserMedal(ctx context.Context, medal, uuid string) error {
	m := make(map[string]interface{}, 0)
	m["uuid"] = uuid
	m[medal] = 1

	err := r.data.DB(ctx).Model(&Medal{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{medal: 1}),
	}).Create(m).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user medal: c(%v)", uuid))
	}
	return nil
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
