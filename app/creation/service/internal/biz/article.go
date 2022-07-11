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
	GetArticleAgreeJudge(ctx context.Context, id int32, uuid string) (bool, error)
	GetArticleCollectJudge(ctx context.Context, id int32, uuid string) (bool, error)
	GetArticleListHot(ctx context.Context, page int32) ([]*ArticleStatistic, error)
	GetArticleStatistic(ctx context.Context, id int32) (*ArticleStatistic, error)
	GetArticleListStatistic(ctx context.Context, ids []int32) ([]*ArticleStatistic, error)
	CreateArticle(ctx context.Context, id int32, uuid string) error
	CreateArticleStatistic(ctx context.Context, id int32, uuid string) error
	CreateArticleDraft(ctx context.Context, uuid string) (int32, error)
	CreateArticleFolder(ctx context.Context, id int32) error
	CreateArticleCache(ctx context.Context, id int32, uuid string) error
	CreateArticleSearch(ctx context.Context, id int32, uuid string) error
	DeleteArticleDraft(ctx context.Context, id int32, uuid string) error
	ArticleDraftMark(ctx context.Context, id int32, uuid string) error
	SendArticle(ctx context.Context, id int32, uuid string) (*ArticleDraft, error)
	SendDraftToMq(ctx context.Context, draft *ArticleDraft) error
	SendArticleToMq(ctx context.Context, article *Article, mode string) error
	SetArticleAgree(ctx context.Context, id int32, uuid string) error
	SetArticleView(ctx context.Context, id int32, uuid string) error
	SetArticleUserCollect(ctx context.Context, id, collectionsId int32, userUuid string) error
	SetArticleCollect(ctx context.Context, id int32, uuid string) error
	SetArticleAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error
	SetArticleViewToCache(ctx context.Context, id int32, uuid string) error
	SetArticleCollectToCache(ctx context.Context, id int32, uuid, userUuid string) error
	CancelArticleAgree(ctx context.Context, id int32, uuid string) error
	CancelArticleAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error
	CancelArticleUserCollect(ctx context.Context, id int32, userUuid string) error
	CancelArticleCollect(ctx context.Context, id int32, uuid string) error
	CancelArticleCollectFromCache(ctx context.Context, id int32, uuid, userUuid string) error
	SendArticleStatisticToMq(ctx context.Context, uuid, mode string) error
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

func (r *ArticleUseCase) CreateArticle(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.DeleteArticleDraft(ctx, id, uuid)
		if err != nil {
			return v1.ErrorDeleteArticleDraftFailed("delete article draft failed: %s", err.Error())
		}

		err = r.repo.CreateArticle(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateArticleFailed("create article failed: %s", err.Error())
		}

		err = r.repo.CreateArticleStatistic(ctx, id, uuid)
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

func (r *ArticleUseCase) CreateArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	err := r.repo.CreateArticleCache(ctx, id, uuid)
	if err != nil {
		return v1.ErrorCreateArticleCacheFailed("create article cache failed: %s", err.Error())
	}

	err = r.repo.CreateArticleSearch(ctx, id, uuid)
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

func (r *ArticleUseCase) ArticleDraftMark(ctx context.Context, id int32, uuid string) error {
	err := r.repo.ArticleDraftMark(ctx, id, uuid)
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

func (r *ArticleUseCase) SendArticle(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		draft, err := r.repo.SendArticle(ctx, id, uuid)
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

func (r *ArticleUseCase) SetArticleAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetArticleAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetArticleAgreeFailed("set article agree failed: %s", err.Error())
		}
		err = r.repo.SetArticleAgreeToCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorSetArticleAgreeFailed("set article agree to cache failed: %s", err.Error())
		}
		err = r.repo.SendArticleStatisticToMq(ctx, uuid, "agree")
		if err != nil {
			return v1.ErrorSendToMqFailed("set article agree to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) SetArticleView(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetArticleView(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetArticleViewFailed("set article view failed: %s", err.Error())
		}
		err = r.repo.SetArticleViewToCache(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetArticleViewFailed("set article view to cache failed: %s", err.Error())
		}
		err = r.repo.SendArticleStatisticToMq(ctx, uuid, "view")
		if err != nil {
			return v1.ErrorSendToMqFailed("set article view to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) SetArticleCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetArticleUserCollect(ctx, id, collectionsId, userUuid)
		if err != nil {
			return v1.ErrorSetArticleViewFailed("set article user collect failed: %s", err.Error())
		}
		err = r.repo.SetArticleCollect(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetArticleViewFailed("set article collect failed: %s", err.Error())
		}
		err = r.repo.SetArticleCollectToCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorSetArticleViewFailed("set article collect to cache failed: %s", err.Error())
		}
		err = r.repo.SendArticleStatisticToMq(ctx, uuid, "collect")
		if err != nil {
			return v1.ErrorSendToMqFailed("set article collect to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) CancelArticleAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelArticleAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCancelArticleAgreeFailed("cancel article agree failed: %s", err.Error())
		}
		err = r.repo.CancelArticleAgreeFromCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorCancelArticleAgreeFailed("cancel article agree from cache failed: %s", err.Error())
		}
		err = r.repo.SendArticleStatisticToMq(ctx, uuid, "agree_cancel")
		if err != nil {
			return v1.ErrorSendToMqFailed("set article agree to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) CancelArticleCollect(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelArticleUserCollect(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorCancelArticleCollectFailed("cancel article user collect failed: %s", err.Error())
		}
		err = r.repo.CancelArticleCollect(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCancelArticleCollectFailed("cancel article collect failed: %s", err.Error())
		}
		err = r.repo.CancelArticleCollectFromCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorCancelArticleCollectFailed("cancel article collect to cache failed: %s", err.Error())
		}
		err = r.repo.SendArticleStatisticToMq(ctx, uuid, "collect_cancel")
		if err != nil {
			return v1.ErrorSendToMqFailed("set article collect to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) ArticleStatisticJudge(ctx context.Context, id int32, uuid string) (*ArticleStatisticJudge, error) {
	agree, err := r.repo.GetArticleAgreeJudge(ctx, id, uuid)
	if err != nil {
		return nil, v1.ErrorArticleStatisticJudgeFailed("get article statistic judge failed: %s", err.Error())
	}
	collect, err := r.repo.GetArticleCollectJudge(ctx, id, uuid)
	if err != nil {
		return nil, v1.ErrorArticleStatisticJudgeFailed("get article statistic judge failed: %s", err.Error())
	}
	return &ArticleStatisticJudge{
		Agree:   agree,
		Collect: collect,
	}, nil
}
