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
		go r.setTalkToCache("talk", talk)
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
		go r.setTalkHotToCache("talk_hot", talk)
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
		go r.setUserTalkListToCache("user_talk_list_"+uuid, talk)
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
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk list visitor from cache: key(%s), page(%v)", "user_talk_list_", page))
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
		go r.setUserTalkListToCache("user_talk_list_visitor_"+uuid, talk)
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
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user talk list visitor from cache: key(%s), page(%v)", "user_talk_list_visitor_", page))
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
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range talk {
			z = append(z, &redis.Z{
				Score:  float64(item.TalkId),
				Member: strconv.Itoa(int(item.TalkId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		pipe.Expire(context.Background(), key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk to cache: talk(%v)", talk)
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
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range talk {
			z = append(z, &redis.Z{
				Score:  float64(item.Agree),
				Member: strconv.Itoa(int(item.TalkId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk to cache: talk(%v)", talk)
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
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range talk {
			z = append(z, &redis.Z{
				Score:  float64(item.TalkId),
				Member: strconv.Itoa(int(item.TalkId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk to cache: talk(%v)", talk)
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
		})
	}

	if len(list) != 0 {
		go r.setTalkListStatisticToCache(list)
	}

	return nil
}

func (r *talkRepo) setTalkListStatisticToCache(commentList []*TalkStatistic) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, item := range commentList {
			key := "talk_" + strconv.Itoa(int(item.TalkId))
			pipe.HSetNX(context.Background(), key, "uuid", item.Uuid)
			pipe.HSetNX(context.Background(), key, "agree", item.Agree)
			pipe.HSetNX(context.Background(), key, "comment", item.Comment)
			pipe.HSetNX(context.Background(), key, "collect", item.Collect)
			pipe.HSetNX(context.Background(), key, "view", item.View)
			pipe.Expire(context.Background(), key, time.Hour*8)
		}
		return nil
	})

	if err != nil {
		r.log.Errorf("fail to set talk statistic to cache, err(%s)", err.Error())
	}
}

func (r *talkRepo) GetTalkStatistic(ctx context.Context, id int32) (*biz.TalkStatistic, error) {
	var statistic *biz.TalkStatistic
	key := "talk_" + strconv.Itoa(int(id))
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		statistic, err = r.getTalkStatisticFromCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return statistic, nil
	}

	statistic, err = r.getTalkStatisticFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	go r.setTalkStatisticToCache(key, statistic)

	return statistic, nil
}

func (r *talkRepo) getTalkStatisticFromCache(ctx context.Context, key string) (*biz.TalkStatistic, error) {
	statistic, err := r.data.redisCli.HMGet(ctx, key, "uuid", "agree", "collect", "view", "comment").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk statistic form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0}
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
	err := r.data.redisCli.HMSet(context.Background(), key, "uuid", statistic.Uuid, "agree", statistic.Agree, "collect", statistic.Collect, "view", statistic.View, "comment", statistic.Comment).Err()
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
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%s)", uuid))
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

func (r *talkRepo) GetTalkAuth(ctx context.Context, id int32) (int32, error) {
	talk := &Talk{}
	err := r.data.db.WithContext(ctx).Where("talk_id = ?", id).First(talk).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get talk auth from db: id(%v)", id))
	}
	return talk.Auth, nil
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
	var incrBy = redis.NewScript(`
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

                    if auth == 2 then
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
	keys := []string{"talk", "talk_hot", "leaderboard", "talk_" + ids, "talk_collect_" + ids, "user_talk_list_" + uuid, "user_talk_list_visitor_" + uuid, "user_creation_info_" + uuid, "user_creation_info_visitor_" + uuid}
	values := []interface{}{ids + "%" + uuid, ids + "%" + uuid + "%talk", auth}
	_, err := incrBy.Run(ctx, r.data.redisCli, keys, values...).Result()
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

func (r *talkRepo) CreateTalkCache(ctx context.Context, id, auth int32, uuid string) error {
	exists := make([]int32, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Exists(ctx, "talk")
		pipe.Exists(ctx, "talk_hot")
		pipe.Exists(ctx, "leaderboard")
		pipe.Exists(ctx, "user_talk_list_"+uuid)
		pipe.Exists(ctx, "user_talk_list_visitor_"+uuid)
		pipe.Exists(ctx, "user_creation_info_"+uuid)
		pipe.Exists(ctx, "user_creation_info_visitor_"+uuid)
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to check if talk exist from cache: id(%v),uuid(%s)", id, uuid))
	}

	for _, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		exists = append(exists, int32(exist))
	}

	ids := strconv.Itoa(int(id))
	_, err = r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSetNX(ctx, "talk_"+ids, "uuid", uuid)
		pipe.HSetNX(ctx, "talk_"+ids, "agree", 0)
		pipe.HSetNX(ctx, "talk_"+ids, "collect", 0)
		pipe.HSetNX(ctx, "talk_"+ids, "view", 0)
		pipe.HSetNX(ctx, "talk_"+ids, "comment", 0)

		if exists[3] == 1 {
			pipe.ZAddNX(ctx, "user_talk_list_"+uuid, &redis.Z{
				Score:  float64(id),
				Member: ids + "%" + uuid,
			})
		}

		if exists[5] == 1 {
			pipe.HIncrBy(ctx, "user_creation_info_"+uuid, "talk", 1)
		}

		if auth == 2 {
			return nil
		}

		if exists[0] == 1 {
			pipe.ZAddNX(ctx, "talk", &redis.Z{
				Score:  float64(id),
				Member: ids + "%" + uuid,
			})
		}

		if exists[1] == 1 {
			pipe.ZAddNX(ctx, "talk_hot", &redis.Z{
				Score:  0,
				Member: ids + "%" + uuid,
			})
		}

		if exists[2] == 1 {
			pipe.ZAddNX(ctx, "leaderboard", &redis.Z{
				Score:  0,
				Member: ids + "%" + uuid + "%talk",
			})
		}

		if exists[4] == 1 {
			pipe.ZAddNX(ctx, "user_talk_list_visitor_"+uuid, &redis.Z{
				Score:  0,
				Member: ids + "%" + uuid + "%talk",
			})
		}

		if exists[6] == 1 {
			pipe.HIncrBy(ctx, "user_creation_info_visitor_"+uuid, "talk", 1)
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create(update) talk cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) UpdateTalkCache(ctx context.Context, id, auth int32, uuid string) error {
	return r.CreateTalkCache(ctx, id, auth, uuid)
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

func (r *talkRepo) SetTalkAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "talk_"+ids, "agree", 1)
		pipe.ZIncrBy(ctx, "talk_hot", 1, ids+"%"+uuid)
		pipe.ZIncrBy(ctx, "leaderboard", 1, ids+"%"+uuid+"%talk")
		pipe.SAdd(ctx, "talk_agree_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add talk agree to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
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

func (r *talkRepo) CancelTalkAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "talk_"+ids, "agree", -1)
		pipe.ZIncrBy(ctx, "talk_hot", -1, ids+"%"+uuid)
		pipe.ZIncrBy(ctx, "leaderboard", -1, ids+"%"+uuid+"%talk")
		pipe.SRem(ctx, "talk_agree_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to cancel talk agree from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
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
		r.log.Errorf("fail to add talk agree to cache: id(%v), uuid(%s)", id, uuid)
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

func (r *talkRepo) SendTalkStatisticToMq(ctx context.Context, uuid, mode string) error {
	achievement := map[string]interface{}{}
	achievement["uuid"] = uuid
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

func (r *talkRepo) SetTalkCollect(ctx context.Context, id int32, uuid string) error {
	ts := TalkStatistic{}
	err := r.data.DB(ctx).Model(&ts).Where("talk_id = ? and uuid = ?", id, uuid).Update("collect", gorm.Expr("collect + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk collect: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) SetTalkCollectToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "talk_"+ids, "collect", 1)
		pipe.SAdd(ctx, "talk_collect_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add talk collect to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
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
	err := r.data.DB(ctx).Model(ts).Where("talk_id = ? and uuid = ?", id, uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel talk collect: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) CancelTalkCollectFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "talk_"+ids, "collect", -1)
		pipe.SRem(ctx, "talk_collect_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to cancel talk collect from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
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
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to check if talk statistic exist from cache: id(%v), uuid(%s)", id, uuid)
	}

	if exist == 0 {
		return nil
	}

	_, err = r.data.redisCli.HIncrBy(ctx, key, "comment", 1).Result()
	if err != nil {
		r.log.Errorf("fail to add talk comment to cache: id(%v), uuid(%s)", id, uuid)
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
	var incrBy = redis.NewScript(`
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
	_, err := incrBy.Run(ctx, r.data.redisCli, keys).Result()
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
