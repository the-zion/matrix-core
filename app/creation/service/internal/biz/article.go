package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type ArticleRepo interface {
	GetLastArticleDraft(ctx context.Context, uuid string) (*ArticleDraft, error)
	GetArticle(ctx context.Context, id int32) (*Article, error)
	GetArticleList(ctx context.Context, page int32) ([]*Article, error)
	GetArticleDraftList(ctx context.Context, uuid string) ([]*ArticleDraft, error)
	GetArticleAgreeJudge(ctx context.Context, id int32, uuid string) (bool, error)
	GetArticleCollectJudge(ctx context.Context, id int32, uuid string) (bool, error)
	GetArticleListHot(ctx context.Context, page int32) ([]*ArticleStatistic, error)
	GetColumnArticleList(ctx context.Context, id int32) ([]*Article, error)
	GetArticleCount(ctx context.Context, uuid string) (int32, error)
	GetArticleCountVisitor(ctx context.Context, uuid string) (int32, error)
	GetUserArticleList(ctx context.Context, page int32, uuid string) ([]*Article, error)
	GetUserArticleListVisitor(ctx context.Context, page int32, uuid string) ([]*Article, error)
	GetArticleStatistic(ctx context.Context, id int32) (*ArticleStatistic, error)
	GetArticleListStatistic(ctx context.Context, ids []int32) ([]*ArticleStatistic, error)

	CreateArticle(ctx context.Context, id, auth int32, uuid string) error
	CreateArticleStatistic(ctx context.Context, id, auth int32, uuid string) error
	CreateArticleDraft(ctx context.Context, uuid string) (int32, error)
	CreateArticleFolder(ctx context.Context, id int32, uuid string) error
	CreateArticleCache(ctx context.Context, id, auth int32, uuid string) error
	CreateArticleSearch(ctx context.Context, id int32, uuid string) error

	EditArticleCos(ctx context.Context, id int32, uuid string) error
	EditArticleSearch(ctx context.Context, id int32, uuid string) error
	UpdateArticleCache(ctx context.Context, id, auth int32, uuid string) error

	DeleteArticle(ctx context.Context, id int32, uuid string) error
	DeleteArticleStatistic(ctx context.Context, id int32, uuid string) error
	DeleteArticleCache(ctx context.Context, id int32, uuid string) error
	DeleteArticleSearch(ctx context.Context, id int32, uuid string) error
	DeleteArticleDraft(ctx context.Context, id int32, uuid string) error

	FreezeArticleCos(ctx context.Context, id int32, uuid string) error
	ArticleDraftMark(ctx context.Context, id int32, uuid string) error

	SendArticle(ctx context.Context, id int32, uuid string) (*ArticleDraft, error)
	SendReviewToMq(ctx context.Context, review *ArticleReview) error
	SendArticleToMq(ctx context.Context, article *Article, mode string) error
	SendArticleStatisticToMq(ctx context.Context, uuid, mode string) error

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

func (r *ArticleUseCase) CreateArticle(ctx context.Context, id, auth int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.DeleteArticleDraft(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateArticleFailed("delete article draft failed: %s", err.Error())
		}

		err = r.repo.CreateArticle(ctx, id, auth, uuid)
		if err != nil {
			return v1.ErrorCreateArticleFailed("create article failed: %s", err.Error())
		}

		err = r.repo.CreateArticleStatistic(ctx, id, auth, uuid)
		if err != nil {
			return v1.ErrorCreateArticleFailed("create article statistic failed: %s", err.Error())
		}

		err = r.repo.SendArticleToMq(ctx, &Article{
			ArticleId: id,
			Uuid:      uuid,
			Auth:      auth,
		}, "create_article_cache_and_search")
		if err != nil {
			return v1.ErrorCreateArticleFailed("create article to mq failed: %s", err.Error())
		}

		return nil
	})
}

func (r *ArticleUseCase) EditArticle(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.SendArticleToMq(ctx, &Article{
		ArticleId: id,
		Auth:      auth,
		Uuid:      uuid,
	}, "edit_article_cos_and_search")
	if err != nil {
		return v1.ErrorEditArticleFailed("edit article to mq failed: %s", err.Error())
	}
	return nil
}

