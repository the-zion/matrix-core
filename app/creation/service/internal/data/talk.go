package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var _ biz.TalkRepo = (*talkRepo)(nil)

type talkRepo struct {
	data *Data
	log  *log.Helper
}

func NewTalkRepo(data *Data, logger log.Logger) biz.TalkRepo {
	return &talkRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "creation/data/talk")),
	}
}

func (r *talkRepo) GetTalk(ctx context.Context, id int32) (*biz.Talk, error) {
	talk := &Talk{}
	err := r.data.db.WithContext(ctx).Select("uuid").Where("talk_id = ?", id).First(talk).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk from db: id(%v)", id))
	}
	return &biz.Talk{
		TalkId: id,
		Uuid:   talk.Uuid,
	}, nil
}

func (r *talkRepo) GetTalkList(ctx context.Context, page int32) ([]*biz.Talk, error) {
	talk, err := r.getTalkFromCache(ctx, page)
	if err != nil {
		return nil, err
	}

	size := len(talk)
	if size != 0 {
		return talk, nil
	}

	talk, err = r.getTalkFromDB(ctx, page)
	if err != nil {
		return nil, err
	}

	size = len(talk)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setTalkToCache("talk", talk)
		})()
	}
	return talk, nil
}

func (r *talkRepo) GetTalkListHot(ctx context.Context, page int32) ([]*biz.TalkStatistic, error) {
	talk, err := r.getTalkHotFromCache(ctx, page)
	if err != nil {
		return nil, err
	}

	size := len(talk)
	if size != 0 {
		return talk, nil
	}

	talk, err = r.GetTalkHotFromDB(ctx, page)
	if err != nil {
		return nil, err
	}

	size = len(talk)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setTalkHotToCache("talk_hot", talk)
		})()
	}
	return talk, nil
}

func (r *talkRepo) GetUserTalkList(ctx context.Context, page int32, uuid string) ([]*biz.Talk, error) {
	talk, err := r.getUserTalkListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(talk)
	if size != 0 {
		return talk, nil
	}

	talk, err = r.getUserTalkListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(talk)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserTalkListToCache("user_talk_list_"+uuid, talk)
		})()
	}
	return talk, nil
}

func (r *talkRepo) getUserTalkListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "user_talk_list_"+uuid, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk list visitor from cache: key(%s), page(%v)", "user_talk_list_"+uuid, page))
	}

	talk := make([]*biz.Talk, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		talk = append(talk, &biz.Talk{
			TalkId: int32(id),
			Uuid:   member[1],
		})
	}
	return talk, nil
}

func (r *talkRepo) getUserTalkListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Talk, 0)
	err := r.data.db.WithContext(ctx).Select("talk_id", "uuid").Where("uuid = ?", uuid).Order("talk_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk from db: page(%v), uuid(%s)", page, uuid))
	}

	talk := make([]*biz.Talk, 0, len(list))
	for _, item := range list {
		talk = append(talk, &biz.Talk{
			TalkId: item.TalkId,
			Uuid:   item.Uuid,
		})
	}
	return talk, nil
}

func (r *talkRepo) GetUserTalkListVisitor(ctx context.Context, page int32, uuid string) ([]*biz.Talk, error) {
	talk, err := r.getUserTalkListVisitorFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(talk)
	if size != 0 {
		return talk, nil
	}

	talk, err = r.getUserTalkListVisitorFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(talk)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserTalkListToCache("user_talk_list_visitor_"+uuid, talk)
		})()
	}
	return talk, nil
}

func (r *talkRepo) getUserTalkListVisitorFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "user_talk_list_visitor_"+uuid, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk list visitor from cache: key(%s), page(%v)", "user_talk_list_visitor_"+uuid, page))
	}

	talk := make([]*biz.Talk, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		talk = append(talk, &biz.Talk{
			TalkId: int32(id),
			Uuid:   member[1],
		})
	}
	return talk, nil
}

func (r *talkRepo) getUserTalkListVisitorFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Talk, 0)
	err := r.data.db.WithContext(ctx).Select("talk_id", "uuid").Where("uuid = ? and auth = ?", uuid, 1).Order("talk_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk visitor from db: page(%v), uuid(%s)", page, uuid))
	}

	talk := make([]*biz.Talk, 0, len(list))
	for _, item := range list {
		talk = append(talk, &biz.Talk{
			TalkId: item.TalkId,
			Uuid:   item.Uuid,
		})
	}
	return talk, nil
}

func (r *talkRepo) setUserTalkListToCache(key string, talk []*biz.Talk) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(talk))
		for _, item := range talk {
			z = append(z, &redis.Z{
				Score:  float64(item.TalkId),
				Member: strconv.Itoa(int(item.TalkId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk to cache: talk(%v), err(%v)", talk, err)
	}
}

func (r *talkRepo) GetTalkCount(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Talk{}).Where("uuid = ?", uuid).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get talk count from db: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *talkRepo) GetTalkCountVisitor(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Talk{}).Where("uuid = ? and auth = ?", uuid, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get talk count from db: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *talkRepo) setTalkHotToCache(key string, talk []*biz.TalkStatistic) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(talk))
		for _, item := range talk {
			z = append(z, &redis.Z{
				Score:  float64(item.Agree),
				Member: strconv.Itoa(int(item.TalkId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk to cache: talk(%v), err(%v)", talk, err)
	}
}

func (r *talkRepo) getTalkHotFromCache(ctx context.Context, page int32) ([]*biz.TalkStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "talk_hot", index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk hot from cache: key(%s), page(%v)", "talk_hot", page))
	}

	talk := make([]*biz.TalkStatistic, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		talk = append(talk, &biz.TalkStatistic{
			TalkId: int32(id),
			Uuid:   member[1],
		})
	}
	return talk, nil
}

func (r *talkRepo) GetTalkHotFromDB(ctx context.Context, page int32) ([]*biz.TalkStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*TalkStatistic, 0)
	err := r.data.db.WithContext(ctx).Select("talk_id", "uuid", "agree").Where("auth", 1).Order("agree desc, talk_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk statistic from db: page(%v)", page))
	}

	talk := make([]*biz.TalkStatistic, 0, len(list))
	for _, item := range list {
		talk = append(talk, &biz.TalkStatistic{
			TalkId: item.TalkId,
			Uuid:   item.Uuid,
			Agree:  item.Agree,
		})
	}
	return talk, nil
}

func (r *talkRepo) getTalkFromCache(ctx context.Context, page int32) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "talk", index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk from cache: key(%s), page(%v)", "talk", page))
	}

	talk := make([]*biz.Talk, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		talk = append(talk, &biz.Talk{
			TalkId: int32(id),
			Uuid:   member[1],
		})
	}
	return talk, nil
}

