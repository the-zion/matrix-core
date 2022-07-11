package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var _ biz.ArticleRepo = (*articleRepo)(nil)

type articleRepo struct {
	data *Data
	log  *log.Helper
}

func NewArticleRepo(data *Data, logger log.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "creation/data/article")),
	}
}

func (r *articleRepo) GetLastArticleDraft(ctx context.Context, uuid string) (*biz.ArticleDraft, error) {
	draft := &ArticleDraft{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Last(draft).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("draft not found from db", fmt.Sprintf("uuid(%s)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%s)", uuid))
	}
	return &biz.ArticleDraft{
		Id:     int32(draft.ID),
		Status: draft.Status,
	}, nil
}

func (r *articleRepo) GetArticleList(ctx context.Context, page int32) ([]*biz.Article, error) {
	article, err := r.getArticleFromCache(ctx, page)
	if err != nil {
		return nil, err
	}

	size := len(article)
	if size != 0 {
		return article, nil
	}

	article, err = r.getArticleFromDB(ctx, page)
	if err != nil {
		return nil, err
	}

	size = len(article)
	if size != 0 {
		go r.setArticleToCache("article", article)
	}
	return article, nil
}

func (r *articleRepo) GetArticleAgreeJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	ids := strconv.Itoa(int(id))
	judge, err := r.data.redisCli.SIsMember(ctx, "article_agree_"+ids, uuid).Result()
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("fail to judge article agree member: id(%v), uuid(%s)", id, uuid))
	}
	return judge, nil
}

func (r *articleRepo) GetArticleCollectJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	ids := strconv.Itoa(int(id))
	judge, err := r.data.redisCli.SIsMember(ctx, "article_collect_"+ids, uuid).Result()
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("fail to judge article collect member: id(%v), uuid(%s)", id, uuid))
	}
	return judge, nil
}

func (r *articleRepo) getArticleFromDB(ctx context.Context, page int32) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Article, 0)
	err := r.data.db.WithContext(ctx).Order("article_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article from db: page(%v)", page))
	}

	article := make([]*biz.Article, 0)
	for _, item := range list {
		article = append(article, &biz.Article{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
		})
	}
	return article, nil
}

func (r *articleRepo) GetArticleListHot(ctx context.Context, page int32) ([]*biz.ArticleStatistic, error) {
	article, err := r.getArticleHotFromCache(ctx, page)
	if err != nil {
		return nil, err
	}

	size := len(article)
	if size != 0 {
		return article, nil
	}

	article, err = r.getArticleHotFromDB(ctx, page)
	if err != nil {
		return nil, err
	}

	size = len(article)
	if size != 0 {
		go r.setArticleHotToCache("article_hot", article)
	}
	return article, nil
}

func (r *articleRepo) getArticleHotFromDB(ctx context.Context, page int32) ([]*biz.ArticleStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*ArticleStatistic, 0)
	err := r.data.db.WithContext(ctx).Order("agree desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article statistic from db: page(%v)", page))
	}

	article := make([]*biz.ArticleStatistic, 0)
	for _, item := range list {
		article = append(article, &biz.ArticleStatistic{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
			Agree:     item.Agree,
		})
	}
	return article, nil
}

func (r *articleRepo) GetArticleStatistic(ctx context.Context, id int32) (*biz.ArticleStatistic, error) {
	var statistic *biz.ArticleStatistic
	key := "article_" + strconv.Itoa(int(id))
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		statistic, err = r.getArticleStatisticFromCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return statistic, nil
	}

	statistic, err = r.getArticleStatisticFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	go r.setArticleStatisticToCache(key, statistic)

	return statistic, nil
}

func (r *articleRepo) getArticleStatisticFromDB(ctx context.Context, id int32) (*biz.ArticleStatistic, error) {
	as := &ArticleStatistic{}
	err := r.data.db.WithContext(ctx).Where("article_id = ?", id).First(as).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get statistic from db: id(%v)", id))
	}
	return &biz.ArticleStatistic{
		Uuid:    as.Uuid,
		Agree:   as.Agree,
		Collect: as.Collect,
		View:    as.View,
		Comment: as.Comment,
	}, nil
}

func (r *articleRepo) GetArticleListStatistic(ctx context.Context, ids []int32) ([]*biz.ArticleStatistic, error) {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, id := range ids {
			pipe.HMGet(ctx, "article_"+strconv.Itoa(int(id)), "agree", "collect", "view", "comment")
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article list statistic from cache: ids(%v)", ids))
	}

	statistic := make([]*biz.ArticleStatistic, 0)
	for index, item := range cmd {
		val := []int32{0, 0, 0, 0}
		for _index, count := range item.(*redis.SliceCmd).Val() {
			if count == nil {
				break
			}
			num, err := strconv.ParseInt(count.(string), 10, 32)
			if err != nil {
				return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
			}
			val[_index] = int32(num)
		}
		statistic = append(statistic, &biz.ArticleStatistic{
			ArticleId: ids[index],
			Agree:     val[0],
			Collect:   val[1],
			View:      val[2],
			Comment:   val[3],
		})
	}
	return statistic, nil
}

