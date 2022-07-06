package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	creationV1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

var _ biz.ArticleRepo = (*articleRepo)(nil)
var _ biz.CreationRepo = (*creationRepo)(nil)

type articleRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

type creationRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewArticleRepo(data *Data, logger log.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/article")),
		sg:   &singleflight.Group{},
	}
}

func NewCreationRepo(data *Data, logger log.Logger) biz.CreationRepo {
	return &creationRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/creation")),
		sg:   &singleflight.Group{},
	}
}

func (r *creationRepo) GetLeaderBoard(ctx context.Context) ([]*biz.LeaderBoard, error) {
	result, err, _ := r.sg.Do("leader_board", func() (interface{}, error) {
		replyBoard := make([]*biz.LeaderBoard, 0)
		reply, err := r.data.cc.GetLeaderBoard(ctx, &emptypb.Empty{})
		if err != nil {
			return nil, err
		}

		for _, item := range reply.Board {
			replyBoard = append(replyBoard, &biz.LeaderBoard{
				Id:   item.Id,
				Uuid: item.Uuid,
				Mode: item.Mode,
			})
		}
		return replyBoard, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.LeaderBoard), nil
}

func (r *articleRepo) GetLastArticleDraft(ctx context.Context, uuid string) (*biz.ArticleDraft, error) {
	reply, err := r.data.cc.GetLastArticleDraft(ctx, &creationV1.GetLastArticleDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return &biz.ArticleDraft{
		Id:     reply.Id,
		Status: reply.Status,
	}, nil
}

func (r *articleRepo) GetArticleList(ctx context.Context, page int32) ([]*biz.Article, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("article_page_%s", strconv.Itoa(int(page))), func() (interface{}, error) {
		reply := make([]*biz.Article, 0)
		articleList, err := r.data.cc.GetArticleList(ctx, &creationV1.GetArticleListReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range articleList.Article {
			reply = append(reply, &biz.Article{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Article), nil
}

func (r *articleRepo) GetArticleListHot(ctx context.Context, page int32) ([]*biz.Article, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("article_page_hot_%s", strconv.Itoa(int(page))), func() (interface{}, error) {
		reply := make([]*biz.Article, 0)
		articleList, err := r.data.cc.GetArticleListHot(ctx, &creationV1.GetArticleListHotReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range articleList.Article {
			reply = append(reply, &biz.Article{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Article), nil
}

func (r *articleRepo) GetArticleStatistic(ctx context.Context, id int32) (*biz.ArticleStatistic, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("article_statistic_%s", strconv.Itoa(int(id))), func() (interface{}, error) {
		statistic, err := r.data.cc.GetArticleStatistic(ctx, &creationV1.GetArticleStatisticReq{
			Id: id,
		})
		if err != nil {
			return nil, err
		}
		return &biz.ArticleStatistic{
			Uuid:    statistic.Uuid,
			Agree:   statistic.Agree,
			Collect: statistic.Collect,
			View:    statistic.View,
			Comment: statistic.Comment,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.ArticleStatistic), nil
}

func (r *articleRepo) GetArticleListStatistic(ctx context.Context, ids []int32) ([]*biz.ArticleStatistic, error) {
	reply := make([]*biz.ArticleStatistic, 0)
	statisticList, err := r.data.cc.GetArticleListStatistic(ctx, &creationV1.GetArticleListStatisticReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range statisticList.Count {
		reply = append(reply, &biz.ArticleStatistic{
			Id:      item.Id,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (r *articleRepo) GetArticleDraftList(ctx context.Context, uuid string) ([]*biz.ArticleDraft, error) {
	reply := make([]*biz.ArticleDraft, 0)
	draftList, err := r.data.cc.GetArticleDraftList(ctx, &creationV1.GetArticleDraftListReq{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range draftList.Draft {
		reply = append(reply, &biz.ArticleDraft{
			Id: item.Id,
		})
	}
	return reply, nil
}

func (r *articleRepo) CreateArticleDraft(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.cc.CreateArticleDraft(ctx, &creationV1.CreateArticleDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Id, nil
}

func (r *articleRepo) ArticleDraftMark(ctx context.Context, uuid string, id int32) error {
	_, err := r.data.cc.ArticleDraftMark(ctx, &creationV1.ArticleDraftMarkReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) SendArticle(ctx context.Context, uuid string, id int32) error {
	_, err := r.data.cc.SendArticle(ctx, &creationV1.SendArticleReq{
		Uuid: uuid,
		Id:   id,
	})
	if err != nil {
		return err
	}
	return nil
}
