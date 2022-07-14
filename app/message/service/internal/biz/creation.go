package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strconv"
)

type CreationRepo interface {
	ToReviewCreateArticle(id int32, uuid string) error
	ToReviewEditArticle(id int32, uuid string) error
	ArticleCreateReviewPass(ctx context.Context, id int32, uuid string) error
	ArticleEditReviewPass(ctx context.Context, id int32, uuid string) error
	CreateArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error
	EditArticleCosAndSearch(ctx context.Context, id int32, uuid string) error
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

func (r *CreationUseCase) ToReviewCreateArticle(id int32, uuid string) error {
	return r.repo.ToReviewCreateArticle(id, uuid)
}

func (r *CreationUseCase) ToReviewEditArticle(id int32, uuid string) error {
	return r.repo.ToReviewEditArticle(id, uuid)
}

func (r *CreationUseCase) ArticleCreateReview(ctx context.Context, tr *TextReview) error {
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
		r.log.Info("article create review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.ArticleCreateReviewPass(ctx, int32(aid), uuid)
	} else {
		r.log.Info("article create review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CreationUseCase) ArticleEditReview(ctx context.Context, tr *TextReview) error {
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
		r.log.Info("article edit review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.ArticleEditReviewPass(ctx, int32(aid), uuid)
	} else {
		r.log.Info("article edit review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CreationUseCase) CreateArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return r.repo.CreateArticleCacheAndSearch(ctx, id, uuid)
}

func (r *CreationUseCase) EditArticleCosAndSearch(ctx context.Context, id int32, uuid string) error {
	return r.repo.EditArticleCosAndSearch(ctx, id, uuid)
}