func (r *articleRepo) GetArticleDraftList(ctx context.Context, uuid string) ([]*biz.ArticleDraft, error) {
	reply := make([]*biz.ArticleDraft, 0)
	draftList := make([]*ArticleDraft, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 3).Order("id desc").Find(&draftList).Error
	//err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Order("id desc").Find(&draftList).Error
	if err != nil {
		return reply, errors.Wrapf(err, fmt.Sprintf("fail to get draft list : uuid(%s)", uuid))
	}
	for _, item := range draftList {
		reply = append(reply, &biz.ArticleDraft{
			Id: int32(item.ID),
		})
	}
	return reply, nil
}

func (r *articleRepo) CreateArticle(ctx context.Context, id int32, uuid string) error {
	article := &Article{
		ArticleId: id,
		Uuid:      uuid,
	}
	err := r.data.DB(ctx).Select("ArticleId", "Uuid").Create(article).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a article: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) CreateArticleDraft(ctx context.Context, uuid string) (int32, error) {
	draft := &ArticleDraft{
		Uuid: uuid,
	}
	err := r.data.DB(ctx).Select("Uuid").Create(draft).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to create an article draft: uuid(%s)", uuid))
	}
	return int32(draft.ID), nil
}

func (r *articleRepo) CreateArticleFolder(ctx context.Context, id int32) error {
	name := "article/" + strconv.Itoa(int(id)) + "/"
	_, err := r.data.cosCli.Object.Put(ctx, name, strings.NewReader(""), nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create an article folder: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) CreateArticleStatistic(ctx context.Context, id int32, uuid string) error {
	as := &ArticleStatistic{
		ArticleId: id,
		Uuid:      uuid,
	}
	err := r.data.DB(ctx).Select("ArticleId", "Uuid").Create(as).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a article statistic: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) CreateArticleCache(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZAddNX(ctx, "article", &redis.Z{
			Score:  float64(id),
			Member: ids + "%" + uuid,
		})
		pipe.ZAddNX(ctx, "article_hot", &redis.Z{
			Score:  0,
			Member: ids + "%" + uuid,
		})
		pipe.ZAddNX(ctx, "leaderboard", &redis.Z{
			Score:  0,
			Member: ids + "%" + uuid + "%article",
		})
		pipe.HSetNX(ctx, "article_"+ids, "uuid", uuid)
		pipe.HSetNX(ctx, "article_"+ids, "agree", 0)
		pipe.HSetNX(ctx, "article_"+ids, "collect", 0)
		pipe.HSetNX(ctx, "article_"+ids, "view", 0)
		pipe.HSetNX(ctx, "article_"+ids, "comment", 0)
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create article cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) CreateArticleSearch(ctx context.Context, id int32, uuid string) error {
	return nil
}

func (r *articleRepo) DeleteArticleDraft(ctx context.Context, id int32, uuid string) error {
	ad := &ArticleDraft{}
	ad.ID = uint(id)
	err := r.data.DB(ctx).Where("uuid = ?", uuid).Delete(ad).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete article draft: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) ArticleDraftMark(ctx context.Context, id int32, uuid string) error {
	err := r.data.db.WithContext(ctx).Model(&ArticleDraft{}).Where("id = ? and uuid = ?", id, uuid).Update("status", 3).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 3: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) SendArticle(ctx context.Context, id int32, uuid string) (*biz.ArticleDraft, error) {
	ad := &ArticleDraft{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&ArticleDraft{}).Where("id = ? and uuid = ? and status = ?", id, uuid, 3).Updates(ad).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 2: uuid(%s), id(%v)", uuid, id))
	}
	return &biz.ArticleDraft{
		Uuid: uuid,
		Id:   id,
	}, nil
}

func (r *articleRepo) SendDraftToMq(ctx context.Context, draft *biz.ArticleDraft) error {
	data, err := json.Marshal(draft)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "article_draft",
		Body:  data,
	}
	msg.WithKeys([]string{draft.Uuid})
	_, err = r.data.articleDraftMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send draft to mq: %v", err))
	}
	return nil
}

func (r *articleRepo) SendArticleToMq(ctx context.Context, article *biz.Article, mode string) error {
	articleMap := map[string]interface{}{}
	articleMap["uuid"] = article.Uuid
	articleMap["id"] = article.ArticleId
	articleMap["mode"] = mode

	data, err := json.Marshal(articleMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "article",
		Body:  data,
	}
	msg.WithKeys([]string{article.Uuid})
	_, err = r.data.articleDraftMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send article to mq: %v", article))
	}
	return nil
}

