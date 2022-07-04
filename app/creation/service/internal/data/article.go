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
	err := r.data.db.Where("uuid = ?", uuid).Last(draft).Error
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
	article, err := r.getArticleFromCache(ctx, "article", page)
	if err != nil {
		return nil, err
	}

	size := len(article)
	if size != 0 {
		return article, nil
	}
	return nil, nil
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

func (r *articleRepo) CreateArticle(ctx context.Context, uuid string, id int32) error {
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

func (r *articleRepo) CreateArticleStatistic(ctx context.Context, uuid string, id int32) error {
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

func (r *articleRepo) CreateArticleCache(ctx context.Context, uuid string, id int32) error {
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

func (r *articleRepo) CreateArticleSearch(ctx context.Context, uuid string, id int32) error {
	return nil
}

func (r *articleRepo) DeleteArticleDraft(ctx context.Context, uuid string, id int32) error {
	ad := &ArticleDraft{}
	ad.ID = uint(id)
	err := r.data.DB(ctx).Where("uuid = ?", uuid).Delete(ad).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete article draft: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) ArticleDraftMark(ctx context.Context, uuid string, id int32) error {
	err := r.data.db.WithContext(ctx).Model(&ArticleDraft{}).Where("id = ? and uuid = ?", id, uuid).Update("status", 3).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 3: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) SendArticle(ctx context.Context, uuid string, id int32) (*biz.ArticleDraft, error) {
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

func (r *articleRepo) getArticleFromCache(ctx context.Context, key string, page int32) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index, index+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article from cache: key(%s), page(%v)", key, page))
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