func (r *talkRepo) getTalkFromDB(ctx context.Context, page int32) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Talk, 0)
	err := r.data.db.WithContext(ctx).Select("talk_id", "uuid").Where("auth", 1).Order("talk_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk from db: page(%v)", page))
	}

	talk := make([]*biz.Talk, 0, len(list))
	for _, item := range list {
		talk = append(talk, &biz.Talk{
			TalkId: item.TalkId,
			Uuid:   item.Uuid,
		})
	}
	return talk, nil
}

func (r *talkRepo) setTalkToCache(key string, talk []*biz.Talk) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(talk))
		for _, item := range talk {
			z = append(z, &redis.Z{
				Score:  float64(item.TalkId),
				Member: strconv.Itoa(int(item.TalkId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk to cache: talk(%v), err(%v)", talk, err)
	}
}

func (r *talkRepo) GetTalkListStatistic(ctx context.Context, ids []int32) ([]*biz.TalkStatistic, error) {
	exists, unExists, err := r.talkListStatisticExist(ctx, ids)
	if err != nil {
		return nil, err
	}

	talkListStatistic := make([]*biz.TalkStatistic, 0, cap(exists))
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(exists) == 0 {
			return nil
		}
		return r.getTalkListStatisticFromCache(ctx, exists, &talkListStatistic)
	}))
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(unExists) == 0 {
			return nil
		}
		return r.getTalkListStatisticFromDb(ctx, unExists, &talkListStatistic)
	}))

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	return talkListStatistic, nil
}

func (r *talkRepo) talkListStatisticExist(ctx context.Context, ids []int32) ([]int32, []int32, error) {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range ids {
			pipe.Exists(ctx, "talk_"+strconv.Itoa(int(item)))
		}
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, fmt.Sprintf("fail to check if talk statistic exist from cache: ids(%v)", ids))
	}

	exists := make([]int32, 0, len(cmd))
	unExists := make([]int32, 0, len(cmd))
	for index, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		if exist == 1 {
			exists = append(exists, ids[index])
		} else {
			unExists = append(unExists, ids[index])
		}
	}
	return exists, unExists, nil
}

func (r *talkRepo) getTalkListStatisticFromCache(ctx context.Context, exists []int32, talkListStatistic *[]*biz.TalkStatistic) error {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, id := range exists {
			pipe.HMGet(ctx, "talk_"+strconv.Itoa(int(id)), "agree", "collect", "view", "comment")
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get talk list statistic from cache: ids(%v)", exists))
	}

	for index, item := range cmd {
		val := []int32{0, 0, 0, 0}
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
		*talkListStatistic = append(*talkListStatistic, &biz.TalkStatistic{
			TalkId:  exists[index],
			Agree:   val[0],
			Collect: val[1],
			View:    val[2],
			Comment: val[3],
		})
	}
	return nil
}

func (r *talkRepo) getTalkListStatisticFromDb(ctx context.Context, unExists []int32, talkListStatistic *[]*biz.TalkStatistic) error {
	list := make([]*TalkStatistic, 0, cap(unExists))
	err := r.data.db.WithContext(ctx).Select("talk_id", "uuid", "agree", "comment", "collect", "view", "auth").Where("talk_id IN ?", unExists).Find(&list).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get talk statistic list from db: ids(%v)", unExists))
	}

	for _, item := range list {
		*talkListStatistic = append(*talkListStatistic, &biz.TalkStatistic{
			TalkId:  item.TalkId,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Comment: item.Comment,
			Collect: item.Collect,
			View:    item.View,
			Auth:    item.Auth,
		})
	}

	if len(list) != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setTalkListStatisticToCache(list)
		})()
	}

	return nil
}

func (r *talkRepo) setTalkListStatisticToCache(commentList []*TalkStatistic) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range commentList {
			key := "talk_" + strconv.Itoa(int(item.TalkId))
			pipe.HSetNX(ctx, key, "uuid", item.Uuid)
			pipe.HSetNX(ctx, key, "agree", item.Agree)
			pipe.HSetNX(ctx, key, "comment", item.Comment)
			pipe.HSetNX(ctx, key, "collect", item.Collect)
			pipe.HSetNX(ctx, key, "view", item.View)
			pipe.HSetNX(ctx, key, "auth", item.Auth)
			pipe.Expire(ctx, key, time.Minute*30)
		}
		return nil
	})

	if err != nil {
		r.log.Errorf("fail to set talk statistic to cache, err(%s)", err.Error())
	}
}

func (r *talkRepo) GetTalkStatistic(ctx context.Context, id int32, uuid string) (*biz.TalkStatistic, error) {
	var statistic *biz.TalkStatistic
	key := "talk_" + strconv.Itoa(int(id))
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		statistic, err = r.getTalkStatisticFromCache(ctx, key, uuid)
		if err != nil {
			return nil, err
		}
		return statistic, nil
	}

	statistic, err = r.getTalkStatisticFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
	go r.data.Recover(newCtx, func(ctx context.Context) {
		r.setTalkStatisticToCache(key, statistic)
	})()

	if statistic.Auth == 2 && statistic.Uuid != uuid {
		return nil, errors.Errorf("fail to get talk statistic from cache: no auth")
	}

	return statistic, nil
}

func (r *talkRepo) getTalkStatisticFromCache(ctx context.Context, key, uuid string) (*biz.TalkStatistic, error) {
	statistic, err := r.data.redisCli.HMGet(ctx, key, "uuid", "agree", "collect", "view", "comment", "auth").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk statistic form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0, 0}
	for _index, count := range statistic[1:] {
		if count == nil {
			break
		}
		num, err := strconv.ParseInt(count.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
		}
		val[_index] = int32(num)
	}
	if val[4] == 2 && statistic[0].(string) != uuid {
		return nil, errors.Errorf("fail to get talk statistic from cache: no auth")
	}
	return &biz.TalkStatistic{
		Uuid:    statistic[0].(string),
		Agree:   val[0],
		Collect: val[1],
		View:    val[2],
		Comment: val[3],
	}, nil
}

func (r *talkRepo) getTalkStatisticFromDB(ctx context.Context, id int32) (*biz.TalkStatistic, error) {
	as := &TalkStatistic{}
	err := r.data.db.WithContext(ctx).Select("uuid", "agree", "collect", "view", "comment").Where("talk_id = ?", id).First(as).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get statistic from db: id(%v)", id))
	}
	return &biz.TalkStatistic{
		Uuid:    as.Uuid,
		Agree:   as.Agree,
		Collect: as.Collect,
		View:    as.View,
		Comment: as.Comment,
	}, nil
}

func (r *talkRepo) setTalkStatisticToCache(key string, statistic *biz.TalkStatistic) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(context.Background(), key, "uuid", statistic.Uuid, "agree", statistic.Agree, "collect", statistic.Collect, "view", statistic.View, "comment", statistic.Comment)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk statistic to cache, err(%s)", err.Error())
	}
}

