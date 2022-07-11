package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"strconv"
	"strings"
)

var _ biz.CreationRepo = (*creationRepo)(nil)

type creationRepo struct {
	data *Data
	log  *log.Helper
}

func NewCreationRepo(data *Data, logger log.Logger) biz.CreationRepo {
	return &creationRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "creation/data/creation")),
	}
}

func (r *creationRepo) GetLeaderBoard(ctx context.Context) ([]*biz.LeaderBoard, error) {
	return r.getLeaderBoardFromCache(ctx)
}

func (r *creationRepo) getLeaderBoardFromCache(ctx context.Context) ([]*biz.LeaderBoard, error) {
	list, err := r.data.redisCli.ZRevRange(ctx, "leaderboard", 0, 9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "fail to get leader board from cache")
	}

	board := make([]*biz.LeaderBoard, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		board = append(board, &biz.LeaderBoard{
			Id:   int32(id),
			Uuid: member[1],
			Mode: member[2],
		})
	}
	return board, nil
}

func (r *creationRepo) GetCollectArticle(ctx context.Context, id, page int32) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collect, 0)
	err := r.data.db.WithContext(ctx).Where("collections_id = ? and mode = ?", id, 1).Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collect article from db: id(%v), page(%v)", id, page))
	}

	article := make([]*biz.Article, 0)
	for _, item := range list {
		article = append(article, &biz.Article{
			ArticleId: item.CreationsId,
			Uuid:      item.Uuid,
		})
	}
	return article, nil
}

func (r *creationRepo) GetCollectArticleCount(ctx context.Context, id int32) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collect{}).Where("collections_id = ? and mode = ?", id, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collect article count: id(%v)", id))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCollections(ctx context.Context, uuid string, page int32) ([]*biz.Collections, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collections, 0)
	handle := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Offset(index * 10).Limit(10).Find(&list)
	err := handle.Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collections: uuid(%s)", uuid))
	}
	collections := make([]*biz.Collections, 0)
	for _, item := range list {
		collections = append(collections, &biz.Collections{
			Id:        int32(item.ID),
			Name:      item.Name,
			Introduce: item.Introduce,
		})
	}
	return collections, nil
}

func (r *creationRepo) GetCollectionsByVisitor(ctx context.Context, uuid string, page int32) ([]*biz.Collections, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collections, 0)
	handle := r.data.db.WithContext(ctx).Where("uuid = ? and auth = ?", uuid, 1).Offset(index * 10).Limit(10).Find(&list)
	err := handle.Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collections visitor: uuid(%s)", uuid))
	}
	collections := make([]*biz.Collections, 0)
	for _, item := range list {
		collections = append(collections, &biz.Collections{
			Id:        int32(item.ID),
			Name:      item.Name,
			Introduce: item.Introduce,
		})
	}
	return collections, nil
}

func (r *creationRepo) GetCollectionsCount(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collections{}).Where("uuid = ?", uuid).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collections count: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collections{}).Where("uuid = ? and auth = ?", uuid, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collections visitor count: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *creationRepo) CreateCollections(ctx context.Context, uuid, name, introduce string, auth int32) error {
	collect := &Collections{
		Uuid:      uuid,
		Name:      name,
		Introduce: introduce,
		Auth:      auth,
	}
	err := r.data.DB(ctx).Create(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create an collections: uuid(%v)", uuid))
	}
	return nil
}
