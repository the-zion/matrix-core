package biz

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strconv"
)

type CreationRepo interface {
	ToReviewArticle(mode string, msg *primitive.MessageExt) error
	ArticleDraftReviewPass(ctx context.Context, uuid string, id int32) error
	CreateArticleCacheAndSearch(ctx context.Context, uuid string, id int32) error
}

type CreationUseCase struct {
	repo CreationRepo
	log  *log.Helper
}

func NewCreationUseCase(repo CreationRepo, logger log.Logger) *CreationUseCase {
	return &CreationUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "message/biz/creationUseCase")),
	}
}

func (r *CreationUseCase) ToReviewArticleDraft(msg *primitive.MessageExt) error {
	return r.repo.ToReviewArticle("article_draft", msg)
}

func (r *CreationUseCase) ArticleDraftReview(ctx context.Context, tr *TextReview) error {
	var err error
	uuid := tr.CosHeaders["X-Cos-Meta-Uuid"]
	if uuid == "" {
		r.log.Info("uuid not exist，%v", tr)
		return nil
	}

	id := tr.CosHeaders["X-Cos-Meta-Id"]
	if id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	aid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	if tr.State != "Success" {
		r.log.Info("article draft review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.ArticleDraftReviewPass(ctx, uuid, int32(aid))
	} else {
		r.log.Info("article draft review not pass，%v", tr)
		//err = r.repo.ProfileReviewNotPass(ctx, tr.CosHeaders["x-cos-meta-uuid"])
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CreationUseCase) CreateArticleCacheAndSearch(ctx context.Context, uuid string, id int32) error {
	return r.repo.CreateArticleCacheAndSearch(ctx, uuid, id)
}