func (r *talkRepo) GetLastTalkDraft(ctx context.Context, uuid string) (*biz.TalkDraft, error) {
	draft := &TalkDraft{}
	err := r.data.db.WithContext(ctx).Select("id", "status").Where("uuid = ? and status = ?", uuid, 1).Last(draft).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("talk draft not found from db", fmt.Sprintf("uuid(%s)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get last talk draft: uuid(%s)", uuid))
	}
	return &biz.TalkDraft{
		Id:     int32(draft.ID),
		Status: draft.Status,
	}, nil
}

func (r *talkRepo) GetTalkSearch(ctx context.Context, page int32, search, time string) ([]*biz.TalkSearch, int32, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	reply := make([]*biz.TalkSearch, 0)

	var buf bytes.Buffer
	var body map[string]interface{}
	if search != "" {
		body = map[string]interface{}{
			"from":    index * 10,
			"size":    10,
			"_source": []string{"update", "tags", "cover", "uuid"},
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{"multi_match": map[string]interface{}{
							"query":  search,
							"fields": []string{"text", "tags", "title"},
						}},
					},
					"filter": map[string]interface{}{
						"range": map[string]interface{}{
							"update": map[string]interface{}{},
						},
					},
				},
			},
			"highlight": map[string]interface{}{
				"fields": map[string]interface{}{
					"text": map[string]interface{}{
						"fragment_size":       300,
						"number_of_fragments": 1,
						"no_match_size":       300,
						"pre_tags":            "<span style='color:red'>",
						"post_tags":           "</span>",
					},
					"title": map[string]interface{}{
						"pre_tags":      "<span style='color:red'>",
						"post_tags":     "</span>",
						"no_match_size": 100,
					},
				},
			},
		}
	} else {
		body = map[string]interface{}{
			"from":    index * 10,
			"size":    10,
			"_source": []string{"update", "tags", "cover", "uuid"},
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{"match_all": map[string]interface{}{}},
					},
					"filter": map[string]interface{}{
						"range": map[string]interface{}{
							"update": map[string]interface{}{},
						},
					},
				},
			},
			"highlight": map[string]interface{}{
				"fields": map[string]interface{}{
					"text": map[string]interface{}{
						"fragment_size":       300,
						"number_of_fragments": 1,
						"no_match_size":       300,
						"pre_tags":            "<span style='color:red'>",
						"post_tags":           "</span>",
					},
					"title": map[string]interface{}{
						"pre_tags":      "<span style='color:red'>",
						"post_tags":     "</span>",
						"no_match_size": 100,
					},
				},
			},
			"sort": []map[string]interface{}{
				{"_id": "desc"},
			},
		}
	}

	switch time {
	case "1day":
		body["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].(map[string]interface{})["range"].(map[string]interface{})["update"].(map[string]interface{})["gte"] = "now-1d"
		break
	case "1week":
		body["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].(map[string]interface{})["range"].(map[string]interface{})["update"].(map[string]interface{})["gte"] = "now-1w"
		break
	case "1month":
		body["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].(map[string]interface{})["range"].(map[string]interface{})["update"].(map[string]interface{})["gte"] = "now-1M"
		break
	case "1year":
		body["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].(map[string]interface{})["range"].(map[string]interface{})["update"].(map[string]interface{})["gte"] = "now-1y"
		break
	}
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, 0, errors.Wrapf(err, fmt.Sprintf("error encoding query: page(%v), search(%s), time(%s)", page, search, time))
	}

	res, err := r.data.elasticSearch.es.Search(
		r.data.elasticSearch.es.Search.WithContext(ctx),
		r.data.elasticSearch.es.Search.WithIndex("talk"),
		r.data.elasticSearch.es.Search.WithBody(&buf),
		r.data.elasticSearch.es.Search.WithTrackTotalHits(true),
		r.data.elasticSearch.es.Search.WithPretty(),
	)
	if err != nil {
		return nil, 0, errors.Wrapf(err, fmt.Sprintf("error getting response from es: page(%v), search(%s), time(%s)", page, search, time))
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, 0, errors.Wrapf(err, fmt.Sprintf("error parsing the response body: page(%v), search(%s), time(%s)", page, search, time))
		} else {
			return nil, 0, errors.Errorf(fmt.Sprintf("error search talk from  es: reason(%v), page(%v), search(%s), time(%s)", e, page, search, time))
		}
	}

	result := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, errors.Wrapf(err, fmt.Sprintf("error parsing the response body: page(%v), search(%s), time(%s)", page, search, time))
	}

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		id, err := strconv.ParseInt(hit.(map[string]interface{})["_id"].(string), 10, 32)
		if err != nil {
			return nil, 0, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: page(%v), search(%s), time(%s)", page, search, time))
		}

		talk := &biz.TalkSearch{
			Id:     int32(id),
			Tags:   hit.(map[string]interface{})["_source"].(map[string]interface{})["tags"].(string),
			Update: hit.(map[string]interface{})["_source"].(map[string]interface{})["update"].(string),
			Cover:  hit.(map[string]interface{})["_source"].(map[string]interface{})["cover"].(string),
			Uuid:   hit.(map[string]interface{})["_source"].(map[string]interface{})["uuid"].(string),
		}

		if text, ok := hit.(map[string]interface{})["highlight"].(map[string]interface{})["text"]; ok {
			talk.Text = text.([]interface{})[0].(string)
		}

		if title, ok := hit.(map[string]interface{})["highlight"].(map[string]interface{})["title"]; ok {
			talk.Title = title.([]interface{})[0].(string)
		}

		reply = append(reply, talk)
	}
	res.Body.Close()
	return reply, int32(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)), nil
}

func (r *talkRepo) GetUserTalkAgree(ctx context.Context, uuid string) (map[int32]bool, error) {
	exist, err := r.userTalkAgreeExist(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if exist == 1 {
		return r.getUserTalkAgreeFromCache(ctx, uuid)
	} else {
		return r.getUserTalkAgreeFromDb(ctx, uuid)
	}
}

func (r *talkRepo) userTalkAgreeExist(ctx context.Context, uuid string) (int32, error) {
	exist, err := r.data.redisCli.Exists(ctx, "user_talk_agree_"+uuid).Result()
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to check if user talk agree exist from cache: uuid(%s)", uuid))
	}
	return int32(exist), nil
}

func (r *talkRepo) getUserTalkAgreeFromCache(ctx context.Context, uuid string) (map[int32]bool, error) {
	key := "user_talk_agree_" + uuid
	agreeSet, err := r.data.redisCli.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk agree from cache: uuid(%s)", uuid))
	}

	agreeMap := make(map[int32]bool, 0)
	for _, item := range agreeSet {
		id, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s), uuid(%s), key(%s)", id, uuid, key))
		}
		agreeMap[int32(id)] = true
	}
	return agreeMap, nil
}

func (r *talkRepo) getUserTalkAgreeFromDb(ctx context.Context, uuid string) (map[int32]bool, error) {
	list := make([]*TalkAgree, 0)
	err := r.data.db.WithContext(ctx).Select("talk_id").Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk agree from db: uuid(%s)", uuid))
	}

	agreeMap := make(map[int32]bool, 0)
	for _, item := range list {
		agreeMap[item.TalkId] = true
	}
	if len(list) != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserTalkAgreeToCache(uuid, list)
		})()
	}
	return agreeMap, nil
}

