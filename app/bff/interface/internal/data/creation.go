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
		reply, err := r.data.cc.GetLeaderBoard(ctx, &emptypb.Empty{})
		if err != nil {
			return nil, err
		}
		replyBoard := make([]*biz.LeaderBoard, 0, len(reply.Board))
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

func (r *creationRepo) GetCollectArticleList(ctx context.Context, id, page int32) ([]*biz.Article, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("collect_article_list_page_%v_%v", id, page), func() (interface{}, error) {
		articleList, err := r.data.cc.GetCollectArticleList(ctx, &creationV1.GetCollectArticleListReq{
			Id:   id,
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Article, 0, len(articleList.Article))
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

func (r *creationRepo) GetCollectTalkList(ctx context.Context, id, page int32) ([]*biz.Talk, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("collect_talk_list_page_%v_%v", id, page), func() (interface{}, error) {
		talkList, err := r.data.cc.GetCollectTalkList(ctx, &creationV1.GetCollectTalkListReq{
			Id:   id,
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Talk, 0, len(talkList.Talk))
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

func (r *creationRepo) GetCollectColumnList(ctx context.Context, id, page int32) ([]*biz.Column, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("collect_column_list_page_%v_%v", id, page), func() (interface{}, error) {
		columnList, err := r.data.cc.GetCollectColumnList(ctx, &creationV1.GetCollectColumnListReq{
			Id:   id,
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Column, 0, len(columnList.Column))
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

func (r *creationRepo) GetLastCollectionsDraft(ctx context.Context, uuid string) (*biz.CollectionsDraft, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_last_collections_draft_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetLastCollectionsDraft(ctx, &creationV1.GetLastCollectionsDraftReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return &biz.CollectionsDraft{
			Id:     reply.Id,
			Status: reply.Status,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.CollectionsDraft), nil
}

func (r *creationRepo) GetCollectionsContentReview(ctx context.Context, page int32, uuid string) ([]*biz.CreationContentReview, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_collections_content_review_%s_%v", uuid, page), func() (interface{}, error) {
		reviewReply, err := r.data.cc.GetCollectionsContentReview(ctx, &creationV1.GetCollectionsContentReviewReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.CreationContentReview, 0, len(reviewReply.Review))
		for _, item := range reviewReply.Review {
			reply = append(reply, &biz.CreationContentReview{
				Id:         item.Id,
				CreationId: item.CreationId,
				Title:      item.Title,
				Kind:       item.Kind,
				Uuid:       item.Uuid,
				CreateAt:   item.CreateAt,
				JobId:      item.JobId,
				Label:      item.Label,
				Result:     item.Result,
				Section:    item.Section,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CreationContentReview), nil
}

func (r *creationRepo) GetUserTimeLineListVisitor(ctx context.Context, page int32, uuid string) ([]*biz.TimeLine, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_timeline_%s_%v", uuid, page), func() (interface{}, error) {
		timeline, err := r.data.cc.GetUserTimeLineList(ctx, &creationV1.GetUserTimeLineListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.TimeLine, 0, len(timeline.Timeline))
		for _, item := range timeline.Timeline {
			reply = append(reply, &biz.TimeLine{
				Id:         item.Id,
				Uuid:       item.Uuid,
				CreationId: item.CreationId,
				Mode:       item.Mode,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.TimeLine), nil
}

func (r *creationRepo) GetCollections(ctx context.Context, id int32, uuid string) (*biz.Collections, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_collections_%v", id), func() (interface{}, error) {
		reply, err := r.data.cc.GetCollections(ctx, &creationV1.GetCollectionsReq{
			Id:   id,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return &biz.Collections{
			Uuid:    reply.Uuid,
			Auth:    reply.Auth,
			Article: reply.Article,
			Column:  reply.Column,
			Talk:    reply.Talk,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.Collections), nil
}

func (r *creationRepo) GetCollectionListInfo(ctx context.Context, collectionsList []*biz.Collections) ([]*biz.Collections, error) {
	ids := make([]int32, 0, len(collectionsList))
	for _, item := range collectionsList {
		ids = append(ids, item.Id)
	}
	collectionsListInfo, err := r.data.cc.GetCollectionListInfo(ctx, &creationV1.GetCollectionListInfoReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	reply := make([]*biz.Collections, 0, len(collectionsListInfo.Collections))
	for _, item := range collectionsListInfo.Collections {
		reply = append(reply, &biz.Collections{
			Id: item.Id,
		})
	}
	return reply, nil
}

func (r *creationRepo) GetCollectionsList(ctx context.Context, uuid string, page int32) ([]*biz.Collections, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_collections_list_%s_%v", uuid, page), func() (interface{}, error) {
		reply, err := r.data.cc.GetCollectionsList(ctx, &creationV1.GetCollectionsListReq{
			Uuid: uuid,
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		collections := make([]*biz.Collections, 0, len(reply.Collections))
		for _, item := range reply.Collections {
			collections = append(collections, &biz.Collections{
				Id: item.Id,
			})
		}
		return collections, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Collections), nil
}

func (r *creationRepo) GetCollectionsListAll(ctx context.Context, uuid string) ([]*biz.Collections, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_collections_list_all_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetCollectionsListAll(ctx, &creationV1.GetCollectionsListAllReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		collections := make([]*biz.Collections, 0, len(reply.Collections))
		for _, item := range reply.Collections {
			collections = append(collections, &biz.Collections{
				Id: item.Id,
			})
		}
		return collections, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Collections), nil
}

func (r *creationRepo) GetCollectionsListByVisitor(ctx context.Context, uuid string, page int32) ([]*biz.Collections, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_collections_list_by_visitor_%s_%v", uuid, page), func() (interface{}, error) {
		reply, err := r.data.cc.GetCollectionsListByVisitor(ctx, &creationV1.GetCollectionsListReq{
			Uuid: uuid,
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		collections := make([]*biz.Collections, 0, len(reply.Collections))
		for _, item := range reply.Collections {
			collections = append(collections, &biz.Collections{
				Id: item.Id,
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
		creation.Collect = reply.Collect
		creation.Subscribe = reply.Subscribe
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

func (r *creationRepo) CreateCollectionsDraft(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.cc.CreateCollectionsDraft(ctx, &creationV1.CreateCollectionsDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Id, nil
}

func (r *creationRepo) SendCollections(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.cc.SendCollections(ctx, &creationV1.SendCollectionsReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SendCollectionsEdit(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.cc.SendCollectionsEdit(ctx, &creationV1.SendCollectionsEditReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
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
	result, err, _ := r.sg.Do(fmt.Sprintf("get_last_article_draft_%s", uuid), func() (interface{}, error) {
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
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.ArticleDraft), nil
}

func (r *articleRepo) GetArticleList(ctx context.Context, page int32) ([]*biz.Article, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("article_page_%v", page), func() (interface{}, error) {
		articleList, err := r.data.cc.GetArticleList(ctx, &creationV1.GetArticleListReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Article, 0, len(articleList.Article))
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
		articleList, err := r.data.cc.GetArticleListHot(ctx, &creationV1.GetArticleListHotReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Article, 0, len(articleList.Article))
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

func (r *articleRepo) GetUserArticleListAll(ctx context.Context, uuid string) ([]*biz.Article, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("article_all_%s", uuid), func() (interface{}, error) {
		articleList, err := r.data.cc.GetUserArticleListAll(ctx, &creationV1.GetUserArticleListAllReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Article, 0, len(articleList.Article))
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
		articleList, err := r.data.cc.GetColumnArticleList(ctx, &creationV1.GetColumnArticleListReq{
			Id: id,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Article, 0, len(articleList.Article))
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
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_article_list_%s_%v", uuid, page), func() (interface{}, error) {
		articleList, err := r.data.cc.GetUserArticleList(ctx, &creationV1.GetUserArticleListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Article, 0, len(articleList.Article))
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

func (r *articleRepo) GetUserArticleListVisitor(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_article_by_visitor_%s_%v", uuid, page), func() (interface{}, error) {
		articleList, err := r.data.cc.GetUserArticleListVisitor(ctx, &creationV1.GetUserArticleListVisitorReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Article, 0, len(articleList.Article))
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

func (r *articleRepo) GetArticleStatistic(ctx context.Context, id int32, uuid string) (*biz.ArticleStatistic, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("article_statistic_%v", id), func() (interface{}, error) {
		statistic, err := r.data.cc.GetArticleStatistic(ctx, &creationV1.GetArticleStatisticReq{
			Id:   id,
			Uuid: uuid,
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

func (r *articleRepo) GetUserArticleAgree(ctx context.Context, uuid string) (map[int32]bool, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("user_article_agree_"+uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetUserArticleAgree(ctx, &creationV1.GetUserArticleAgreeReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return reply.Agree, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[int32]bool), nil
}

func (r *articleRepo) GetUserArticleCollect(ctx context.Context, uuid string) (map[int32]bool, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("user_article_collect_"+uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetUserArticleCollect(ctx, &creationV1.GetUserArticleCollectReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return reply.Collect, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[int32]bool), nil
}

func (r *articleRepo) GetArticleListStatistic(ctx context.Context, articleList []*biz.Article) ([]*biz.ArticleStatistic, error) {
	ids := make([]int32, 0, len(articleList))
	for _, item := range articleList {
		ids = append(ids, item.Id)
	}
	statisticList, err := r.data.cc.GetArticleListStatistic(ctx, &creationV1.GetArticleListStatisticReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	reply := make([]*biz.ArticleStatistic, 0, len(statisticList.Count))
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
	result, err, _ := r.sg.Do(fmt.Sprintf("get_article_draft_list_%s", uuid), func() (interface{}, error) {
		draftList, err := r.data.cc.GetArticleDraftList(ctx, &creationV1.GetArticleDraftListReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.ArticleDraft, 0, len(draftList.Draft))
		for _, item := range draftList.Draft {
			reply = append(reply, &biz.ArticleDraft{
				Id: item.Id,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.ArticleDraft), nil
}

func (r *articleRepo) GetArticleSearch(ctx context.Context, page int32, search, time string) ([]*biz.Article, int32, error) {
	searchReply, err := r.data.cc.GetArticleSearch(ctx, &creationV1.GetArticleSearchReq{
		Page:   page,
		Search: search,
		Time:   time,
	})
	if err != nil {
		return nil, 0, err
	}
	reply := make([]*biz.Article, 0, len(searchReply.List))
	for _, item := range searchReply.List {
		reply = append(reply, &biz.Article{
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

func (r *articleRepo) GetArticleImageReview(ctx context.Context, page int32, uuid string) ([]*biz.CreationImageReview, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_article_image_review_%s_%v", uuid, page), func() (interface{}, error) {
		reviewReply, err := r.data.cc.GetArticleImageReview(ctx, &creationV1.GetArticleImageReviewReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.CreationImageReview, 0, len(reviewReply.Review))
		for _, item := range reviewReply.Review {
			reply = append(reply, &biz.CreationImageReview{
				Id:         item.Id,
				CreationId: item.CreationId,
				Kind:       item.Kind,
				Uid:        item.Uid,
				Uuid:       item.Uuid,
				CreateAt:   item.CreateAt,
				JobId:      item.JobId,
				Url:        item.Url,
				Label:      item.Label,
				Result:     item.Result,
				Score:      item.Score,
				Category:   item.Category,
				SubLabel:   item.SubLabel,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CreationImageReview), nil
}

func (r *articleRepo) GetArticleContentReview(ctx context.Context, page int32, uuid string) ([]*biz.CreationContentReview, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_article_content_review_%s_%v", uuid, page), func() (interface{}, error) {
		reviewReply, err := r.data.cc.GetArticleContentReview(ctx, &creationV1.GetArticleContentReviewReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.CreationContentReview, 0, len(reviewReply.Review))
		for _, item := range reviewReply.Review {
			reply = append(reply, &biz.CreationContentReview{
				Id:         item.Id,
				CreationId: item.CreationId,
				Title:      item.Title,
				Kind:       item.Kind,
				Uuid:       item.Uuid,
				CreateAt:   item.CreateAt,
				JobId:      item.JobId,
				Label:      item.Label,
				Result:     item.Result,
				Section:    item.Section,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CreationContentReview), nil
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
	result, err, _ := r.sg.Do(fmt.Sprintf("article_statistic_judge_%s_%v", uuid, id), func() (interface{}, error) {
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
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.ArticleStatisticJudge), nil
}

func (r *talkRepo) GetTalkList(ctx context.Context, page int32) ([]*biz.Talk, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("talk_page_%v", page), func() (interface{}, error) {
		talkList, err := r.data.cc.GetTalkList(ctx, &creationV1.GetTalkListReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Talk, 0, len(talkList.Talk))
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
		talkList, err := r.data.cc.GetTalkListHot(ctx, &creationV1.GetTalkListHotReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Talk, 0, len(talkList.Talk))
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
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_talk_list_%s_%v", uuid, page), func() (interface{}, error) {
		talkList, err := r.data.cc.GetUserTalkList(ctx, &creationV1.GetUserTalkListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Talk, 0, len(talkList.Talk))
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

func (r *talkRepo) GetUserTalkListVisitor(ctx context.Context, page int32, uuid string) ([]*biz.Talk, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_talk_by_visitor_%s_%v", uuid, page), func() (interface{}, error) {
		talkList, err := r.data.cc.GetUserTalkListVisitor(ctx, &creationV1.GetUserTalkListVisitorReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Talk, 0, len(talkList.Talk))
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
	ids := make([]int32, 0, len(talkList))
	for _, item := range talkList {
		ids = append(ids, item.Id)
	}
	statisticList, err := r.data.cc.GetTalkListStatistic(ctx, &creationV1.GetTalkListStatisticReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	reply := make([]*biz.TalkStatistic, 0, len(statisticList.Count))
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

func (r *talkRepo) GetTalkStatistic(ctx context.Context, id int32, uuid string) (*biz.TalkStatistic, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("talk_statistic_%v", id), func() (interface{}, error) {
		statistic, err := r.data.cc.GetTalkStatistic(ctx, &creationV1.GetTalkStatisticReq{
			Id:   id,
			Uuid: uuid,
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
	result, err, _ := r.sg.Do(fmt.Sprintf("get_last_talk_draft_%s", uuid), func() (interface{}, error) {
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
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.TalkDraft), nil
}

func (r *talkRepo) GetTalkSearch(ctx context.Context, page int32, search, time string) ([]*biz.Talk, int32, error) {
	searchReply, err := r.data.cc.GetTalkSearch(ctx, &creationV1.GetTalkSearchReq{
		Page:   page,
		Search: search,
		Time:   time,
	})
	if err != nil {
		return nil, 0, err
	}
	reply := make([]*biz.Talk, 0, len(searchReply.List))
	for _, item := range searchReply.List {
		reply = append(reply, &biz.Talk{
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

func (r *talkRepo) GetUserTalkAgree(ctx context.Context, uuid string) (map[int32]bool, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("user_talk_agree_"+uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetUserTalkAgree(ctx, &creationV1.GetUserTalkAgreeReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return reply.Agree, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[int32]bool), nil
}

func (r *talkRepo) GetUserTalkCollect(ctx context.Context, uuid string) (map[int32]bool, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("user_talk_collect_"+uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetUserTalkCollect(ctx, &creationV1.GetUserTalkCollectReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return reply.Collect, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[int32]bool), nil
}

func (r *talkRepo) GetTalkImageReview(ctx context.Context, page int32, uuid string) ([]*biz.CreationImageReview, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_talk_image_review_%s_%v", uuid, page), func() (interface{}, error) {
		reviewReply, err := r.data.cc.GetTalkImageReview(ctx, &creationV1.GetTalkImageReviewReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.CreationImageReview, 0, len(reviewReply.Review))
		for _, item := range reviewReply.Review {
			reply = append(reply, &biz.CreationImageReview{
				Id:         item.Id,
				CreationId: item.CreationId,
				Kind:       item.Kind,
				Uid:        item.Uid,
				Uuid:       item.Uuid,
				CreateAt:   item.CreateAt,
				JobId:      item.JobId,
				Url:        item.Url,
				Label:      item.Label,
				Result:     item.Result,
				Score:      item.Score,
				Category:   item.Category,
				SubLabel:   item.SubLabel,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CreationImageReview), nil
}

func (r *talkRepo) GetTalkContentReview(ctx context.Context, page int32, uuid string) ([]*biz.CreationContentReview, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_talk_content_review_%s_%v", uuid, page), func() (interface{}, error) {
		reviewReply, err := r.data.cc.GetTalkContentReview(ctx, &creationV1.GetTalkContentReviewReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.CreationContentReview, 0, len(reviewReply.Review))
		for _, item := range reviewReply.Review {
			reply = append(reply, &biz.CreationContentReview{
				Id:         item.Id,
				CreationId: item.CreationId,
				Title:      item.Title,
				Kind:       item.Kind,
				Uuid:       item.Uuid,
				CreateAt:   item.CreateAt,
				JobId:      item.JobId,
				Label:      item.Label,
				Result:     item.Result,
				Section:    item.Section,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CreationContentReview), nil
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
	result, err, _ := r.sg.Do(fmt.Sprintf("talk_statistic_judge_%s_%v", uuid, id), func() (interface{}, error) {
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
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.TalkStatisticJudge), nil
}

func (r *columnRepo) GetLastColumnDraft(ctx context.Context, uuid string) (*biz.ColumnDraft, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_last_column_draft_%s", uuid), func() (interface{}, error) {
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
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.ColumnDraft), nil
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

func (r *columnRepo) SubscribeColumn(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.SubscribeColumn(ctx, &creationV1.SubscribeColumnReq{
		Id:   id,
		Uuid: uuid,
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
		columnList, err := r.data.cc.GetColumnList(ctx, &creationV1.GetColumnListReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Column, 0, len(columnList.Column))
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
		columnList, err := r.data.cc.GetColumnListHot(ctx, &creationV1.GetColumnListHotReq{
			Page: page,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Column, 0, len(columnList.Column))
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
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_column_list_%s_%v", uuid, page), func() (interface{}, error) {
		columnList, err := r.data.cc.GetUserColumnList(ctx, &creationV1.GetUserColumnListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Column, 0, len(columnList.Column))
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

func (r *columnRepo) GetUserColumnListVisitor(ctx context.Context, page int32, uuid string) ([]*biz.Column, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_column_by_visitor_%s_%v", uuid, page), func() (interface{}, error) {
		columnList, err := r.data.cc.GetUserColumnListVisitor(ctx, &creationV1.GetUserColumnListVisitorReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Column, 0, len(columnList.Column))
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
	ids := make([]int32, 0, len(columnList))
	for _, item := range columnList {
		ids = append(ids, item.Id)
	}
	statisticList, err := r.data.cc.GetColumnListStatistic(ctx, &creationV1.GetColumnListStatisticReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	reply := make([]*biz.ColumnStatistic, 0, len(statisticList.Count))
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

func (r *columnRepo) GetColumnStatistic(ctx context.Context, id int32, uuid string) (*biz.ColumnStatistic, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("column_statistic_%v", id), func() (interface{}, error) {
		statistic, err := r.data.cc.GetColumnStatistic(ctx, &creationV1.GetColumnStatisticReq{
			Id:   id,
			Uuid: uuid,
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

func (r *columnRepo) GetSubscribeList(ctx context.Context, page int32, uuid string) ([]*biz.Column, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_subscribe_list_%s_%v", uuid, page), func() (interface{}, error) {
		subscribeList, err := r.data.cc.GetSubscribeList(ctx, &creationV1.GetSubscribeListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.Column, 0, len(subscribeList.Subscribe))
		for _, item := range subscribeList.Subscribe {
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
	subscribeList, err := r.data.cc.GetColumnSubscribes(ctx, &creationV1.GetColumnSubscribesReq{
		Ids:  ids,
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	reply := make([]*biz.Subscribe, 0, len(subscribeList.Subscribes))
	for _, item := range subscribeList.Subscribes {
		reply = append(reply, &biz.Subscribe{
			ColumnId: item.Id,
			Status:   item.SubscribeJudge,
		})
	}
	return reply, nil
}

func (r *columnRepo) GetColumnSearch(ctx context.Context, page int32, search, time string) ([]*biz.Column, int32, error) {
	searchReply, err := r.data.cc.GetColumnSearch(ctx, &creationV1.GetColumnSearchReq{
		Page:   page,
		Search: search,
		Time:   time,
	})
	if err != nil {
		return nil, 0, err
	}
	reply := make([]*biz.Column, 0, len(searchReply.List))
	for _, item := range searchReply.List {
		reply = append(reply, &biz.Column{
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

func (r *columnRepo) GetUserColumnAgree(ctx context.Context, uuid string) (map[int32]bool, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("user_column_agree_"+uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetUserColumnAgree(ctx, &creationV1.GetUserColumnAgreeReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return reply.Agree, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[int32]bool), nil
}

func (r *columnRepo) GetUserColumnCollect(ctx context.Context, uuid string) (map[int32]bool, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("user_column_collect_"+uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetUserColumnCollect(ctx, &creationV1.GetUserColumnCollectReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return reply.Collect, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[int32]bool), nil
}

func (r *columnRepo) GetUserSubscribeColumn(ctx context.Context, uuid string) (map[int32]bool, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("user_subscribe_column_"+uuid), func() (interface{}, error) {
		reply, err := r.data.cc.GetUserSubscribeColumn(ctx, &creationV1.GetUserSubscribeColumnReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return reply.Subscribe, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[int32]bool), nil
}

func (r *columnRepo) GetColumnImageReview(ctx context.Context, page int32, uuid string) ([]*biz.CreationImageReview, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_column_image_review_%s_%v", uuid, page), func() (interface{}, error) {
		reviewReply, err := r.data.cc.GetColumnImageReview(ctx, &creationV1.GetColumnImageReviewReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.CreationImageReview, 0, len(reviewReply.Review))
		for _, item := range reviewReply.Review {
			reply = append(reply, &biz.CreationImageReview{
				Id:         item.Id,
				CreationId: item.CreationId,
				Kind:       item.Kind,
				Uid:        item.Uid,
				Uuid:       item.Uuid,
				CreateAt:   item.CreateAt,
				JobId:      item.JobId,
				Url:        item.Url,
				Label:      item.Label,
				Result:     item.Result,
				Score:      item.Score,
				Category:   item.Category,
				SubLabel:   item.SubLabel,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CreationImageReview), nil
}

func (r *columnRepo) GetColumnContentReview(ctx context.Context, page int32, uuid string) ([]*biz.CreationContentReview, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_column_content_review_%s_%v", uuid, page), func() (interface{}, error) {
		reviewReply, err := r.data.cc.GetColumnContentReview(ctx, &creationV1.GetColumnContentReviewReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		reply := make([]*biz.CreationContentReview, 0, len(reviewReply.Review))
		for _, item := range reviewReply.Review {
			reply = append(reply, &biz.CreationContentReview{
				Id:         item.Id,
				CreationId: item.CreationId,
				Title:      item.Title,
				Kind:       item.Kind,
				Uuid:       item.Uuid,
				CreateAt:   item.CreateAt,
				JobId:      item.JobId,
				Label:      item.Label,
				Result:     item.Result,
				Section:    item.Section,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CreationContentReview), nil
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
	fmt.Println("before get news")
	newsList, err := r.data.cc.GetNews(ctx, &creationV1.GetNewsReq{
		Page: page,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("after get news")
	reply := make([]*biz.News, 0, len(newsList.News))
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