func (r *ArticleUseCase) DeleteArticle(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.DeleteArticle(ctx, id, uuid)
		if err != nil {
			return v1.ErrorDeleteArticleFailed("delete article failed: %s", err.Error())
		}

		err = r.repo.DeleteArticleStatistic(ctx, id, uuid)
		if err != nil {
			return v1.ErrorDeleteArticleFailed("delete article statistic failed: %s", err.Error())
		}

		err = r.repo.SendArticleToMq(ctx, &Article{
			ArticleId: id,
			Uuid:      uuid,
		}, "delete_article_cache_and_search")
		if err != nil {
			return v1.ErrorEditArticleFailed("delete article to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) CreateArticleCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.CreateArticleCache(ctx, id, auth, uuid)
	if err != nil {
		return v1.ErrorCreateArticleFailed("create article cache failed: %s", err.Error())
	}

	if auth == 2 {
		return nil
	}

	err = r.repo.CreateArticleSearch(ctx, id, uuid)
	if err != nil {
		return v1.ErrorCreateArticleFailed("create article search failed: %s", err.Error())
	}
	return nil
}

func (r *ArticleUseCase) EditArticleCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.UpdateArticleCache(ctx, id, auth, uuid)
	if err != nil {
		return v1.ErrorEditArticleFailed("edit article cache failed: %s", err.Error())
	}

	err = r.repo.EditArticleCos(ctx, id, uuid)
	if err != nil {
		return v1.ErrorEditArticleFailed("edit article cache failed: %s", err.Error())
	}

	err = r.repo.EditArticleSearch(ctx, id, uuid)
	if err != nil {
		return v1.ErrorEditArticleFailed("edit article search failed: %s", err.Error())
	}
	return nil
}

func (r *ArticleUseCase) DeleteArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	err := r.repo.DeleteArticleCache(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteArticleFailed("delete article cache failed: %s", err.Error())
	}

	err = r.repo.FreezeArticleCos(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteArticleFailed("freeze article cos failed: %s", err.Error())
	}

	err = r.repo.DeleteArticleSearch(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteArticleFailed("delete article search failed: %s", err.Error())
	}
	return nil
}