func (r *talkRepo) setUserTalkAgreeToCache(uuid string, agreeList []*TalkAgree) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0, len(agreeList))
		key := "user_talk_agree_" + uuid
		for _, item := range agreeList {
			set = append(set, item.TalkId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user talk agree to cache: uuid(%s), agreeList(%v), error(%s) ", uuid, agreeList, err.Error())
	}
}

func (r *talkRepo) GetUserTalkCollect(ctx context.Context, uuid string) (map[int32]bool, error) {
	exist, err := r.userTalkCollectExist(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if exist == 1 {
		return r.getUserTalkCollectFromCache(ctx, uuid)
	} else {
		return r.getUserTalkCollectFromDb(ctx, uuid)
	}
}

func (r *talkRepo) userTalkCollectExist(ctx context.Context, uuid string) (int32, error) {
	exist, err := r.data.redisCli.Exists(ctx, "user_talk_collect_"+uuid).Result()
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to check if user talk collect exist from cache: uuid(%s)", uuid))
	}
	return int32(exist), nil
}

func (r *talkRepo) getUserTalkCollectFromCache(ctx context.Context, uuid string) (map[int32]bool, error) {
	key := "user_talk_collect_" + uuid
	collectSet, err := r.data.redisCli.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk collect from cache: uuid(%s)", uuid))
	}

	collectMap := make(map[int32]bool, 0)
	for _, item := range collectSet {
		id, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s), uuid(%s), key(%s)", id, uuid, key))
		}
		collectMap[int32(id)] = true
	}
	return collectMap, nil
}

func (r *talkRepo) getUserTalkCollectFromDb(ctx context.Context, uuid string) (map[int32]bool, error) {
	list := make([]*TalkCollect, 0)
	err := r.data.db.WithContext(ctx).Select("talk_id").Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk collect from db: uuid(%s)", uuid))
	}

	collectMap := make(map[int32]bool, 0)
	for _, item := range list {
		collectMap[item.TalkId] = true
	}
	if len(list) != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserTalkCollectToCache(uuid, list)
		})()
	}
	return collectMap, nil
}

func (r *talkRepo) setUserTalkCollectToCache(uuid string, collectList []*TalkCollect) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0, len(collectList))
		key := "user_talk_collect_" + uuid
		for _, item := range collectList {
			set = append(set, item.TalkId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user talk collect to cache: uuid(%s), collectList(%v), error(%s) ", uuid, collectList, err.Error())
	}
}

func (r *talkRepo) GetTalkAuth(ctx context.Context, id int32) (int32, error) {
	talk := &Talk{}
	err := r.data.db.WithContext(ctx).Select("auth").Where("talk_id = ?", id).First(talk).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get talk auth from db: id(%v)", id))
	}
	return talk.Auth, nil
}

func (r *talkRepo) GetCollectionsIdFromTalkCollect(ctx context.Context, id int32) (int32, error) {
	collect := &Collect{}
	err := r.data.db.WithContext(ctx).Select("collections_id").Where("creations_id = ? and mode = ?", id, 3).First(collect).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collections id  from db: creationsId(%v)", id))
	}
	return collect.CollectionsId, nil
}

func (r *talkRepo) GetTalkImageReview(ctx context.Context, page int32, uuid string) ([]*biz.ImageReview, error) {
	key := "talk_image_irregular_" + uuid
	review, err := r.getTalkImageReviewFromCache(ctx, page, key)
	if err != nil {
		return nil, err
	}

	size := len(review)
	if size != 0 {
		return review, nil
	}

	review, err = r.getTalkImageReviewFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(review)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setTalkImageReviewToCache(key, review)
		})()
	}
	return review, nil
}

func (r *talkRepo) getTalkImageReviewFromCache(ctx context.Context, page int32, key string) ([]*biz.ImageReview, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.LRange(ctx, key, index*20, index*20+19).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk image irregular list from cache: key(%s), page(%v)", key, page))
	}

	review := make([]*biz.ImageReview, 0, len(list))
	for _index, item := range list {
		var imageReview = &biz.ImageReview{}
		err = imageReview.UnmarshalJSON([]byte(item))
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("json unmarshal error: imageReview(%v)", item))
		}
		review = append(review, &biz.ImageReview{
			Id:         int32(_index+1) + (page-1)*20,
			CreationId: imageReview.CreationId,
			Kind:       imageReview.Kind,
			Uid:        imageReview.Uid,
			Uuid:       imageReview.Uuid,
			CreateAt:   imageReview.CreateAt,
			JobId:      imageReview.JobId,
			Url:        imageReview.Url,
			Label:      imageReview.Label,
			Result:     imageReview.Result,
			Category:   imageReview.Category,
			SubLabel:   imageReview.SubLabel,
			Score:      imageReview.Score,
		})
	}
	return review, nil
}

func (r *talkRepo) getTalkImageReviewFromDB(ctx context.Context, page int32, uuid string) ([]*biz.ImageReview, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*TalkReview, 0)
	err := r.data.db.WithContext(ctx).Select("talk_id", "kind", "uid", "uuid", "job_id", "created_at", "url", "label", "result", "category", "sub_label", "score").Where("uuid", uuid).Order("id desc").Offset(index * 20).Limit(20).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk image review from db: page(%v), uuid(%s)", page, uuid))
	}

	review := make([]*biz.ImageReview, 0, len(list))
	for _index, item := range list {
		review = append(review, &biz.ImageReview{
			Id:         int32(_index+1) + (page-1)*20,
			CreationId: item.TalkId,
			Kind:       item.Kind,
			Uid:        item.Uid,
			Uuid:       item.Uuid,
			JobId:      item.JobId,
			CreateAt:   item.CreatedAt.Format("2006-01-02"),
			Url:        item.Url,
			Label:      item.Label,
			Result:     item.Result,
			Category:   item.Category,
			SubLabel:   item.SubLabel,
			Score:      item.Score,
		})
	}
	return review, nil
}

func (r *talkRepo) setTalkImageReviewToCache(key string, review []*biz.ImageReview) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		list := make([]interface{}, 0, len(review))
		for _, item := range review {
			m, err := item.MarshalJSON()
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("fail to marshal avatar review: imageReview(%v)", review))
			}
			list = append(list, m)
		}
		pipe.RPush(ctx, key, list...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk image review to cache: imageReview(%v), err(%v)", review, err)
	}
}

func (r *talkRepo) GetTalkContentReview(ctx context.Context, page int32, uuid string) ([]*biz.TextReview, error) {
	key := "talk_content_irregular_" + uuid
	review, err := r.getTalkContentReviewFromCache(ctx, page, key)
	if err != nil {
		return nil, err
	}

	size := len(review)
	if size != 0 {
		return review, nil
	}

	review, err = r.getTalkContentReviewFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(review)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setTalkContentReviewToCache(key, review)
		})()
	}
	return review, nil
}

