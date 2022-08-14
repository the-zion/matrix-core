package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	creationV1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ biz.ArticleRepo = (*articleRepo)(nil)
var _ biz.TalkRepo = (*talkRepo)(nil)
var _ biz.CreationRepo = (*creationRepo)(nil)
var _ biz.ColumnRepo = (*columnRepo)(nil)
var _ biz.NewsRepo = (*newsRepo)(nil)

type articleRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

type talkRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

type creationRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

type columnRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

type newsRepo struct {
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

func NewTalkRepo(data *Data, logger log.Logger) biz.TalkRepo {
	return &talkRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/talk")),
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

func NewColumnRepo(data *Data, logger log.Logger) biz.ColumnRepo {
	return &columnRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/column")),
		sg:   &singleflight.Group{},
	}
}

func NewNewsRepo(data *Data, logger log.Logger) biz.NewsRepo {
	return &newsRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/news")),
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

func (r *creationRepo) GetCollectArticle(ctx context.Context, id, page int32) ([]*biz.Article, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("collect_article_page_%v_%v", id, page), func() (interface{}, error) {
		reply := make([]*biz.Article, 0)
		articleList, err := r.data.cc.GetCollectArticle(ctx, &creationV1.GetCollectArticleReq{
			Id:   id,
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

func (r *creationRepo) GetCollectArticleCount(ctx context.Context, id int32) (int32, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("collect_article_count_%v", id), func() (interface{}, error) {
		reply, err := r.data.cc.GetCollectArticleCount(ctx, &creationV1.GetCollectArticleCountReq{
			Id: id,
		})
		if err != nil {
			return nil, err
		}
		return reply.Count, nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func (r *creationRepo) GetCollectTalk(ctx context.Context, id, page int32) ([]*biz.Talk, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("collect_talk_page_%v_%v", id, page), func() (interface{}, error) {
		reply := make([]*biz.Talk, 0)
		talkList, err := r.data.cc.GetCollectTalk(ctx, &creationV1.GetCollectTalkReq{
			Id:   id,
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range talkList.Talk {
			reply = append(reply, &biz.Talk{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Talk), nil
}

func (r *creationRepo) GetCollectTalkCount(ctx context.Context, id int32) (int32, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("collect_talk_count_%v", id), func() (interface{}, error) {
		reply, err := r.data.cc.GetCollectTalkCount(ctx, &creationV1.GetCollectTalkCountReq{
			Id: id,
		})
		if err != nil {
			return nil, err
		}
		return reply.Count, nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func (r *creationRepo) GetCollectColumn(ctx context.Context, id, page int32) ([]*biz.Column, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("collect_column_page_%v_%v", id, page), func() (interface{}, error) {
		reply := make([]*biz.Column, 0)
		columnList, err := r.data.cc.GetCollectColumn(ctx, &creationV1.GetCollectColumnReq{
			Id:   id,
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range columnList.Column {
			reply = append(reply, &biz.Column{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Column), nil
}

func (r *creationRepo) GetCollectColumnCount(ctx context.Context, id int32) (int32, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("collect_column_count_%v", id), func() (interface{}, error) {
		reply, err := r.data.cc.GetCollectColumnCount(ctx, &creationV1.GetCollectColumnCountReq{
			Id: id,
		})
		if err != nil {
			return nil, err
		}
		return reply.Count, nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func (r *creationRepo) GetCollection(ctx context.Context, id int32, uuid string) (*biz.Collections, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_collection_%v", id), func() (interface{}, error) {
		reply, err := r.data.cc.GetCollection(ctx, &creationV1.GetCollectionReq{
			Id:   id,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return &biz.Collections{
			Uuid:      reply.Uuid,
			Name:      reply.Name,
			Introduce: reply.Introduce,
			Auth:      reply.Auth,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.Collections), nil
}

func (r *creationRepo) GetCollectionListInfo(ctx context.Context, collectionsList []*biz.Collections) ([]*biz.Collections, error) {
	ids := make([]int32, 0)
	for _, item := range collectionsList {
		ids = append(ids, item.Id)
	}
	reply := make([]*biz.Collections, 0)
	collectionsListInfo, err := r.data.cc.GetCollectionListInfo(ctx, &creationV1.GetCollectionListInfoReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range collectionsListInfo.Collections {
		reply = append(reply, &biz.Collections{
			Id:        item.Id,
			Name:      item.Name,
			Introduce: item.Introduce,
		})
	}
	return reply, nil
}

func (r *creationRepo) GetCollections(ctx context.Context, uuid string, page int32) ([]*biz.Collections, error) {
	collections := make([]*biz.Collections, 0)
	reply, err := r.data.cc.GetCollections(ctx, &creationV1.GetCollectionsReq{
		Uuid: uuid,
		Page: page,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range reply.Collections {
		collections = append(collections, &biz.Collections{
			Id:        item.Id,
			Name:      item.Name,
			Introduce: item.Introduce,
		})
	}
	return collections, nil
}

func (r *creationRepo) GetCollectionsByVisitor(ctx context.Context, uuid string, page int32) ([]*biz.Collections, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_collections_by_visitor_%s_%v", uuid, page), func() (interface{}, error) {
		collections := make([]*biz.Collections, 0)
		reply, err := r.data.cc.GetCollectionsByVisitor(ctx, &creationV1.GetCollectionsReq{
			Uuid: uuid,
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range reply.Collections {
			collections = append(collections, &biz.Collections{
				Id:        item.Id,
				Name:      item.Name,
				Introduce: item.Introduce,
			})
		}
		return collections, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Collections), nil
}

func (r *creationRepo) GetCollectionsCount(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.cc.GetCollectionsCount(ctx, &creationV1.GetCollectionsCountReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Count, nil
}

func (r *creationRepo) GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_collections_visitor_count_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetCollectionsVisitorCount(ctx, &creationV1.GetCollectionsCountReq{
			Uuid: uuid,
		})
		if err != nil {
			return 0, err
		}
		return reply.Count, nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func (r *creationRepo) GetCreationUser(ctx context.Context, uuid string) (*biz.CreationUser, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_creation_user_%s", uuid), func() (interface{}, error) {
		creation := &biz.CreationUser{}
		reply, err := r.data.cc.GetCreationUser(ctx, &creationV1.GetCreationUserReq{
			Uuid: uuid,
		})
		if err != nil {
			return 0, err
		}
		creation.Article = reply.Article
		creation.Talk = reply.Talk
		creation.Collections = reply.Collections
		creation.Column = reply.Column
		return creation, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.CreationUser), nil
}

func (r *creationRepo) GetCreationUserVisitor(ctx context.Context, uuid string) (*biz.CreationUser, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_creation_user_visitor_%s", uuid), func() (interface{}, error) {
		creation := &biz.CreationUser{}
		reply, err := r.data.cc.GetCreationUserVisitor(ctx, &creationV1.GetCreationUserReq{
			Uuid: uuid,
		})
		if err != nil {
			return 0, err
		}
		creation.Article = reply.Article
		creation.Talk = reply.Talk
		creation.Collections = reply.Collections
		creation.Column = reply.Column
		return creation, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.CreationUser), nil
}

func (r *creationRepo) CreateCollections(ctx context.Context, uuid, name, introduce string, auth int32) error {
	_, err := r.data.cc.CreateCollections(ctx, &creationV1.CreateCollectionsReq{
		Uuid:      uuid,
		Name:      name,
		Introduce: introduce,
		Auth:      auth,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) EditCollections(ctx context.Context, id int32, uuid, name, introduce string, auth int32) error {
	_, err := r.data.cc.EditCollections(ctx, &creationV1.EditCollectionsReq{
		Id:        id,
		Uuid:      uuid,
		Name:      name,
		Introduce: introduce,
		Auth:      auth,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) DeleteCollections(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.DeleteCollections(ctx, &creationV1.DeleteCollectionsReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
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
	result, err, _ := r.sg.Do(fmt.Sprintf("article_page_%v", page), func() (interface{}, error) {
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
	result, err, _ := r.sg.Do(fmt.Sprintf("article_page_hot_%v", page), func() (interface{}, error) {
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

func (r *articleRepo) GetColumnArticleList(ctx context.Context, id int32) ([]*biz.Article, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_column_article_list_%v", id), func() (interface{}, error) {
		reply := make([]*biz.Article, 0)
		articleList, err := r.data.cc.GetColumnArticleList(ctx, &creationV1.GetColumnArticleListReq{
			Id: id,
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

func (r *articleRepo) GetArticleCount(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.cc.GetArticleCount(ctx, &creationV1.GetArticleCountReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Count, nil
}

func (r *articleRepo) GetArticleCountVisitor(ctx context.Context, uuid string) (int32, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_article_visitor_count_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetArticleCountVisitor(ctx, &creationV1.GetArticleCountVisitorReq{
			Uuid: uuid,
		})
		if err != nil {
			return 0, err
		}
		return reply.Count, nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func (r *articleRepo) GetUserArticleList(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	reply := make([]*biz.Article, 0)
	articleList, err := r.data.cc.GetUserArticleList(ctx, &creationV1.GetUserArticleListReq{
		Page: page,
		Uuid: uuid,
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
}

func (r *articleRepo) GetUserArticleListVisitor(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_article_by_visitor_%s_%v", uuid, page), func() (interface{}, error) {
		reply := make([]*biz.Article, 0)
		articleList, err := r.data.cc.GetUserArticleListVisitor(ctx, &creationV1.GetUserArticleListVisitorReq{
			Page: page,
			Uuid: uuid,
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
	result, err, _ := r.sg.Do(fmt.Sprintf("article_statistic_%v", id), func() (interface{}, error) {
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

func (r *articleRepo) GetArticleListStatistic(ctx context.Context, articleList []*biz.Article) ([]*biz.ArticleStatistic, error) {
	ids := make([]int32, 0)
	for _, item := range articleList {
		ids = append(ids, item.Id)
	}
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

func (r *articleRepo) GetArticleSearch(ctx context.Context, page int32, search, time string) ([]*biz.ArticleSearch, int32, error) {
	reply := make([]*biz.ArticleSearch, 0)
	searchReply, err := r.data.cc.GetArticleSearch(ctx, &creationV1.GetArticleSearchReq{
		Page:   page,
		Search: search,
		Time:   time,
	})
	if err != nil {
		return nil, 0, err
	}
	for _, item := range searchReply.List {
		reply = append(reply, &biz.ArticleSearch{
			Id:     item.Id,
			Tags:   item.Tags,
			Title:  item.Title,
			Text:   item.Text,
			Uuid:   item.Uuid,
			Cover:  item.Cover,
			Update: item.Update,
		})
	}
	return reply, searchReply.Total, nil
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

func (r *articleRepo) ArticleDraftMark(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.ArticleDraftMark(ctx, &creationV1.ArticleDraftMarkReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) SendArticle(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.cc.SendArticle(ctx, &creationV1.SendArticleReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) SendArticleEdit(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.cc.SendArticleEdit(ctx, &creationV1.SendArticleEditReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) DeleteArticle(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.DeleteArticle(ctx, &creationV1.DeleteArticleReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) SetArticleAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetArticleAgree(ctx, &creationV1.SetArticleAgreeReq{
		Uuid:     uuid,
		Id:       id,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) SetArticleView(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.SetArticleView(ctx, &creationV1.SetArticleViewReq{
		Uuid: uuid,
		Id:   id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) SetArticleCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetArticleCollect(ctx, &creationV1.SetArticleCollectReq{
		CollectionsId: collectionsId,
		Uuid:          uuid,
		UserUuid:      userUuid,
		Id:            id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) CancelArticleAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelArticleAgree(ctx, &creationV1.CancelArticleAgreeReq{
		Uuid:     uuid,
		Id:       id,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) CancelArticleCollect(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelArticleCollect(ctx, &creationV1.CancelArticleCollectReq{
		Uuid:     uuid,
		UserUuid: userUuid,
		Id:       id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) ArticleStatisticJudge(ctx context.Context, id int32, uuid string) (*biz.ArticleStatisticJudge, error) {
	reply, err := r.data.cc.ArticleStatisticJudge(ctx, &creationV1.ArticleStatisticJudgeReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return &biz.ArticleStatisticJudge{
		Agree:   reply.Agree,
		Collect: reply.Collect,
	}, nil
}

func (r *talkRepo) GetTalkList(ctx context.Context, page int32) ([]*biz.Talk, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("talk_page_%v", page), func() (interface{}, error) {
		reply := make([]*biz.Talk, 0)
		talkList, err := r.data.cc.GetTalkList(ctx, &creationV1.GetTalkListReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range talkList.Talk {
			reply = append(reply, &biz.Talk{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Talk), nil
}

func (r *talkRepo) GetTalkListHot(ctx context.Context, page int32) ([]*biz.Talk, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("talk_page_hot_%v", page), func() (interface{}, error) {
		reply := make([]*biz.Talk, 0)
		talkList, err := r.data.cc.GetTalkListHot(ctx, &creationV1.GetTalkListHotReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range talkList.Talk {
			reply = append(reply, &biz.Talk{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Talk), nil
}

func (r *talkRepo) GetUserTalkList(ctx context.Context, page int32, uuid string) ([]*biz.Talk, error) {
	reply := make([]*biz.Talk, 0)
	talkList, err := r.data.cc.GetUserTalkList(ctx, &creationV1.GetUserTalkListReq{
		Page: page,
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range talkList.Talk {
		reply = append(reply, &biz.Talk{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (r *talkRepo) GetUserTalkListVisitor(ctx context.Context, page int32, uuid string) ([]*biz.Talk, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_talk_by_visitor_%s_%v", uuid, page), func() (interface{}, error) {
		reply := make([]*biz.Talk, 0)
		talkList, err := r.data.cc.GetUserTalkListVisitor(ctx, &creationV1.GetUserTalkListVisitorReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range talkList.Talk {
			reply = append(reply, &biz.Talk{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Talk), nil
}

func (r *talkRepo) GetTalkCount(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.cc.GetTalkCount(ctx, &creationV1.GetTalkCountReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Count, nil
}

func (r *talkRepo) GetTalkCountVisitor(ctx context.Context, uuid string) (int32, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_talk_visitor_count_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetTalkCountVisitor(ctx, &creationV1.GetTalkCountVisitorReq{
			Uuid: uuid,
		})
		if err != nil {
			return 0, err
		}
		return reply.Count, nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func (r *talkRepo) GetTalkListStatistic(ctx context.Context, talkList []*biz.Talk) ([]*biz.TalkStatistic, error) {
	ids := make([]int32, 0)
	for _, item := range talkList {
		ids = append(ids, item.Id)
	}
	reply := make([]*biz.TalkStatistic, 0)
	statisticList, err := r.data.cc.GetTalkListStatistic(ctx, &creationV1.GetTalkListStatisticReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range statisticList.Count {
		reply = append(reply, &biz.TalkStatistic{
			Id:      item.Id,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (r *talkRepo) GetTalkStatistic(ctx context.Context, id int32) (*biz.TalkStatistic, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("talk_statistic_%v", id), func() (interface{}, error) {
		statistic, err := r.data.cc.GetTalkStatistic(ctx, &creationV1.GetTalkStatisticReq{
			Id: id,
		})
		if err != nil {
			return nil, err
		}
		return &biz.TalkStatistic{
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
	return result.(*biz.TalkStatistic), nil
}

func (r *talkRepo) GetLastTalkDraft(ctx context.Context, uuid string) (*biz.TalkDraft, error) {
	reply, err := r.data.cc.GetLastTalkDraft(ctx, &creationV1.GetLastTalkDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return &biz.TalkDraft{
		Id:     reply.Id,
		Status: reply.Status,
	}, nil
}

func (r *talkRepo) GetTalkSearch(ctx context.Context, page int32, search, time string) ([]*biz.TalkSearch, int32, error) {
	reply := make([]*biz.TalkSearch, 0)
	searchReply, err := r.data.cc.GetTalkSearch(ctx, &creationV1.GetTalkSearchReq{
		Page:   page,
		Search: search,
		Time:   time,
	})
	if err != nil {
		return nil, 0, err
	}
	for _, item := range searchReply.List {
		reply = append(reply, &biz.TalkSearch{
			Id:     item.Id,
			Tags:   item.Tags,
			Title:  item.Title,
			Text:   item.Text,
			Uuid:   item.Uuid,
			Cover:  item.Cover,
			Update: item.Update,
		})
	}
	return reply, searchReply.Total, nil
}

func (r *talkRepo) CreateTalkDraft(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.cc.CreateTalkDraft(ctx, &creationV1.CreateTalkDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Id, nil
}

func (r *talkRepo) SendTalk(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.cc.SendTalk(ctx, &creationV1.SendTalkReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *talkRepo) SendTalkEdit(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.cc.SendTalkEdit(ctx, &creationV1.SendTalkEditReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *talkRepo) DeleteTalk(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.DeleteTalk(ctx, &creationV1.DeleteTalkReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *talkRepo) SetTalkAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetTalkAgree(ctx, &creationV1.SetTalkAgreeReq{
		Uuid:     uuid,
		Id:       id,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *talkRepo) CancelTalkAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelTalkAgree(ctx, &creationV1.CancelTalkAgreeReq{
		Uuid:     uuid,
		Id:       id,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *talkRepo) CancelTalkCollect(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelTalkCollect(ctx, &creationV1.CancelTalkCollectReq{
		Uuid:     uuid,
		UserUuid: userUuid,
		Id:       id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *talkRepo) SetTalkView(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.SetTalkView(ctx, &creationV1.SetTalkViewReq{
		Uuid: uuid,
		Id:   id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *talkRepo) SetTalkCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetTalkCollect(ctx, &creationV1.SetTalkCollectReq{
		CollectionsId: collectionsId,
		Uuid:          uuid,
		UserUuid:      userUuid,
		Id:            id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *talkRepo) TalkStatisticJudge(ctx context.Context, id int32, uuid string) (*biz.TalkStatisticJudge, error) {
	reply, err := r.data.cc.TalkStatisticJudge(ctx, &creationV1.TalkStatisticJudgeReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return &biz.TalkStatisticJudge{
		Agree:   reply.Agree,
		Collect: reply.Collect,
	}, nil
}

func (r *columnRepo) GetLastColumnDraft(ctx context.Context, uuid string) (*biz.ColumnDraft, error) {
	reply, err := r.data.cc.GetLastColumnDraft(ctx, &creationV1.GetLastColumnDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return &biz.ColumnDraft{
		Id:     reply.Id,
		Status: reply.Status,
	}, nil
}

func (r *columnRepo) CreateColumnDraft(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.cc.CreateColumnDraft(ctx, &creationV1.CreateColumnDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Id, nil
}

func (r *columnRepo) SubscribeColumn(ctx context.Context, id int32, author, uuid string) error {
	_, err := r.data.cc.SubscribeColumn(ctx, &creationV1.SubscribeColumnReq{
		Id:     id,
		Author: author,
		Uuid:   uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) CancelSubscribeColumn(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.CancelSubscribeColumn(ctx, &creationV1.CancelSubscribeColumnReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) SubscribeJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	reply, err := r.data.cc.SubscribeJudge(ctx, &creationV1.SubscribeJudgeReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return false, err
	}
	return reply.Subscribe, nil
}

func (r *columnRepo) GetColumnList(ctx context.Context, page int32) ([]*biz.Column, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("column_page_%v", page), func() (interface{}, error) {
		reply := make([]*biz.Column, 0)
		columnList, err := r.data.cc.GetColumnList(ctx, &creationV1.GetColumnListReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range columnList.Column {
			reply = append(reply, &biz.Column{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Column), nil
}

func (r *columnRepo) GetColumnListHot(ctx context.Context, page int32) ([]*biz.Column, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("column_page_hot_%v", page), func() (interface{}, error) {
		reply := make([]*biz.Column, 0)
		columnList, err := r.data.cc.GetColumnListHot(ctx, &creationV1.GetColumnListHotReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range columnList.Column {
			reply = append(reply, &biz.Column{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Column), nil
}

func (r *columnRepo) GetUserColumnList(ctx context.Context, page int32, uuid string) ([]*biz.Column, error) {
	reply := make([]*biz.Column, 0)
	columnList, err := r.data.cc.GetUserColumnList(ctx, &creationV1.GetUserColumnListReq{
		Page: page,
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range columnList.Column {
		reply = append(reply, &biz.Column{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (r *columnRepo) GetUserColumnListVisitor(ctx context.Context, page int32, uuid string) ([]*biz.Column, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_column_by_visitor_%s_%v", uuid, page), func() (interface{}, error) {
		reply := make([]*biz.Column, 0)
		columnList, err := r.data.cc.GetUserColumnListVisitor(ctx, &creationV1.GetUserColumnListVisitorReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range columnList.Column {
			reply = append(reply, &biz.Column{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Column), nil
}

func (r *columnRepo) GetColumnCount(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.cc.GetColumnCount(ctx, &creationV1.GetColumnCountReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Count, nil
}

func (r *columnRepo) GetColumnCountVisitor(ctx context.Context, uuid string) (int32, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_column_visitor_count_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetColumnCountVisitor(ctx, &creationV1.GetColumnCountVisitorReq{
			Uuid: uuid,
		})
		if err != nil {
			return 0, err
		}
		return reply.Count, nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func (r *columnRepo) GetColumnListStatistic(ctx context.Context, columnList []*biz.Column) ([]*biz.ColumnStatistic, error) {
	ids := make([]int32, 0)
	for _, item := range columnList {
		ids = append(ids, item.Id)
	}
	reply := make([]*biz.ColumnStatistic, 0)
	statisticList, err := r.data.cc.GetColumnListStatistic(ctx, &creationV1.GetColumnListStatisticReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range statisticList.Count {
		reply = append(reply, &biz.ColumnStatistic{
			Id:      item.Id,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
		})
	}
	return reply, nil
}

func (r *columnRepo) GetColumnStatistic(ctx context.Context, id int32) (*biz.ColumnStatistic, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("column_statistic_%v", id), func() (interface{}, error) {
		statistic, err := r.data.cc.GetColumnStatistic(ctx, &creationV1.GetColumnStatisticReq{
			Id: id,
		})
		if err != nil {
			return nil, err
		}
		return &biz.ColumnStatistic{
			Uuid:    statistic.Uuid,
			Agree:   statistic.Agree,
			Collect: statistic.Collect,
			View:    statistic.View,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.ColumnStatistic), nil
}

func (r *columnRepo) GetSubscribeList(ctx context.Context, page int32, uuid string) ([]*biz.Subscribe, error) {
	reply := make([]*biz.Subscribe, 0)
	subscribeList, err := r.data.cc.GetSubscribeList(ctx, &creationV1.GetSubscribeListReq{
		Page: page,
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range subscribeList.Subscribe {
		reply = append(reply, &biz.Subscribe{
			ColumnId: item.Id,
			AuthorId: item.Uuid,
		})
	}
	return reply, nil
}

func (r *columnRepo) GetSubscribeListCount(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.cc.GetSubscribeListCount(ctx, &creationV1.GetSubscribeListCountReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Count, nil
}

func (r *columnRepo) GetColumnSubscribes(ctx context.Context, uuid string, ids []int32) ([]*biz.Subscribe, error) {
	reply := make([]*biz.Subscribe, 0)
	subscribeList, err := r.data.cc.GetColumnSubscribes(ctx, &creationV1.GetColumnSubscribesReq{
		Ids:  ids,
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range subscribeList.Subscribes {
		reply = append(reply, &biz.Subscribe{
			ColumnId: item.Id,
			Status:   item.SubscribeJudge,
		})
	}
	return reply, nil
}

func (r *columnRepo) GetColumnSearch(ctx context.Context, page int32, search, time string) ([]*biz.ColumnSearch, int32, error) {
	reply := make([]*biz.ColumnSearch, 0)
	searchReply, err := r.data.cc.GetColumnSearch(ctx, &creationV1.GetColumnSearchReq{
		Page:   page,
		Search: search,
		Time:   time,
	})
	if err != nil {
		return nil, 0, err
	}
	for _, item := range searchReply.List {
		reply = append(reply, &biz.ColumnSearch{
			Id:        item.Id,
			Tags:      item.Tags,
			Name:      item.Name,
			Introduce: item.Introduce,
			Uuid:      item.Uuid,
			Cover:     item.Cover,
			Update:    item.Update,
		})
	}
	return reply, searchReply.Total, nil
}

func (r *columnRepo) SendColumn(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.cc.SendColumn(ctx, &creationV1.SendColumnReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) SendColumnEdit(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.cc.SendColumnEdit(ctx, &creationV1.SendColumnEditReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) DeleteColumn(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.DeleteColumn(ctx, &creationV1.DeleteColumnReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) ColumnStatisticJudge(ctx context.Context, id int32, uuid string) (*biz.ColumnStatisticJudge, error) {
	reply, err := r.data.cc.ColumnStatisticJudge(ctx, &creationV1.ColumnStatisticJudgeReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return &biz.ColumnStatisticJudge{
		Agree:   reply.Agree,
		Collect: reply.Collect,
	}, nil
}

func (r *columnRepo) SetColumnAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetColumnAgree(ctx, &creationV1.SetColumnAgreeReq{
		Uuid:     uuid,
		Id:       id,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) CancelColumnAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelColumnAgree(ctx, &creationV1.CancelColumnAgreeReq{
		Uuid:     uuid,
		Id:       id,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) SetColumnCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetColumnCollect(ctx, &creationV1.SetColumnCollectReq{
		CollectionsId: collectionsId,
		Uuid:          uuid,
		UserUuid:      userUuid,
		Id:            id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) CancelColumnCollect(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelColumnCollect(ctx, &creationV1.CancelColumnCollectReq{
		Uuid:     uuid,
		UserUuid: userUuid,
		Id:       id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) SetColumnView(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.SetColumnView(ctx, &creationV1.SetColumnViewReq{
		Uuid: uuid,
		Id:   id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) AddColumnIncludes(ctx context.Context, id, articleId int32, uuid string) error {
	_, err := r.data.cc.AddColumnIncludes(ctx, &creationV1.AddColumnIncludesReq{
		Id:        id,
		ArticleId: articleId,
		Uuid:      uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) DeleteColumnIncludes(ctx context.Context, id, articleId int32, uuid string) error {
	_, err := r.data.cc.DeleteColumnIncludes(ctx, &creationV1.DeleteColumnIncludesReq{
		Id:        id,
		ArticleId: articleId,
		Uuid:      uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *newsRepo) GetNews(ctx context.Context, page int32) ([]*biz.News, error) {
	reply := make([]*biz.News, 0)
	newsList, err := r.data.cc.GetNews(ctx, &creationV1.GetNewsReq{
		Page: page,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range newsList.News {
		reply = append(reply, &biz.News{
			Id:     item.Id,
			Update: item.Update,
			Title:  item.Title,
			Author: item.Author,
			Text:   item.Text,
			Tags:   item.Tags,
			Cover:  item.Cover,
			Url:    item.Url,
		})
	}
	return reply, nil
}