func (r *articleRepo) SetArticleAgree(ctx context.Context, id int32, uuid string) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and uuid = ?", id, uuid).Update("agree", gorm.Expr("agree + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article agree: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) SetArticleAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "article_"+ids, "agree", 1)
		pipe.ZIncrBy(ctx, "article_hot", 1, ids+"%"+uuid)
		pipe.ZIncrBy(ctx, "leaderboard", 1, ids+"%"+uuid+"%article")
		pipe.SAdd(ctx, "article_agree_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add article agree to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *articleRepo) SetArticleView(ctx context.Context, id int32, uuid string) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and uuid = ?", id, uuid).Update("view", gorm.Expr("view + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article view: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) SetArticleViewToCache(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "article_"+ids, "view", 1)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add article agree to cache: id(%v), uuid(%s)", id, uuid)
	}
	return nil
}

func (r *articleRepo) SetArticleUserCollect(ctx context.Context, id, collectionsId int32, userUuid string) error {
	collect := &Collect{
		CollectionsId: collectionsId,
		Uuid:          userUuid,
		CreationsId:   id,
		Mode:          1,
	}
	err := r.data.DB(ctx).Create(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to collect an article: article_id(%v), collectionsId(%v), userUuid(%s)", id, collectionsId, userUuid))
	}
	return nil
}

func (r *articleRepo) SetArticleCollect(ctx context.Context, id int32, uuid string) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and uuid = ?", id, uuid).Update("collect", gorm.Expr("collect + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article collect: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) SetArticleCollectToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "article_"+ids, "collect", 1)
		pipe.SAdd(ctx, "article_collect_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add article collect to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *articleRepo) CancelArticleAgree(ctx context.Context, id int32, uuid string) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and uuid = ?", id, uuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article agree: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) CancelArticleAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "article_"+ids, "agree", -1)
		pipe.ZIncrBy(ctx, "article_hot", -1, ids+"%"+uuid)
		pipe.ZIncrBy(ctx, "leaderboard", -1, ids+"%"+uuid+"%article")
		pipe.SRem(ctx, "article_agree_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to cancel article agree from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *articleRepo) CancelArticleUserCollect(ctx context.Context, id int32, userUuid string) error {
	collect := &Collect{}
	err := r.data.DB(ctx).Where("creations_id = ? and uuid = ?", id, userUuid).Delete(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article collect: article_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) CancelArticleCollect(ctx context.Context, id int32, uuid string) error {
	as := &ArticleStatistic{}
	err := r.data.DB(ctx).Model(as).Where("article_id = ? and uuid = ?", id, uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article collect: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) CancelArticleCollectFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "article_"+ids, "collect", -1)
		pipe.SRem(ctx, "article_collect_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to cancel article collect from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *articleRepo) SendArticleStatisticToMq(ctx context.Context, uuid, mode string) error {
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
		return errors.Wrapf(err, fmt.Sprintf("fail to send article statistic to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *articleRepo) getArticleFromCache(ctx context.Context, page int32) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "article", index*10, index+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article from cache: key(%s), page(%v)", "article", page))
	}

	article := make([]*biz.Article, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		article = append(article, &biz.Article{
			ArticleId: int32(id),
			Uuid:      member[1],
		})
	}
	return article, nil
}

func (r *articleRepo) getArticleHotFromCache(ctx context.Context, page int32) ([]*biz.ArticleStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "article_hot", index*10, index+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article hot from cache: key(%s), page(%v)", "article_hot", page))
	}

	article := make([]*biz.ArticleStatistic, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		article = append(article, &biz.ArticleStatistic{
			ArticleId: int32(id),
			Uuid:      member[1],
		})
	}
	return article, nil
}

func (r *articleRepo) setArticleToCache(key string, article []*biz.Article) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(item.ArticleId),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article to cache: article(%v)", article)
	}
}

func (r *articleRepo) setArticleHotToCache(key string, article []*biz.ArticleStatistic) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(item.Agree),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article to cache: article(%v)", article)
	}
}

func (r *articleRepo) getArticleStatisticFromCache(ctx context.Context, key string) (*biz.ArticleStatistic, error) {
	statistic, err := r.data.redisCli.HMGet(ctx, key, "uuid", "agree", "collect", "view", "comment").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article statistic form cache: key(%s)", key))
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
	return &biz.ArticleStatistic{
		Uuid:    statistic[0].(string),
		Agree:   val[0],
		Collect: val[1],
		View:    val[2],
		Comment: val[3],
	}, nil
}

func (r *articleRepo) setArticleStatisticToCache(key string, statistic *biz.ArticleStatistic) {
	err := r.data.redisCli.HMSet(context.Background(), key, "uuid", statistic.Uuid, "agree", statistic.Agree, "collect", statistic.Collect, "view", statistic.View, "comment", statistic.Comment).Err()
	if err != nil {
		r.log.Errorf("fail to set article statistic to cache, err(%s)", err.Error())
	}
}
