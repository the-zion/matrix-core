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
	ArticleCreateReviewPass(ctx context.Context, id, auth int32, uuid string) error
	ArticleEditReviewPass(ctx context.Context, id, auth int32, uuid string) error
	CreateArticleCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error
	EditArticleCosAndSearch(ctx context.Context, id, auth int32, uuid string) error
	DeleteArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error

	ToReviewCreateTalk(id int32, uuid string) error
	ToReviewEditTalk(id int32, uuid string) error
	TalkCreateReviewPass(ctx context.Context, id, auth int32, uuid string) error
	CreateTalkCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error
	TalkEditReviewPass(ctx context.Context, id, auth int32, uuid string) error
	EditTalkCosAndSearch(ctx context.Context, id, auth int32, uuid string) error
	DeleteTalkCacheAndSearch(ctx context.Context, id int32, uuid string) error

	ToReviewCreateColumn(id int32, uuid string) error
	ToReviewEditColumn(id int32, uuid string) error
	ColumnCreateReviewPass(ctx context.Context, id, auth int32, uuid string) error
	ColumnEditReviewPass(ctx context.Context, id, auth int32, uuid string) error
	CreateColumnCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error
	EditColumnCosAndSearch(ctx context.Context, id, auth int32, uuid string) error
	DeleteColumnCacheAndSearch(ctx context.Context, id int32, uuid string) error
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

	auths := tr.CosHeaders["X-Cos-Meta-Auth"]
	if auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	aid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	auth, err := strconv.ParseInt(auths, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	if tr.State != "Success" {
		r.log.Info("article create review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.ArticleCreateReviewPass(ctx, int32(aid), int32(auth), uuid)
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

	auths := tr.CosHeaders["X-Cos-Meta-Auth"]
	if auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	aid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	auth, err := strconv.ParseInt(auths, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	if tr.State != "Success" {
		r.log.Info("article edit review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.ArticleEditReviewPass(ctx, int32(aid), int32(auth), uuid)
	} else {
		r.log.Info("article edit review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CreationUseCase) CreateArticleCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.CreateArticleCacheAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) EditArticleCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.EditArticleCosAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) DeleteArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return r.repo.DeleteArticleCacheAndSearch(ctx, id, uuid)
}

func (r *CreationUseCase) ToReviewCreateTalk(id int32, uuid string) error {
	return r.repo.ToReviewCreateTalk(id, uuid)
}

func (r *CreationUseCase) ToReviewEditTalk(id int32, uuid string) error {
	return r.repo.ToReviewEditTalk(id, uuid)
}

func (r *CreationUseCase) TalkCreateReview(ctx context.Context, tr *TextReview) error {
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

	auths := tr.CosHeaders["X-Cos-Meta-Auth"]
	if auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	aid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	auth, err := strconv.ParseInt(auths, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	if tr.State != "Success" {
		r.log.Info("talk create review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.TalkCreateReviewPass(ctx, int32(aid), int32(auth), uuid)
	} else {
		r.log.Info("talk create review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CreationUseCase) TalkEditReview(ctx context.Context, tr *TextReview) error {
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

	auths := tr.CosHeaders["X-Cos-Meta-Auth"]
	if auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	aid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	auth, err := strconv.ParseInt(auths, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	if tr.State != "Success" {
		r.log.Info("talk edit review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.TalkEditReviewPass(ctx, int32(aid), int32(auth), uuid)
	} else {
		r.log.Info("talk edit review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CreationUseCase) CreateTalkCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.CreateTalkCacheAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) EditTalkCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.EditTalkCosAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) DeleteTalkCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return r.repo.DeleteTalkCacheAndSearch(ctx, id, uuid)
}

func (r *CreationUseCase) ToReviewCreateColumn(id int32, uuid string) error {
	return r.repo.ToReviewCreateColumn(id, uuid)
}

func (r *CreationUseCase) ToReviewEditColumn(id int32, uuid string) error {
	return r.repo.ToReviewEditColumn(id, uuid)
}

func (r *CreationUseCase) ColumnCreateReview(ctx context.Context, tr *TextReview) error {
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

	auths := tr.CosHeaders["X-Cos-Meta-Auth"]
	if auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	aid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	auth, err := strconv.ParseInt(auths, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	if tr.State != "Success" {
		r.log.Info("column create review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.ColumnCreateReviewPass(ctx, int32(aid), int32(auth), uuid)
	} else {
		r.log.Info("column create review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CreationUseCase) ColumnEditReview(ctx context.Context, tr *TextReview) error {
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

	auths := tr.CosHeaders["X-Cos-Meta-Auth"]
	if auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	aid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	auth, err := strconv.ParseInt(auths, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	if tr.State != "Success" {
		r.log.Info("column edit review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.ColumnEditReviewPass(ctx, int32(aid), int32(auth), uuid)
	} else {
		r.log.Info("column edit review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CreationUseCase) CreateColumnCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.CreateColumnCacheAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) EditColumnCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.EditColumnCosAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) DeleteColumnCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return r.repo.DeleteColumnCacheAndSearch(ctx, id, uuid)
}