func (r *talkRepo) getTalkContentReviewFromCache(ctx context.Context, page int32, key string) ([]*biz.TextReview, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.LRange(ctx, key, index*20, index*20+19).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk content irregular list from cache: key(%s), page(%v)", key, page))
	}

	review := make([]*biz.TextReview, 0, len(list))
	for _index, item := range list {
		var textReview = &biz.TextReview{}
		err = textReview.UnmarshalJSON([]byte(item))
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("json unmarshal error: contentReview(%v)", item))
		}
		review = append(review, &biz.TextReview{
			Id:         int32(_index+1) + (page-1)*20,
			CreationId: textReview.CreationId,
			Title:      textReview.Title,
			Kind:       textReview.Kind,
			Uuid:       textReview.Uuid,
			CreateAt:   textReview.CreateAt,
			JobId:      textReview.JobId,
			Label:      textReview.Label,
			Result:     textReview.Result,
			Section:    textReview.Section,
		})
	}
	return review, nil
}

func (r *talkRepo) getTalkContentReviewFromDB(ctx context.Context, page int32, uuid string) ([]*biz.TextReview, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*TalkContentReview, 0)
	err := r.data.db.WithContext(ctx).Select("talk_id", "kind", "title", "uuid", "job_id", "created_at", "label", "result", "section").Where("uuid", uuid).Order("id desc").Offset(index * 20).Limit(20).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk content review from db: page(%v), uuid(%s)", page, uuid))
	}

	review := make([]*biz.TextReview, 0, len(list))
	for _index, item := range list {
		review = append(review, &biz.TextReview{
			Id:         int32(_index+1) + (page-1)*20,
			CreationId: item.TalkId,
			Kind:       item.Kind,
			Title:      item.Title,
			Uuid:       item.Uuid,
			JobId:      item.JobId,
			CreateAt:   item.CreatedAt.Format("2006-01-02"),
			Label:      item.Label,
			Result:     item.Result,
			Section:    item.Section,
		})
	}
	return review, nil
}

func (r *talkRepo) setTalkContentReviewToCache(key string, review []*biz.TextReview) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		list := make([]interface{}, 0, len(review))
		for _, item := range review {
			m, err := item.MarshalJSON()
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("fail to marshal avatar review: contentReview(%v)", review))
			}
			list = append(list, m)
		}
		pipe.RPush(ctx, key, list...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk content review to cache: contentReview(%v), err(%v)", review, err)
	}
}

func (r *talkRepo) CreateTalkDraft(ctx context.Context, uuid string) (int32, error) {
	draft := &TalkDraft{
		Uuid: uuid,
	}
	err := r.data.DB(ctx).Select("Uuid").Create(draft).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to create a talk draft: uuid(%s)", uuid))
	}
	return int32(draft.ID), nil
}

func (r *talkRepo) CreateTalkFolder(ctx context.Context, id int32, uuid string) error {
	name := "talk/" + uuid + "/" + strconv.Itoa(int(id)) + "/"
	_, err := r.data.cosCli.Object.Put(ctx, name, strings.NewReader(""), nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a talk folder: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) SendTalk(ctx context.Context, id int32, uuid string) (*biz.TalkDraft, error) {
	td := &TalkDraft{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&TalkDraft{}).Where("id = ? and uuid = ? and status = ?", id, uuid, 1).Updates(td).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 2: uuid(%s), id(%v)", uuid, id))
	}
	return &biz.TalkDraft{
		Uuid: uuid,
		Id:   id,
	}, nil
}

func (r *talkRepo) SendReviewToMq(ctx context.Context, review *biz.TalkReview) error {
	data, err := review.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send review to mq: %v", err))
	}
	return nil
}

func (r *talkRepo) DeleteTalkDraft(ctx context.Context, id int32, uuid string) error {
	td := &TalkDraft{}
	td.ID = uint(id)
	err := r.data.DB(ctx).Where("uuid = ?", uuid).Delete(td).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete talk draft: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) CreateTalk(ctx context.Context, id, auth int32, uuid string) error {
	talk := &Talk{
		TalkId: id,
		Uuid:   uuid,
		Auth:   auth,
	}
	err := r.data.DB(ctx).Select("TalkId", "Uuid", "Auth").Create(talk).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a talk: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) DeleteTalk(ctx context.Context, id int32, uuid string) error {
	talk := &Talk{}
	err := r.data.DB(ctx).Where("talk_id = ? and uuid = ?", id, uuid).Delete(talk).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete a talk: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *talkRepo) DeleteTalkStatistic(ctx context.Context, id int32, uuid string) error {
	statistic := &TalkStatistic{}
	err := r.data.DB(ctx).Where("talk_id = ? and uuid = ?", id, uuid).Delete(statistic).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete a talk statistic: uuid(%s)", uuid))
	}
	return nil
}

func (r *talkRepo) DeleteTalkCache(ctx context.Context, id, auth int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	keys := []string{"talk", "talk_hot", "leaderboard", "talk_" + ids, "talk_collect_" + ids, "user_talk_list_" + uuid, "user_talk_list_visitor_" + uuid, "creation_user_" + uuid, "creation_user_visitor_" + uuid}
	values := []interface{}{ids + "%" + uuid, ids + "%" + uuid + "%talk", auth}
	_, err := r.data.redisCli.EvalSha(ctx, "723a24fff880154f9a85292efb2152330b72d994", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete talk cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *talkRepo) FreezeTalkCos(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	key := "talk/" + uuid + "/" + ids + "/content"
	_, err := r.data.cosCli.Object.Delete(ctx, key)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to freeze talk: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *talkRepo) AddCreationUserTalk(ctx context.Context, uuid string, auth int32) error {
	cu := &CreationUser{
		Uuid: uuid,
		Talk: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"talk": gorm.Expr("talk + ?", 1)}),
	}).Create(cu).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation talk: uuid(%v)", uuid))
	}

	if auth == 2 {
		return nil
	}

	cuv := &CreationUserVisitor{
		Uuid: uuid,
		Talk: 1,
	}
	err = r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"talk": gorm.Expr("talk + ?", 1)}),
	}).Create(cuv).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation talk visitor: uuid(%v)", uuid))
	}

	return nil
}

func (r *talkRepo) CreateTalkStatistic(ctx context.Context, id, auth int32, uuid string) error {
	ts := &TalkStatistic{
		TalkId: id,
		Uuid:   uuid,
		Auth:   auth,
	}
	err := r.data.DB(ctx).Select("TalkId", "Uuid", "Auth").Create(ts).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a talk statistic: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) SendTalkToMq(ctx context.Context, talk *biz.Talk, mode string) error {
	talkMap := &biz.SendTalkMap{
		Uuid: talk.Uuid,
		Id:   talk.TalkId,
		Auth: talk.Auth,
		Mode: mode,
	}
	data, err := talkMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{talk.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk to mq: %v", talk))
	}
	return nil
}

