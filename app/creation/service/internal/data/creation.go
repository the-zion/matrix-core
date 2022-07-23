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
	err := r.data.db.WithContext(ctx).Where("collections_id = ? and mode = ? and status = ?", id, 1, 1).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
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
	err := r.data.db.WithContext(ctx).Model(&Collect{}).Where("collections_id = ? and mode = ? and status = ?", id, 1, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collect article count from db: id(%v)", id))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCollectTalk(ctx context.Context, id, page int32) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collect, 0)
	err := r.data.db.WithContext(ctx).Where("collections_id = ? and mode = ? and status = ?", id, 2, 1).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collect talk from db: id(%v), page(%v)", id, page))
	}

	talk := make([]*biz.Talk, 0)
	for _, item := range list {
		talk = append(talk, &biz.Talk{
			TalkId: item.CreationsId,
			Uuid:   item.Uuid,
		})
	}
	return talk, nil
}

func (r *creationRepo) GetCollectTalkCount(ctx context.Context, id int32) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collect{}).Where("collections_id = ? and mode = ? and status = ?", id, 2, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collect talk count from db: id(%v)", id))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCollectColumn(ctx context.Context, id, page int32) ([]*biz.Column, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collect, 0)
	err := r.data.db.WithContext(ctx).Where("collections_id = ? and mode = ? and status = ?", id, 3, 1).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collect column from db: id(%v), page(%v)", id, page))
	}

	column := make([]*biz.Column, 0)
	for _, item := range list {
		column = append(column, &biz.Column{
			ColumnId: item.CreationsId,
			Uuid:     item.Uuid,
		})
	}
	return column, nil
}

func (r *creationRepo) GetCollectColumnCount(ctx context.Context, id int32) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collect{}).Where("collections_id = ? and mode = ? and status = ?", id, 3, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collect column count from db: id(%v)", id))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCollection(ctx context.Context, id int32, uuid string) (*biz.Collections, error) {
	collections := &Collections{}
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(collections).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collection from db: id(%v)", id))
	}
	if collections.Auth == 2 && collections.Uuid != uuid {
		return nil, errors.Errorf("fail to get collection: no auth")
	}
	return &biz.Collections{
		Uuid:      collections.Uuid,
		Name:      collections.Name,
		Introduce: collections.Introduce,
		Auth:      collections.Auth,
	}, nil
}

func (r *creationRepo) GetCollectCount(ctx context.Context, id int32) (int64, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collect{}).Where("collections_id = ? and status = ?", id, 1).Limit(1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collect count from db: collections_id(%v)", id))
	}
	return count, nil
}

func (r *creationRepo) GetCollections(ctx context.Context, uuid string, page int32) ([]*biz.Collections, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collections, 0)
	handle := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Order("id desc").Offset(index * 10).Limit(10).Find(&list)
	err := handle.Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collections from db: uuid(%s)", uuid))
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
	handle := r.data.db.WithContext(ctx).Where("uuid = ? and auth = ?", uuid, 1).Order("id desc").Offset(index * 10).Limit(10).Find(&list)
	err := handle.Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collections visitor from db: uuid(%s)", uuid))
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
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collections count from db: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collections{}).Where("uuid = ? and auth = ?", uuid, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collections visitor count from db: uuid(%s)", uuid))
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
		return errors.Wrapf(err, fmt.Sprintf("fail to create an collection from db: uuid(%v)", uuid))
	}
	return nil
}

func (r *creationRepo) EditCollections(ctx context.Context, id int32, uuid, name, introduce string, auth int32) error {
	collect := &Collections{
		Name:      name,
		Introduce: introduce,
		Auth:      auth,
	}
	err := r.data.db.WithContext(ctx).Model(&Collections{}).Where("id = ? and uuid = ?", id, uuid).Updates(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to edit an collection from db: id(%v), uuid(%v)", id, uuid))
	}
	return nil
}

func (r *creationRepo) DeleteCollections(ctx context.Context, id int32, uuid string) error {
	collect := &Collections{}
	err := r.data.DB(ctx).Where("id = ? and uuid = ?", id, uuid).Delete(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete an collection from db: id(%v), uuid(%v)", id, uuid))
	}
	return nil
}
