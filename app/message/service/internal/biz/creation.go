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
	CreateArticleDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error
	EditArticleCosAndSearch(ctx context.Context, id, auth int32, uuid string) error
	DeleteArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error
	SetArticleViewDbAndCache(ctx context.Context, id int32, uuid string) error
	SetArticleAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	SetArticleCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error
	CancelArticleAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	CancelArticleCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error

	ToReviewCreateTalk(id int32, uuid string) error
	ToReviewEditTalk(id int32, uuid string) error
	TalkCreateReviewPass(ctx context.Context, id, auth int32, uuid string) error
	CreateTalkDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error
	TalkEditReviewPass(ctx context.Context, id, auth int32, uuid string) error
	EditTalkCosAndSearch(ctx context.Context, id, auth int32, uuid string) error
	DeleteTalkCacheAndSearch(ctx context.Context, id int32, uuid string) error
	SetTalkViewDbAndCache(ctx context.Context, id int32, uuid string) error
	SetTalkAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	CancelTalkAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	SetTalkCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error
	CancelTalkCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error

	ToReviewCreateColumn(id int32, uuid string) error
	ToReviewEditColumn(id int32, uuid string) error
	ColumnCreateReviewPass(ctx context.Context, id, auth int32, uuid string) error
	ColumnEditReviewPass(ctx context.Context, id, auth int32, uuid string) error
	CreateColumnDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error
	EditColumnCosAndSearch(ctx context.Context, id, auth int32, uuid string) error
	DeleteColumnCacheAndSearch(ctx context.Context, id int32, uuid string) error
	SetColumnAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	SetColumnViewDbAndCache(ctx context.Context, id int32, uuid string) error
	CancelColumnAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	SetColumnCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error
	CancelColumnCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	AddColumnIncludesDbAndCache(ctx context.Context, id, articleId int32, uuid string) error
	DeleteColumnIncludesDbAndCache(ctx context.Context, id, articleId int32, uuid string) error
	SetColumnSubscribeDbAndCache(ctx context.Context, id int32, uuid string) error
	CancelColumnSubscribeDbAndCache(ctx context.Context, id int32, uuid string) error

	ToReviewCreateCollections(id int32, uuid string) error
	ToReviewEditCollections(id int32, uuid string) error
	CollectionsCreateReviewPass(ctx context.Context, id, auth int32, uuid string) error
	CollectionsEditReviewPass(ctx context.Context, id, auth int32, uuid string) error
	CreateCollectionsDbAndCache(ctx context.Context, id, auth int32, uuid string) error
	EditCollectionsCos(ctx context.Context, id, auth int32, uuid string) error
	DeleteCollectionsCache(ctx context.Context, id int32, uuid string) error

	AddCreationComment(ctx context.Context, createId, createType int32, uuid string)
	GetCreationUser(ctx context.Context, uuid string) (int32, int32, int32, error)
}

type CreationUseCase struct {
	repo CreationRepo
	jwt  Jwt
	log  *log.Helper
}

func NewCreationUseCase(repo CreationRepo, jwt Jwt, logger log.Logger) *CreationUseCase {
	return &CreationUseCase{
		repo: repo,
		jwt:  jwt,
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
	var token, id, auths string
	var ok bool

	if token, ok = tr.CosHeaders["X-Cos-Meta-Token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", tr)
		return nil
	}

	if id, ok = tr.CosHeaders["X-Cos-Meta-Id"]; !ok || id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	if auths, ok = tr.CosHeaders["X-Cos-Meta-Auth"]; !ok || auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
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
	var token, id, auths string
	var ok bool

	if token, ok = tr.CosHeaders["X-Cos-Meta-Token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", tr)
		return nil
	}

	if id, ok = tr.CosHeaders["X-Cos-Meta-Id"]; !ok || id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	if auths, ok = tr.CosHeaders["X-Cos-Meta-Auth"]; !ok || auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
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