func (r *talkRepo) CreateTalkCache(ctx context.Context, id, auth int32, uuid, mode string) error {
	ids := strconv.Itoa(int(id))
	talkStatistic := "talk_" + ids
	talk := "talk"
	talkHot := "talk_hot"
	leaderboard := "leaderboard"
	userTalkList := "user_talk_list_" + uuid
	userTalkListVisitor := "user_talk_list_visitor_" + uuid
	creationUser := "creation_user_" + uuid
	creationUserVisitor := "creation_user_visitor_" + uuid
	keys := []string{talkStatistic, talk, talkHot, leaderboard, userTalkList, userTalkListVisitor, creationUser, creationUserVisitor}
	values := []interface{}{uuid, auth, id, ids + "%" + uuid, mode, ids + "%" + uuid + "%talk"}
	_, err := r.data.redisCli.EvalSha(ctx, "eb5d94a2a972340b5ac93769524a51911f9c032f", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create(update) talk cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) UpdateTalkCache(ctx context.Context, id, auth int32, uuid, mode string) error {
	return r.CreateTalkCache(ctx, id, auth, uuid, mode)
}

func (r *talkRepo) EditTalkCos(ctx context.Context, id int32, uuid string) error {
	err := r.EditTalkCosContent(ctx, id, uuid)
	if err != nil {
		return err
	}

	err = r.EditTalkCosIntroduce(ctx, id, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *talkRepo) SetTalkAgree(ctx context.Context, id int32, uuid string) error {
	ts := TalkStatistic{}
	err := r.data.DB(ctx).Model(&ts).Where("talk_id = ? and uuid = ?", id, uuid).Update("agree", gorm.Expr("agree + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk agree: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) SetUserTalkAgree(ctx context.Context, id int32, userUuid string) error {
	ta := &TalkAgree{
		TalkId: id,
		Uuid:   userUuid,
		Status: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(ta).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add user talk agree: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) SetTalkAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	hotKey := fmt.Sprintf("talk_hot")
	statisticKey := fmt.Sprintf("talk_%v", id)
	boardKey := fmt.Sprintf("leaderboard")
	userKey := fmt.Sprintf("user_talk_agree_%s", userUuid)
	keys := []string{hotKey, statisticKey, boardKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), fmt.Sprintf("%v%s%s%s", id, "%", uuid, "%talk"), id}
	_, err := r.data.redisCli.EvalSha(ctx, "7c94d6c820881b0b69ec949e5565ea0557be65ae", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add user talk agree to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid))
	}
	return nil
}

func (r *talkRepo) CancelTalkAgree(ctx context.Context, id int32, uuid string) error {
	ts := TalkStatistic{}
	err := r.data.DB(ctx).Model(&ts).Where("talk_id = ? and uuid = ?", id, uuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel talk agree: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) CancelUserTalkAgree(ctx context.Context, id int32, userUuid string) error {
	ta := TalkAgree{}
	err := r.data.DB(ctx).Model(&ta).Where("talk_id = ? and uuid = ?", id, userUuid).Update("status", 2).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user talk agree: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) CancelTalkAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	hotKey := "talk_hot"
	boardKey := "leaderboard"
	statisticKey := fmt.Sprintf("talk_%v", id)
	userKey := fmt.Sprintf("user_talk_agree_%s", userUuid)

	keys := []string{hotKey, boardKey, statisticKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), fmt.Sprintf("%v%s%s%s", id, "%", uuid, "%talk"), id}
	_, err := r.data.redisCli.EvalSha(ctx, "c2c68100b9b682b94faff1a9b182159921d0b8e4", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel talk agree from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid))
	}
	return nil
}

func (r *talkRepo) SetTalkView(ctx context.Context, id int32, uuid string) error {
	ts := TalkStatistic{}
	err := r.data.DB(ctx).Model(&ts).Where("talk_id = ? and uuid = ?", id, uuid).Update("view", gorm.Expr("view + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk view: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) SetTalkViewToCache(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	keys := []string{"talk_" + ids}
	values := []interface{}{1}
	_, err := r.data.redisCli.EvalSha(ctx, "6abf0f5d9afbfe7ed552297f59b9d898c6b06747", keys, values...).Result()

	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk view to cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *talkRepo) GetTalkAgreeJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	ids := strconv.Itoa(int(id))
	judge, err := r.data.redisCli.SIsMember(ctx, "talk_agree_"+ids, uuid).Result()
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("fail to judge talk agree member: id(%v), uuid(%s)", id, uuid))
	}
	return judge, nil
}

func (r *talkRepo) GetTalkCollectJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	ids := strconv.Itoa(int(id))
	judge, err := r.data.redisCli.SIsMember(ctx, "talk_collect_"+ids, uuid).Result()
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("fail to judge talk collect member: id(%v), uuid(%s)", id, uuid))
	}
	return judge, nil
}

func (r *talkRepo) EditTalkCosContent(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	name := "talk/" + uuid + "/" + ids + "/content-edit"
	dest := "talk/" + uuid + "/" + ids + "/content"
	sourceURL := fmt.Sprintf("%s/%s", r.data.cosCli.BaseURL.BucketURL.Host, name)
	_, _, err := r.data.cosCli.Object.Copy(ctx, dest, sourceURL, nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to copy talk from edit to content: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) EditTalkCosIntroduce(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	name := "talk/" + uuid + "/" + ids + "/introduce-edit"
	dest := "talk/" + uuid + "/" + ids + "/introduce"
	sourceURL := fmt.Sprintf("%s/%s", r.data.cosCli.BaseURL.BucketURL.Host, name)
	_, _, err := r.data.cosCli.Object.Copy(ctx, dest, sourceURL, nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to copy talk from edit to content: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) SendTalkStatisticToMq(ctx context.Context, uuid, userUuid, mode string) error {
	achievement := &biz.SendTalkStatisticMap{
		Uuid:     uuid,
		UserUuid: userUuid,
		Mode:     mode,
	}
	data, err := achievement.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk statistic to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *talkRepo) SendScoreToMq(ctx context.Context, score int32, uuid, mode string) error {
	scoreMap := &biz.SendScoreMap{
		Uuid:  uuid,
		Score: score,
		Mode:  mode,
	}
	data, err := scoreMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send score to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *talkRepo) SendStatisticToMq(ctx context.Context, id, collectionsId int32, uuid, userUuid, mode string) error {
	statisticMap := &biz.SendStatisticMap{
		Id:            id,
		CollectionsId: collectionsId,
		Uuid:          uuid,
		UserUuid:      userUuid,
		Mode:          mode,
	}
	data, err := statisticMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk statistic to mq: map(%v)", statisticMap))
	}
	return nil
}

func (r *talkRepo) SendTalkImageIrregularToMq(ctx context.Context, review *biz.ImageReview) error {
	data, err := review.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk image review to mq: %v", err))
	}
	return nil
}

