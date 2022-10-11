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
	err := r.data.db.WithContext(ctx).Where("talk_id = ?", id).First(talk).Error
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
		go r.data.Recover(context.Background(), func(ctx context.Context) {
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
		go r.data.Recover(context.Background(), func(ctx context.Context) {
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
		go r.data.Recover(context.Background(), func(ctx context.Context) {
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

	talk := make([]*biz.Talk, 0)
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
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Order("talk_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk from db: page(%v), uuid(%s)", page, uuid))
	}

	talk := make([]*biz.Talk, 0)
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
		go r.data.Recover(context.Background(), func(ctx context.Context) {
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

	talk := make([]*biz.Talk, 0)
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
	err := r.data.db.WithContext(ctx).Where("uuid = ? and auth = ?", uuid, 1).Order("talk_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk visitor from db: page(%v), uuid(%s)", page, uuid))
	}

	talk := make([]*biz.Talk, 0)
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
		z := make([]*redis.Z, 0)
		for _, item := range talk {
			z = append(z, &redis.Z{
				Score:  float64(item.TalkId),
				Member: strconv.Itoa(int(item.TalkId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Hour*8)
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
		z := make([]*redis.Z, 0)
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

	talk := make([]*biz.TalkStatistic, 0)
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
	err := r.data.db.WithContext(ctx).Where("auth", 1).Order("agree desc, talk_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk statistic from db: page(%v)", page))
	}

	talk := make([]*biz.TalkStatistic, 0)
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

	talk := make([]*biz.Talk, 0)
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
	err := r.data.db.WithContext(ctx).Where("auth", 1).Order("talk_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk from db: page(%v)", page))
	}

	talk := make([]*biz.Talk, 0)
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
		z := make([]*redis.Z, 0)
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
	talkListStatistic := make([]*biz.TalkStatistic, 0)
	exists, unExists, err := r.talkListStatisticExist(ctx, ids)
	if err != nil {
		return nil, err
	}

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
	exists := make([]int32, 0)
	unExists := make([]int32, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range ids {
			pipe.Exists(ctx, "talk_"+strconv.Itoa(int(item)))
		}
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, fmt.Sprintf("fail to check if talk statistic exist from cache: ids(%v)", ids))
	}

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
	list := make([]*TalkStatistic, 0)
	err := r.data.db.WithContext(ctx).Where("talk_id IN ?", unExists).Find(&list).Error
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
		go r.data.Recover(context.Background(), func(ctx context.Context) {
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
			pipe.Expire(ctx, key, time.Hour*8)
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

	go r.data.Recover(context.Background(), func(ctx context.Context) {
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
	err := r.data.db.WithContext(ctx).Where("talk_id = ?", id).First(as).Error
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
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk statistic to cache, err(%s)", err.Error())
	}
}

func (r *talkRepo) GetLastTalkDraft(ctx context.Context, uuid string) (*biz.TalkDraft, error) {
	draft := &TalkDraft{}
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Last(draft).Error
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

		reply = append(reply, &biz.TalkSearch{
			Id:     int32(id),
			Tags:   hit.(map[string]interface{})["_source"].(map[string]interface{})["tags"].(string),
			Update: hit.(map[string]interface{})["_source"].(map[string]interface{})["update"].(string),
			Cover:  hit.(map[string]interface{})["_source"].(map[string]interface{})["cover"].(string),
			Uuid:   hit.(map[string]interface{})["_source"].(map[string]interface{})["uuid"].(string),
			Text:   hit.(map[string]interface{})["highlight"].(map[string]interface{})["text"].([]interface{})[0].(string),
			Title:  hit.(map[string]interface{})["highlight"].(map[string]interface{})["title"].([]interface{})[0].(string),
		})
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
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk agree from db: uuid(%s)", uuid))
	}

	agreeMap := make(map[int32]bool, 0)
	for _, item := range list {
		agreeMap[item.TalkId] = true
	}
	if len(list) != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setUserTalkAgreeToCache(uuid, list)
		})()
	}
	return agreeMap, nil
}

func (r *talkRepo) setUserTalkAgreeToCache(uuid string, agreeList []*TalkAgree) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0)
		key := "user_talk_agree_" + uuid
		for _, item := range agreeList {
			set = append(set, item.TalkId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Hour*8)
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
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk collect from db: uuid(%s)", uuid))
	}

	collectMap := make(map[int32]bool, 0)
	for _, item := range list {
		collectMap[item.TalkId] = true
	}
	if len(list) != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setUserTalkCollectToCache(uuid, list)
		})()
	}
	return collectMap, nil
}