func (r *CreationUseCase) CreateArticleDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.CreateArticleDbCacheAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) EditArticleCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.EditArticleCosAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) DeleteArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return r.repo.DeleteArticleCacheAndSearch(ctx, id, uuid)
}

func (r *CreationUseCase) SetArticleViewDbAndCache(ctx context.Context, id int32, uuid string) error {
	return r.repo.SetArticleViewDbAndCache(ctx, id, uuid)
}

func (r *CreationUseCase) SetArticleAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.SetArticleAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CreationUseCase) CancelArticleAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.CancelArticleAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CreationUseCase) SetArticleCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	return r.repo.SetArticleCollectDbAndCache(ctx, id, collectionsId, uuid, userUuid)
}

func (r *CreationUseCase) CancelArticleCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.CancelArticleCollectDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CreationUseCase) ToReviewCreateTalk(id int32, uuid string) error {
	return r.repo.ToReviewCreateTalk(id, uuid)
}

func (r *CreationUseCase) ToReviewEditTalk(id int32, uuid string) error {
	return r.repo.ToReviewEditTalk(id, uuid)
}

func (r *CreationUseCase) TalkCreateReview(ctx context.Context, tr *TextReview) error {
	var err error
	var token, id, auths string
	var ok bool

	if token, ok = tr.CosHeaders["X-Cos-Meta-Token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", tr)
		return nil
	}

	if id, ok = tr.CosHeaders["X-Cos-Meta-Id"]; !ok || id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	if auths, ok = tr.CosHeaders["X-Cos-Meta-Auth"]; !ok || auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
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
	var token, id, auths string
	var ok bool

	if token, ok = tr.CosHeaders["X-Cos-Meta-Token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", tr)
		return nil
	}

	if id, ok = tr.CosHeaders["X-Cos-Meta-Id"]; !ok || id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	if auths, ok = tr.CosHeaders["X-Cos-Meta-Auth"]; !ok || auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
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

func (r *CreationUseCase) CreateTalkDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.CreateTalkDbCacheAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) EditTalkCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.EditTalkCosAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) DeleteTalkCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return r.repo.DeleteTalkCacheAndSearch(ctx, id, uuid)
}

func (r *CreationUseCase) SetTalkViewDbAndCache(ctx context.Context, id int32, uuid string) error {
	return r.repo.SetTalkViewDbAndCache(ctx, id, uuid)
}

func (r *CreationUseCase) SetTalkAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.SetTalkAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CreationUseCase) CancelTalkAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.CancelTalkAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CreationUseCase) SetTalkCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	return r.repo.SetTalkCollectDbAndCache(ctx, id, collectionsId, uuid, userUuid)
}

func (r *CreationUseCase) CancelTalkCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.CancelTalkCollectDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CreationUseCase) ToReviewCreateColumn(id int32, uuid string) error {
	return r.repo.ToReviewCreateColumn(id, uuid)
}

func (r *CreationUseCase) ToReviewEditColumn(id int32, uuid string) error {
	return r.repo.ToReviewEditColumn(id, uuid)
}

func (r *CreationUseCase) ColumnCreateReview(ctx context.Context, tr *TextReview) error {
	var err error
	var token, id, auths string
	var ok bool

	if token, ok = tr.CosHeaders["X-Cos-Meta-Token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", tr)
		return nil
	}

	if id, ok = tr.CosHeaders["X-Cos-Meta-Id"]; !ok || id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	if auths, ok = tr.CosHeaders["X-Cos-Meta-Auth"]; !ok || auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
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
	var token, id, auths string
	var ok bool

	if token, ok = tr.CosHeaders["X-Cos-Meta-Token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", tr)
		return nil
	}

	if id, ok = tr.CosHeaders["X-Cos-Meta-Id"]; !ok || id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	if auths, ok = tr.CosHeaders["X-Cos-Meta-Auth"]; !ok || auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
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

func (r *CreationUseCase) CreateColumnDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.CreateColumnDbCacheAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) EditColumnCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.EditColumnCosAndSearch(ctx, id, auth, uuid)
}