func (r *talkRepo) SendTalkContentIrregularToMq(ctx context.Context, review *biz.TextReview) error {
	data, err := review.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk content review to mq: %v", err))
	}
	return nil
}

func (r *talkRepo) SetTalkUserCollect(ctx context.Context, id, collectionsId int32, userUuid string) error {
	collect := &Collect{
		CollectionsId: collectionsId,
		Uuid:          userUuid,
		CreationsId:   id,
		Mode:          2,
		Status:        1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to collect a talk: talk_id(%v), collectionsId(%v), userUuid(%s)", id, collectionsId, userUuid))
	}
	return nil
}

func (r *talkRepo) SetTalkCollectToCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	statisticKey := fmt.Sprintf("talk_%v", id)
	collectKey := fmt.Sprintf("collections_%v_talk", collectionsId)
	collectionsKey := fmt.Sprintf("collections_%v", collectionsId)
	creationKey := fmt.Sprintf("creation_user_%s", userUuid)
	userKey := fmt.Sprintf("user_talk_collect_%s", userUuid)
	keys := []string{statisticKey, collectKey, collectionsKey, creationKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid)}
	_, err := r.data.redisCli.EvalSha(ctx, "1feecf8c0a6bb8c20f300d1cbb00e4d0e1d2c5a4", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk collect to cache: id(%v), collectionsId(%v), uuid(%s), userUuid(%s)", id, collectionsId, uuid, userUuid))
	}
	return nil
}

func (r *talkRepo) SetTalkCollect(ctx context.Context, id int32, uuid string) error {
	ts := TalkStatistic{}
	err := r.data.DB(ctx).Model(&ts).Where("talk_id = ? and uuid = ?", id, uuid).Update("collect", gorm.Expr("collect + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk collect: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) SetUserTalkAgreeToCache(ctx context.Context, id int32, userUuid string) error {
	keys := []string{"user_talk_agree_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := r.data.redisCli.EvalSha(ctx, "1c622b6aa57cb9c250a7d728d5639901e2a75fbd", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user talk agree to cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) SetUserTalkCollectToCache(ctx context.Context, id int32, userUuid string) error {
	keys := []string{"user_talk_collect_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := r.data.redisCli.EvalSha(ctx, "16415c44ae0544f8c7f85d841521813d41c35994", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user talk collect to cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) SetCollectionsTalkCollect(ctx context.Context, id, collectionsId int32, userUuid string) error {
	collect := &Collect{
		CollectionsId: collectionsId,
		Uuid:          userUuid,
		CreationsId:   id,
		Mode:          3,
		Status:        1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to collect an talk: talk_id(%v), collectionsId(%v), userUuid(%s)", id, collectionsId, userUuid))
	}
	return nil
}

func (r *talkRepo) SetCollectionTalk(ctx context.Context, collectionsId int32, userUuid string) error {
	c := Collections{}
	err := r.data.DB(ctx).Model(&c).Where("collections_id = ? and uuid = ?", collectionsId, userUuid).Update("talk", gorm.Expr("talk + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add collection talk collect: collectionsId(%v), userUuid(%s)", collectionsId, userUuid))
	}
	return nil
}

func (r *talkRepo) SetUserTalkCollect(ctx context.Context, id int32, userUuid string) error {
	tc := &TalkCollect{
		TalkId: id,
		Uuid:   userUuid,
		Status: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(tc).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add user talk collect: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) SetCreationUserCollect(ctx context.Context, userUuid string) error {
	cu := &CreationUser{
		Uuid:    userUuid,
		Collect: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"collect": gorm.Expr("collect + ?", 1)}),
	}).Create(cu).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation collect: uuid(%v)", userUuid))
	}
	return nil
}

func (r *talkRepo) SetTalkImageIrregular(ctx context.Context, review *biz.ImageReview) (*biz.ImageReview, error) {
	ar := &TalkReview{
		TalkId:   review.CreationId,
		Kind:     review.Kind,
		Uid:      review.Uid,
		Uuid:     review.Uuid,
		JobId:    review.JobId,
		Url:      review.Url,
		Label:    review.Label,
		Result:   review.Result,
		Category: review.Category,
		SubLabel: review.SubLabel,
		Score:    review.Score,
	}
	err := r.data.DB(ctx).Select("TalkId", "Kind", "Uid", "Uuid", "JobId", "Url", "Label", "Result", "Category", "SubLabel", "Score").Create(ar).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to add talk image review record: review(%v)", review))
	}
	review.Id = int32(ar.ID)
	review.CreateAt = time.Now().Format("2006-01-02")
	return review, nil
}

func (r *talkRepo) SetTalkImageIrregularToCache(ctx context.Context, review *biz.ImageReview) error {
	marshal, err := review.MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set talk image irregular to json: json.Marshal(%v)", review))
	}
	keys := []string{"talk_image_irregular_" + review.Uuid}
	values := []interface{}{marshal}
	_, err = r.data.redisCli.EvalSha(ctx, "8f6205011a2b264278a7c5bc0a2bcd1006ac6e5d", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set talk image irregular to cache: review(%v)", review))
	}
	return nil
}

func (r *talkRepo) SetTalkContentIrregular(ctx context.Context, review *biz.TextReview) (*biz.TextReview, error) {
	ar := &TalkContentReview{
		TalkId:  review.CreationId,
		Title:   review.Title,
		Kind:    review.Kind,
		Uuid:    review.Uuid,
		JobId:   review.JobId,
		Label:   review.Label,
		Result:  review.Result,
		Section: review.Section,
	}
	err := r.data.DB(ctx).Select("TalkId", "Title", "Kind", "Uuid", "JobId", "Label", "Result", "Section").Create(ar).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to add talk content review record: review(%v)", review))
	}
	review.Id = int32(ar.ID)
	review.CreateAt = time.Now().Format("2006-01-02")
	return review, nil
}

func (r *talkRepo) SetTalkContentIrregularToCache(ctx context.Context, review *biz.TextReview) error {
	marshal, err := review.MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set talk content irregular to json: json.Marshal(%v)", review))
	}
	keys := []string{"talk_content_irregular_" + review.Uuid}
	values := []interface{}{marshal}
	_, err = r.data.redisCli.EvalSha(ctx, "8f6205011a2b264278a7c5bc0a2bcd1006ac6e5d", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set talk content irregular to cache: review(%v)", review))
	}
	return nil
}