func (r *ArticleUseCase) CreateArticleDraft(ctx context.Context, uuid string) (int32, error) {
	var id int32
	err := r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		id, err = r.repo.CreateArticleDraft(ctx, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("create article draft failed: %s", err.Error())
		}

		err = r.repo.CreateArticleFolder(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("create article folder failed: %s", err.Error())
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
		return nil, v1.ErrorGetArticleListFailed("get article hot list failed: %s", err.Error())
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetColumnArticleList(ctx context.Context, id int32) ([]*Article, error) {
	articleList, err := r.repo.GetColumnArticleList(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetArticleListFailed("get column article list failed: %s", err.Error())
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetArticleCount(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetArticleCount(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get article count failed: %s", err.Error())
	}
	return count, nil
}

func (r *ArticleUseCase) GetArticleCountVisitor(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetArticleCountVisitor(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get article count visitor failed: %s", err.Error())
	}
	return count, nil
}

func (r *ArticleUseCase) GetUserArticleList(ctx context.Context, page int32, uuid string) ([]*Article, error) {
	articleList, err := r.repo.GetUserArticleList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetArticleListFailed("get user article list failed: %s", err.Error())
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetUserArticleListVisitor(ctx context.Context, page int32, uuid string) ([]*Article, error) {
	articleList, err := r.repo.GetUserArticleListVisitor(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetArticleListFailed("get user article list visitor failed: %s", err.Error())
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetArticleStatistic(ctx context.Context, id int32) (*ArticleStatistic, error) {
	statistic, err := r.repo.GetArticleStatistic(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetStatisticFailed("get article statistic failed: %s", err.Error())
	}
	return statistic, nil
}

func (r *ArticleUseCase) GetArticleListStatistic(ctx context.Context, ids []int32) ([]*ArticleStatistic, error) {
	statisticList, err := r.repo.GetArticleListStatistic(ctx, ids)
	if err != nil {
		return nil, v1.ErrorGetStatisticFailed("get article list statistic failed: %s", err.Error())
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
			return v1.ErrorCreateArticleFailed("send article failed: %s", err.Error())
		}
		err = r.repo.SendReviewToMq(ctx, &ArticleReview{
			Uuid: draft.Uuid,
			Id:   draft.Id,
			Mode: "create",
		})
		if err != nil {
			return v1.ErrorCreateArticleFailed("send create review to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) SendArticleEdit(ctx context.Context, id int32, uuid string) error {
	article, err := r.repo.GetArticle(ctx, id)
	if err != nil {
		return v1.ErrorGetArticleFailed("get article failed: %s", err.Error())
	}

	if article.Uuid != uuid {
		return v1.ErrorEditArticleFailed("send article edit failed: no auth")
	}

	err = r.repo.SendReviewToMq(ctx, &ArticleReview{
		Uuid: article.Uuid,
		Id:   article.ArticleId,
		Mode: "edit",
	})
	if err != nil {
		return v1.ErrorEditArticleFailed("send edit review to mq failed: %s", err.Error())
	}
	return nil
}

func (r *ArticleUseCase) SetArticleAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetArticleAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set article agree failed: %s", err.Error())
		}
		err = r.repo.SetArticleAgreeToCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set article agree to cache failed: %s", err.Error())
		}
		err = r.repo.SendArticleStatisticToMq(ctx, uuid, "agree")
		if err != nil {
			return v1.ErrorSetAgreeFailed("set article agree to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) SetArticleView(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetArticleView(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetViewFailed("set article view failed: %s", err.Error())
		}
		err = r.repo.SetArticleViewToCache(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetViewFailed("set article view to cache failed: %s", err.Error())
		}
		err = r.repo.SendArticleStatisticToMq(ctx, uuid, "view")
		if err != nil {
			return v1.ErrorSetViewFailed("set article view to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) SetArticleCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetArticleUserCollect(ctx, id, collectionsId, userUuid)
		if err != nil {
			return v1.ErrorSetCollectFailed("set article user collect failed: %s", err.Error())
		}
		err = r.repo.SetArticleCollect(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetCollectFailed("set article collect failed: %s", err.Error())
		}
		err = r.repo.SetArticleCollectToCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorSetCollectFailed("set article collect to cache failed: %s", err.Error())
		}
		err = r.repo.SendArticleStatisticToMq(ctx, uuid, "collect")
		if err != nil {
			return v1.ErrorSetCollectFailed("set article collect to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) CancelArticleAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelArticleAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel article agree failed: %s", err.Error())
		}
		err = r.repo.CancelArticleAgreeFromCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel article agree from cache failed: %s", err.Error())
		}
		err = r.repo.SendArticleStatisticToMq(ctx, uuid, "agree_cancel")
		if err != nil {
			return v1.ErrorCancelAgreeFailed("set article agree to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) CancelArticleCollect(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelArticleUserCollect(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorCancelCollectFailed("cancel article user collect failed: %s", err.Error())
		}
		err = r.repo.CancelArticleCollect(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCancelCollectFailed("cancel article collect failed: %s", err.Error())
		}
		err = r.repo.CancelArticleCollectFromCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorCancelCollectFailed("cancel article collect to cache failed: %s", err.Error())
		}
		err = r.repo.SendArticleStatisticToMq(ctx, uuid, "collect_cancel")
		if err != nil {
			return v1.ErrorCancelCollectFailed("set article collect to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ArticleUseCase) ArticleStatisticJudge(ctx context.Context, id int32, uuid string) (*ArticleStatisticJudge, error) {
	agree, err := r.repo.GetArticleAgreeJudge(ctx, id, uuid)
	if err != nil {
		return nil, v1.ErrorGetStatisticJudgeFailed("get article statistic judge failed: %s", err.Error())
	}
	collect, err := r.repo.GetArticleCollectJudge(ctx, id, uuid)
	if err != nil {
		return nil, v1.ErrorGetStatisticJudgeFailed("get article statistic judge failed: %s", err.Error())
	}
	return &ArticleStatisticJudge{
		Agree:   agree,
		Collect: collect,
	}, nil
}