func (r *CreationUseCase) DeleteColumnCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return r.repo.DeleteColumnCacheAndSearch(ctx, id, uuid)
}

func (r *CreationUseCase) SetColumnViewDbAndCache(ctx context.Context, id int32, uuid string) error {
	return r.repo.SetColumnViewDbAndCache(ctx, id, uuid)
}

func (r *CreationUseCase) SetColumnAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.SetColumnAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CreationUseCase) CancelColumnAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.CancelColumnAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CreationUseCase) SetColumnCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	return r.repo.SetColumnCollectDbAndCache(ctx, id, collectionsId, uuid, userUuid)
}

func (r *CreationUseCase) CancelColumnCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.CancelColumnCollectDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CreationUseCase) AddColumnIncludesDbAndCache(ctx context.Context, id, articleId int32, uuid string) error {
	return r.repo.AddColumnIncludesDbAndCache(ctx, id, articleId, uuid)
}

func (r *CreationUseCase) DeleteColumnIncludesDbAndCache(ctx context.Context, id, articleId int32, uuid string) error {
	return r.repo.DeleteColumnIncludesDbAndCache(ctx, id, articleId, uuid)
}

func (r *CreationUseCase) SetColumnSubscribeDbAndCache(ctx context.Context, id int32, uuid string) error {
	return r.repo.SetColumnSubscribeDbAndCache(ctx, id, uuid)
}

func (r *CreationUseCase) CancelColumnSubscribeDbAndCache(ctx context.Context, id int32, uuid string) error {
	return r.repo.CancelColumnSubscribeDbAndCache(ctx, id, uuid)
}

func (r *CreationUseCase) CollectionsCreateReview(ctx context.Context, tr *TextReview) error {
	var err error
	var token, id, auths string
	var ok bool

	if token, ok = tr.CosHeaders["X-Cos-Meta-Token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", tr)
		return nil
	}

	if id, ok = tr.CosHeaders["X-Cos-Meta-Id"]; !ok || id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	if auths, ok = tr.CosHeaders["X-Cos-Meta-Auth"]; !ok || auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
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
		r.log.Info("collections create review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.CollectionsCreateReviewPass(ctx, int32(aid), int32(auth), uuid)
	} else {
		r.log.Info("collections create review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CreationUseCase) CollectionsEditReview(ctx context.Context, tr *TextReview) error {
	var err error
	var token, id, auths string
	var ok bool

	if token, ok = tr.CosHeaders["X-Cos-Meta-Token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", tr)
		return nil
	}

	if id, ok = tr.CosHeaders["X-Cos-Meta-Id"]; !ok || id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	if auths, ok = tr.CosHeaders["X-Cos-Meta-Auth"]; !ok || auths == "" {
		r.log.Info("auth not exist，%v", tr)
		return nil
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
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
		r.log.Info("collections edit review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.CollectionsEditReviewPass(ctx, int32(aid), int32(auth), uuid)
	} else {
		r.log.Info("collections edit review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CreationUseCase) CreateCollectionsDbAndCache(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.CreateCollectionsDbAndCache(ctx, id, auth, uuid)
}

func (r *CreationUseCase) EditCollectionsCos(ctx context.Context, id, auth int32, uuid string) error {
	return r.repo.EditCollectionsCos(ctx, id, auth, uuid)
}

func (r *CreationUseCase) DeleteCollectionsCache(ctx context.Context, id int32, uuid string) error {
	return r.repo.DeleteCollectionsCache(ctx, id, uuid)
}

func (r *CreationUseCase) AddCreationComment(ctx context.Context, createId, createType int32, uuid string) {
	r.repo.AddCreationComment(ctx, createId, createType, uuid)
}

func (r *CreationUseCase) ToReviewCreateCollections(id int32, uuid string) error {
	return r.repo.ToReviewCreateCollections(id, uuid)
}

func (r *CreationUseCase) ToReviewEditCollections(id int32, uuid string) error {
	return r.repo.ToReviewEditCollections(id, uuid)
}
