package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type ArticleRepo interface {
	GetLastArticleDraft(ctx context.Context, uuid string) (*ArticleDraft, error)
	GetArticleDraftList(ctx context.Context, uuid string) ([]*ArticleDraft, error)
	GetArticleList(ctx context.Context, page int32) ([]*Article, error)
	GetArticleListHot(ctx context.Context, page int32) ([]*ArticleStatistic, error)
	GetArticleStatistic(ctx context.Context, id int32) (*ArticleStatistic, error)
	GetArticleListStatistic(ctx context.Context, ids []int32) ([]*ArticleStatistic, error)
	CreateArticle(ctx context.Context, uuid string, id int32) error
	CreateArticleStatistic(ctx context.Context, uuid string, id int32) error
	CreateArticleDraft(ctx context.Context, uuid string) (int32, error)
	CreateArticleFolder(ctx context.Context, id int32) error
	CreateArticleCache(ctx context.Context, uuid string, id int32) error
	CreateArticleSearch(ctx context.Context, uuid string, id int32) error
	DeleteArticleDraft(ctx context.Context, uuid string, id int32) error
	ArticleDraftMark(ctx context.Context, uuid string, id int32) error
	SendArticle(ctx context.Context, uuid string, id int32) (*ArticleDraft, error)
	SendDraftToMq(ctx context.Context, draft *ArticleDraft) error
	SendArticleToMq(ctx context.Context, article *Article, mode string) error
}
type ArticleUseCase struct {
	repo ArticleRepo
	tm   Transaction
	log  *log.Helper
}

func NewArticleUseCase(repo ArticleRepo, tm Transaction, logger log.Logger) *ArticleUseCase {
	return &ArticleUseCase{
		repo: repo,
		tm:   tm,
		log:  log.NewHelper(log.With(logger, "module", "creation/biz/articleUseCase")),
	}
}

func (r *ArticleUseCase) GetLastArticleDraft(ctx context.Context, uuid string) (*ArticleDraft, error) {
	draft, err := r.repo.GetLastArticleDraft(ctx, uuid)
	if kerrors.IsNotFound(err) {
		return nil, v1.ErrorRecordNotFound("last draft not found: %s", err.Error())
	}
	if err != nil {
		return nil, v1.ErrorGetArticleDraftFailed("get last draft failed: %s", err.Error())
	}
	return draft, nil
}

func (r *ArticleUseCase) CreateArticle(ctx context.Context, uuid string, id int32) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.DeleteArticleDraft(ctx, uuid, id)
		if err != nil {
			return v1.ErrorDeleteArticleDraftFailed("delete article draft failed: %s", err.Error())
		}

		err = r.repo.CreateArticle(ctx, uuid, id)
		if err != nil {
			return v1.ErrorCreateArticleFailed("create article failed: %s", err.Error())
		}

		err = r.repo.CreateArticleStatistic(ctx, uuid, id)
		if err != nil {
			return v1.ErrorCreateArticleStatisticFailed("create article statistic failed: %s", err.Error())
		}

		err = r.repo.SendArticleToMq(ctx, &Article{
			ArticleId: id,
			Uuid:      uuid,
		}, "create_article_cache_and_search")
		if err != nil {
			return v1.ErrorSendToMqFailed("create article to mq failed: %s", err.Error())
		}

		return nil
	})
}

func (r *ArticleUseCase) CreateArticleCacheAndSearch(ctx context.Context, uuid string, id int32) error {
	err := r.repo.CreateArticleCache(ctx, uuid, id)
	if err != nil {
		return v1.ErrorCreateArticleCacheFailed("create article cache failed: %s", err.Error())
	}

	err = r.repo.CreateArticleSearch(ctx, uuid, id)
	if err != nil {
		return v1.ErrorCreateArticleSearchFailed("create article search failed: %s", err.Error())
	}
	return nil
}

func (r *ArticleUseCase) CreateArticleDraft(ctx context.Context, uuid string) (int32, error) {
	var id int32
	err := r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		id, err = r.repo.CreateArticleDraft(ctx, uuid)
		if err != nil {
			return v1.ErrorCreateArticleDraftFailed("create article draft failed: %s", err.Error())
		}

		err = r.repo.CreateArticleFolder(ctx, id)
		if err != nil {
			return v1.ErrorCreateArticleFolderFailed("create article folder failed: %s", err.Error())
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ArticleUseCase) ArticleDraftMark(ctx context.Context, uuid string, id int32) error {
	err := r.repo.ArticleDraftMark(ctx, uuid, id)
	if err != nil {
		return v1.ErrorDraftMarkFailed("mark draft failed: %s", err.Error())
	}
	return nil
}

func (r *ArticleUseCase) GetArticleList(ctx context.Context, page int32) ([]*Article, error) {
	articleList, err := r.repo.GetArticleList(ctx, page)
	if err != nil {
		return nil, v1.ErrorGetArticleListFailed("get article list failed: %s", err.Error())
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetArticleListHot(ctx context.Context, page int32) ([]*ArticleStatistic, error) {
	articleList, err := r.repo.GetArticleListHot(ctx, page)
	if err != nil {
		return nil, v1.ErrorGetArticleListHotFailed("get article hot list failed: %s", err.Error())
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetArticleStatistic(ctx context.Context, id int32) (*ArticleStatistic, error) {
	statistic, err := r.repo.GetArticleStatistic(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetArticleStatisticFailed("get article statistic failed: %s", err.Error())
	}
	return statistic, nil
}

func (r *ArticleUseCase) GetArticleListStatistic(ctx context.Context, ids []int32) ([]*ArticleStatistic, error) {
	statisticList, err := r.repo.GetArticleListStatistic(ctx, ids)
	if err != nil {
		return nil, v1.ErrorGetArticleListStatisticFailed("get article list statistic failed: %s", err.Error())
	}
	return statisticList, nil
}

func (r *ArticleUseCase) GetArticleDraftList(ctx context.Context, uuid string) ([]*ArticleDraft, error) {
	draftList, err := r.repo.GetArticleDraftList(ctx, uuid)
	if err != nil {
		return draftList, v1.ErrorGetDraftListFailed("get article draft list failed: %s", err.Error())
	}
	return draftList, nil
}

func (r *ArticleUseCase) SendArticle(ctx context.Context, uuid string, id int32) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		draft, err := r.repo.SendArticle(ctx, uuid, id)
		if err != nil {
			return v1.ErrorSendArticleFailed("send article failed: %s", err.Error())
		}
		err = r.repo.SendDraftToMq(ctx, draft)
		if err != nil {
			return v1.ErrorSendToMqFailed("send draft to mq failed: %s", err.Error())
		}
		return nil
	})
}