func (r *talkRepo) CancelTalkUserCollect(ctx context.Context, id int32, userUuid string) error {
	collect := &Collect{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&Collect{}).Where("creations_id = ? and mode = ? and uuid = ?", id, 2, userUuid).Updates(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel talk collect: talk_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) CancelTalkCollect(ctx context.Context, id int32, uuid string) error {
	ts := &TalkStatistic{}
	err := r.data.DB(ctx).Model(ts).Where("talk_id = ? and uuid = ? and collect > 0", id, uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel talk collect: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) CancelTalkCollectFromCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	statisticKey := fmt.Sprintf("talk_%v", id)
	collectKey := fmt.Sprintf("collections_%v_talk", collectionsId)
	collectionsKey := fmt.Sprintf("collections_%v", collectionsId)
	creationKey := fmt.Sprintf("creation_user_%s", userUuid)
	userKey := fmt.Sprintf("user_talk_collect_%s", userUuid)
	keys := []string{statisticKey, collectKey, collectionsKey, creationKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), id}
	_, err := r.data.redisCli.EvalSha(ctx, "c54e39b45bc53e7cd70ad2d28fc2eb4e68625784", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel talk collect from cache: id(%v), collectionsId(%v), uuid(%s), userUuid(%s)", id, collectionsId, uuid, userUuid))
	}
	return nil
}

func (r *talkRepo) CancelUserTalkAgreeFromCache(ctx context.Context, id int32, userUuid string) error {
	keys := []string{"user_talk_agree_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := r.data.redisCli.EvalSha(ctx, "d7d8a559ccbad93b64d6e4ef0b9130d2b5992224", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user talk agree from cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) CancelUserTalkCollectFromCache(ctx context.Context, id int32, userUuid string) error {
	keys := []string{"user_talk_collect_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := r.data.redisCli.EvalSha(ctx, "9aab0f65688566d3d54b7851570412ad816a80bb", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user talk agree from cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) CancelCollectionsTalkCollect(ctx context.Context, id int32, userUuid string) error {
	collect := &Collect{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&Collect{}).Where("creations_id = ? and mode = ? and uuid = ?", id, 3, userUuid).Updates(collect).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel talk collect: talk_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) CancelUserTalkCollect(ctx context.Context, id int32, userUuid string) error {
	tc := &TalkCollect{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&TalkCollect{}).Where("talk_id = ? and uuid = ?", id, userUuid).Updates(tc).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user talk collect: talk_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) CancelCollectionTalk(ctx context.Context, id int32, uuid string) error {
	collections := &Collections{}
	err := r.data.DB(ctx).Model(collections).Where("collections_id = ? and uuid = ? and talk > 0", id, uuid).Update("talk", gorm.Expr("talk - ?", 1)).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel collections talk: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) ReduceCreationUserCollect(ctx context.Context, uuid string) error {
	cu := &CreationUser{}
	err := r.data.DB(ctx).Model(cu).Where("uuid = ? and collect > 0", uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user collect: uuid(%v)", uuid))
	}
	return nil
}

func (r *talkRepo) CreateTalkSearch(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	key := "talk/" + uuid + "/" + ids + "/search"
	resp, err := r.data.cosCli.Object.Get(ctx, key, &cos.ObjectGetOptions{})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get talk from cos: id(%v), uuid(%s)", id, uuid))
	}

	talk, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to read request body: id(%v), uuid(%s)", id, uuid))
	}

	resp.Body.Close()

	req := esapi.IndexRequest{
		Index:      "talk",
		DocumentID: strconv.Itoa(int(id)),
		Body:       bytes.NewReader(talk),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, r.data.elasticSearch.es)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Error getting talk search create response: id(%v), uuid(%s)", id, uuid))
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("error parsing the response body: id(%v), uuid(%s)", id, uuid))
		} else {
			return errors.Errorf(fmt.Sprintf("error indexing document to es: reason(%v), id(%v), uuid(%s)", e, id, uuid))
		}
	}
	return nil
}

func (r *talkRepo) AddTalkComment(ctx context.Context, id int32) error {
	ts := TalkStatistic{}
	err := r.data.DB(ctx).Model(&ts).Where("talk_id = ?", id).Update("comment", gorm.Expr("comment + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk comment: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) AddTalkCommentToCache(ctx context.Context, id int32, uuid string) error {
	key := "talk_" + strconv.Itoa(int(id))
	keys := []string{key}
	var values []interface{}
	_, err := r.data.redisCli.EvalSha(ctx, "17cef96499553bacd3611f03630542e2b475a42c", keys, values...).Result()
	if err != nil {
		r.log.Errorf("fail to add talk comment to cache: id(%v), uuid(%s), err(%v)", id, uuid, err)
	}
	return nil
}

func (r *talkRepo) ReduceTalkComment(ctx context.Context, id int32) error {
	ts := TalkStatistic{}
	err := r.data.DB(ctx).Model(&ts).Where("talk_id = ? and comment > 0", id).Update("comment", gorm.Expr("comment - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce talk comment: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) ReduceTalkCommentToCache(ctx context.Context, id int32, uuid string) error {
	key := "talk_" + strconv.Itoa(int(id))
	keys := []string{key}
	_, err := r.data.redisCli.EvalSha(ctx, "e9942f8b94d123b097211f27b5bf8fc896b1fa45", keys).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce talk comment to cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *talkRepo) ReduceCreationUserTalk(ctx context.Context, auth int32, uuid string) error {
	cu := CreationUser{}
	err := r.data.DB(ctx).Model(&cu).Where("uuid = ? and talk > 0", uuid).Update("talk", gorm.Expr("talk - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user talk: uuid(%s)", uuid))
	}

	if auth == 2 {
		return nil
	}

	cuv := CreationUserVisitor{}
	err = r.data.DB(ctx).Model(&cuv).Where("uuid = ? and talk > 0", uuid).Update("talk", gorm.Expr("talk - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user talk visitor: uuid(%s)", uuid))
	}
	return nil
}

func (r *talkRepo) EditTalkSearch(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	key := "talk/" + uuid + "/" + ids + "/search"
	resp, err := r.data.cosCli.Object.Get(ctx, key, &cos.ObjectGetOptions{})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get talk from cos: id(%v), uuid(%s)", id, uuid))
	}

	talk, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to read request body: id(%v), uuid(%s)", id, uuid))
	}

	resp.Body.Close()

	req := esapi.IndexRequest{
		Index:      "talk",
		DocumentID: strconv.Itoa(int(id)),
		Body:       bytes.NewReader(talk),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, r.data.elasticSearch.es)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Error getting talk search edit response: id(%v), uuid(%s)", id, uuid))
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("error parsing the response body: id(%v), uuid(%s)", id, uuid))
		} else {
			return errors.Errorf(fmt.Sprintf("error indexing document to es: reason(%v), id(%v), uuid(%s)", e, id, uuid))
		}
	}
	return nil
}

func (r *talkRepo) DeleteTalkSearch(ctx context.Context, id int32, uuid string) error {
	req := esapi.DeleteRequest{
		Index:      "talk",
		DocumentID: strconv.Itoa(int(id)),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, r.data.elasticSearch.es)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Error getting talk search delete response: id(%v), uuid(%s)", id, uuid))
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("error parsing the response body: id(%v), uuid(%s)", id, uuid))
		} else {
			return errors.Errorf(fmt.Sprintf("error delete document to es: reason(%v), id(%v), uuid(%s)", e, id, uuid))
		}
	}
	return nil
}