func (r *talkRepo) setUserTalkCollectToCache(uuid string, collectList []*TalkCollect) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0)
		key := "user_talk_collect_" + uuid
		for _, item := range collectList {
			set = append(set, item.TalkId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user talk collect to cache: uuid(%s), collectList(%v), error(%s) ", uuid, collectList, err.Error())
	}
}

func (r *talkRepo) GetTalkAuth(ctx context.Context, id int32) (int32, error) {
	talk := &Talk{}
	err := r.data.db.WithContext(ctx).Where("talk_id = ?", id).First(talk).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get talk auth from db: id(%v)", id))
	}
	return talk.Auth, nil
}

func (r *talkRepo) GetCollectionsIdFromTalkCollect(ctx context.Context, id int32) (int32, error) {
	collect := &Collect{}
	err := r.data.db.WithContext(ctx).Where("creations_id = ? and mode = ?", id, 3).First(collect).Error
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
		go r.data.Recover(context.Background(), func(ctx context.Context) {
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

	review := make([]*biz.ImageReview, 0)
	for _index, item := range list {
		var imageReview = &biz.ImageReview{}
		err = json.Unmarshal([]byte(item), imageReview)
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
	err := r.data.db.WithContext(ctx).Where("uuid", uuid).Order("id desc").Offset(index * 20).Limit(20).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk image review from db: page(%v), uuid(%s)", page, uuid))
	}

	review := make([]*biz.ImageReview, 0)
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
		list := make([]interface{}, 0)
		for _, item := range review {
			m, err := json.Marshal(item)
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("fail to marshal avatar review: imageReview(%v)", review))
			}
			list = append(list, m)
		}
		pipe.RPush(ctx, key, list...)
		pipe.Expire(ctx, key, time.Hour*8)
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
		go r.data.Recover(context.Background(), func(ctx context.Context) {
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

	review := make([]*biz.TextReview, 0)
	for _index, item := range list {
		var textReview = &biz.TextReview{}
		err = json.Unmarshal([]byte(item), textReview)
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
	err := r.data.db.WithContext(ctx).Where("uuid", uuid).Order("id desc").Offset(index * 20).Limit(20).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk content review from db: page(%v), uuid(%s)", page, uuid))
	}

	review := make([]*biz.TextReview, 0)
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
		list := make([]interface{}, 0)
		for _, item := range review {
			m, err := json.Marshal(item)
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("fail to marshal avatar review: contentReview(%v)", review))
			}
			list = append(list, m)
		}
		pipe.RPush(ctx, key, list...)
		pipe.Expire(ctx, key, time.Hour*8)
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
	data, err := json.Marshal(review)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "talk_review",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.talkReviewMqPro.producer.SendSync(ctx, msg)
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
	var script = redis.NewScript(`
					local key1 = KEYS[1]
					local key2 = KEYS[2]
					local key3 = KEYS[3]
					local key4 = KEYS[4]
					local key5 = KEYS[5]
					local key6 = KEYS[6]
					local key7 = KEYS[7]
					local key8 = KEYS[8]
					local key9 = KEYS[9]

                    local member = ARGV[1]
					local leadMember = ARGV[2]
					local auth = ARGV[3]

                    redis.call("ZREM", key1, member)
					redis.call("ZREM", key2, member)
					redis.call("ZREM", key3, leadMember)
					redis.call("DEL", key4)
					redis.call("DEL", key5)

                    redis.call("ZREM", key6, member)

                    local exist = redis.call("EXISTS", key8)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key8, "talk"))
						if number > 0 then
  							redis.call("HINCRBY", key8, "talk", -1)
						end
					end

                    if auth == "2" then
						return 0
					end

					redis.call("ZREM", key7, member)

					local exist = redis.call("EXISTS", key9)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key9, "talk"))
						if number > 0 then
  							redis.call("HINCRBY", key9, "talk", -1)
						end
					end

					return 0
	`)
	keys := []string{"talk", "talk_hot", "leaderboard", "talk_" + ids, "talk_collect_" + ids, "user_talk_list_" + uuid, "user_talk_list_visitor_" + uuid, "creation_user_" + uuid, "creation_user_visitor_" + uuid}
	values := []interface{}{ids + "%" + uuid, ids + "%" + uuid + "%talk", auth}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	talkMap := map[string]interface{}{}
	talkMap["uuid"] = talk.Uuid
	talkMap["id"] = talk.TalkId
	talkMap["auth"] = talk.Auth
	talkMap["mode"] = mode

	data, err := json.Marshal(talkMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "talk",
		Body:  data,
	}
	msg.WithKeys([]string{talk.Uuid})
	_, err = r.data.talkMqPro.producer.SendSync(ctx, msg)
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
	var script = redis.NewScript(`
					local talkStatistic = KEYS[1]
					local talk = KEYS[2]
					local talkHot = KEYS[3]
					local leaderboard = KEYS[4]
					local userTalkList = KEYS[5]
					local userTalkListVisitor = KEYS[6]
					local creationUser = KEYS[7]
					local creationUserVisitor = KEYS[8]

					local uuid = ARGV[1]
					local auth = ARGV[2]
					local id = ARGV[3]
					local member = ARGV[4]
					local mode = ARGV[5]
					local member2 = ARGV[6]

					local userTalkListExist = redis.call("EXISTS", userTalkList)
					local creationUserExist = redis.call("EXISTS", creationUser)
					local talkExist = redis.call("EXISTS", talk)
					local talkHotExist = redis.call("EXISTS", talkHot)
					local leaderboardExist = redis.call("EXISTS", leaderboard)
					local userTalkListVisitorExist = redis.call("EXISTS", userTalkListVisitor)
					local creationUserVisitorExist = redis.call("EXISTS", creationUserVisitor)

					redis.call("HSETNX", talkStatistic, "uuid", uuid)
					redis.call("HSETNX", talkStatistic, "agree", 0)
					redis.call("HSETNX", talkStatistic, "collect", 0)
					redis.call("HSETNX", talkStatistic, "view", 0)
					redis.call("HSETNX", talkStatistic, "comment", 0)
					redis.call("HSETNX", talkStatistic, "auth", auth)
					redis.call("EXPIRE", talkStatistic, 28800)

					if userTalkListExist == 1 then
						redis.call("ZADD", userTalkList, id, member)
					end

					if (creationUserExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUser, "talk", 1)
					end

					if auth == "2" then
						return 0
					end

					if talkExist == 1 then
						redis.call("ZADD", talk, id, member)
					end

					if talkHotExist == 1 then
						redis.call("ZADD", talkHot, id, member)
					end

					if leaderboardExist == 1 then
						redis.call("ZADD", leaderboard, 0, member2)
					end

					if userTalkListVisitorExist == 1 then
						redis.call("ZADD", userTalkListVisitor, id, member)
					end

					if (creationUserVisitorExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUserVisitor, "talk", 1)
					end
					return 0
	`)
	keys := []string{talkStatistic, talk, talkHot, leaderboard, userTalkList, userTalkListVisitor, creationUser, creationUserVisitor}
	values := []interface{}{uuid, auth, id, ids + "%" + uuid, mode, ids + "%" + uuid + "%talk"}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	var script = redis.NewScript(`
					local hotKey = KEYS[1]
					local statisticKey = KEYS[2]
					local boardKey = KEYS[3]
					local userKey = KEYS[4]

					local member1 = ARGV[1]
					local member2 = ARGV[2]
					local id = ARGV[3]

					local hotKeyExist = redis.call("EXISTS", hotKey)
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					local boardKeyExist = redis.call("EXISTS", boardKey)
					local userKeyExist = redis.call("EXISTS", userKey)

					if hotKeyExist == 1 then
						redis.call("ZINCRBY", hotKey, 1, member1)
					end

					if statisticKeyExist == 1 then
						redis.call("HINCRBY", statisticKey, "talk", 1)
					end

					if boardKeyExist == 1 then
						redis.call("ZINCRBY", boardKey, 1, member2)
					end

					if userKeyExist == 1 then
						redis.call("SADD", userKey, id)
					end
					return 0
	`)
	keys := []string{hotKey, statisticKey, boardKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), fmt.Sprintf("%v%s%s%s", id, "%", uuid, "%talk"), id}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
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

	var script = redis.NewScript(`
					local hotKey = KEYS[1]
                    local member = ARGV[1]
					local hotKeyExist = redis.call("EXISTS", hotKey)
					if hotKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", hotKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", hotKey, -1, member)
						end
					end

					local boardKey = KEYS[2]
                    local member = ARGV[2]
					local boardKeyExist = redis.call("EXISTS", boardKey)
					if boardKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", boardKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", boardKey, -1, member)
						end
					end

					local statisticKey = KEYS[3]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "agree", -1)
						end
					end

					local userKey = KEYS[4]
					local commentId = ARGV[3]
					redis.call("SREM", userKey, commentId)
					return 0
	`)
	keys := []string{hotKey, boardKey, statisticKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), fmt.Sprintf("%v%s%s%s", id, "%", uuid, "%talk"), id}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "talk_"+ids, "view", 1)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add talk agree to cache: id(%v), uuid(%s), err(%v)", id, uuid, err)
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
	achievement := map[string]interface{}{}
	achievement["uuid"] = uuid
	achievement["userUuid"] = userUuid
	achievement["mode"] = mode

	data, err := json.Marshal(achievement)
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
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk statistic to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *talkRepo) SendScoreToMq(ctx context.Context, score int32, uuid, mode string) error {
	scoreMap := map[string]interface{}{}
	scoreMap["uuid"] = uuid
	scoreMap["score"] = score
	scoreMap["mode"] = mode

	data, err := json.Marshal(scoreMap)
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
		return errors.Wrapf(err, fmt.Sprintf("fail to send score to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *talkRepo) SendStatisticToMq(ctx context.Context, id, collectionsId int32, uuid, userUuid, mode string) error {
	statisticMap := map[string]interface{}{}
	statisticMap["id"] = id
	statisticMap["collectionsId"] = collectionsId
	statisticMap["uuid"] = uuid
	statisticMap["userUuid"] = userUuid
	statisticMap["mode"] = mode

	data, err := json.Marshal(statisticMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "talk",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.talkMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk statistic to mq: map(%v)", statisticMap))
	}
	return nil
}

func (r *talkRepo) SendTalkImageIrregularToMq(ctx context.Context, review *biz.ImageReview) error {
	m := make(map[string]interface{}, 0)
	m["creation_id"] = review.CreationId
	m["score"] = review.Score
	m["result"] = review.Result
	m["kind"] = review.Kind
	m["uid"] = review.Uid
	m["uuid"] = review.Uuid
	m["job_id"] = review.JobId
	m["label"] = review.Label
	m["category"] = review.Category
	m["sub_label"] = review.SubLabel
	m["mode"] = review.Mode
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "talk",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.talkMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk image review to mq: %v", err))
	}
	return nil
}

func (r *talkRepo) SendTalkContentIrregularToMq(ctx context.Context, review *biz.TextReview) error {
	m := make(map[string]interface{}, 0)
	m["creation_id"] = review.CreationId
	m["result"] = review.Result
	m["uuid"] = review.Uuid
	m["job_id"] = review.JobId
	m["label"] = review.Label
	m["title"] = review.Title
	m["kind"] = review.Kind
	m["section"] = review.Section
	m["mode"] = review.Mode
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "talk",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.talkMqPro.producer.SendSync(ctx, msg)
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
	var script = redis.NewScript(`
					local statisticKey = KEYS[1]
					local collectKey = KEYS[2]
					local collectionsKey = KEYS[3]
					local creationKey = KEYS[4]
					local userKey = KEYS[5]

					local member = ARGV[1]

					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					local collectKeyExist = redis.call("EXISTS", collectKey)
					local collectionsKeyExist = redis.call("EXISTS", collectionsKey)
					local creationKeyExist = redis.call("EXISTS", creationKey)
					local userKeyExist = redis.call("EXISTS", userKey)
					
					if statisticKeyExist == 1 then
						redis.call("HINCRBY", statisticKey, "collect", 1)
					end

					if collectKeyExist == 1 then
						redis.call("ZADD", collectKey, 0, member)
					end

					if collectionsKeyExist == 1 then
						redis.call("HINCRBY", collectionsKey, "talk", 1)
					end

					if creationKeyExist == 1 then
						redis.call("HINCRBY", creationKey, "collect", 1)
					end

					if userKeyExist == 1 then
						redis.call("SADD", userKey, id)
					end
					return 0
	`)
	keys := []string{statisticKey, collectKey, collectionsKey, creationKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid)}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
						redis.call("SADD", key, change)
					end
					return 0
	`)
	keys := []string{"user_talk_agree_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user talk agree to cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) SetUserTalkCollectToCache(ctx context.Context, id int32, userUuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`)
	keys := []string{"user_talk_collect_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	marshal, err := json.Marshal(review)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set talk image irregular to json: json.Marshal(%v)", review))
	}
	var script = redis.NewScript(`
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`)
	keys := []string{"talk_image_irregular_" + review.Uuid}
	values := []interface{}{marshal}
	_, err = script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	marshal, err := json.Marshal(review)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set talk content irregular to json: json.Marshal(%v)", review))
	}
	var script = redis.NewScript(`
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`)
	keys := []string{"talk_content_irregular_" + review.Uuid}
	values := []interface{}{marshal}
	_, err = script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	var script = redis.NewScript(`
					local statisticKey = KEYS[1]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "collect", -1)
						end
					end

					local collectKey = KEYS[2]
					local talkMember = ARGV[1]
					redis.call("ZREM", collectKey, talkMember)

					local collectionsKey = KEYS[3]
					local collectionsKeyExist = redis.call("EXISTS", collectionsKey)
					if collectionsKeyExist == 1 then
						local number = tonumber(redis.call("HGET", collectionsKey, "talk"))
						if number > 0 then
  							redis.call("HINCRBY", collectionsKey, "talk", -1)
						end
					end

					local creationKey = KEYS[4]
					local creationKeyExist = redis.call("EXISTS", creationKey)
					if creationKeyExist == 1 then
						local number = tonumber(redis.call("HGET", creationKey, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", creationKey, "collect", -1)
						end
					end

					local userKey = KEYS[5]
					local talkId = ARGV[2]
					redis.call("SREM", userKey, talkId)

					return 0
	`)
	keys := []string{statisticKey, collectKey, collectionsKey, creationKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), id}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel talk collect from cache: id(%v), collectionsId(%v), uuid(%s), userUuid(%s)", id, collectionsId, uuid, userUuid))
	}
	return nil
}

func (r *talkRepo) CancelUserTalkAgreeFromCache(ctx context.Context, id int32, userUuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SREM", key, change)
					end
					return 0
	`)
	keys := []string{"user_talk_agree_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user talk agree from cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *talkRepo) CancelUserTalkCollectFromCache(ctx context.Context, id int32, userUuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
						redis.call("SREM", key, change)
					end
					return 0
	`)
	keys := []string{"user_talk_collect_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	var script = redis.NewScript(`
					local key = KEYS[1]
					local keyExist = redis.call("EXISTS", key)
					if keyExist == 1 then
						redis.call("HINCRBY", key, "comment", 1)
					end
					return 0
	`)
	keys := []string{key}
	var values []interface{}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	var script = redis.NewScript(`
					local key = KEYS[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
						local number = tonumber(redis.call("HGET", key, "comment"))
						if number > 0 then
  							redis.call("HINCRBY", key, "comment", -1)
						end
					end
					return 0
	`)
	keys := []string{key}
	_, err := script.Run(ctx, r.data.redisCli, keys).Result()
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
